package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed scripts_vault/vol.scpt
var volScript []byte

//go:embed scripts_vault/art.scpt
var artScript []byte

//go:embed scripts_vault/music.scpt
var musicScript []byte

func setup() {
	if _, err := os.Stat("scripts/art.scpt"); err == nil {
		fmt.Printf("Skipping setup\n")
		return
	} else {
		fmt.Printf("Launching setup\n")
		os.Mkdir("scripts", 0755)
		os.Create("scripts/art.scpt")
		os.Create("scripts/music.scpt")
		os.Create("scripts/vol.scpt")
		os.WriteFile("scripts/art.scpt", artScript, 0744)
		os.WriteFile("scripts/music.scpt", musicScript, 0744)
		os.WriteFile("scripts/vol.scpt", volScript, 0744)
	}
}

func getPlayer() Player {
	var player Player

	cmd := exec.Command("osascript", "./scripts/music.scpt")

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

func getVol() string {
	cmd := exec.Command("osascript", "./scripts/vol.scpt")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}

	return strings.Split(strings.Split(out.String(), ",")[0], ":")[1]
}
