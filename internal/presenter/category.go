package presenter

import (
	"fmt"
	"strconv"

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
	return map[string]string{"ID": strconv.Itoa(int(c.ID)), "Title": c.Title}
}

// FromMap parses map[string]string to model.Category. It doesn't handles ID field.
func (p *Category) FromMap(m map[string]string) (*model.Category, error) {
	if err := checkKeys(m, []string{"Title"}); err != nil {
		return nil, fmt.Errorf("checkKeys: %w", err)
	}

	id, err := getID(m)
	if err != nil {
		return nil, fmt.Errorf("getID: %w", err)
	}

	return &model.Category{ID: id, Title: m["Title"]}, nil
}
