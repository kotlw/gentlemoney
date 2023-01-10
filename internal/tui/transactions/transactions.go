package transactions

import (
	"strings"
	"time"

	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/tui/ext"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// View is a transactions view.
type View struct {
	*tview.Pages

	service   *service.Service
	presenter *presenter.Presenter

	table       *ext.Table
	createForm  *ext.Form
	updateForm  *ext.Form
	deleteModal *tview.Modal
	errorModal  *tview.Modal
}

// New returns new transactions view.
func New(service *service.Service, presenter *presenter.Presenter) *View {
	v := &View{
		Pages: tview.NewPages(),

		service:   service,
		presenter: presenter,
	}

	dataProvider := NewDataProvider(v.service, v.presenter)

	// table
	cols := []string{"Date", "Account", "Category", "Amount", "Currency", "Note"}
	v.table = ext.NewTable(cols, dataProvider).SetOrder("Date", true).Refresh()
	v.AddPage("table", v.table, true, true)

	// create form
	v.createForm = v.newForm("Create Transaction", v.submitCreateForm, v.hideCreateForm, dataProvider)
	v.AddPage("createForm", ext.WrapIntoModal(v.createForm, 40, 15), true, false)

	// update form
	v.updateForm = v.newForm("Update Transaction", v.submitUpdateForm, v.hideUpdateForm, dataProvider)
	v.AddPage("updateForm", ext.WrapIntoModal(v.updateForm, 40, 15), true, false)

	// delete modal
	v.deleteModal = ext.NewAskModal("Are you sure?", v.submitDeleteModal, v.hideDeleteModal)
	v.AddPage("deleteModal", v.deleteModal, true, false)

	// error modal
	v.errorModal = ext.NewErrorModal(v.hideError)
	v.AddPage("errorModal", v.errorModal, true, false)

	return v
}

// ModalHasFocus returns true if any of modal is currently on focus.
func (v *View) ModalHasFocus() bool {
	for _, modal := range []tview.Primitive{v.createForm, v.updateForm, v.deleteModal, v.errorModal} {
		if modal.HasFocus() {
			return true
		}
	}
	return false
}

// InputHandler returns the handler for this primitive.
func (v *View) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if v.table.HasFocus() {
			if event.Rune() == 'c' {
				v.showCreateForm()
			}

			if event.Rune() == 'u' {
				v.showUpdateForm()
			}

			if event.Rune() == 'd' {
				v.showDeleteModal()
			}

			// if none of keys has pressed use standard table input handler.
			if handler := v.table.InputHandler(); handler != nil {
				handler(event, setFocus)

				return
			}
		}

		// give control to the child view.
		for _, modal := range []tview.Primitive{v.createForm, v.updateForm, v.deleteModal, v.errorModal} {
			if modal.HasFocus() {
				if handler := modal.InputHandler(); handler != nil {
					handler(event, setFocus)

					return
				}
			}
		}

	})
}

// newForm returns new form with corresponding transaction fields.
func (v *View) newForm(title string, submit func(), cancel func(), dataProvider *DataProvider) *ext.Form {
	form := tview.NewForm().
		AddInputField("Date", "", 0, nil, nil).
		AddDropDown("Category", nil, 0, nil).
		AddDropDown("Account", nil, 0, nil).
		AddInputField("Amount", "", 0, nil, nil).
		AddInputField("Note", "", 0, nil, nil).
		AddButton(strings.Split(title, " ")[0], submit).
		AddButton("Cancel", cancel)

	form.SetBorder(true)
	form.SetTitle(title)
	form.SetCancelFunc(cancel)

	return ext.NewForm(form, dataProvider)
}

// showCreateForm shows create form with initialized empty fields.
func (v *View) showCreateForm() {
	d := time.Now().Format("2006-01-02")
	m := map[string]string{"Date": d, "Account": "", "Category": "", "Amount": "", "Note": ""}

	v.createForm.SetFields(m)
	v.Pages.ShowPage("createForm")
}

// hideCreateForm hides create form.
func (v *View) hideCreateForm() {
	v.Pages.HidePage("createForm")
}

// isValidCreateForm checks if all necessary fields are filled.
func (v *View) isValidCreateForm(m map[string]string) bool {
	value, ok := m["Account"]
	if !ok || value == "" {
		v.showError("Can't create transaction without account.")
		return false
	}
	value, ok = m["Category"]
	if !ok || value == "" {
		v.showError("Can't create transaction without category.")
		return false
	}
	value, ok = m["Amount"]
	if !ok || value == "" {
		v.showError("Can't create transaction without amount.")
		return false
	}

	return true
}

// submitCreateForm create form submit handler.
func (v *View) submitCreateForm() {
	m := v.createForm.GetFields()
	if !v.isValidCreateForm(m) {
		return
	}

	tr, err := v.presenter.Transaction().FromMap(m)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Transaction().Insert(tr); err != nil {
		v.showError("Error insert transaction: \n" + err.Error())
		return
	}

	v.table.Refresh()
	v.hideCreateForm()
}

// showUpdateForm shows update form with initialized with selected transaction fields.
func (v *View) showUpdateForm() {
	m := v.table.GetSelectedRef()
	v.updateForm.SetFields(m)
	v.Pages.ShowPage("updateForm")
}

// hideUpdateForm hides update form.
func (v *View) hideUpdateForm() {
	v.Pages.HidePage("updateForm")
}

// submitUpdateForm update form submit handler.
func (v *View) submitUpdateForm() {
	m := v.updateForm.GetFields()
	if !v.isValidCreateForm(m) {
		return
	}

	ref := v.table.GetSelectedRef()
	m["ID"] = ref["ID"]

	tr, err := v.presenter.Transaction().FromMap(m)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Transaction().Update(tr); err != nil {
		v.showError("Error update transaction: \n" + err.Error())
		return
	}

	v.table.Refresh()
	v.hideUpdateForm()
}

// showDeleteModal shows delete modal.
func (v *View) showDeleteModal() {
	v.Pages.ShowPage("deleteModal")
}

// hideDeleteModal hides delete modal.
func (v *View) hideDeleteModal() {
	v.Pages.HidePage("deleteModal")
}

// submitDeleteModal delete modal submit handler.
func (v *View) submitDeleteModal() {
	ref := v.table.GetSelectedRef()
	tr, err := v.presenter.Transaction().FromMap(ref)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Transaction().Delete(tr); err != nil {
		v.showError("Error update transaction: \n" + err.Error())
		return
	}

	v.table.Refresh()
	v.hideDeleteModal()
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
