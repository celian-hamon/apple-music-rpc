package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/api/admin/search"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/flytam/filenamify"
	"github.com/hugolgst/rich-go/client"
)

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

func (player *Player) setRpc() {
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
			State:      "by " + player.Artist,
			LargeImage: cover,
			LargeText:  player.Album,
			SmallImage: "https://media.giphy.com/media/3o6gDP9oLOGtBMMBSU/giphy.gif",
			SmallText:  player.State,
			Timestamps: &client.Timestamps{
				Start: &now,
				End:   &end,
			},
			Buttons: []*client.Button{
				{
					Label: "Volume : " + outVol,
					Url:   "https://github.com/celian-hamon/apple-music-rpc",
				},
			},
		})

	} else if player.State == "paused" {
		now = time.Now()
		client.SetActivity(client.Activity{
			Details:    player.State + " " + player.Title,
			State:      "by " + player.Artist,
			LargeImage: cover,
			LargeText:  player.Album,
			SmallImage: "https://media.giphy.com/media/UUQ7nCPmqF9mY04Co0/giphy.gif",
			SmallText:  player.State,
			Timestamps: &client.Timestamps{
				Start: &now,
			},
			Buttons: []*client.Button{
				{
					Label: "Volume : " + outVol,
					Url:   "https://github.com/celian-hamon/apple-music-rpc",
				},
			},
		})
	} else {
		client.SetActivity(client.Activity{
			Details: player.State,
		})
	}
}

func (player *Player) getArtwork() {

	cmd := exec.Command("osascript", "./scripts/art.scpt")

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
