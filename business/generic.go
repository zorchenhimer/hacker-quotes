package business

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/zorchenhimer/hacker-quotes/database"
	"github.com/zorchenhimer/hacker-quotes/models"
)

type generic struct {
	db database.DB
}

func NewGeneric(db database.DB) (HackerQuotes, error) {
	return &generic{db: db}, nil
}

func (g *generic) Random() (string, error) {
	definite := rand.Int() % 2 == 0
	hasAdj := rand.Int() % 2 == 0
	plural := rand.Int() % 2 == 0

	np, err := g.nounPhrase(definite, hasAdj, plural)
	if err != nil {
		return "", err
	}

	//fmt.Printf("(%s) definite: %t; hasAdj: %t; plural: %t\n", np, definite, hasAdj, plural)

	sb := strings.Builder{}
	sb.WriteString(np)

	return sb.String(), nil
}

func (g *generic) Format(format string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

func (g *generic) nounPhrase(definite, hasAdj, plural bool) (string, error){
	adj := ""
	var err error
	if hasAdj {
		adj, err = g.randomAdjective()
		if err != nil {
			return "", err
		}
	}

	noun, err := g.randomNoun(plural)
	if err != nil {
		return "", err
	}

	phrase := adj
	if phrase != "" {
		phrase += " " + noun
	} else {
		phrase = noun
	}
	
	if definite && !plural {
		//fmt.Println("[nounPhrase] definite && !plural")
		return "the " + phrase, nil
	}

	if !plural {
		//fmt.Println("[nounPhrase] !plural")
		return ana(phrase), nil
	}

	return phrase, nil
}

func (g *generic) randomAdjective() (string, error) {
	ids, err := g.db.GetAdjectiveIds()
	if err != nil {
		return "", fmt.Errorf("[adj] get IDs error: %v", err)
	}

	if len(ids) <= 0 {
		return "", fmt.Errorf("No adjective IDs returned from database")
	}

	rid := int(rand.Int63n(int64(len(ids))))
	//fmt.Printf("[adj] len(ids): %d; rid: %d; %d\n", len(ids), rid, ids[rid])

	adj, err := g.db.GetAdjective(ids[rid])
	if err != nil {
		return "", fmt.Errorf("[adj] ID: %d; %v", ids[rid], err)
	}

	return adj.Word, nil
}

func (g *generic) randomNoun(plural bool) (string, error) {
	ids, err := g.db.GetNounIds()
	if err != nil {
		return "", fmt.Errorf("[noun] get IDs error: %v", err)
	}

	if len(ids) <= 0 {
		return "", fmt.Errorf("No noun IDs returned from database")
	}

	rid := int(rand.Int63n(int64(len(ids))))
	//fmt.Printf("[noun] len(ids): %d; rid: %d; ID: %d\n", len(ids), rid, ids[rid])

	noun, err := g.db.GetNoun(ids[rid])
	if err != nil {
		return "", fmt.Errorf("[noun] ID: %d; %v", ids[rid], err)
	}

	if plural {
		return noun.Plural(), nil
	}
	return noun.Word, nil
}

func (g *generic) randomVerb() (string, error) {
	ids, err := g.db.GetVerbIds()
	if err != nil {
		return "", fmt.Errorf("[verb] get IDs error: %v", err)
	}

	if len(ids) <= 0 {
		return "", fmt.Errorf("No verb IDs returned from database")
	}

	rid := int(rand.Int63n(int64(len(ids))))
	verb, err := g.db.GetVerb(ids[rid])
	if err != nil {
		return "", fmt.Errorf("[verb] ID: %d; %v", ids[rid], err)
	}

	return verb.Word, nil
}

func (g *generic) InitData(filename string) error {
	fmt.Printf("Initializing database with data in %q\n", filename)
	if g.db == nil {
		return fmt.Errorf("databse is nil!")
	}

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	data := map[string][][]string{}
	if err = json.Unmarshal(raw, &data); err != nil {
		return err
	}

	rawadj, ok := data["adjectives"]
	if !ok {
		return fmt.Errorf("Missing adjectives key in data")
	}


	adjectives := []models.Adjective{}
	for _, adj := range rawadj {
		t, word := adj[0], adj[1]
		a := models.Adjective{Word: word}
		if strings.Contains(t, "a") {
			a.Absolute = true
		}
		if strings.Contains(t, "e") {
			a.AppendEst = true
		}
		if strings.Contains(t, "m") {
			a.AppendMore = true
		}

		adjectives = append(adjectives, a)
	}

	rawnoun, ok := data["nouns"]
	if !ok {
		return fmt.Errorf("Missing nouns key in data")
	}

	nouns := []models.Noun{}
	for _, noun := range rawnoun {
		t, word := noun[0], noun[1]
		n := models.Noun{Word: word}

		if strings.Contains(t, "m") {
			n.Multiple = true
		}

		if strings.Contains(t, "b") {
			n.Begin = true
		}

		if strings.Contains(t, "e") {
			n.End = true
		}

		if strings.Contains(t, "a") {
			n.Alone = true
		}

		if strings.Contains(t, "r") {
			n.Regular = true
		}

		nouns = append(nouns, n)
	}

	rawverbs, ok := data["verbs"]
	if !ok {
		return fmt.Errorf("Missing verbs key in data")
	}

	verbs := []models.Verb{}
	for _, word := range rawverbs {
		v := models.Verb{Word: word[1]}
		if strings.Contains(word[0], "r") {
			v.Regular = true
		}

		verbs = append(verbs, v)
	}

	return g.db.InitData(adjectives, nouns, verbs, nil)
}

// Prepend "a", "an" or nothing to a phrase
func ana(phrase string) string {
	//fmt.Printf("[ana] phrase[0]: %s; %q\n", string(phrase[0]), phrase)
	if strings.ContainsAny(string(phrase[0]), "aeiou") {
		return "an " + phrase
	}

	return "a " + phrase
}
