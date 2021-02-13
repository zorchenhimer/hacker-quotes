package models

type Verb struct {
	Id int

	Regular bool

	// Indefinite form
	Word string

	ConjugationsPast    []Conjugate
	ConjugationsPresent []Conjugate
	ConjugationsFuture  []Conjugate
}

type ConjugateType int
const (
	CT_I ConjugateType = iota
	CT_You
	CT_It
	CT_We
	CT_They
)

type Conjugate struct {
	Type ConjugateType
	Form string
}
