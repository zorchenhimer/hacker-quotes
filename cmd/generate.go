package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/zorchenhimer/hacker-quotes"
	"github.com/zorchenhimer/hacker-quotes/database"
	"github.com/zorchenhimer/hacker-quotes/files"
)

func main() {
	var count int
	var format string

	flag.IntVar(&count, "c", 1, "Number of sentences to generate")
	flag.StringVar(&format, "f", "", "Custom format to use when generating sentences")
	flag.Parse()

	db, err := database.New("sqlite", "file:db.sqlite?mode=memory")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hq, err := hacker.NewEnglish(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := bytes.NewReader(files.WordLists)
	err = hq.InitData(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("")
	for i := 0; i < count; i++ {
		var sentence string
		var err error

		if format != "" {
			sentence, err = hq.HackThis(format)
		} else {
			sentence, err = hq.Hack()
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(">>", sentence, "<<")
	}
	fmt.Println("")
}
