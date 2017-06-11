package main

import (
	"github.com/zither-oss/zither-backend/player"
)

func main() {
	remote, err := player.Connect("localhost", "6600")
	if err != nil {
		panic(err)
	}
	remote.Add(remote.ItemByURI("spotify:track:4uLU6hMCjMI75M1A2tKUQC"))
	remote.Play()
}
