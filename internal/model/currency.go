package model

// Currency is a model of account currency field.
type Currency struct {
	ID           int64
	Abbreviation string
	IsMain       bool
}

// NewEmptyCurrency returns an empty Currency. This function for consistancy with NewEmptyAccount
// and NewEmptyTransaction.
func NewEmptyCurrency() *Currency {
	return &Currency{}
}
