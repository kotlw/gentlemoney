package settings

import (
	"strings"

	"github.com/kotlw/gentlemoney/internal/tui/ext"

	"github.com/rivo/tview"
)

// newForm returns new form with corresponding currency fields.
func (v *View) newCurrencyForm(title string, submit func(), cancel func(), dataProvider *CurrencyDataProvider) *ext.Form {
	form := tview.NewForm().
		AddInputField("Abbreviation", "", 0, nil, nil).
		AddInputField("ExchangeRate", "", 0, nil, nil).
		AddButton(strings.Split(title, " ")[0], submit).
		AddButton("Cancel", cancel)

	form.SetBorder(true)
	form.SetTitle(title)
	form.SetCancelFunc(cancel)

	return ext.NewForm(form, dataProvider)
}

// showCurrencyCreateForm shows currency create form with initialized empty fields.
func (v *View) showCurrencyCreateForm() {
	v.currencyCreateForm.SetFields(map[string]string{"Abbreviation": "", "ExchangeRate": ""})
	v.Pages.ShowPage("currencyCreateForm")
}

// hideCurrencyCreateForm hides currency create form.
func (v *View) hideCurrencyCreateForm() {
	v.Pages.HidePage("currencyCreateForm")
	v.app.SetFocus(v.currencyTable)
}

// isValidCreateForm checks if all necessary fields are filled.
func (v *View) isValidCurrencyCreateForm(m map[string]string) bool {
	value, ok := m["Abbreviation"]
	if !ok || value == "" {
		v.showError("Can't create currency without abbreviation.")
		return false
	}

	return true
}

// submitCurrencyCreateFormcurrency create form submit handler.
func (v *View) submitCurrencyCreateForm() {
	m := v.currencyCreateForm.GetFields()
	if !v.isValidCurrencyCreateForm(m) {
		return
	}

	c, err := v.presenter.Currency().FromMap(m)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Currency().Insert(c); err != nil {
		v.showError("Error insert transaction: \n" + err.Error())
		return
	}

	v.currencyTable.Refresh()
	v.hideCurrencyCreateForm()
}

// showCurrencyUpdateForm shows update form with initialized with selected currency fields.
func (v *View) showCurrencyUpdateForm() {
	m := v.currencyTable.GetSelectedRef()
	v.currencyUpdateForm.SetFields(m)
	v.Pages.ShowPage("currencyUpdateForm")
}

// hideCurrencyUpdateForm hides update form.
func (v *View) hideCurrencyUpdateForm() {
	v.Pages.HidePage("currencyUpdateForm")
	v.app.SetFocus(v.currencyTable)
}

// submitCurrencyUpdateForm update form submit handler.
func (v *View) submitCurrencyUpdateForm() {
	m := v.currencyUpdateForm.GetFields()
	if !v.isValidCurrencyCreateForm(m) {
		return
	}

	ref := v.currencyTable.GetSelectedRef()
	m["ID"] = ref["ID"]

	c, err := v.presenter.Currency().FromMap(m)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Currency().Update(c); err != nil {
		v.showError("Error update transaction: \n" + err.Error())
		return
	}

	v.currencyTable.Refresh()
	v.hideCurrencyCreateForm()
}

// showCurrencyDeleteModal shows delete modal.
func (v *View) showCurrencyDeleteModal() {
	v.Pages.ShowPage("deleteCurrencyModal")
}

// hideCurrencyDeleteModal hides delete modal.
func (v *View) hideCurrencyDeleteModal() {
	v.Pages.HidePage("deleteCurrencyModal")
	v.app.SetFocus(v.currencyTable)
}

// submitCurrencyDeleteModaldelete modal submit handler.
func (v *View) submitCurrencyDeleteModal() {
	ref := v.currencyTable.GetSelectedRef()
	c, err := v.presenter.Currency().FromMap(ref)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Currency().Delete(c); err != nil {
		v.showError("Error update transaction: \n" + err.Error())
		return
	}

	v.currencyTable.Refresh()
	v.hideCurrencyDeleteModal()
}
