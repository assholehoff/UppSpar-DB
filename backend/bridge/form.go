package bridge

import (
	"UppSpar/backend"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	midget "github.com/assholehoff/fyne-midget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type Form struct {
	Container *container.Scroll
	Check     Checks
	Entry     Entries
	Label     Labels
	Radio     Radios
	Select    Selects
	Value     Labels
}

func NewForm(b *backend.Backend, w fyne.Window) *Form {
	c := container.New(layout.NewFormLayout())

	f := &Form{
		Container: container.NewVScroll(c),
		Check:     make(Checks),
		Entry:     make(Entries),
		Label:     make(Labels),
		Radio:     make(Radios),
		Select:    make(Selects),
		Value:     make(Labels),
	}

	for _, key := range Combine(
		CategoryFormCheckKeys,
		ItemFormCheckKeys,
		ManufacturerFormCheckKeys,
		ModelFormCheckKeys,
	) {
		f.Check[key] = ttw.NewCheck("Template", func(b bool) {})
	}

	for _, key := range Combine(
		CategoryFormEntryKeys,
		ItemFormEntryKeys,
		ManufacturerFormEntryKeys,
		ModelFormEntryKeys,
	) {
		f.Entry[key] = midget.NewEntry()
	}

	for _, key := range Combine(
		CategoryFormLabelKeys,
		ItemFormLabelKeys,
		ManufacturerFormLabelKeys,
		ModelFormLabelKeys,
	) {
		f.Label[key] = ttw.NewLabel("Template label")
	}

	for _, key := range Combine(
		CategoryFormRadioKeys,
		ItemFormRadioKeys,
		ManufacturerFormRadioKeys,
		ModelFormRadioKeys,
	) {
		f.Radio[key] = widget.NewRadioGroup([]string{}, func(s string) {})
	}

	for _, key := range Combine(
		CategoryFormSelectKeys,
		ItemFormSelectKeys,
		ManufacturerFormSelectKeys,
		ModelFormSelectKeys,
	) {
		f.Select[key] = ttw.NewSelect([]string{}, func(s string) {})
	}

	for _, key := range Combine(
		CategoryFormValuesKeys,
		ItemFormValueKeys,
		ManufacturerFormValuesKeys,
		ModelFormValuesKeys,
	) {
		f.Value[key] = ttw.NewLabel("Template value label text")
	}

	return f
}

func (f *Form) Clear() {
	f.Check.Uncheck()
	f.Entry.Clear()
	f.Radio.Uncheck()
	f.Select.Clear()
	f.Value.Clear()
}
func (f *Form) Disable() {
	f.Check.Disable()
	f.Entry.Disable()
	f.Radio.Disable()
	f.Select.Disable()
}
func (f *Form) Enable() {
	f.Check.Enable()
	f.Entry.Enable()
	f.Radio.Enable()
	f.Select.Enable()
}
func (f *Form) LoadItem(id backend.ItemID) {
	f.Entry.Bind(Sieve(id.Item().Bindings(), ItemFormEntryKeys))
}
func (f *Form) LoadMfr(id backend.MfrID) {
	f.Entry.Bind(Sieve(id.Manufacturer().Bindings(), ManufacturerFormEntryKeys))
}
func (f *Form) LoadModel(id backend.ModelID) {
	f.Entry.Bind(Sieve(id.Model().Bindings(), ModelFormEntryKeys))
}
func (f *Form) LoadCategory(id backend.CatID) {
	f.Entry.Bind(Sieve(id.Category().Bindings(), CategoryFormEntryKeys))
}
