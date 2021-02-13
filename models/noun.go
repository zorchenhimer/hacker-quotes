package models

import (
	"strings"
)

type Noun struct {
	Id int

	Multiple bool

	Begin bool
	End bool
	Alone bool

	Regular bool

	Word string
}

func (n Noun) Plural() string {
	suffixes := []string{
		"s",
		"x",
		"sh",
		"ch",
		"ss",
	}

	for _, sfx := range suffixes {
		if strings.HasSuffix(n.Word, sfx) {
			return n.Word + "es"
		}
	}

	if strings.HasSuffix(n.Word, "y") && !strings.ContainsAny(string(n.Word[len(n.Word)-2]), "aeiou") {
		return n.Word[:len(n.Word)-1] + "ies"
	}

	return n.Word + "s"
}

