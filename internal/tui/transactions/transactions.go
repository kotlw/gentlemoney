package transactions

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/tui/ext"

	"github.com/rivo/tview"
)

type View struct {
	*tview.Pages

	service     *service.Service
	presenter   *presenter.Presenter
	table       *ext.Table
	createForm  *ext.Form
	updateForm  *ext.Form
	deleteModal *tview.Modal
	errorModal  *tview.Modal
}

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
	v.createForm = v.NewForm("Create Transaction", v.CreateFormSubmit, v.HideCreateForm, dataProvider)
	v.AddPage("createForm", ext.WrapIntoModal(v.createForm, 40, 15), true, false)

	// update form
	v.updateForm = v.NewForm("Update Transaction", v.UpdateFormSubmit, v.HideUpdateForm, dataProvider)
	v.AddPage("updateForm", ext.WrapIntoModal(v.updateForm, 40, 15), true, false)

	// delete modal
	v.deleteModal = ext.NewAskModal("Are you sure?", v.DeleteModalSubmit, v.HideDeleteModal)
	v.AddPage("deleteModal", v.deleteModal, true, false)

	// error modal
	v.errorModal = ext.NewErrorModal(v.HideError)
	v.AddPage("errorModal", v.errorModal, true, false)

	return v
}

func (v *View) ModalHasFocus() bool {
	return v.createForm.HasFocus()
}

func (v *View) NewForm(title string, submit func(), cancel func(), dataProvider *DataProvider) *ext.Form {
	form := ext.NewForm(dataProvider).
		AddInputField("Date", "", 0, nil, nil).
		AddDropDown("Category", nil, 0, nil).
		AddDropDown("Account", nil, 0, nil).
		AddInputField("Amount", "", 0, nil, nil).
		AddInputField("Note", "", 0, nil, nil).
		AddButton("Create", submit).
		AddButton("Cancel", cancel)

	form.SetBorder(true)
	form.SetTitle(title)
	form.SetCancelFunc(cancel)

	return form
}

func (v *View) ShowCreateForm() {
	d := time.Now().Format("2006-01-02")
	m := map[string]string{"Date": d, "Account": "", "Category": "", "Amount": "", "Note": ""}

	v.createForm.SetFields(m)
	v.Pages.ShowPage("createForm")
}

func (v *View) HideCreateForm() {
	v.Pages.HidePage("createForm")
}

func (v *View) CreateFormSubmit() {
	m := v.createForm.GetFields()
	value, ok := m["Account"]
	if !ok || value == "" {
		v.ShowError("Can't create transaction without account.")
		return
	}
	value, ok = m["Category"]
	if !ok || value == "" {
		v.ShowError("Can't create transaction without category.")
		return
	}

	tr, err := v.presenter.Transaction().FromMap(m)
	if err != nil {
		v.ShowError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Transaction().Insert(tr); err != nil {
		v.ShowError("Error insert transaction: \n" + err.Error())
		return
	}

	v.table.Refresh()
	v.HideCreateForm()
}

func (v *View) ShowUpdateForm() {
	m := v.table.GetSelectedRef()
	v.updateForm.SetFields(m)
	v.Pages.ShowPage("updateForm")
}

func (v *View) HideUpdateForm() {
	v.Pages.HidePage("updateForm")
}

func (v *View) UpdateFormSubmit() {
	ref := v.table.GetSelectedRef()
	tr, err := v.presenter.Transaction().FromMap(ref)
	if err != nil {
		v.ShowError("Error parse form: \n" + err.Error())
		return
	}

	id, err := strconv.Atoi(ref["ID"])
	if err != nil {
		v.ShowError("Internal error: something wrong with getting transaction ID from reference.")
	}
	tr.ID = int64(id)

	if err := v.service.Transaction().Update(tr); err != nil {
		v.ShowError("Error update transaction: \n" + err.Error())
		return
	}

	v.table.Refresh()
	v.HideUpdateForm()
}

func (v *View) ShowDeleteModal() {
	v.Pages.ShowPage("deleteModal")
}

func (v *View) HideDeleteModal() {
	v.Pages.HidePage("deleteModal")
}

func (v *View) DeleteModalSubmit() {
	ref := v.table.GetSelectedRef()
	tr, err := v.presenter.Transaction().FromMap(ref)
	if err != nil {
		v.ShowError("Error parse form: \n" + err.Error())
		return
	}

	id, err := strconv.Atoi(ref["ID"])
	if err != nil {
		v.ShowError("Internal error: something wrong with getting transaction ID from reference.")
	}
	tr.ID = int64(id)

	if err := v.service.Transaction().Delete(tr); err != nil {
		v.ShowError("Error update transaction: \n" + err.Error())
		return
	}

	v.table.Refresh()
	v.HideDeleteModal()
}

func (v *View) ShowError(text string) {
	v.errorModal.SetText(text)
	v.Pages.ShowPage("errorModal")
}

func (v *View) HideError() {
	v.Pages.HidePage("errorModal")
}

func (v *View) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if v.table.HasFocus() {
			if event.Rune() == 'c' {
				v.ShowCreateForm()
			}

			if event.Rune() == 'u' {
				v.ShowUpdateForm()
			}

			if event.Rune() == 'd' {
				v.ShowDeleteModal()
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
