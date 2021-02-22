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

type InitialData struct {
	Adjectives [][]string
	Nouns [][]string
	Verbs [][]string
	Pronouns [][]string
	Sentences []string
}

type english struct {
	db database.DB

	currentPronoun string
}

func NewEnglish(db database.DB) (HackerQuotes, error) {
	return &english{db: db}, nil
}

/*
	Sentence format

	bare words
	{word_type:options}
	{{word_type:new word:properties}}

	[rng:random,choice,list]		// one word is chosen from the comma-separated list

	// Randomly hide/show the enclosed definition
	[hide:{word_type:options}]
	[hide:{{word_type:new word:options}}]

	{pronoun} can't {verb:i,present} {noun_phrase}, it {verb:it,future} {noun_phrase}!
	{verb:you,present} {noun_phrase:indefinite}, then you can {verb:you,present} {noun_phrase:definite}!
	{noun_phrase} {verb}. With {noun_phrase:indifinite,noadj,compound}!

	{pronoun} {{verb:need:from_pronoun,present}} to {verb:i,present} {noun_phrase:definite}!

	{pronoun} [rng:can,may] {verb:i,present} {noun_phrase} [hide:{noun_phrase:singlular}], it {verb:it,future} {noun_phrase:definite}.
*/

func (g *english) Hack() (string, error) {
	//var fmtString string = `{verb:you,present} {noun_phrase:definite} then you can {verb:you,present} {noun_phrase:definite}!`
	//var str string = `{pronoun} {{verb:need:from_pronoun,present}} to {verb:i,present} {noun_phrase:definite}!`
	//var str string = `{pronoun} [rng:can,may] {verb:i,present} {noun_phrase}[hide: {noun_phrase:singlular}], it {verb:it,future} {noun_phrase:definite}.`
	str, err := g.randomSentence()
	if err != nil {
		return "", err
	}

	return g.HackThis(str)
}

func (g *english) HackThis(fmtString string) (string, error) {
	var idx int
	var err error
	var nidx int
	output := &strings.Builder{}
	g.currentPronoun = ""

	for idx < len(fmtString) {
		if fmtString[idx] == '{' {
			nidx, err = g.consumeWord(fmtString, idx, output)
			if err != nil {
				return "", err
			}
			idx = nidx
			continue

		} else if fmtString[idx] == '[' {
			nidx, err = g.consumeRng(fmtString, idx, output)
			if err != nil {
				return "", err
			}
			idx = nidx
			continue
		}

		nidx, err = g.consumeRaw(fmtString, idx, output)
		if err != nil {
			return "", err
		}
		idx = nidx
	}

	return toCap(output.String()), nil
}

func (g *english) consumeRng(fmtString string, idx int, output *strings.Builder) (int, error) {
	idx++
	def := strings.Index(fmtString[idx:], ":")
	if def == -1 {
		return 0, fmt.Errorf("Missing command separator in RNG block")
	}
	def += idx

	end := strings.Index(fmtString[idx:], "]")
	if end == -1 {
		return 0, fmt.Errorf("Unclosed RNG block starting at offset %d", idx-1)
	}
	end += idx

	if def > end {
		return 0, fmt.Errorf("Missing command separator in RNG block")
	}

	switch fmtString[idx:def] {
	case "rng":
		choices := strings.Split(fmtString[def+1:end], ",")
		ridx := rand.Intn(len(choices))
		output.WriteString(choices[ridx])

	case "hide":
		if rand.Int() % 2 == 0 {
			return end+1, nil
		}

		var newEnd int = def+1
		var err error
		for newEnd < end {
			if fmtString[newEnd] == '{' {
				newEnd, err = g.consumeWord(fmtString, newEnd, output)
				if err != nil {
					return 0, err
				}
			} else {
				newEnd, err = g.consumeRaw(fmtString, newEnd, output)
			}
		}

	default:
		return 0, fmt.Errorf("RNG type %q is not implemented", fmtString[idx:def])
	}

	return end+1, nil
}

func (g *english) consumeRaw(fmtString string, idx int, output *strings.Builder) (int, error) {
	end := strings.IndexAny(fmtString[idx:], "{[")
	if end == -1 {
		output.WriteString(fmtString[idx:len(fmtString)])
		return len(fmtString), nil
	}

	output.WriteString(fmtString[idx:end+idx])
	return idx+end, nil
}

