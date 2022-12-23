package model

// Category is a model of transaction category field.
type Category struct {
	ID    int64
	Title string
}

// NewCategory returns an empty Category.
func NewCategory() *Category {
	return &Category{}	
}
