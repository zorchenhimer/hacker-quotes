package hacker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/zorchenhimer/hacker-quotes/database"
	"github.com/zorchenhimer/hacker-quotes/models"
)

type english struct {
	db database.DB
}

func NewEnglish(db database.DB) (HackerQuotes, error) {
	return &english{db: db}, nil
}

/*

	Sentence format

	{word_type:options}
	{{word_type:new word:properties}}

	{pronoun} can't {verb:i,present} {noun_phrase}, it {verb:it,future} {noun_phrase}!
	{verb:you,present} {noun_phrase:definite}, than you can {verb:you,present} {noun_phrase:definite}!
	{noun_phrase} {verb}. With {noun_phrase:indifinite,noadj,compound}!
*/

func (g *english) Hack() (string, error) {
	sb := strings.Builder{}

	invert := rand.Int() % 2 == 0
	plural := rand.Int() % 2 == 0

	pn, err := g.randomPronoun(plural)
	if err != nil {
		return "", err
	}

	sb.WriteString(pn)
	sb.WriteString(" can't ")

	v, err := g.randomVerb(models.CT_I, models.CM_Present, invert)
	if err != nil {
		return "", err
	}

	sb.WriteString(v)
	sb.WriteString(" ")

	definite := rand.Int() % 2 == 0
	hasAdj := rand.Int() % 2 == 0
	plural = rand.Int() % 2 == 0
	compound := rand.Int() % 2 == 0

	np, err := g.nounPhrase(definite, hasAdj, plural, compound)
	if err != nil {
		return "", err
	}

	sb.WriteString(np)
	sb.WriteString(", it ")

	v2, err := g.randomVerb(models.CT_It, models.CM_Future, invert)
	if err != nil {
		return "", err
	}

	sb.WriteString(v2)
	sb.WriteString(" ")

	definite = rand.Int() % 2 == 0
	hasAdj = rand.Int() % 2 == 0
	plural = rand.Int() % 2 == 0
	compound = rand.Int() % 2 == 0

	np2, err := g.nounPhrase(definite, hasAdj, plural, compound)
	if err != nil {
		return "", err
	}

	sb.WriteString(np2)
	sb.WriteString("!")

	return toCap(sb.String()), nil
}

func (g *english) Hack_t1() (string, error) {
	sb := strings.Builder{}
	invert := false

	v, err := g.randomVerb(models.CT_You, models.CM_Present, invert)
	if err != nil {
		return "", err
	}

	sb.WriteString(toCap(v))
	sb.WriteString(" ")

	hasAdj := rand.Int() % 2 == 0
	plural := rand.Int() % 2 == 0
	compound := rand.Int() % 2 == 0

	np, err := g.nounPhrase(true, hasAdj, plural, compound)
	if err != nil {
		return "", err
	}

	sb.WriteString(np)
	sb.WriteString(", then you can ")

	v2, err := g.randomVerb(models.CT_You, models.CM_Present, invert)
	if err != nil {
		return "", err
	}

	sb.WriteString(v2)
	sb.WriteString(" ")

	hasAdj = rand.Int() % 2 == 0
	plural = rand.Int() % 2 == 0
	compound = rand.Int() % 2 == 0

	np2, err := g.nounPhrase(true, hasAdj, plural, compound)
	if err != nil {
		return "", err
	}
	sb.WriteString(np2)
	sb.WriteString("!")

	return sb.String(), err
}

func (g *english) Hack_t0() (string, error) {
	definite := rand.Int() % 2 == 0
	hasAdj := rand.Int() % 2 == 0
	plural := rand.Int() % 2 == 0
	compound := rand.Int() % 2 == 0

	np, err := g.nounPhrase(definite, hasAdj, plural, compound)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	sb.WriteString(toCap(np))

	ctime := models.CM_Present
	ctype := models.CT_It
	invert := false	// TODO: implement this

	if plural {
		ctype = models.CT_They
	}

	v, err := g.randomVerb(ctype, ctime, invert)
	if err != nil {
		return "", err
	}

	sb.WriteString(" ")
	sb.WriteString(v)

	definite = rand.Int() % 2 == 0
	hasAdj = rand.Int() % 2 == 0
	plural = rand.Int() % 2 == 0

	np2, err := g.nounPhrase(definite, hasAdj, plural, false)
	if err != nil {
		return "", err
	}

	sb.WriteString(" ")
	sb.WriteString(np2)
	sb.WriteString(". With ")

	plural = rand.Int() % 2 == 0

	np3, err := g.nounPhrase(false, false, plural, true)
	if err != nil {
		return "", err
	}
	sb.WriteString(np3)
	sb.WriteString("!")

	return sb.String(), nil
}

func (g *english) Format(format string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

func (g *english) nounPhrase(definite, hasAdj, plural, compound bool) (string, error){
	adj := ""
	var err error
	if hasAdj {
		adj, err = g.randomAdjective()
		if err != nil {
			return "", err
		}
	}

	noun, err := g.randomNoun(plural, compound)
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

func (g *english) randomAdjective() (string, error) {
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

func (g *english) randomNoun(plural, compound bool) (string, error) {
	var ids []int
	var err error
	if compound {
		ids, err = g.db.GetNounIds(true, true, true)
		if err != nil {
			return "", fmt.Errorf("[noun] get IDs error: %v", err)
		}
	} else {
		ids, err = g.db.GetNounIds(true, false, false)
		if err != nil {
			return "", fmt.Errorf("[noun] get IDs error: %v", err)
		}
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

func (g *english) randomVerb(ctype models.ConjugationType, ctime models.ConjugationTime, invert bool) (string, error) {
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

	return verb.Conjugate(ctype, ctime, invert), nil
}

func (g *english) randomPronoun(plural bool) (string, error) {
	ids, err := g.db.GetPronounIds(plural)
	if err != nil {
		return "", fmt.Errorf("[pronoun] get IDs error: %v", err)
	}

	if len(ids) <= 0 {
		return "", fmt.Errorf("No pronoun IDs returned from database")
	}

	rid := int(rand.Int63n(int64(len(ids))))
	pronoun, err := g.db.GetPronoun(ids[rid])
	if err != nil {
		return "", fmt.Errorf("[pronoun] ID: %d; %v", ids[rid], err)
	}

	return pronoun.Word, nil
}

func (g *english) InitData(filename string) error {
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

	rawpronouns, ok := data["pronouns"]
	if !ok {
		return fmt.Errorf("Missing pronouns key in data")
	}

	pronouns := []models.Pronoun{}
	for _, word := range rawpronouns {
		p := models.Pronoun{Word: word[1]}
		if strings.Contains(word[0], "p") {
			p.Plural = true
		}
		pronouns = append(pronouns, p)
	}

	return g.db.InitData(adjectives, nouns, verbs, pronouns, nil)
}

// Prepend "a", "an" or nothing to a phrase
func ana(phrase string) string {
	//fmt.Printf("[ana] phrase[0]: %s; %q\n", string(phrase[0]), phrase)
	if strings.ContainsAny(string(phrase[0]), "aeiou") {
		return "an " + phrase
	}

	return "a " + phrase
}

func toCap(words string) string {
	return strings.ToUpper(string(words[0])) + words[1:]
}
