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
func NewForm(dataProvider FormDataProvider) *Form {
	return &Form{
		Form: tview.NewForm(),

		inputFields:  make(map[string]*tview.InputField),
		dropDowns:    make(map[string]*tview.DropDown),
		dataProvider: dataProvider,
	}
}

// AddInputField adds an input field to the form. For more information read tview.Form.AddInputField doc.
func (f *Form) AddInputField(label, value string, fieldWidth int, accept func(textToCheck string, lastChar rune) bool, changed func(text string)) *Form {
	f.Form.AddInputField(label, value, fieldWidth, accept, changed)
	f.inputFields[label] = f.Form.GetFormItemByLabel(label).(*tview.InputField)
	return f
}

// AddDropDown adds a drop-down element to the form. For more information read tview.Form.AddDropDown doc.
func (f *Form) AddDropDown(label string, options []string, initialOption int, selected func(option string, optionIndex int)) *Form {
	f.Form.AddDropDown(label, options, initialOption, selected)
	f.dropDowns[label] = f.Form.GetFormItemByLabel(label).(*tview.DropDown)
	return f
}

// AddButton adds a new button to the form. For more information read tview.Form.AddButton doc.
func (f *Form) AddButton(label string, selected func()) *Form {
	f.Form.AddButton(label, selected)
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
