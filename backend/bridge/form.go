package bridge

import (
	"UppSpar/backend"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var (
	CategoryFormCheckKeys  = []string{}
	CategoryFormEntryKeys  = []string{}
	CategoryFormLabelKeys  = []string{}
	CategoryFormRadioKeys  = []string{}
	CategoryFormSelectKeys = []string{}
	CategoryFormValuesKeys = []string{}

	ItemFormCheckKeys = []string{
		"New",
		"Working",
	}
	ItemFormEntryKeys = []string{
		"Name",
		"Price",
		"Vat",
		"ImgURL1",
		"ImgURL2",
		"ImgURL3",
		"ImgURL4",
		"ImgURL5",
		"SpecsURL",
		"LongDesc",
		"Manufacturer",
		"ModelDesc",
		"ModelURL",
		"Notes",
		"Width",
		"Height",
		"Depth",
		"Volume",
		"Weight",
	}
	ItemFormLabelKeys = []string{
		"ItemID",
		"Name",
		"Category",
		"Currency",
		"Price",
		"Vat",
		"ImgURL1",
		"ImgURL2",
		"ImgURL3",
		"ImgURL4",
		"ImgURL5",
		"SpecsURL",
		"AddDesc",
		"LongDesc",
		"Manufacturer",
		"ModelName",
		"ModelDesc",
		"ModelURL",
		"Notes",
		"Dimensions",
		"Width",
		"Height",
		"Depth",
		"Volume",
		"Weight",
		"Status",
		"DateCreated",
		"DateModified",
		"Condition",
		"Functionality",
	}
	ItemFormRadioKeys = []string{
		"Tested",
	}
	ItemFormSelectKeys = []string{
		"Category",
		"Manufacturer",
		"ModelName",
		"LengthUnit",
		"VolumeUnit",
		"WeightUnit",
		"Status",
	}
	ItemFormValueKeys = []string{
		"ItemID",
		"AddDesc",
		"LongDesc",
		"DateCreated",
		"DateModified",
	}

	ManufacturerFormCheckKeys  = []string{}
	ManufacturerFormEntryKeys  = []string{}
	ManufacturerFormLabelKeys  = []string{}
	ManufacturerFormRadioKeys  = []string{}
	ManufacturerFormSelectKeys = []string{}
	ManufacturerFormValuesKeys = []string{}

	ModelFormCheckKeys  = []string{}
	ModelFormEntryKeys  = []string{}
	ModelFormLabelKeys  = []string{}
	ModelFormRadioKeys  = []string{}
	ModelFormSelectKeys = []string{}
	ModelFormValuesKeys = []string{}
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

func NewFormView() *Form {
	c := container.New(layout.NewFormLayout())
	return &Form{
		Container: container.NewVScroll(c),
		Checks:    make(Checks),
		Entries:   make(Entries),
		Labels:    make(Labels),
		Radios:    make(Radios),
		Selects:   make(Selects),
		Values:    make(Labels),
	}
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
