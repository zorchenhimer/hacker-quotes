package models

import (
	"fmt"
	"strings"
)

type Verb struct {
	Id int

	Regular bool

	// Indefinite form
	Word string

	//ConjugationsPast    []Conjugate
	//ConjugationsPresent []Conjugate
	//ConjugationsFuture  []Conjugate
}

//type Conjugate struct {
//	Type ConjugateType
//	Form string
//}

type ConjugationType int
const (
	CT_Unknown ConjugationType = iota
	CT_I
	CT_You
	CT_It
	CT_We
	CT_They
)

func (ct ConjugationType) String() string {
	switch ct {
	case CT_I:
		return "I"
	case CT_It:
		return "It"
	case CT_You:
		return "You"
	case CT_We:
		return "We"
	case CT_They:
		return "They"
	default:
		return "Unknown"
	}
}

type ConjugationTime int
const (
	CM_Unknown ConjugationTime = iota
	CM_Present
	CM_Past
	CM_Future
)

func (cm ConjugationTime) String() string {
	switch cm {
	case CM_Present:
		return "Present"
	case CM_Past:
		return "Past"
	case CM_Future:
		return "Future"
	default:
		return "Unknown"
	}
}

// TODO: implement inverted
func (v Verb) Conjugate(conjugation ConjugationType, time ConjugationTime, inverted bool ) string {
	if !v.Regular {
		fmt.Printf("[verb] Irregular verbs not implemented! %q\n", v.Word)
		return v.Word
	}

	//fmt.Printf("[verb] word: %s; type: %s; time: %s\n", v.Word, conjugation, time)

	switch time {
	case CM_Present:
		if conjugation == CT_It {
			sfx := []string{"s", "x", "sh", "ch", "ss"}
			for _, s := range sfx {
				if strings.HasSuffix(v.Word, s) {
					return v.Word + "es"
				}
			}

			if strings.HasSuffix(v.Word, "y") && !strings.ContainsAny(string(v.Word[len(v.Word)-2]), "aeiou") {
				return v.Word[:len(v.Word)-1] + "ies"
			}

			return v.Word + "s"

		} else {
			return v.Word
		}

	case CM_Past:
		if strings.HasSuffix(v.Word, "e") {
			return v.Word + "d"
		}

		if strings.HasSuffix(v.Word, "y") && !strings.ContainsAny(string(v.Word[len(v.Word)-2]), "aeiou") {
			return v.Word[:len(v.Word)-1] + "ied"
		}

		return v.Word + "ed"

	case CM_Future:
		return "will " + v.Word
	}

	fmt.Printf("[verb] Unknown ConjugationTime: %d\n", conjugation)
	return v.Word
}
