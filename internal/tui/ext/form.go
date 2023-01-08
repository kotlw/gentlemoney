package ext

import (
	"sort"

	"github.com/rivo/tview"
)

type FormDataProvider interface {
	GetDropDownOptions(label string) []string
}

type Form struct {
	*tview.Form

	inputFields  map[string]*tview.InputField
	dropDowns    map[string]*tview.DropDown
	dataProvider FormDataProvider
}

func NewForm(dataProvider FormDataProvider) *Form {
	return &Form{
		Form: tview.NewForm(),

		inputFields:  make(map[string]*tview.InputField),
		dropDowns:    make(map[string]*tview.DropDown),
		dataProvider: dataProvider,
	}
}

func (f *Form) AddInputField(label, value string, fieldWidth int, accept func(textToCheck string, lastChar rune) bool, changed func(text string)) *Form {
	f.Form.AddInputField(label, value, fieldWidth, accept, changed)
	f.inputFields[label] = f.Form.GetFormItemByLabel(label).(*tview.InputField)
	return f
}

func (f *Form) AddDropDown(label string, options []string, initialOption int, selected func(option string, optionIndex int)) *Form {
	f.Form.AddDropDown(label, options, initialOption, selected)
	f.dropDowns[label] = f.Form.GetFormItemByLabel(label).(*tview.DropDown)
	return f
}

func (f *Form) AddButton(label string, selected func()) *Form {
	f.Form.AddButton(label, selected)
	return f
}

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
