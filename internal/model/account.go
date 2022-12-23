package model

// Account is a model of transaction account field.
type Account struct {
	ID          int64
	Name        string
	Currency    *Currency
}

// NewAccount returns an empty Account with nested structures.
func NewAccount() *Account {
	return &Account{Currency: NewCurrency()}
}
