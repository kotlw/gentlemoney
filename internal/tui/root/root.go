package root

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Root is the root view of the app. It aggregates pages and navBar in flex.
type Root struct {
	*tview.Flex

	navBar   *navBar
	pages    *tview.Pages
	captures []func(*tcell.EventKey)
}

// New returns Root.
func New() *Root {
	root := &Root{
		Flex:   tview.NewFlex().SetDirection(tview.FlexRow),
		navBar: newNavBar(),
		pages:  tview.NewPages(),
	}

	root.AddItem(root.navBar, 1, 1, false)
	root.AddItem(root.pages, 0, 16, true)

	return root
}

// AddPage adds page with corresponding navbar item.
func (r *Root) AddView(key rune, name string, view tview.Primitive) {
	r.navBar.addItem(name, key)
	r.pages.AddPage(name, view, true, true)

	r.captures = append(r.captures, func(event *tcell.EventKey) {
		if event.Rune() == key {
			r.SwitchToView(name)
		}
	})
}

// SwitchToView switches to page by name and highlights corresponting navbar item.
func (l *Root) SwitchToView(name string) {
	l.navBar.highlightItem(name)
	l.pages.SwitchToPage(name)
}

// Capture is a func with defined conditions to switch to certain page by defined key.
// Should be used with SetInputCapture func.
func (l *Root) Capture(event *tcell.EventKey) *tcell.EventKey {
	for _, cap := range l.captures {
		cap(event)
	}
	return event
}
