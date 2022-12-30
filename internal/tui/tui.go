package tui

import (
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/tui/root"
	"github.com/kotlw/gentlemoney/internal/tui/transactions"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func New(service *service.Service, presenter *presenter.Presenter) *tview.Application {
	root := root.New()
	transactionsView := transactions.New(service, presenter)
	root.AddView('1', "Transactions", transactionsView)
	root.AddView('2', "Settings", tview.NewBox())
	root.SwitchToView("Transactions")

	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		root.Capture(event)

		return event
	})

	app := tview.NewApplication()
	app.SetRoot(root, true).EnableMouse(false)

	return app
}
