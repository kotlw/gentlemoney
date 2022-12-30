package transactions

import (
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
)

type TransactionsDataProvider struct {
	service   *service.Transaction
	presenter *presenter.Transaction
}

func NewTransactionDataProvider(service *service.Transaction, presenter *presenter.Transaction) *TransactionsDataProvider {
	return &TransactionsDataProvider{service: service, presenter: presenter}
}

func (d *TransactionsDataProvider) GetAll() []map[string]string {
	data := d.service.GetAll()

	res := make([]map[string]string, len(data))

	for i, e := range data {
		res[i] = d.presenter.ToMap(e)
	}

	return res
}
