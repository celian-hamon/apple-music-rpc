package main

import (
	"embed"
	"encoding/json"

	"github.com/andybrewer/mack"
)

//go:embed config.json
var jsonFile embed.FS

type Config struct {
	ConfigCloud  ConfigCloud `json:"cloudConfig"`
	DiscordAppId string      `json:"discordAppId"`
	Setup        bool        `json:"setup"`
}

func (Config *Config) loadConfig() {

	bytes, _ := jsonFile.ReadFile("config.json")

	err := json.Unmarshal(bytes, &Config)
	if err != nil {
		alert := mack.AlertOptions{
			Title:   "Apple Music RPC",
			Message: "Config file not found",
			Buttons: "OK",
		}
		mack.AlertBox(alert)
	}
}
