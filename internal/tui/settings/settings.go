package settings

import (
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/tui/ext"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// View is a transactions view.
type View struct {
	*tview.Pages

	tuiApp    *tview.Application
	service   *service.Service
	presenter *presenter.Presenter

	flex *tview.Flex

	categoryTable       *ext.Table
	categoryCreateForm  *ext.Form
	categoryUpdateForm  *ext.Form
	categoryDeleteModal *tview.Modal

	currencyTable       *ext.Table
	currencyCreateForm  *ext.Form
	currencyUpdateForm  *ext.Form
	currencyDeleteModal *tview.Modal

	accountTable       *ext.Table
	accountCreateForm  *ext.Form
	accountUpdateForm  *ext.Form
	accountDeleteModal *tview.Modal

	errorModal *tview.Modal
}

// New returns new settings view.
func New(tuiApp *tview.Application, service *service.Service, presenter *presenter.Presenter) *View {
	v := &View{
		Pages: tview.NewPages(),

		tuiApp:    tuiApp,
		service:   service,
		presenter: presenter,
		flex:      tview.NewFlex(),
	}

	// dataProviders
	categoryDataProvider := NewCategoryDataProvider(service, presenter)
	currencyDataProvider := NewCurrencyDataProvider(service, presenter)
	accountDataProvider := NewAccountDataProvider(service, presenter)

	// table
	v.categoryTable = ext.NewTable([]string{"Title"}, categoryDataProvider).SetOrder("Title", false).Refresh()
	v.currencyTable = ext.NewTable([]string{"Abbreviation"}, currencyDataProvider).SetOrder("Abbreviation", false).Refresh()
	v.accountTable = ext.NewTable([]string{"Name", "Currency"}, accountDataProvider).SetOrder("Name", false).Refresh()
	v.categoryTable.SetTitle("Category")
	v.currencyTable.SetTitle("Currency")
	v.accountTable.SetTitle("Account")
	v.flex.AddItem(v.categoryTable, 0, 1, true)
	v.flex.AddItem(v.currencyTable, 0, 1, false)
	v.flex.AddItem(v.accountTable, 0, 1, false)
	v.AddPage("flex", v.flex, true, true)

	// create form
	v.categoryCreateForm = v.newCategoryForm("Create Category", v.submitCategoryCreateForm, v.hideCategoryCreateForm, categoryDataProvider)
	v.currencyCreateForm = v.newCurrencyForm("Create Currency", v.submitCurrencyCreateForm, v.hideCurrencyCreateForm, currencyDataProvider)
	v.accountCreateForm = v.newAccountForm("Create Account", v.submitAccountCreateForm, v.hideAccountCreateForm, accountDataProvider)
	v.AddPage("categoryCreateForm", ext.WrapIntoModal(v.categoryCreateForm, 40, 7), true, false)
	v.AddPage("currencyCreateForm", ext.WrapIntoModal(v.currencyCreateForm, 40, 7), true, false)
	v.AddPage("accountCreateForm", ext.WrapIntoModal(v.accountCreateForm, 40, 9), true, false)

	// update form
	v.categoryUpdateForm = v.newCategoryForm("Update Category", v.submitCategoryUpdateForm, v.hideCategoryUpdateForm, categoryDataProvider)
	v.currencyUpdateForm = v.newCurrencyForm("Update Currency", v.submitCurrencyUpdateForm, v.hideCurrencyUpdateForm, currencyDataProvider)
	v.accountUpdateForm = v.newAccountForm("Update Account", v.submitAccountUpdateForm, v.hideAccountUpdateForm, accountDataProvider)
	v.AddPage("categoryUpdateForm", ext.WrapIntoModal(v.categoryUpdateForm, 40, 7), true, false)
	v.AddPage("currencyUpdateForm", ext.WrapIntoModal(v.currencyUpdateForm, 40, 7), true, false)
	v.AddPage("accountUpdateForm", ext.WrapIntoModal(v.accountUpdateForm, 40, 9), true, false)

	// delete modal
	v.categoryDeleteModal = ext.NewAskModal("Are you sure?", v.submitCategoryDeleteModal, v.hideCategoryDeleteModal)
	v.currencyDeleteModal = ext.NewAskModal("Are you sure?", v.submitCurrencyDeleteModal, v.hideCurrencyDeleteModal)
	v.accountDeleteModal = ext.NewAskModal("Are you sure?", v.submitAccountDeleteModal, v.hideAccountDeleteModal)
	v.AddPage("categoryDeleteModal", v.categoryDeleteModal, true, false)
	v.AddPage("currencyDeleteModal", v.currencyDeleteModal, true, false)
	v.AddPage("accountDeleteModal", v.accountDeleteModal, true, false)

	// error modal
	v.errorModal = ext.NewErrorModal(v.hideError)
	v.AddPage("errorModal", v.errorModal, true, false)

	return v
}

// Table returns table view by given label.
func (v *View) Table(label string) *ext.Table {
	switch label {
	case "Category":
		return v.categoryTable
	case "Currency":
		return v.currencyTable
	case "Account":
		return v.accountTable
	}
	return nil
}

// ModalHasFocus returns true if any of modal is currently on focus.
func (v *View) ModalHasFocus() bool {
	for _, modal := range []tview.Primitive{
		v.categoryCreateForm, v.categoryUpdateForm, v.categoryDeleteModal,
		v.currencyCreateForm, v.currencyUpdateForm, v.currencyDeleteModal,
		v.accountCreateForm, v.accountUpdateForm, v.accountDeleteModal,
		v.errorModal,
	} {
		if modal.HasFocus() {
			return true
		}
	}
	return false
}

// InputHandler returns the handler for this primitive.
func (v *View) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if v.categoryTable.HasFocus() {
			// table controllers
			switch event.Rune() {
			case 'c':
				v.showCategoryCreateForm()
			case 'u':
        if len(v.categoryTable.GetSelectedRef()) != 0 {
          v.showCategoryUpdateForm()
        } else {
          v.showError("Nothing to update")
        }
			case 'd':
        if len(v.categoryTable.GetSelectedRef()) != 0 {
          v.showCategoryDeleteModal()
        } else {
          v.showError("Nothing to delete")
        }
			}

			// navigation between settings
			switch event.Key() {
			case tcell.KeyTab:
				v.tuiApp.SetFocus(v.currencyTable)
			case tcell.KeyBacktab:
				v.tuiApp.SetFocus(v.accountTable)
			}

			// if none of keys has pressed use standard table input handler.
			if handler := v.categoryTable.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}

		if v.currencyTable.HasFocus() {
			// table controllers
			switch event.Rune() {
			case 'c':
				v.showCurrencyCreateForm()
			case 'u':
        if len(v.currencyTable.GetSelectedRef()) != 0 {
          v.showCurrencyUpdateForm()
        } else {
          v.showError("Nothing to update")
        }
			case 'd':
        if len(v.currencyTable.GetSelectedRef()) != 0 {
          v.showCurrencyDeleteModal()
        } else {
          v.showError("Nothing to delete")
        }
			}

			// navigation between settings
			switch event.Key() {
			case tcell.KeyTab:
				v.tuiApp.SetFocus(v.accountTable)
			case tcell.KeyBacktab:
				v.tuiApp.SetFocus(v.categoryTable)
			}

			// if none of keys has pressed use standard table input handler.
			if handler := v.currencyTable.InputHandler(); handler != nil {
				handler(event, setFocus)

				return
			}
		}

		if v.accountTable.HasFocus() {
			// table controllers
			switch event.Rune() {
			case 'c':
				v.showAccountCreateForm()
			case 'u':
        if len(v.accountTable.GetSelectedRef()) != 0 {
          v.showAccountUpdateForm()
        } else {
          v.showError("Nothing to update")
        }
			case 'd':
        if len(v.accountTable.GetSelectedRef()) != 0 {
          v.showAccountDeleteModal()
        } else {
          v.showError("Nothing to delete")
        }
			}

			// navigation between settings
			switch event.Key() {
			case tcell.KeyTab:
				v.tuiApp.SetFocus(v.categoryTable)
			case tcell.KeyBacktab:
				v.tuiApp.SetFocus(v.currencyTable)
			}

			// if none of keys has pressed use standard table input handler.
			if handler := v.accountTable.InputHandler(); handler != nil {
				handler(event, setFocus)

				return
			}
		}

		// give control to the child view.
		for _, modal := range []tview.Primitive{
			v.categoryCreateForm, v.categoryUpdateForm, v.categoryDeleteModal,
			v.currencyCreateForm, v.currencyUpdateForm, v.currencyDeleteModal,
			v.accountCreateForm, v.accountUpdateForm, v.accountDeleteModal,
			v.errorModal,
		} {
			if modal.HasFocus() {
				if handler := modal.InputHandler(); handler != nil {
					handler(event, setFocus)

					return
				}
			}
		}

	})
}

// showError shows error modal.
func (v *View) showError(text string) {
	v.errorModal.SetText(text)
	v.Pages.ShowPage("errorModal")
}

// hideError hides error modal.
func (v *View) hideError() {
	v.Pages.HidePage("errorModal")
}
