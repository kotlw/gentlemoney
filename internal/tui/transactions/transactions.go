package transactions

import (
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/tui/table"

	"github.com/rivo/tview"
)

type TransactionView struct {
	*tview.Pages

	table *table.Table
}

func New(service *service.Service, presenter *presenter.Presenter) *TransactionView {
    tv := &TransactionView{Pages: tview.NewPages()}

	cols := []string{"Date", "Account", "Category", "Amount", "Currency", "Note"}
	dataProvider := NewTransactionDataProvider(service.Transaction(), presenter.Transaction())
	tv.table = table.New(cols, dataProvider).SetOrder("Date", true).Refresh()

	tv.AddPage("table", tv.table, true, true)

	return tv
}
