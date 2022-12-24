package model

// Any is an interface for using in generic functions.
type Any interface {
	Category | Currency | Account | Transaction
}
