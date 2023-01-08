package presenter

import (
	"fmt"
	"strconv"

	"github.com/kotlw/gentlemoney/internal/model"
)

// Currency presenter contains logic related to UI.
type Currency struct{}

// NewCurrency returns Currency presenter.
func NewCurrency() *Currency {
	return &Currency{}
}

// ToMap converts model.Currency to map[string]string. It doesn't handles ID field.
func (p *Currency) ToMap(c *model.Currency) map[string]string {
	return map[string]string{
		"ID":           strconv.Itoa(int(c.ID)),
		"Abbreviation": c.Abbreviation,
		"ExchangeRate": reprMoney(c.ExchangeRate),
		"IsMain":       p.reprIsMain(c.IsMain),
	}
}

// FromMap parses map[string]string to model.Currency. It doesn't handles ID field.
func (p *Currency) FromMap(m map[string]string) (*model.Currency, error) {
	if err := checkKeys(m, []string{"Abbreviation", "ExchangeRate", "IsMain"}); err != nil {
		return nil, fmt.Errorf("checkKeys: %w", err)
	}

	id, err := getID(m)
	if err != nil {
		return nil, fmt.Errorf("getID: %w", err)
	}

	exchangeRate, err := parseMoney(m["ExchangeRate"])
	if err != nil {
		return nil, fmt.Errorf("parseMoney: %w", err)
	}

	return &model.Currency{
		ID:           id,
		Abbreviation: m["Abbreviation"],
		ExchangeRate: exchangeRate,
		IsMain:       p.parseIsMain(m["IsMain"]),
	}, nil
}

func (p *Currency) reprIsMain(value bool) string {
	if value {
		return "*"
	}
	return ""
}

func (p *Currency) parseIsMain(value string) bool {
	return value == "*"
}
