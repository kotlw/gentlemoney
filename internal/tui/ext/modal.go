package ext

import "github.com/rivo/tview"

// WrapIntoModal wraps primitive into modal with specified width and height.
func WrapIntoModal(primitive tview.Primitive, width, height int) *tview.Flex {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(primitive, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

// NewModal returns new tview.Modal with default Form styles.
func NewModal() *tview.Modal {
	return tview.NewModal().
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor)
}
