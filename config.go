package main

import (
	"embed"
	"encoding/json"

	"github.com/andybrewer/mack"
)

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
