package tui

import (
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/tui/root"

	"github.com/rivo/tview"
)

func New(service *service.Service, presenter *presenter.Presenter) *tview.Application {
	app := tview.NewApplication()
	root := root.New(app, service, presenter)
	app.SetRoot(root, true).EnableMouse(false)

	return app
}
