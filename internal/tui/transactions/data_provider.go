package transactions

import (
	"sort"

	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
)

type DataProvider struct {
	service   *service.Service
	presenter *presenter.Presenter
}

func NewDataProvider(service *service.Service, presenter *presenter.Presenter) *DataProvider {
	return &DataProvider{service: service, presenter: presenter}
}

func (d *DataProvider) GetAll() []map[string]string {
	data := d.service.Transaction().GetAll()

	res := make([]map[string]string, len(data))

	for i, e := range data {
		res[i] = d.presenter.Transaction().ToMap(e)
	}

	return res
}

func (d *DataProvider) GetDropDownOptions(label string) []string {
	switch label {
	case "Account":
		return d.AccountOptions()
	case "Category":
		return d.CategoryOptions()
	}
	return nil
}

func (d *DataProvider) AccountOptions() []string {
	accounts := d.service.Account().GetAll()

	res := make([]string, len(accounts))

	for i, e := range accounts {
		res[i] = e.Name
	}

	sort.Strings(res)

	return res
}

func (d *DataProvider) CategoryOptions() []string {
	categories := d.service.Category().GetAll()

	res := make([]string, len(categories))

	for i, e := range categories {
		res[i] = e.Title
	}

	sort.Strings(res)

	return res
}
