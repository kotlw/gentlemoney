package tui

import (
    "github.com/kotlw/gentlemoney/internal/tui/root"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

func New() *tview.Application {
	root := root.New()
	root.AddView('1', "Transactions", tview.NewBox())
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
