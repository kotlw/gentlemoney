package ext

import "github.com/rivo/tview"

// NewErrorModal returns new tview.Modal with "okay" button. Use cancel variable to set action button.
func NewErrorModal(cancel func()) *tview.Modal {
	return NewModal().
		AddButtons([]string{"ok"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "ok" {
				cancel()
			}
		})
}
