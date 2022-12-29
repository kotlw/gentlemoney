package presenter

import (
	"fmt"

	"github.com/kotlw/gentlemoney/internal/model"
)

// Category presenter contains logic related to UI.
type Category struct{}

// NewCategory returns Category presenter.
func NewCategory() *Category {
	return &Category{}
}

// ToMap converts model.Category to map[string]string. It doesn't handles ID field.
func (p *Category) ToMap(c *model.Category) map[string]string {
	return map[string]string{"Title": c.Title}
}

// FromMap parses map[string]string to model.Category. It doesn't handles ID field.
func (p *Category) FromMap(m map[string]string) (*model.Category, error) {
	if err := checkKeys(m, []string{"Title"}); err != nil {
		return nil, fmt.Errorf("checkKeys: %w", err)
	}
	return &model.Category{Title: m["Title"]}, nil
}
