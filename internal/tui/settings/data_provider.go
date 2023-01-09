package settings

import (
	"sort"

	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
)

// CategoryDataProvider implements ext.TableDataProvider for interaction with categories.
type CategoryDataProvider struct {
	service   *service.Service
	presenter *presenter.Presenter
}

// NewCategoryDataProvider returns new CategoryDataProvider.
func NewCategoryDataProvider(service *service.Service, presenter *presenter.Presenter) *CategoryDataProvider {
	return &CategoryDataProvider{service: service, presenter: presenter}
}

// GetAll returns slice of maps which represents category struct.
func (d *CategoryDataProvider) GetAll() []map[string]string {
	data := d.service.Category().GetAll()

	res := make([]map[string]string, len(data))

	for i, e := range data {
		res[i] = d.presenter.Category().ToMap(e)
	}

	return res
}

// GetDropDownOptions returns dropdown obtions for given label.
func (d *CategoryDataProvider) GetDropDownOptions(label string) []string {
	return nil
}

// CurrencyDataProvider implements ext.TableDataProvider for interaction with currencies.
type CurrencyDataProvider struct {
	service   *service.Service
	presenter *presenter.Presenter
}

// NewCurrencyDataProvider returns new CurrencyDataProvider.
func NewCurrencyDataProvider(service *service.Service, presenter *presenter.Presenter) *CurrencyDataProvider {
	return &CurrencyDataProvider{service: service, presenter: presenter}
}

// GetAll returns slice of maps which represents currency struct.
func (d *CurrencyDataProvider) GetAll() []map[string]string {
	data := d.service.Currency().GetAll()

	res := make([]map[string]string, len(data))

	for i, e := range data {
		res[i] = d.presenter.Currency().ToMap(e)
	}

	return res
}

// GetDropDownOptions returns dropdown obtions for given label.
func (d *CurrencyDataProvider) GetDropDownOptions(label string) []string {
	return nil
}

// AccountDataProvider implements ext.TableDataProvider and ext.FromDataProvider for interaction with accounts.
type AccountDataProvider struct {
	service   *service.Service
	presenter *presenter.Presenter
}

// NewAccountDataProvider returns new AccountDataProvider.
func NewAccountDataProvider(service *service.Service, presenter *presenter.Presenter) *AccountDataProvider {
	return &AccountDataProvider{service: service, presenter: presenter}
}

// GetAll returns slice of maps which represents account struct.
func (d *AccountDataProvider) GetAll() []map[string]string {
	data := d.service.Account().GetAll()

	res := make([]map[string]string, len(data))

	for i, e := range data {
		res[i] = d.presenter.Account().ToMap(e)
	}

	return res
}

// GetDropDownOptions returns dropdown obtions for given label.
func (d *AccountDataProvider) GetDropDownOptions(label string) []string {
	switch label {
	case "Currency":
		return d.currencyOptions()
	}
	return nil
}

// currencyOptions returns currency dropdown options.
func (d *AccountDataProvider) currencyOptions() []string {
	currencies := d.service.Currency().GetAll()

	res := make([]string, len(currencies))

	for i, e := range currencies {
		res[i] = e.Abbreviation
	}

	sort.Strings(res)

	return res
}
