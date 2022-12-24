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

// NewEmptyTransaction returns an empty Transaction with non nil nested structures. The purpose of
// this func to avoid erros when calling nested fields when they points to nil.
func NewEmptyTransaction() *Transaction {
    return &Transaction{Account: NewEmptyAccount(), Category: NewEmptyCategory()}
}