func (g *english) parseVerbOptions(optlist []string) (models.ConjugationType, models.ConjugationTime, bool) {
		var ct models.ConjugationType = models.ConjugationType(rand.Intn(5) + 1)
		if sliceContains(optlist, "i") {
			ct = models.CT_I
		} else if sliceContains(optlist, "you") {
			ct = models.CT_You
		} else if sliceContains(optlist, "it") {
			ct = models.CT_It
		} else if sliceContains(optlist, "we") {
			ct = models.CT_We
		} else if sliceContains(optlist, "they") {
			ct = models.CT_They
		}

		var cm models.ConjugationTime = models.ConjugationTime(rand.Intn(3) + 1)
		if sliceContains(optlist, "present") {
			cm = models.CM_Present
		} else if sliceContains(optlist, "past") {
			cm = models.CM_Past
		} else if sliceContains(optlist, "future") {
			cm = models.CM_Future
		}

		var invert bool = rand.Int() % 2 == 0
		if sliceContains(optlist, "invert") {
			invert = true
		}

		//var regular bool = true
		//if sliceContains(optlist, "irregular") {
		//	regular = false
		//}

		if sliceContains(optlist, "from_pronoun") {
			switch g.currentPronoun {
			case "he", "she", "it":
				ct = models.CT_It
			case "i":
				ct = models.CT_I
			case "you":
				ct = models.CT_You
			case "we":
				ct = models.CT_We
			case "they":
				ct = models.CT_They
			}
		}

		return ct, cm, invert
}

func (g *english) parseNounOptions(optlist []string) (bool, bool) {
	var plural bool = rand.Int() % 2 == 0
	if sliceContains(optlist, "plural") {
		plural = true
	} else if sliceContains(optlist, "singular") {
		plural = false
	}

	var compound bool = rand.Int() % 2 == 0
	if sliceContains(optlist, "compound") {
		compound = true
	} else if sliceContains(optlist, "simple") {
		compound = false
	}

	return plural, compound
}

func (g *english) consumeNewWord(fmtString string, idx int, output *strings.Builder) (int, error) {
	idx += 2
	var wordtype string
	var word string
	var options string

	end := strings.Index(fmtString[idx:], "}")
	if end == -1 {
		return 0, fmt.Errorf("[consumeNewWord] Unclosed definition starting at %d", idx)
	}
	end += idx

	// Find the start of the new word
	def := strings.Index(fmtString[idx:], ":")
	if def == -1 {
		return 0, fmt.Errorf("[consumeNewWord] Word definition missing word")
	}
	def += idx
	wordtype = fmtString[idx:def]

	// Find the start of the options
	opts := strings.Index(fmtString[def+1:], ":")

	if opts > -1 {
		opts += def+2
		options = fmtString[opts:end]
		word = fmtString[def+1:opts-1]
	} else {
		word = fmtString[def+1:end]
	}

	//fmt.Printf("idx:%d def:%d opts:%d end:%d type:%q word:%q opts:%q\n",
	//	idx, def, opts, end,
	//	wordtype,
	//	word,
	//	options)

	optlist := strings.Split(options, ",")
	var formatted string
	switch wordtype {
	case "verb":
		ct, cm, invert := g.parseVerbOptions(optlist)

		v := models.Verb{
			Regular: true,
			Word: word,
		}

		formatted = v.Conjugate(ct, cm, invert)

	case "noun":
		plural, _ := g.parseNounOptions(optlist)

		if plural {
			n := models.Noun{
				Regular: true,
				Word: word,
			}
			formatted = n.Plural()
		} else {
			formatted = word
		}

	default:
		// Pronouns and adjectives aren't implemented because they have no modifications to apply.
		return 0, fmt.Errorf("Word type %q not implemented creating in-sentence", wordtype)
	}

	output.WriteString(formatted)
	return end+2, nil
}

