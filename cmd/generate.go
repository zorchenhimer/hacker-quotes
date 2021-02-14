package main

import (
	"os"
	"fmt"

	"github.com/zorchenhimer/hacker-quotes"
	"github.com/zorchenhimer/hacker-quotes/database"
)

func main() {
	fmt.Println("len(os.Args):", len(os.Args))

	var count int = 1
	if len(os.Args) == 2 {
		fmt.Sscanf(os.Args[1], "%d", &count)
	}

	db, err := database.New("sqlite", "file:db.sqlite?mode=memory")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hq, err := hacker.NewGeneric(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = hq.InitData("word_lists.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if count == 1 {
		sentence, err := hq.Hack()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("\n>>", sentence, "<<\n")
	} else {
		fmt.Println("")
		for i := 0; i < count; i++ {
			sentence, err := hq.Hack()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(">>", sentence, "<<")
		}
	}
}
