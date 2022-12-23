package model

// Currency is a model of account currency field.
type Currency struct {
	ID           int64
	Abbreviation string
	ExchangeRate int64
	IsMain       bool
}

// NewCurrency returns an empty Currency.
func NewCurrency() *Currency {
    return &Currency{}
}
