package main

import (
	"fmt"
	"os"

	"github.com/zorchenhimer/hacker-quotes"
)

func main() {
	server, err := hacker.New("settings.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server.Hack()
}
