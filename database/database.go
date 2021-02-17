package database

import (
	"fmt"
	"os"

	"github.com/zorchenhimer/hacker-quotes/models"
)

type DbType string
const (
	DB_Json DbType = "json"
	DB_PostgresSQL DbType = "pgsql"
	DB_SQLite DbType = "sqlite"
)

type DB interface {
	AddAdjective(word models.Adjective) error
	AddNoun(word models.Noun) error
	AddVerb(word models.Verb) error
	//AddPronoun(word models.Verb) error

	RemoveAdjective(id int) error
	RemoveNoun(id int) error
	RemoveVerb(id int) error
	//RemovePronoun(id int) error

	GetAdjectiveIds() ([]int, error)
	GetNounIds(begin, end, alone bool) ([]int, error)
	GetVerbIds() ([]int, error)
	GetPronounIds(plural bool) ([]int, error)
	GetSentenceIds() ([]int, error)

	GetAdjective(id int) (*models.Adjective, error)
	GetNoun(id int) (*models.Noun, error)
	GetVerb(id int) (*models.Verb, error)
	GetPronoun(id int) (*models.Pronoun, error)
	GetSentence(id int) (string, error)

	InitData([]models.Adjective, []models.Noun, []models.Verb, []models.Pronoun, []string) error
	IsNew() bool
	Close()
}

type dbInit func(connectionString string) (DB, error)
var registered map[DbType]dbInit

func New(databaseType DbType, connectionString string) (DB, error) {
	f, ok := registered[databaseType]
	if !ok {
		return nil, fmt.Errorf("Unregistered database type: %s", databaseType)
	}

	return f(connectionString)
}

func register(databaseType DbType, initFunc dbInit) {
	if registered == nil {
		registered = make(map[DbType]dbInit)
	}
	if _, exists := registered[databaseType]; exists {
		panic(fmt.Sprintf("Unable to register database with type %s: already exists.", databaseType))
	}

	registered[databaseType] = initFunc
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
