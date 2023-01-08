package ext

import "github.com/rivo/tview"

// NewAskModal returns new tview.Modal with text and "Yes/No" buttons. Use submit and cancel args
// to set actions on buttons click.
func NewAskModal(question string, submit, cancel func()) *tview.Modal {
	return NewModal().
		SetText(question).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			switch buttonLabel {
			case "Yes":
				submit()
			case "No":
				cancel()
			}
		})
}
