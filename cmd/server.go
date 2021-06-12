package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/zorchenhimer/hacker-quotes"
	"github.com/zorchenhimer/hacker-quotes/database"
	"github.com/zorchenhimer/hacker-quotes/frontend"
)

//go:embed settings_default.json
var defaultSettings []byte

func main() {
	var s *settings
	var err error

	if fileExists("settings.json") {
		fmt.Println("settings.json found")
		s, err = loadSettings("settings.json")
	} else {
		fmt.Println("settings.json not found.  Using default values.")
		s, err = loadRawSettings(defaultSettings)
	}

	if err != nil {
		fmt.Printf("Unable to load settings: %s\n", err)
		os.Exit(1)
	}

	db, err := database.New(s.DatabaseType, s.ConnectionString)
	if err != nil {
		fmt.Printf("Unable to load database type %s: %s\n", s.DatabaseType, err)
		os.Exit(1)
	}

	hack, err := hacker.NewEnglish(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if db.IsNew() {
		fmt.Println("database is new")
		err = hack.InitData("word_lists.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("database isn't new")
	}

	web, err := frontend.New(hack)
	if err != nil {
		fmt.Printf("Unable to load frontend: %s\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	//mux.Handle("/api", api)
	mux.Handle("/", web)

	hs := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error running HTTP server:", err)
	}
}

type settings struct {
	DatabaseType database.DbType
	ConnectionString string

	HttpAddr string
}

func loadRawSettings(raw []byte) (*settings, error) {
	s := &settings{}
	if err := json.Unmarshal(raw, s); err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %s", err)
	}

	return s, nil
}

func loadSettings(filename string) (*settings, error) {
	if !fileExists(filename) {
		return nil, fmt.Errorf("%q doesn't exist", filename)
		//return &settings{
		//	HttpAddr: ":8080",
		//}, nil
	}

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s", err)
	}

	return loadRawSettings(raw)
}

// fileExists returns whether the given file or directory exists or not.
// Taken from https://stackoverflow.com/a/10510783
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
