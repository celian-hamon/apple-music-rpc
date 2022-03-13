package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/andybrewer/mack"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin/search"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/flytam/filenamify"
	"github.com/hugolgst/rich-go/client"
)

type Config struct {
	ConfigCloud  ConfigCloud `json:"cloudConfig"`
	DiscordAppId string      `json:"discordAppId"`
}

type ConfigCloud struct {
	CloudId          string `json:"cloudId"`
	CloudToken       string `json:"cloudToken"`
	CloudTokenSecret string `json:"cloudTokenSecret"`
}

type Player struct {
	State    string
	Artist   string
	ID       string
	Title    string
	Album    string
	Duration string
	Position string
	Kind     string
	AlbumUrl string
}

//go:embed "music.scpt"
var musicScript string

func getPlayer() Player {
	var player Player

	cmd := exec.Command("osascript", "-e", musicScript)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}

	content := strings.Split(out.String(), ", ")
	player.State = content[0]
	player.Artist = content[1]
	player.Title = content[2]
	player.Album = content[3]
	player.Duration = content[5]
	player.Position = content[6]
	player.Kind = content[7]
	player.ID = content[8]

	return player
}

//go:embed "art.scpt"
var artScript string

func (player *Player) getArtwork() {

	cmd := exec.Command("osascript", "-e", artScript)

	cmd.Run()

	var resp *uploader.UploadResult
	output, _ := filenamify.Filenamify(player.Album, filenamify.Options{})
	output = strings.Replace(output, " ", "_", -1)
	os.Rename("tmp.jpg", output+".jpg")

	resp, err := cld.Upload.Upload(context.Background(), output+".jpg", uploader.UploadParams{
		UseFilename:    true,
		UniqueFilename: false,
	})

	if err != nil {
		fmt.Println(err)
	}
	player.AlbumUrl = resp.SecureURL
	os.Remove(output + ".jpg")
}

//go:embed config.json
var jsonFile embed.FS

func loadConfig() Config {

	bytes, _ := jsonFile.ReadFile("config.json")
	var payload Config
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		alert := mack.AlertOptions{
			Title:   "Apple Music RPC",
			Message: "Config file not found",
			Buttons: "OK",
		}
		mack.AlertBox(alert)
	}

	return payload
}

//go:embed "vol.scpt"
var volScript string

func getVol() string {
	cmd := exec.Command("osascript", "-e", volScript)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
	liste := strings.Split(out.String(), ",")

	return liste[0][len(liste[0])-2:]
}

var cld *cloudinary.Cloudinary

func main() {
	config := loadConfig()
	lastState := ""
	lastTitle := ""
	err := client.Login(config.DiscordAppId)
	if err != nil {
		alert := mack.AlertOptions{
			Title:   "Apple Music RPC",
			Message: "Cant connect to discord",
			Buttons: "OK",
		}
		mack.AlertBox(alert)
		panic(err)
	}

	cld, err = cloudinary.NewFromParams(config.ConfigCloud.CloudId, config.ConfigCloud.CloudToken, config.ConfigCloud.CloudTokenSecret)
	if err != nil {
		alert := mack.AlertOptions{
			Title:   "Apple Music RPC",
			Message: "Cant connect to cloud",
			Buttons: "OK",
		}
		mack.AlertBox(alert)
		panic(err)
	}

	mack.Notify("Service up and Running", "Apple Music RPC")
	for {
		player := getPlayer()
		if player.State == "playing" || player.State == "paused" {
			if player.Title != lastTitle {
				fmt.Printf("------------------------------------------\nNow listening to : %s\nOn the album     : %s\nBy               : %s \n------------------------------------------\n", player.Title, player.Album, player.Artist)
				go setRpc(player)
				lastState = player.State
				lastTitle = player.Title
			} else if player.State != lastState {
				go setRpc(player)
				lastState = player.State
				lastTitle = player.Title
				fmt.Printf("State Change now : %s\n", player.State)
			}
			time.Sleep(time.Second * 1)
		}
	}
}

func setRpc(player Player) {
	playerPos, _ := strconv.ParseFloat(player.Position, 64)
	playerDur, _ := strconv.ParseFloat(player.Duration, 64)

	now := time.Now().Add(time.Second * -time.Duration(playerPos))
	end := time.Now().Add(time.Second * time.Duration(playerDur-playerPos))

	outVol := getVol()

	searchQuery := search.Query{
		Expression: "resource_type:image",
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: 30,
	}

	searchResult, _ := cld.Admin.Search(context.Background(), searchQuery)

	filename, _ := filenamify.Filenamify(player.Album, filenamify.Options{})
	filename = strings.Replace(filename, " ", "_", -1)

	for _, asset := range searchResult.Assets {
		fmt.Println("Searching for : '" + filename + "'")
		if strings.Contains(asset.PublicID, filename) {
			player.AlbumUrl = asset.SecureURL
			break
		}
	}

	if player.AlbumUrl == "" {
		player.getArtwork()
	}

	var cover string
	if player.AlbumUrl != "" {
		cover = player.AlbumUrl
	}

	if player.State == "playing" {
		client.SetActivity(client.Activity{
			Details:    player.State + " " + player.Title,
			State:      "by " + player.Artist + " on " + player.Album,
			LargeImage: cover,
			LargeText:  player.Title,
			SmallImage: "https://media.giphy.com/media/3o6gDP9oLOGtBMMBSU/giphy.gif",
			SmallText:  player.State,
			Timestamps: &client.Timestamps{
				Start: &now,
				End:   &end,
			},
			Buttons: []*client.Button{
				{
					Label: "(: Volume : " + outVol + " :)",
					Url:   "https://github.com/celian-hamon",
				},
				{
					Label: "(: made by cece :)",
					Url:   "https://github.com/celian-hamon",
				},
			},
		})

	} else if player.State == "paused" {
		now = time.Now()
		client.SetActivity(client.Activity{
			Details:    player.State + " " + player.Title,
			State:      "by " + player.Artist + " on " + player.Album,
			LargeImage: cover,
			LargeText:  player.Title,
			SmallImage: "https://media.giphy.com/media/UUQ7nCPmqF9mY04Co0/giphy.gif",
			SmallText:  player.State,
			Timestamps: &client.Timestamps{
				Start: &now,
			},
		})
	} else {
		client.SetActivity(client.Activity{
			Details: player.State,
		})
	}
}
