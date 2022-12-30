package root

import (
	"fmt"

	"github.com/rivo/tview"
)

// navBar is a menu bar which shown on top of the screen.
type navBar struct {
	*tview.TextView
}

// newNavBar returns new navBar.
func newNavBar() *navBar {
	return &navBar{
		tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWrap(false),
	}
}

// addItem adds item into navbar, which could be highlighted by it's label. Hint should be used as
// a key to select the desired page.
func (tb *navBar) addItem(label string, hint rune) {
	fmt.Fprintf(tb, `  ["%s"]%c %s[""] `, label, hint, label)
}

// highlightItem highlights item corresponding to currently selected page.
func (tb *navBar) highlightItem(name string) {
	tb.Highlight(name)
	tb.ScrollToHighlight()
}
