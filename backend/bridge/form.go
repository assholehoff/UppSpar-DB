package bridge

import (
	"UppSpar/backend"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
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

func NewItemForm(b *backend.Backend, w fyne.Window) *Form {
	initItemStringMaps()

	f := &Form{
		Check:  make(Checks),
		Entry:  make(Entries),
		Label:  make(Labels),
		Radio:  make(Radios),
		Select: make(Selects),
		Value:  make(Labels),
	}

	for _, key := range Combine(
		CategoryFormCheckKeys,
		ItemFormCheckKeys,
	) {
		f.Check[key] = ttw.NewCheck("Template", func(b bool) {})
	}

	for _, key := range Combine(
		CategoryFormEntryKeys,
		ItemFormEntryKeys,
	) {
		f.Entry[key] = midget.NewEntry()
	}
	f.Entry["ModelDesc"].MultiLine = true
	f.Entry["ModelDesc"].SetMinRowsVisible(5)
	f.Entry["ModelDesc"].Wrapping = fyne.TextWrapWord
	f.Entry["LongDesc"].MultiLine = true
	f.Entry["LongDesc"].SetMinRowsVisible(5)
	f.Entry["LongDesc"].Wrapping = fyne.TextWrapWord
	f.Entry["Notes"].MultiLine = true
	f.Entry["Notes"].SetMinRowsVisible(5)
	f.Entry["Notes"].Wrapping = fyne.TextWrapWord

	for _, key := range Combine(
		CategoryFormLabelKeys,
		ItemFormLabelKeys,
	) {
		f.Label[key] = ttw.NewLabel(key + " label")
	}
	f.Label.Set(ItemFormLabelStrings)

	for _, key := range Combine(
		CategoryFormRadioKeys,
		ItemFormRadioKeys,
	) {
		f.Radio[key] = widget.NewRadioGroup([]string{}, func(s string) {})
	}

	for _, key := range Combine(
		CategoryFormSelectKeys,
		ItemFormSelectKeys,
	) {
		f.Select[key] = ttw.NewSelect([]string{}, func(s string) {})
	}
	f.Select["Status"].SetOptions(b.Metadata.ItemStatusList())
	f.Select["LengthUnit"].SetOptions([]string{"mm", "cm", "dm", "m"})
	f.Select["VolumeUnit"].SetOptions([]string{"ml", "cl", "dl", "l"})
	f.Select["WeightUnit"].SetOptions([]string{"g", "hg", "kg"})
	b.Metadata.Categories.AddListener(binding.NewDataListener(func() {
		cats, _ := b.Metadata.Categories.Get()
		f.Select["Category"].SetOptions(cats)
	}))
	b.Metadata.MfrNameList.AddListener(binding.NewDataListener(func() {
		manufacturers, _ := b.Metadata.MfrNameList.Get()
		f.Select["Manufacturer"].SetOptions(manufacturers)
	}))

	for _, key := range Combine(
		CategoryFormValuesKeys,
		ItemFormValueKeys,
	) {
		f.Value[key] = ttw.NewLabel(key + " template")
	}
	f.Value["LongDesc"].Wrapping = fyne.TextWrapWord
	f.Value["AddDesc"].Wrapping = fyne.TextWrapWord
	f.Value.Set(ItemFormValueStrings)

	f.Value["DateCreated"].Hide()
	f.Value["DateModified"].Hide()
	f.Value["AddDesc"].Hide()
	f.Value["LongDesc"].Hide()

	f.Container = container.NewVScroll(f.itemContainer())

	f.Check.Disable()
	f.Entry.Disable()
	f.Select.Disable()

	return f
}
func NewProductForm(b *backend.Backend, w fyne.Window) *Form {
	initProductStringMaps()

	f := &Form{
		Check:  make(Checks),
		Entry:  make(Entries),
		Label:  make(Labels),
		Radio:  make(Radios),
		Select: make(Selects),
		Value:  make(Labels),
	}

	for _, key := range Combine(ManufacturerFormCheckKeys, ModelFormCheckKeys) {
		f.Check[key] = ttw.NewCheck("Template", func(b bool) {})
	}

	for _, key := range Combine(ManufacturerFormEntryKeys, ModelFormEntryKeys) {
		f.Entry[key] = midget.NewEntry()
	}
	f.Entry["ModelDesc"].MultiLine = true
	f.Entry["ModelDesc"].SetMinRowsVisible(5)
	f.Entry["ModelDesc"].Wrapping = fyne.TextWrapWord

	for _, key := range Combine(ManufacturerFormLabelKeys, ModelFormLabelKeys) {
		f.Label[key] = ttw.NewLabel(key + " label")
	}
	f.Label.Set(ProductFormLabelStrings)

	for _, key := range Combine(ManufacturerFormSelectKeys, ModelFormSelectKeys) {
		f.Select[key] = ttw.NewSelect([]string{}, func(s string) {})
	}
	f.Select["Status"].SetOptions(b.Metadata.ItemStatusList())
	f.Select["LengthUnit"].SetOptions([]string{"mm", "cm", "dm", "m"})
	f.Select["VolumeUnit"].SetOptions([]string{"ml", "cl", "dl", "l"})
	f.Select["WeightUnit"].SetOptions([]string{"g", "hg", "kg"})
	b.Metadata.Categories.AddListener(binding.NewDataListener(func() {
		cats, _ := b.Metadata.Categories.Get()
		f.Select["Category"].SetOptions(cats)
	}))
	b.Metadata.MfrNameList.AddListener(binding.NewDataListener(func() {
		manufacturers, _ := b.Metadata.MfrNameList.Get()
		f.Select["Manufacturer"].SetOptions(manufacturers)
	}))

	for _, key := range Combine(ManufacturerFormValuesKeys, ModelFormValuesKeys) {
		f.Value[key] = ttw.NewLabel(key + " value template")
	}

	f.Container = container.NewVScroll(f.productContainer())

	f.Check.Disable()
	f.Entry.Disable()
	f.Select.Disable()

	return f
}

func (f *Form) Clear() {
	f.Check.Unbind()
	f.Entry.Unbind()
	f.Select.Unbind()
	f.Value.Unbind()

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
func (f *Form) LoadItem(b *backend.Backend, id backend.ItemID) {
	var enabled, disabled []string

	f.Clear()

	f.Value["DateCreated"].Hide()
	f.Value["DateModified"].Hide()
	f.Value["AddDesc"].Hide()
	f.Value["LongDesc"].Hide()

	/* Label widgets */
	f.Value["ItemID"].Bind(id.Item().ItemIDString)
	f.Value["DateCreated"].Bind(id.Item().DateCreated)
	f.Value["DateModified"].Bind(id.Item().DateModified)
	f.Value["AddDesc"].Bind(id.Item().AddDesc)
	f.Value["LongDesc"].Bind(id.Item().LongDesc)

	f.Value["DateCreated"].Show()
	f.Value["DateModified"].Show()
	f.Value["AddDesc"].Show()
	f.Value["LongDesc"].Show()

	/* Entry widgets */
	f.Entry["Name"].Bind(id.Item().Name)
	f.Entry["Price"].Bind(id.Item().PriceString)
	f.Entry["Vat"].Bind(id.Item().VatString)
	f.Entry["ImgURL1"].Bind(id.Item().ImgURL1)
	f.Entry["ImgURL2"].Bind(id.Item().ImgURL2)
	f.Entry["ImgURL3"].Bind(id.Item().ImgURL3)
	f.Entry["ImgURL4"].Bind(id.Item().ImgURL4)
	f.Entry["ImgURL5"].Bind(id.Item().ImgURL5)
	f.Entry["SpecsURL"].Bind(id.Item().SpecsURL)
	f.Entry["ModelDesc"].Bind(id.Item().ModelDesc)
	f.Entry["ModelURL"].Bind(id.Item().ModelURL)
	f.Entry["Notes"].Bind(id.Item().Notes)
	f.Entry["Width"].Bind(id.Item().WidthString)
	f.Entry["Height"].Bind(id.Item().HeightString)
	f.Entry["Depth"].Bind(id.Item().DepthString)
	f.Entry["Volume"].Bind(id.Item().VolumeString)
	f.Entry["Weight"].Bind(id.Item().WeightString)

	for _, key := range ItemFormEntryKeys {
		f.Entry[key].Enable()
	}

	/* Select widgets */
	f.Select["Status"].Bind(id.Item().ItemStatus)
	f.Select["Category"].Bind(id.Item().Category)
	f.Select["Manufacturer"].Bind(id.Item().Manufacturer)
	f.Select["ModelName"].Bind(id.Item().ModelName)
	f.Select["LengthUnit"].Bind(id.Item().LengthUnit)
	f.Select["VolumeUnit"].Bind(id.Item().VolumeUnit)
	f.Select["WeightUnit"].Bind(id.Item().WeightUnit)

	id.Item().Manufacturer.AddListener(binding.NewDataListener(func() {
		models := func() []string {
			b.Metadata.GetModelIDs(id.Item().MfrID)
			var names []string
			ids := id.Item().MfrID.Children()
			for _, id := range ids {
				name, _ := id.Name()
				names = append(names, name)
			}
			return names
		}()
		f.Select["ModelName"].SetOptions(models)
	}))

	/* This step is needed because child categories have spaces prepended to them in the select list */
	cat, _ := id.Item().Category.Get()
	f.Select["Category"].SetSelectedIndex(b.Metadata.GetListItemIDForCategory(cat))

	id.Item().CatID.Category().Config["ShowPrice"].AddListener(binding.NewDataListener(func() {
		p, _ := id.Item().CatID.Category().Config["ShowPrice"].Get()
		if p {
			disabled = slices.DeleteFunc(disabled, func(s string) bool {
				switch s {
				case "Currency", "Price":
					return true
				default:
					return false
				}
			})
			f.show([]string{"Currency", "Price"})
		} else {
			disabled = Combine(disabled, []string{"Currency", "Price"})
			f.hide(disabled)
		}
	}))

	enabled = []string{"Status", "Category", "Manufacturer", "ModelName", "LengthUnit", "VolumeUnit", "WeightUnit"}

	f.enable(enabled)
}
func (f *Form) LoadMfr(id backend.MfrID) {
	f.Clear()

	f.Entry["Name"].Bind(id.Manufacturer().Name)
}
func (f *Form) LoadModel(id backend.ModelID) {
	f.Clear()
}
func (f *Form) LoadCategory(id backend.CatID) {
	f.Entry.Unbind()
	f.Entry.Clear()
	f.Entry.Bind(Sieve(id.Category().Bindings(), CategoryFormEntryKeys))
}

func (f *Form) itemContainer() *fyne.Container {
	return container.New(
		layout.NewFormLayout(),
		layout.NewSpacer(), container.NewHBox(f.Label["DateCreated"], f.Value["DateCreated"]),
		layout.NewSpacer(), container.NewHBox(f.Label["DateModified"], f.Value["DateModified"]),
		f.Label["ItemID"], container.NewHBox(f.Value["ItemID"], f.Select["Status"]),
		f.Label["Name"], f.Entry["Name"],
		f.Label["Category"], f.Select["Category"],
		f.Label["Manufacturer"], f.Select["Manufacturer"],
		f.Label["ModelName"], f.Select["ModelName"],
		f.Label["ModelDesc"], f.Entry["ModelDesc"],
		f.Label["ModelURL"], f.Entry["ModelURL"],
		f.Label["Dimensions"], f.dimbox(),
		layout.NewSpacer(), f.massbox(),
		f.Label["Price"], container.NewBorder(nil, nil, nil, f.Label["Currency"], f.Entry["Price"]),
		f.Label["ImgURL1"], f.Entry["ImgURL1"],
		f.Label["Notes"], f.Entry["Notes"],
		layout.NewSpacer(), widget.NewLabel(" "),
		layout.NewSpacer(), widget.NewRichTextFromMarkdown(`### `+lang.L("Preview")),
		f.Label["LongDesc"], f.Value["LongDesc"],
		f.Label["AddDesc"], f.Value["AddDesc"],
	)
}

func (f *Form) productContainer() *fyne.Container {
	return container.New(
		layout.NewFormLayout(),
		f.Label["Name"], f.Entry["Name"],
		f.Label["Manufacturer"], f.Select["Manufacturer"],
		f.Label["Category"], f.Select["Category"],
		f.Label["ModelDesc"], f.Entry["ModelDesc"],
		f.Label["ImgURL1"], f.Entry["ImgURL1"],
		f.Label["ImgURL2"], f.Entry["ImgURL2"],
		f.Label["ImgURL3"], f.Entry["ImgURL3"],
		f.Label["ImgURL4"], f.Entry["ImgURL4"],
		f.Label["ImgURL5"], f.Entry["ImgURL5"],
		f.Label["SpecsURL"], f.Entry["SpecsURL"],
		f.Label["ModelURL"], f.Entry["ModelURL"],
		f.Label["Dimensions"], f.dimbox(),
		layout.NewSpacer(), f.massbox(),
		// create X new objects from product -button
	)
}

func (f *Form) dimbox() *fyne.Container {
	return container.NewGridWithRows(1,
		container.NewBorder(nil, nil, f.Label["Width"], nil, f.Entry["Width"]),
		container.NewBorder(nil, nil, f.Label["Height"], nil, f.Entry["Height"]),
		container.NewBorder(nil, nil, f.Label["Depth"], nil, f.Entry["Depth"]),
		f.Select["LengthUnit"],
	)
}

func (f *Form) massbox() *fyne.Container {
	return container.NewGridWithRows(1,
		container.NewBorder(nil, nil, f.Label["Volume"], f.Select["VolumeUnit"], f.Entry["Volume"]),
		container.NewBorder(nil, nil, f.Label["Weight"], f.Select["WeightUnit"], f.Entry["Weight"]),
	)
}
func (f *Form) disable(disabled []string) {
	if len(disabled) < 1 {
		return
	}
	for key, val := range f.Check {
		if slices.Contains(disabled, key) {
			val.Disable()
		}
	}
	for key, val := range f.Entry {
		if slices.Contains(disabled, key) {
			val.Disable()
		}
	}
	for key, val := range f.Radio {
		if slices.Contains(disabled, key) {
			val.Disable()
		}
	}
	for key, val := range f.Select {
		if slices.Contains(disabled, key) {
			val.Disable()
		}
	}
}
func (f *Form) enable(enabled []string) {
	if len(enabled) < 1 {
		return
	}
	for key, val := range f.Check {
		if slices.Contains(enabled, key) {
			val.Enable()
		}
	}
	for key, val := range f.Entry {
		if slices.Contains(enabled, key) {
			val.Enable()
		}
	}
	for key, val := range f.Radio {
		if slices.Contains(enabled, key) {
			val.Enable()
		}
	}
	for key, val := range f.Select {
		if slices.Contains(enabled, key) {
			val.Enable()
		}
	}
}
func (f *Form) hide(disabled []string) {
	if len(disabled) < 1 {
		return
	}
	for key, val := range f.Check {
		if slices.Contains(disabled, key) {
			val.Hide()
		}
	}
	for key, val := range f.Entry {
		if slices.Contains(disabled, key) {
			val.Hide()
		}
	}
	for key, val := range f.Label {
		if slices.Contains(disabled, key) {
			val.Hide()
		}
	}
	for key, val := range f.Radio {
		if slices.Contains(disabled, key) {
			val.Hide()
		}
	}
	for key, val := range f.Select {
		if slices.Contains(disabled, key) {
			val.Hide()
		}
	}
	for key, val := range f.Value {
		if slices.Contains(disabled, key) {
			val.Hide()
		}
	}
}
func (f *Form) show(enabled []string) {
	if len(enabled) < 1 {
		return
	}
	for key, val := range f.Check {
		if slices.Contains(enabled, key) {
			val.Show()
		}
	}
	for key, val := range f.Entry {
		if slices.Contains(enabled, key) {
			val.Show()
		}
	}
	for key, val := range f.Label {
		if slices.Contains(enabled, key) {
			val.Show()
		}
	}
	for key, val := range f.Radio {
		if slices.Contains(enabled, key) {
			val.Show()
		}
	}
	for key, val := range f.Select {
		if slices.Contains(enabled, key) {
			val.Show()
		}
	}
	for key, val := range f.Value {
		if slices.Contains(enabled, key) {
			val.Show()
		}
	}
}
