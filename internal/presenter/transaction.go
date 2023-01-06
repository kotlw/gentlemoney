package presenter

import (
	"fmt"
	"time"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/service"
)

// Transaction presenter contains logic related to UI.
type Transaction struct {
	accountService  *service.Account
	categoryService *service.Category
}

// NewTransaction returns Transaction presenter.
func NewTransaction(accountService *service.Account, categoryService *service.Category) *Transaction {
	return &Transaction{accountService: accountService, categoryService: categoryService}
}

// ToMap converts model.Transaction to map[string]string. It doesn't handles ID field.
func (p *Transaction) ToMap(t *model.Transaction) map[string]string {
	return map[string]string{
		"Date":     t.Date.Format("2006-01-02"),
		"Account":  t.Account.Name,
		"Category": t.Category.Title,
		"Amount":   p.reprAmount(t.Amount),
		"Currency": t.Account.Currency.Abbreviation,
		"Note":     t.Note,
	}
}

// FromMap parses map[string]string to model.Transaction. It doesn't handles ID field.
func (p *Transaction) FromMap(m map[string]string) (*model.Transaction, error) {
	if err := checkKeys(m, []string{"Date", "Account", "Category", "Amount", "Note"}); err != nil {
		return nil, fmt.Errorf("checkKeys: %w", err)
	}

	date, err := time.Parse("2006-01-02", m["Date"])
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}

	amount, err := parseMoney(m["Amount"])
	if err != nil {
		return nil, fmt.Errorf("parseMoney: %w", err)
	}

	return &model.Transaction{
		Date:     date,
		Account:  p.accountService.GetByName(m["Account"]),
		Category: p.categoryService.GetByTitle(m["Category"]),
		Amount:   amount,
		Note:     m["Note"],
	}, nil
}

func (*Transaction) reprAmount(value int64) string {
	sign := ""
	if value > 0 {
		sign = "+"
	}
	return sign + reprMoney(value)
}
