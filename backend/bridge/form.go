package bridge

import (
	"UppSpar/backend"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type Form struct {
	Container *container.Scroll
	Checks    Checks
	Entries   Entries
	Labels    Labels
	Radios    Radios
	Selects   Selects
	Values    Labels
}

func NewForm(b *backend.Backend, w fyne.Window) *Form {
	c := container.New(layout.NewFormLayout())
	f := &Form{
		Container: container.NewVScroll(c),
		Checks:    make(Checks),
		Entries:   make(Entries),
		Labels:    make(Labels),
		Radios:    make(Radios),
		Selects:   make(Selects),
		Values:    make(Labels),
	}
	return f
}

func (f *Form) Clear() {
	f.Checks.Uncheck()
	f.Entries.Clear()
	f.Radios.Uncheck()
	f.Selects.Clear()
	f.Values.Clear()
}
func (f *Form) Disable() {
	f.Checks.Disable()
	f.Entries.Disable()
	f.Radios.Disable()
	f.Selects.Disable()
}
func (f *Form) Enable() {
	f.Checks.Enable()
	f.Entries.Enable()
	f.Radios.Enable()
	f.Selects.Enable()
}
func (f *Form) LoadItem(id backend.ItemID) {
	f.Entries.Bind(Sieve(id.Item().Bindings(), ItemFormEntryKeys))
}
func (f *Form) LoadMfr(id backend.MfrID) {
	f.Entries.Bind(Sieve(id.Manufacturer().Bindings(), ManufacturerFormEntryKeys))
}
func (f *Form) LoadModel(id backend.ModelID) {
	f.Entries.Bind(Sieve(id.Model().Bindings(), ModelFormEntryKeys))
}
func (f *Form) LoadCategory(id backend.CatID) {
	f.Entries.Bind(Sieve(id.Category().Bindings(), CategoryFormEntryKeys))
}
