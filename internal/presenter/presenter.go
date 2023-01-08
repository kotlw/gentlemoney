package presenter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kotlw/gentlemoney/internal/service"
)

// Presenter is a facade structure which aggregates all Presenters. It is used for convenience.
type Presenter struct {
	category    *Category
	currency    *Currency
	account     *Account
	transaction *Transaction
}

// New returns new Presenter.
func New(service *service.Service) *Presenter {
	return &Presenter{
		category:    NewCategory(),
		currency:    NewCurrency(),
		account:     NewAccount(service.Currency()),
		transaction: NewTransaction(service.Account(), service.Category()),
	}
}

// Category returns category presenter.
func (p *Presenter) Category() *Category {
	return p.category
}

// Currency returns category presenter.
func (p *Presenter) Currency() *Currency {
	return p.currency
}

// Account returns category presenter.
func (p *Presenter) Account() *Account {
	return p.account
}

// Transaction returns category presenter.
func (p *Presenter) Transaction() *Transaction {
	return p.transaction
}

// checkKeys checks if all given keys are exist.
func checkKeys(m map[string]string, keys []string) error {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return fmt.Errorf(`key "` + k + `" is missing`)
		}
	}
	return nil
}

// reprMoney converts int64 to string money format 0.00. Last 0-2 digits always will be after ".",
// so value 1 becomes to "0.01", -12 to "-0.12", 0 to "0.00".
func reprMoney(value int64) string {
	res := strconv.Itoa(int(value))
	sign := ""

	// check if the value is negative and store the sign.
	if res[0] == '-' {
		res = res[1:]
		sign = "-"
	}

	// handle edge cases
	if len(res) == 1 {
		return sign + "0.0" + res
	} else if len(res) == 2 {
		return sign + "0." + res
	}

	return sign + res[:len(res)-2] + "." + res[len(res)-2:]
}

// parseMoney parses strings like "13.41" to int64 value of 1341.
func parseMoney(value string) (int64, error) {
	multiplier := 1
	if !strings.Contains(value, ".") {
		multiplier = 100
	}
	value = strings.Replace(value, ".", "", 1)

	res, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return int64(res) * int64(multiplier), nil
}

// getID returns id from given map, returns 0 if key "ID" is missing in map.
func getID(m map[string]string) (int64, error) {
	idStr, ok := m["ID"]
	if !ok {
		idStr = "0"
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
    return int64(id), nil
}
