package model

// Category is a model of transaction category field.
type Category struct {
	ID    int64
	Title string
}

// NewEmptyCategory returns an empty Category. This function for consistancy with NewEmptyAccount
// and NewEmptyTransaction.
func NewEmptyCategory() *Category {
	return &Category{}
}
