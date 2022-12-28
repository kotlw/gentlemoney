package model

// Account is a model of transaction account field.
type Account struct {
	ID       int64
	Name     string
	Currency *Currency
}

// NewEmptyAccount returns an empty Account with non nil nested structure. The purpose of this func
// to avoid erros when calling nested fields when they points to nil.
func NewEmptyAccount() *Account {
	return &Account{Currency: NewEmptyCurrency()}
}
