package ext

import "github.com/rivo/tview"

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
