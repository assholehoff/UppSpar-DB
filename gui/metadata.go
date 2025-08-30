package gui

import (
	"UppSpar/backend"
	"log"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	midget "github.com/assholehoff/fyne-midget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type metadataView struct {
	category *categoryView
	product  *productView
	tabs     *container.AppTabs
}

func newMetadataView(b *backend.Backend) *metadataView {
	categoryView := newCategoryView(b)
	productView := newProductView(b)

	return &metadataView{
		category: categoryView,
		product:  productView,
		tabs:     newMetadataTabs(categoryView, productView),
	}
}

func newMetadataTabs(c *categoryView, mdl *productView) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem(lang.L("Products"), mdl.container),
		container.NewTabItem(lang.L("Categories"), c.container),
	)
	return tabs
}

type productView struct {
	container *container.Split
	form      *productForm
	tree      *widget.Tree
}

type productForm struct {
	entries *productEntries
	labels  *productLabels
	selects *productSelects
}

type productEntries struct {
	Name     *midget.Entry
	Desc     *midget.Entry
	ImgURL1  *midget.Entry
	ImgURL2  *midget.Entry
	ImgURL3  *midget.Entry
	ImgURL4  *midget.Entry
	ImgURL5  *midget.Entry
	SpecsURL *midget.Entry
	ModelURL *midget.Entry
	Width    *midget.Entry
	Height   *midget.Entry
	Depth    *midget.Entry
	Volume   *midget.Entry
	Weight   *midget.Entry
}
type productLabels struct {
	Name         *widget.Label
	Category     *widget.Label
	Manufacturer *widget.Label
	Desc         *widget.Label
	Dimensions   *widget.Label
	ImgURL1      *widget.Label
	ImgURL2      *widget.Label
	ImgURL3      *widget.Label
	ImgURL4      *widget.Label
	ImgURL5      *widget.Label
	SpecsURL     *widget.Label
	ModelURL     *widget.Label
	Width        *widget.Label
	Height       *widget.Label
	Depth        *widget.Label
	Volume       *widget.Label
	Weight       *widget.Label
}
type productSelects struct {
	Manufacturer *widget.Select
	Category     *widget.Select
	LengthUnit   *widget.Select
	VolumeUnit   *widget.Select
	WeightUnit   *widget.Select
}

