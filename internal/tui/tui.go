package tui

import (
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/tui/root"

	"github.com/rivo/tview"
)

func New(service *service.Service, presenter *presenter.Presenter) *tview.Application {
	root := root.New(service, presenter)
	app := tview.NewApplication()
	app.SetRoot(root, true).EnableMouse(false)

	return app
}
