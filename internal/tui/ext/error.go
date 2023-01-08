package ext

import "github.com/rivo/tview"

func NewErrorModal(cancel func()) *tview.Modal {
	return NewModal().
		AddButtons([]string{"ok"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "ok" {
				cancel()
			}
		})
}
