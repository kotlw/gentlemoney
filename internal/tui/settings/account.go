package settings

import (
	"strings"

	"github.com/kotlw/gentlemoney/internal/tui/ext"

	"github.com/rivo/tview"
)

// newForm returns new form with corresponding account fields.
func (v *View) newAccountForm(title string, submit func(), cancel func(), dataProvider *AccountDataProvider) *ext.Form {
	form := tview.NewForm().
		AddInputField("Name", "", 0, nil, nil).
		AddDropDown("Currency", nil, 0, nil).
		AddButton(strings.Split(title, " ")[0], submit).
		AddButton("Cancel", cancel)

	form.SetBorder(true)
	form.SetTitle(title)
	form.SetCancelFunc(cancel)

	return ext.NewForm(form, dataProvider)
}

// showAccountCreateForm shows account create form with initialized empty fields.
func (v *View) showAccountCreateForm() {
	v.accountCreateForm.SetFields(map[string]string{"Name": "", "Currency": ""})
	v.Pages.ShowPage("accountCreateForm")
}

// hideAccountCreateForm hides account create form.
func (v *View) hideAccountCreateForm() {
	v.Pages.HidePage("accountCreateForm")
	v.tuiApp.SetFocus(v.accountTable)
}

// isValidCreateForm checks if all necessary fields are filled.
func (v *View) isValidAccountCreateForm(m map[string]string) bool {
	value, ok := m["Name"]
	if !ok || value == "" {
		v.showError("Can't create account without name.")
		return false
	}

	value, ok = m["Currency"]
	if !ok || value == "" {
		v.showError("Can't create account without currency.")
		return false
	}

	return true
}

// submitAccountCreateFormaccount create form submit handler.
func (v *View) submitAccountCreateForm() {
	m := v.accountCreateForm.GetFields()
	if !v.isValidAccountCreateForm(m) {
		return
	}

	c, err := v.presenter.Account().FromMap(m)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Account().Insert(c); err != nil {
		v.showError("Error insert transaction: \n" + err.Error())
		return
	}

	v.accountTable.Refresh()
	v.hideAccountCreateForm()
}

// showAccountUpdateForm shows update form with initialized with selected account fields.
func (v *View) showAccountUpdateForm() {
	m := v.accountTable.GetSelectedRef()
	v.accountUpdateForm.SetFields(m)
	v.Pages.ShowPage("accountUpdateForm")
}

// hideAccountUpdateForm hides update form.
func (v *View) hideAccountUpdateForm() {
	v.Pages.HidePage("accountUpdateForm")
	v.tuiApp.SetFocus(v.accountTable)
}

// submitAccountUpdateForm update form submit handler.
func (v *View) submitAccountUpdateForm() {
	m := v.accountUpdateForm.GetFields()
	if !v.isValidAccountCreateForm(m) {
		return
	}

	ref := v.accountTable.GetSelectedRef()
	m["ID"] = ref["ID"]

	c, err := v.presenter.Account().FromMap(m)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Account().Update(c); err != nil {
		v.showError("Error update transaction: \n" + err.Error())
		return
	}

	v.accountTable.Refresh()
	v.hideAccountCreateForm()
}

// showAccountDeleteModal shows delete modal.
func (v *View) showAccountDeleteModal() {
	v.Pages.ShowPage("accountDeleteModal")
}

// hideAccountDeleteModal hides delete modal.
func (v *View) hideAccountDeleteModal() {
	v.Pages.HidePage("accountDeleteModal")
	v.tuiApp.SetFocus(v.accountTable)
}

// submitAccountDeleteModaldelete modal submit handler.
func (v *View) submitAccountDeleteModal() {
	ref := v.accountTable.GetSelectedRef()
	c, err := v.presenter.Account().FromMap(ref)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Account().Delete(c); err != nil {
		v.showError("Error update transaction: \n" + err.Error())
		return
	}

	v.accountTable.Refresh()
	v.hideAccountDeleteModal()
}