func newProductView(b *backend.Backend) *productView {
	pv := &productView{}

	categories, _ := b.Metadata.Categories.Get()
	b.Metadata.Categories.AddListener(binding.NewDataListener(func() {
		categories, _ := b.Metadata.Categories.Get()
		pv.form.selects.Category.SetOptions(categories)
	}))
	manufacturers, _ := b.Metadata.MfrNameList.Get()
	b.Metadata.MfrNameList.AddListener(binding.NewDataListener(func() {
		manufacturers, _ := b.Metadata.MfrNameList.Get()
		manufacturers = append([]string{lang.L("None")}, manufacturers...)
		pv.form.selects.Manufacturer.SetOptions(manufacturers)
	}))
	lengthUnits := []string{"mm", "cm", "dm", "m"}
	volumeUnits := []string{"ml", "cl", "dl", "l"}
	weightUnits := []string{"g", "hg", "kg"}

	createItem := func(branch bool) fyne.CanvasObject {
		if branch {
			return &widget.Label{
				Text:      "Template branch category name",
				TextStyle: fyne.TextStyle{Bold: true},
			}
		}
		return widget.NewLabel("Template leaf category name")
	}
	updateItem := func(di binding.DataItem, branch bool, co fyne.CanvasObject) {
		co.(*widget.Label).Bind(di.(binding.String))
	}
	tree := widget.NewTreeWithData(b.Metadata.ProductTree, createItem, updateItem)
	tree.OnSelected = func(uid widget.TreeNodeID) {
		// productIdString := b.Metadata.GetProductIDFor(uid)
		if strings.HasPrefix(uid, "MFR-") {
			mfrIdString := strings.TrimPrefix(uid, "MFR-")
			MfrID, err := strconv.Atoi(mfrIdString)
			if err != nil {
				log.Printf("strconv.Atoi(%s) error: %s", mfrIdString, err)
				panic(err)
			}
			pv.LoadMfr(backend.MfrID(MfrID))
		} else {
			mdlIdString := strings.TrimPrefix(uid, "MDL-")
			ModelID, err := strconv.Atoi(mdlIdString)
			if err != nil {
				log.Printf("strconv.Atoi(%s) error: %s", mdlIdString, err)
				panic(err)
			}
			pv.LoadModel(backend.ModelID(ModelID))
		}
	}
	// t.OnUnselected = func(uid widget.TreeNodeID) {}

	pv.form = &productForm{
		entries: &productEntries{
			Name:     midget.NewEntry(),
			Desc:     midget.NewEntry(),
			ImgURL1:  midget.NewEntry(),
			ImgURL2:  midget.NewEntry(),
			ImgURL3:  midget.NewEntry(),
			ImgURL4:  midget.NewEntry(),
			ImgURL5:  midget.NewEntry(),
			SpecsURL: midget.NewEntry(),
			ModelURL: midget.NewEntry(),
			Width:    midget.NewEntry(),
			Height:   midget.NewEntry(),
			Depth:    midget.NewEntry(),
			Volume:   midget.NewEntry(),
			Weight:   midget.NewEntry(),
		},
		labels: &productLabels{
			Name:         widget.NewLabel(lang.L("Name")),
			Category:     widget.NewLabel(lang.X("item.form.label.category", "item.form.label.category")),
			Manufacturer: widget.NewLabel(lang.L("Manufacturer")),
			Desc:         widget.NewLabel(lang.X("metadata.product.form.description", "metadata.product.form.description")),
			Dimensions:   widget.NewLabel(lang.L("Dimensions")),
			ImgURL1:      widget.NewLabel(lang.L("Image URL")),
			ImgURL2:      widget.NewLabel(lang.L("Image URL")),
			ImgURL3:      widget.NewLabel(lang.L("Image URL")),
			ImgURL4:      widget.NewLabel(lang.L("Image URL")),
			ImgURL5:      widget.NewLabel(lang.L("Image URL")),
			SpecsURL:     widget.NewLabel(lang.L("Specs URL")),
			ModelURL:     widget.NewLabel(lang.L("Model URL")),
			Width:        widget.NewLabel(lang.L("Width")),
			Height:       widget.NewLabel(lang.L("Height")),
			Depth:        widget.NewLabel(lang.L("Depth")),
			Volume:       widget.NewLabel(lang.L("Volume")),
			Weight:       widget.NewLabel(lang.L("Weight")),
		},
		selects: &productSelects{
			Manufacturer: widget.NewSelect(manufacturers, func(s string) {}),
			Category:     widget.NewSelect(categories, func(s string) {}),
			LengthUnit:   widget.NewSelect(lengthUnits, func(s string) {}),
			VolumeUnit:   widget.NewSelect(volumeUnits, func(s string) {}),
			WeightUnit:   widget.NewSelect(weightUnits, func(s string) {}),
		},
	}
	pv.form.entries.Desc.MultiLine = true
	pv.form.entries.Desc.SetMinRowsVisible(5)
	pv.form.entries.Desc.Wrapping = fyne.TextWrapWord

	dimBox := container.NewGridWithRows(1,
		container.NewBorder(nil, nil, pv.form.labels.Width, nil, pv.form.entries.Width),
		container.NewBorder(nil, nil, pv.form.labels.Height, nil, pv.form.entries.Height),
		container.NewBorder(nil, nil, pv.form.labels.Depth, nil, pv.form.entries.Depth),
		pv.form.selects.LengthUnit,
	)

	massBox := container.NewGridWithRows(1,
		container.NewBorder(nil, nil, pv.form.labels.Volume, pv.form.selects.VolumeUnit, pv.form.entries.Volume),
		container.NewBorder(nil, nil, pv.form.labels.Weight, pv.form.selects.WeightUnit, pv.form.entries.Weight),
	)

	f := container.New(layout.NewFormLayout(),
		pv.form.labels.Name, pv.form.entries.Name,
		pv.form.labels.Manufacturer, pv.form.selects.Manufacturer,
		pv.form.labels.Category, pv.form.selects.Category,
		pv.form.labels.Desc, pv.form.entries.Desc,
		pv.form.labels.ImgURL1, pv.form.entries.ImgURL1,
		pv.form.labels.SpecsURL, pv.form.entries.SpecsURL,
		pv.form.labels.ModelURL, pv.form.entries.ModelURL,
		pv.form.labels.Dimensions, dimBox,
		layout.NewSpacer(), massBox,
	)

	t := container.NewBorder(
		container.NewHBox(
			widget.NewButton(lang.L("New Manufacturer"), func() {
				_, err := b.Metadata.CreateNewManufacturer()
				if err != nil {
					return
				}
				// pv.LoadMfr(id)
			}),
			widget.NewButton(lang.L("New Product"), func() {
				_, err := b.Metadata.CreateNewProduct()
				if err != nil {
					return
				}
				// pv.LoadModel(id)
			}),
		),
		nil, nil, nil, tree,
	)
	pv.container = container.NewHSplit(t, f)
	pv.container.SetOffset(0.25)
	pv.Clear()
	pv.Disable()
	pv.Hide()
	return pv
}

