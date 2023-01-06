package root

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/tui/transactions"

	"github.com/rivo/tview"
)

// Root is the root view of the app. It aggregates pages and navbar in flex.
type Root struct {
	*tview.Flex

	navbar       *tview.TextView
	pages        *tview.Pages
	transactions *transactions.View
}

// New returns Root.
func New(service *service.Service, presenter *presenter.Presenter) *Root {
	root := &Root{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
		navbar: tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWrap(false),
		pages: tview.NewPages(),
	}

	root.AddItem(root.navbar, 1, 1, false)
	root.AddItem(root.pages, 0, 16, true)

	root.transactions = transactions.New(service, presenter)
	root.AddView('1', "Transactions", root.transactions)

	root.SwitchToView("Transactions")

	return root
}

// AddView adds view with corresponding navbar item.
func (r *Root) AddView(hint rune, title string, view tview.Primitive) {
	fmt.Fprintf(r.navbar, `  ["%s"]%c %s[""] `, title, hint, title)
	r.pages.AddPage(title, view, true, false)
}

// SwitchToView switches to page by name and highlights corresponting navbar item.
func (r *Root) SwitchToView(name string) {
	r.navbar.Highlight(name)
	r.navbar.ScrollToHighlight()
	r.pages.SwitchToPage(name)
}

// IsModalOnTop check if modal of any child view is on top.
func (r *Root) IsModalOnTop() bool {
	return r.transactions.ModalHasFocus()
}

// InputHandler returns the handler for this primitive.
func (r *Root) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// modal shouldn't be on top, to perform switch view.
		if !r.IsModalOnTop() {
			if event.Rune() == '1' {
				r.SwitchToView("Transactions")

				return
			}
		}

		// if modal is active all other handlers should be ignored except modal handler.
		for _, view := range []tview.Primitive{r.transactions} {
			if view.HasFocus() {
				// give control to the child view.
				if handler := view.InputHandler(); handler != nil {
					handler(event, setFocus)

					return
				}
			}
		}
	})
}
