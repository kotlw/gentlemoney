package ext

import (
	"sort"

	"github.com/rivo/tview"
)

// FormDataProvider an interface for getting options for dropdown fields.
type FormDataProvider interface {
	GetDropDownOptions(label string) []string
}

// Form an extension for tview.Form. It used to avoid duplication while setting and getting field values.
type Form struct {
	*tview.Form

	inputFields  map[string]*tview.InputField
	dropDowns    map[string]*tview.DropDown
	dataProvider FormDataProvider
}

// NewForm returns new extended Form.
func NewForm(form *tview.Form, dataProvider FormDataProvider) *Form {
	f := &Form{
		Form: form,

		inputFields:  make(map[string]*tview.InputField),
		dropDowns:    make(map[string]*tview.DropDown),
		dataProvider: dataProvider,
	}

	// gather form items into groups
	for i := 0; i < form.GetFormItemCount(); i++ {
		item := form.GetFormItem(i)

		inputField, ok := item.(*tview.InputField)
		if ok {
			f.inputFields[item.GetLabel()] = inputField
		}
		dropDown, ok := item.(*tview.DropDown)
		if ok {
			f.dropDowns[item.GetLabel()] = dropDown
		}
	}

	return f
}

// SetFields sets fields value from map where key is field label and value is a value.
func (f *Form) SetFields(m map[string]string) {
	for label, value := range m {

		field, ok := f.inputFields[label]
		if ok {
			field.SetText(value)
			continue
		}

		dropDown, ok := f.dropDowns[label]
		if ok {
			opts := f.dataProvider.GetDropDownOptions(label)

			index := -1
			if value != "" {
				index = sort.SearchStrings(opts, value)
			}

			dropDown.SetOptions(opts, nil).SetCurrentOption(index)
			continue
		}
	}
	f.Form.SetFocus(0)
}

// GetFields returns fields values as map of strings where the key is field label.
func (f *Form) GetFields() map[string]string {
	res := make(map[string]string)

	for label, field := range f.inputFields {
		res[label] = field.GetText()
	}

	for label, dropDown := range f.dropDowns {
		_, res[label] = dropDown.GetCurrentOption()
	}

	return res
}