func (pv *productView) Clear() {
	pv.form.entries.Name.SetText("")
	pv.form.entries.Desc.SetText("")
	pv.form.entries.ImgURL1.SetText("")
	pv.form.entries.ImgURL2.SetText("")
	pv.form.entries.ImgURL3.SetText("")
	pv.form.entries.ImgURL4.SetText("")
	pv.form.entries.ImgURL5.SetText("")
	pv.form.entries.SpecsURL.SetText("")
	pv.form.entries.ModelURL.SetText("")
	pv.form.entries.Width.SetText("")
	pv.form.entries.Height.SetText("")
	pv.form.entries.Depth.SetText("")
	pv.form.entries.Volume.SetText("")
	pv.form.entries.Weight.SetText("")

	pv.form.labels.Name.SetText(lang.L("Name"))

	pv.form.selects.Manufacturer.ClearSelected()
	pv.form.selects.Category.ClearSelected()
	pv.form.selects.LengthUnit.ClearSelected()
	pv.form.selects.VolumeUnit.ClearSelected()
	pv.form.selects.WeightUnit.ClearSelected()

}
func (pv *productView) Disable() {
	pv.form.entries.Name.Disable()
	pv.form.entries.Desc.Disable()
	pv.form.entries.ImgURL1.Disable()
	pv.form.entries.ImgURL2.Disable()
	pv.form.entries.ImgURL3.Disable()
	pv.form.entries.ImgURL4.Disable()
	pv.form.entries.ImgURL5.Disable()
	pv.form.entries.SpecsURL.Disable()
	pv.form.entries.ModelURL.Disable()
	pv.form.entries.Width.Disable()
	pv.form.entries.Height.Disable()
	pv.form.entries.Depth.Disable()
	pv.form.entries.Volume.Disable()
	pv.form.entries.Weight.Disable()

	pv.form.selects.Manufacturer.Disable()
	pv.form.selects.Category.Disable()
	pv.form.selects.LengthUnit.Disable()
	pv.form.selects.VolumeUnit.Disable()
	pv.form.selects.WeightUnit.Disable()
}
func (pv *productView) Enable() {
	pv.form.entries.Name.Enable()
	pv.form.entries.Desc.Enable()
	pv.form.entries.ImgURL1.Enable()
	pv.form.entries.ImgURL2.Enable()
	pv.form.entries.ImgURL3.Enable()
	pv.form.entries.ImgURL4.Enable()
	pv.form.entries.ImgURL5.Enable()
	pv.form.entries.SpecsURL.Enable()
	pv.form.entries.ModelURL.Enable()
	pv.form.entries.Width.Enable()
	pv.form.entries.Height.Enable()
	pv.form.entries.Depth.Enable()
	pv.form.entries.Volume.Enable()
	pv.form.entries.Weight.Enable()

	pv.form.selects.Manufacturer.Enable()
	pv.form.selects.Category.Enable()
	pv.form.selects.LengthUnit.Enable()
	pv.form.selects.VolumeUnit.Enable()
	pv.form.selects.WeightUnit.Enable()
}
func (pv *productView) Hide() {
	pv.form.entries.Name.Hide()
	pv.form.entries.Desc.Hide()
	pv.form.entries.ImgURL1.Hide()
	pv.form.entries.ImgURL2.Hide()
	pv.form.entries.ImgURL3.Hide()
	pv.form.entries.ImgURL4.Hide()
	pv.form.entries.ImgURL5.Hide()
	pv.form.entries.SpecsURL.Hide()
	pv.form.entries.ModelURL.Hide()
	pv.form.entries.Width.Hide()
	pv.form.entries.Height.Hide()
	pv.form.entries.Depth.Hide()
	pv.form.entries.Volume.Hide()
	pv.form.entries.Weight.Hide()

	pv.form.labels.Name.Hide()
	pv.form.labels.Category.Hide()
	pv.form.labels.Manufacturer.Hide()
	pv.form.labels.Desc.Hide()
	pv.form.labels.Dimensions.Hide()
	pv.form.labels.ImgURL1.Hide()
	pv.form.labels.ImgURL2.Hide()
	pv.form.labels.ImgURL3.Hide()
	pv.form.labels.ImgURL4.Hide()
	pv.form.labels.ImgURL5.Hide()
	pv.form.labels.SpecsURL.Hide()
	pv.form.labels.ModelURL.Hide()
	pv.form.labels.Width.Hide()
	pv.form.labels.Height.Hide()
	pv.form.labels.Depth.Hide()
	pv.form.labels.Volume.Hide()
	pv.form.labels.Weight.Hide()

	pv.form.selects.Manufacturer.Hide()
	pv.form.selects.Category.Hide()
	pv.form.selects.LengthUnit.Hide()
	pv.form.selects.VolumeUnit.Hide()
	pv.form.selects.WeightUnit.Hide()
}
func (pv *productView) LoadMfr(id backend.MfrID) {
	pv.Unbind()
	pv.Clear()
	pv.Hide()

	pv.form.entries.Name.Bind(id.Manufacturer().Name)
	pv.form.entries.Name.Enable()
	pv.form.entries.Name.Show()

	pv.form.labels.Name.SetText(lang.L("Manufacturer"))
	pv.form.labels.Name.Show()
}
func (pv *productView) LoadModel(id backend.ModelID) {
	pv.Hide()
	pv.Unbind()
	pv.Clear()

	pv.form.entries.Name.Bind(id.Model().Name)
	pv.form.entries.Desc.Bind(id.Model().Desc)
	pv.form.entries.ImgURL1.Bind(id.Model().ImgURL1)
	pv.form.entries.ImgURL2.Bind(id.Model().ImgURL2)
	pv.form.entries.ImgURL3.Bind(id.Model().ImgURL3)
	pv.form.entries.ImgURL4.Bind(id.Model().ImgURL4)
	pv.form.entries.ImgURL5.Bind(id.Model().ImgURL5)
	pv.form.entries.SpecsURL.Bind(id.Model().SpecsURL)
	pv.form.entries.ModelURL.Bind(id.Model().ModelURL)
	pv.form.entries.Width.Bind(id.Model().Width)
	pv.form.entries.Height.Bind(id.Model().Height)
	pv.form.entries.Depth.Bind(id.Model().Depth)
	pv.form.entries.Volume.Bind(id.Model().Volume)
	pv.form.entries.Weight.Bind(id.Model().Weight)

	pv.form.labels.Name.SetText(lang.L("Product"))

	pv.form.selects.Category.Bind(id.Model().Category)
	pv.form.selects.Manufacturer.Bind(id.Model().Manufacturer)
	pv.form.selects.LengthUnit.Bind(id.Model().LengthUnit)
	pv.form.selects.VolumeUnit.Bind(id.Model().VolumeUnit)
	pv.form.selects.WeightUnit.Bind(id.Model().WeightUnit)

	if mfrid, _ := id.MfrID(); mfrid != 0 {
		n, _ := mfrid.Name()
		pv.form.selects.Manufacturer.SetSelected(n)
	}

	pv.form.entries.Name.Enable()
	pv.form.entries.Desc.Enable()
	pv.form.entries.ImgURL1.Enable()
	pv.form.entries.ImgURL2.Enable()
	pv.form.entries.ImgURL3.Enable()
	pv.form.entries.ImgURL4.Enable()
	pv.form.entries.ImgURL5.Enable()
	pv.form.entries.SpecsURL.Enable()
	pv.form.entries.ModelURL.Enable()
	pv.form.entries.Width.Enable()
	pv.form.entries.Height.Enable()
	pv.form.entries.Depth.Enable()
	pv.form.entries.Volume.Enable()
	pv.form.entries.Weight.Enable()
	pv.form.selects.Manufacturer.Enable()
	pv.form.selects.Category.Enable()
	pv.form.selects.LengthUnit.Enable()
	pv.form.selects.VolumeUnit.Enable()
	pv.form.selects.WeightUnit.Enable()

	pv.form.entries.Name.Show()
	pv.form.entries.Desc.Show()
	pv.form.entries.ImgURL1.Show()
	pv.form.entries.ImgURL2.Show()
	pv.form.entries.ImgURL3.Show()
	pv.form.entries.ImgURL4.Show()
	pv.form.entries.ImgURL5.Show()
	pv.form.entries.SpecsURL.Show()
	pv.form.entries.ModelURL.Show()
	pv.form.entries.Width.Show()
	pv.form.entries.Height.Show()
	pv.form.entries.Depth.Show()
	pv.form.entries.Volume.Show()
	pv.form.entries.Weight.Show()

	pv.form.labels.Name.Show()
	pv.form.labels.Category.Show()
	pv.form.labels.Manufacturer.Show()
	pv.form.labels.Desc.Show()
	pv.form.labels.Dimensions.Show()
	pv.form.labels.ImgURL1.Show()
	pv.form.labels.ImgURL2.Show()
	pv.form.labels.ImgURL3.Show()
	pv.form.labels.ImgURL4.Show()
	pv.form.labels.ImgURL5.Show()
	pv.form.labels.SpecsURL.Show()
	pv.form.labels.ModelURL.Show()
	pv.form.labels.Width.Show()
	pv.form.labels.Height.Show()
	pv.form.labels.Depth.Show()
	pv.form.labels.Volume.Show()
	pv.form.labels.Weight.Show()

	pv.form.selects.Manufacturer.Show()
	pv.form.selects.Category.Show()
	pv.form.selects.LengthUnit.Show()
	pv.form.selects.VolumeUnit.Show()
	pv.form.selects.WeightUnit.Show()
}
func (pv *productView) Unbind() {
	pv.form.entries.Name.Unbind()
	pv.form.entries.Desc.Unbind()
	pv.form.entries.ImgURL1.Unbind()
	pv.form.entries.ImgURL2.Unbind()
	pv.form.entries.ImgURL3.Unbind()
	pv.form.entries.ImgURL4.Unbind()
	pv.form.entries.ImgURL5.Unbind()
	pv.form.entries.SpecsURL.Unbind()
	pv.form.entries.ModelURL.Unbind()
	pv.form.entries.Width.Unbind()
	pv.form.entries.Height.Unbind()
	pv.form.entries.Depth.Unbind()
	pv.form.entries.Volume.Unbind()
	pv.form.entries.Weight.Unbind()

	pv.form.selects.Manufacturer.Unbind()
	pv.form.selects.Category.Unbind()
	pv.form.selects.LengthUnit.Unbind()
	pv.form.selects.VolumeUnit.Unbind()
	pv.form.selects.WeightUnit.Unbind()
}

