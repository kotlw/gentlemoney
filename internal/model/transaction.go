package model

import (
	"time"
)

// Transaction is a model of transaction which is main entitty of the app.
type Transaction struct {
	ID       int64
	Date     time.Time
	Account  *Account
	Category *Category
	Amount   int64
	Note     string
}

// NewTransaction returns an empty Transaction with nested structures.
func NewTransaction() *Transaction {
    return &Transaction{Account: NewAccount(), Category: NewCategory()}
}
