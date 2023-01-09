package settings

import (
	"strings"

	"github.com/kotlw/gentlemoney/internal/tui/ext"

	"github.com/rivo/tview"
)

// newForm returns new form with corresponding category fields.
func (v *View) newCategoryForm(title string, submit func(), cancel func(), dataProvider *CategoryDataProvider) *ext.Form {
	form := tview.NewForm().
		AddInputField("Title", "", 0, nil, nil).
		AddButton(strings.Split(title, " ")[0], submit).
		AddButton("Cancel", cancel)

	form.SetBorder(true)
	form.SetTitle(title)
	form.SetCancelFunc(cancel)

	return ext.NewForm(form, dataProvider)
}

// showCategoryCreateForm shows category create form with initialized empty fields.
func (v *View) showCategoryCreateForm() {
	v.categoryCreateForm.SetFields(map[string]string{"Title": ""})
	v.Pages.ShowPage("categoryCreateForm")
}

// hideCategoryCreateForm hides category create form.
func (v *View) hideCategoryCreateForm() {
	v.Pages.HidePage("categoryCreateForm")
	v.app.SetFocus(v.categoryTable)
}

// isValidCreateForm checks if all necessary fields are filled.
func (v *View) isValidCategoryCreateForm(m map[string]string) bool {
	value, ok := m["Title"]
	if !ok || value == "" {
		v.showError("Can't create category without title.")
		return false
	}

	return true
}

// submitCategoryCreateFormcategory create form submit handler.
func (v *View) submitCategoryCreateForm() {
	m := v.categoryCreateForm.GetFields()
	if !v.isValidCategoryCreateForm(m) {
		return
	}

	c, err := v.presenter.Category().FromMap(m)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Category().Insert(c); err != nil {
		v.showError("Error insert transaction: \n" + err.Error())
		return
	}

	v.categoryTable.Refresh()
	v.hideCategoryCreateForm()
}

// showCategoryUpdateForm shows update form with initialized with selected category fields.
func (v *View) showCategoryUpdateForm() {
	m := v.categoryTable.GetSelectedRef()
	v.categoryUpdateForm.SetFields(m)
	v.Pages.ShowPage("categoryUpdateForm")
}

// hideCategoryUpdateForm hides update form.
func (v *View) hideCategoryUpdateForm() {
	v.Pages.HidePage("categoryUpdateForm")
	v.app.SetFocus(v.categoryTable)
}

// submitCategoryUpdateForm update form submit handler.
func (v *View) submitCategoryUpdateForm() {
	m := v.categoryUpdateForm.GetFields()
	if !v.isValidCategoryCreateForm(m) {
		return
	}

	ref := v.categoryTable.GetSelectedRef()
	m["ID"] = ref["ID"]

	c, err := v.presenter.Category().FromMap(m)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Category().Update(c); err != nil {
		v.showError("Error update transaction: \n" + err.Error())
		return
	}

	v.categoryTable.Refresh()
	v.hideCategoryCreateForm()
}

// showCategoryDeleteModal shows delete modal.
func (v *View) showCategoryDeleteModal() {
	v.Pages.ShowPage("categoryDeleteModal")
}

// hideCategoryDeleteModal hides delete modal.
func (v *View) hideCategoryDeleteModal() {
	v.Pages.HidePage("categoryDeleteModal")
	v.app.SetFocus(v.categoryTable)
}

// submitCategoryDeleteModaldelete modal submit handler.
func (v *View) submitCategoryDeleteModal() {
	ref := v.categoryTable.GetSelectedRef()
	c, err := v.presenter.Category().FromMap(ref)
	if err != nil {
		v.showError("Error parse form: \n" + err.Error())
		return
	}

	if err := v.service.Category().Delete(c); err != nil {
		v.showError("Error update transaction: \n" + err.Error())
		return
	}

	v.categoryTable.Refresh()
	v.hideCategoryDeleteModal()
}