type categoryChecks struct {
	price  *ttw.Check
	length *ttw.Check
	volume *ttw.Check
	weight *ttw.Check
}

type categoryEntries struct {
	name *midget.Entry
}

type categoryLabels struct {
	name   *midget.Label
	parent *midget.Label
}

type categorySelects struct {
	parent *widget.Select
}

type categoryFormView struct {
	container *fyne.Container
	checks    *categoryChecks
	entries   *categoryEntries
	labels    *categoryLabels
	selects   *categorySelects
}

type categoryView struct {
	container *container.Split
	form      *categoryFormView
	tree      *widget.Tree
	toolbar   *widget.Toolbar
}

func newCategoryView(b *backend.Backend) *categoryView {
	createTreeItem := func(branch bool) fyne.CanvasObject {
		if branch {
		}
		return widget.NewLabel("Template category name")
	}
	updateTreeItem := func(di binding.DataItem, branch bool, co fyne.CanvasObject) {
		v, err := di.(binding.Untyped).Get()
		if err != nil {
			log.Println(err)
			return
		}
		CatID := v.(backend.CatID)
		co.(*widget.Label).Bind(CatID.Category().Name)
	}
	m := &categoryView{
		tree: widget.NewTreeWithData(b.Metadata.CatIDTree, createTreeItem, updateTreeItem),
	}
	m.form = &categoryFormView{
		checks: &categoryChecks{
			price:  ttw.NewCheck(lang.X("metadata.category.check.price", "metadata.category.check.price"), func(b bool) {}),
			length: ttw.NewCheck(lang.X("metadata.category.check.length", "metadata.category.check.length"), func(b bool) {}),
			volume: ttw.NewCheck(lang.X("metadata.category.check.volume", "metadata.category.check.volume"), func(b bool) {}),
			weight: ttw.NewCheck(lang.X("metadata.category.check.weight", "metadata.category.check.weight"), func(b bool) {}),
		},
		entries: &categoryEntries{
			name: midget.NewEntry(),
		},
		labels: &categoryLabels{
			name:   midget.NewLabel(lang.X("metadata.form.name", "metadata.form.name"), "", ""),
			parent: midget.NewLabel(lang.X("metadata.form.parent", "metadata.form.parent"), "", ""),
		},
		selects: &categorySelects{
			parent: widget.NewSelect([]string{lang.L("None")}, func(s string) {}),
		},
	}
	m.form.container = container.New(layout.NewFormLayout(),
		m.form.labels.parent, m.form.selects.parent,
		m.form.labels.name, m.form.entries.name,
		layout.NewSpacer(), m.form.checks.price,
		layout.NewSpacer(), container.NewHBox(
			m.form.checks.length,
			m.form.checks.volume,
			m.form.checks.weight,
		),
	)
	m.tree.OnSelected = func(uid widget.TreeNodeID) {
		b.Metadata.SelectCategory(b.Metadata.GetCatIDForTreeItem(uid))
		m.form.Bind(b, b.Metadata.GetCatIDForTreeItem(uid))
	}
	m.tree.OnUnselected = func(uid widget.TreeNodeID) {
		b.Metadata.UnselectCategory(b.Metadata.GetCatIDForTreeItem(uid))
		m.form.entries.name.Unbind()
	}
	m.toolbar = widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			m.toolbar.Items[0].(*widget.ToolbarAction).Disable()
			go func() {
				b.Metadata.CreateNewCategory()
				fyne.Do(func() {
					time.Sleep(300 * time.Millisecond) // Prevent accidental multiclick
					m.toolbar.Items[0].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
			m.toolbar.Items[1].(*widget.ToolbarAction).Disable()
			go func() {
				b.Metadata.DeleteCategory()
				fyne.Do(func() {
					time.Sleep(300 * time.Millisecond)
					m.toolbar.Items[1].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			m.toolbar.Items[2].(*widget.ToolbarAction).Disable()
			go func() {
				b.Metadata.CopyCategory()
				fyne.Do(func() {
					time.Sleep(300 * time.Millisecond)
					m.toolbar.Items[2].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
	)
	listView := container.NewBorder(m.toolbar, nil, nil, nil, m.tree)
	m.container = container.NewHSplit(listView, m.form.container)
	m.container.SetOffset(0.25)
	return m
}

func (m *categoryFormView) Bind(b *backend.Backend, id backend.CatID) {
	remove := func(s []string, r string) []string {
		for i, v := range s {
			if strings.TrimSpace(v) == r {
				return append(s[:i], s[i+1:]...)
			}
		}
		return s
	}
	var categories []string
	fetchCategories := func() []string {
		categories, _ := b.Metadata.Categories.Get()
		return append([]string{lang.L("None")}, categories...)
	}
	categories = fetchCategories()
	log.Println(categories)
	self, _ := id.Name()
	remove(categories, self)

	m.selects.parent.SetOptions(categories)
	m.entries.name.Bind(id.Category().Name)
	m.checks.price.Bind(id.Category().Config["ShowPrice"])
	m.checks.length.Bind(id.Category().Config["ShowLength"])
	m.checks.volume.Bind(id.Category().Config["ShowVolume"])
	m.checks.weight.Bind(id.Category().Config["ShowWeight"])

	m.selects.parent.Bind(id.Category().Parent)

	p, _ := id.ParentID()
	if p == 0 {
		m.selects.parent.SetSelected(lang.L("None"))
	}

	n, _ := p.Name()
	m.selects.parent.SetSelected(n)

	id.Category().Name.AddListener(binding.NewDataListener(func() { b.Metadata.UpdateCatList() }))
	id.Category().Parent.AddListener(binding.NewDataListener(func() { b.Metadata.UpdateCatList() }))
}
