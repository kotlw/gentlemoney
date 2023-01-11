package transactions

import (
	"sort"

	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
)

// DataProvider implements ext.TableDataProvider and ext.FromDataProvider for interaction with transactions.
type DataProvider struct {
	service   *service.Service
	presenter *presenter.Presenter
}

// NewDataProvider returns new DataProvider.
func NewDataProvider(service *service.Service, presenter *presenter.Presenter) *DataProvider {
	return &DataProvider{service: service, presenter: presenter}
}

// GetAll returns slice of maps which represents transaction struct.
func (d *DataProvider) GetAll() []map[string]string {
	data := d.service.Transaction().GetAll()

	res := make([]map[string]string, len(data))

	for i, e := range data {
		m := d.presenter.Transaction().ToMap(e)
		if m["Amount"][0] == '+' {
			m["Amount"] = "[green]" + m["Amount"] + "[white]"
		}
		if m["Amount"][0] == '-' {
			m["Amount"] = "[red]" + m["Amount"] + "[white]"
		}
		res[i] = m
	}

	return res
}

// GetDropDownOptions returns dropdown obtions for given label.
func (d *DataProvider) GetDropDownOptions(label string) []string {
	switch label {
	case "Account":
		return d.accountOptions()
	case "Category":
		return d.categoryOptions()
	}
	return nil
}

// accountOptions returns account dropdown options.
func (d *DataProvider) accountOptions() []string {
	accounts := d.service.Account().GetAll()

	res := make([]string, len(accounts))

	for i, e := range accounts {
		res[i] = e.Name
	}

	sort.Strings(res)

	return res
}

// categoryOptions returns category dropdown options.
func (d *DataProvider) categoryOptions() []string {
	categories := d.service.Category().GetAll()

	res := make([]string, len(categories))

	for i, e := range categories {
		res[i] = e.Title
	}

	sort.Strings(res)

	return res
}
