package main

import (
	"fmt"
	"time"

	"github.com/andybrewer/mack"
	"github.com/cloudinary/cloudinary-go"
	"github.com/hugolgst/rich-go/client"
)

type ConfigCloud struct {
	CloudId          string `json:"cloudId"`
	CloudToken       string `json:"cloudToken"`
	CloudTokenSecret string `json:"cloudTokenSecret"`
}

var cld *cloudinary.Cloudinary

func main() {
	var config Config
	config.loadConfig()
	setup()
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
				go player.setRpc()
				lastState, lastTitle = player.State, player.Title
				fmt.Printf("------------------------------------------\nNow listening to : %s\nOn the album     : %s\nBy               : %s \n------------------------------------------\n", player.Title, player.Album, player.Artist)
			} else if player.State != lastState {
				go player.setRpc()
				lastState, lastTitle = player.State, player.Title
				fmt.Printf("State Change now : %s\n", player.State)
			}
			time.Sleep(time.Second * 2)
		}
	}
}
