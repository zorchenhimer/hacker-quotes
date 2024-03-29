package hacker

import (
	"io"

	"github.com/zorchenhimer/hacker-quotes/models"
)

type HackerQuotes interface {
	// Hack returns a completely randomized quote.
	Hack() (string, error)

	// Format returns a quote in the given format.
	HackThis(format string) (string, error)

	// InitData populates the underlying database with data from the given json file.
	InitData(reader io.Reader) error
}

type Admin interface {
	AddNoun(word models.Noun) error
	AddVerb(word models.Verb) error

	RemoveNoun(word string) error
	// Word is the indefinite form.
	RemoveVerb(word string) error

	GetNouns() ([]models.Noun, error)
	GetVerbs() ([]models.Verb, error)
}
