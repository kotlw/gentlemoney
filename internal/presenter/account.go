package presenter

import (
	"fmt"
	"strconv"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/service"
)

// Account presenter contains logic related to UI.
type Account struct {
	currencyService *service.Currency
}

// NewAccount returns Account presenter.
func NewAccount(currencyService *service.Currency) *Account {
	return &Account{currencyService: currencyService}
}

// ToMap converts model.Account to map[string]string. It doesn't handles ID field.
func (p *Account) ToMap(a *model.Account) map[string]string {
	return map[string]string{
		"ID":       strconv.Itoa(int(a.ID)),
		"Name":     a.Name,
		"Currency": a.Currency.Abbreviation,
	}
}

// FromMap parses map[string]string to model.Account. It doesn't handles ID field.
func (p *Account) FromMap(m map[string]string) (*model.Account, error) {
	if err := checkKeys(m, []string{"Name", "Currency"}); err != nil {
		return nil, fmt.Errorf("checkKeys: %w", err)
	}

	id, err := getID(m)
	if err != nil {
		return nil, fmt.Errorf("getID: %w", err)
	}

	return &model.Account{
		ID:       id,
		Name:     m["Name"],
		Currency: p.currencyService.GetByAbbreviation(m["Currency"]),
	}, nil
}
