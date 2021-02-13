package hacker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/zorchenhimer/hacker-quotes/database"
)

type settings struct {
	DatabaseType database.DbType
	ConnectionString string

	HttpAddr string
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

	s := &settings{}
	if err = json.Unmarshal(raw, s); err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %s", err)
	}

	return s, nil
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