func (g *english) consumeWord(fmtString string, idx int, output *strings.Builder) (int, error) {
	if fmtString[idx+1] == '{' {
		return g.consumeNewWord(fmtString, idx, output)
	}

	idx++
	var wordtype string
	var options string

	end := strings.Index(fmtString[idx:], "}")
	if end == -1 {
		return 0, fmt.Errorf("[consumeWord] Unclosed definition starting at %d", idx)
	}

	end += idx
	optsStart := strings.Index(fmtString[idx:end], ":")

	if optsStart != -1 {
		options = fmtString[optsStart+idx+1:end]
		wordtype = fmtString[idx:optsStart+idx]
	} else {
		wordtype = fmtString[idx:end]
	}

	if wordtype == "" {
		return 0, fmt.Errorf("[consumeWord] Missing word type at idx: %d", idx)
	}

	opts := strings.Split(options, ",")
	var word string
	var err error

	switch wordtype {
	case "pronoun":
		var plural bool = rand.Int() % 2 == 0
		if sliceContains(opts, "plural") {
			plural = true
		} else if sliceContains(opts, "singular") {
			plural = false
		}

		word, err = g.randomPronoun(plural)
		if err != nil {
			return 0, err
		}
		g.currentPronoun = word

	case "verb":
		ct, cm, invert := g.parseVerbOptions(opts)

		word, err = g.randomVerb(ct, cm, invert)
		if err != nil {
			return 0, err
		}

	case "noun":
		plural, compound := g.parseNounOptions(opts)

		word, err = g.randomNoun(plural, compound)
		if err != nil {
			return 0, err
		}

	case "noun_phrase":
		var definite bool = rand.Int() % 2 == 0
		var hasAdj bool = rand.Int() % 2 == 0
		var plural bool = rand.Int() % 2 == 0
		var compound bool = rand.Int() % 2 == 0

		if sliceContains(opts, "indefinite") {
			definite = false
		} else if sliceContains(opts, "definite") {
			definite = true
		}

		if sliceContains(opts, "noadj") {
			hasAdj = false
		} else if sliceContains(opts, "adj") {
			hasAdj = true
		}

		if sliceContains(opts, "plural") {
			plural = true
		} else if sliceContains(opts, "singular") {
			plural = false
		}

		if sliceContains(opts, "compound") {
			compound = true
		} else if sliceContains(opts, "simple") {
			compound = false
		}

		word, err = g.nounPhrase(definite, hasAdj, plural, compound)
		if err != nil {
			return 0, err
		}

	case "adjective":
		word, err = g.randomAdjective()
		if err != nil {
			return 0, err
		}

	default:
		return 0, fmt.Errorf("[consumeWord] Invalid word type %s at %d", wordtype, idx)
	}

	output.WriteString(word)
	return end+1, nil
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
		return g.ana(phrase), nil
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

	if pronoun.Word == "i" {
		return strings.ToUpper(pronoun.Word), nil
	}
	return pronoun.Word, nil
}

func (g *english) randomSentence() (string, error) {
	ids, err := g.db.GetSentenceIds()
	if err != nil {
		return "", fmt.Errorf("[sentence] get IDs error: %v", err)
	}

	if len(ids) <= 0 {
		return "", fmt.Errorf("[sentence] No sentence IDs returned from database")
	}

	rid := int(rand.Int63n(int64(len(ids))))
	sentence, err := g.db.GetSentence(ids[rid])
	if err != nil {
		return "", fmt.Errorf("[sentence] ID: %d, %v", ids[rid], err)
	}

	return sentence, nil
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

	//data := map[string][]interface{}{}
	data := InitialData{}
	if err = json.Unmarshal(raw, &data); err != nil {
		return err
	}

	if data.Adjectives == nil || len(data.Adjectives) == 0 {
		return fmt.Errorf("Missing Adjectives in input data")
	}

	adjectives := []models.Adjective{}
	for _, adj := range data.Adjectives {
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

	if data.Nouns == nil || len(data.Nouns) == 0 {
		return fmt.Errorf("Missing nouns key in data")
	}

	nouns := []models.Noun{}
	for _, noun := range data.Nouns {
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

	if data.Verbs == nil || len(data.Verbs) == 0 {
		return fmt.Errorf("Missing verbs key in data")
	}

	verbs := []models.Verb{}
	for _, verb := range data.Verbs {
		v := models.Verb{Word: verb[1]}
		if strings.Contains(verb[0], "r") {
			v.Regular = true
		}

		verbs = append(verbs, v)
	}

	if data.Pronouns == nil || len(data.Pronouns) == 0 {
		return fmt.Errorf("Missing pronouns key in data")
	}

	pronouns := []models.Pronoun{}
	for _, pro := range data.Pronouns {
		p := models.Pronoun{Word: pro[1]}
		if strings.Contains(pro[0], "p") {
			p.Plural = true
		}
		pronouns = append(pronouns, p)
	}

	if data.Sentences == nil || len(data.Sentences) == 0 {
		return fmt.Errorf("Missing sentences key in data")
	}

	return g.db.InitData(adjectives, nouns, verbs, pronouns, data.Sentences)
}

// Prepend "a", "an" or nothing to a phrase
func (g *english) ana(phrase string) string {
	//fmt.Printf("[ana] phrase[0]: %s; %q\n", string(phrase[0]), phrase)
	if strings.ContainsAny(string(phrase[0]), "aeiou") {
		return "an " + phrase
	}

	return "a " + phrase
}

// toCap capitalizes the first word of each sentence in the input string.
func toCap(words string) string {
	ret :=  strings.ToUpper(string(words[0])) + words[1:]

	next := strings.Index(words, ". ")
	if next == -1 {
		return ret
	}

	for next+3 < len(words) {
		newnext := strings.Index(words[next+1:], ". ")
		ret = ret[0:next+2] + strings.ToUpper(string(ret[next+2])) + ret[next+3:]

		if newnext == -1 {
			break
		}
		next = newnext + next+1
	}

	return ret
}

func sliceContains(haystack []string, needle string) bool {
	if len(haystack) == 0 {
		return false
	}

	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
