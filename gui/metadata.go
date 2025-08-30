package gui

import (
	"UppSpar/backend"
	"UppSpar/backend/bridge"
	"log"
	"regexp"
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
	entry     bridge.Entries
	label     bridge.Labels
	selects   bridge.Selects
}

var entryKeys = []string{"Name", "Desc", "ImgURL1", "ImgURL2", "ImgURL3", "ImgURL4", "ImgURL5", "SpecsURL", "ModelURL", "Width", "Height", "Depth", "Volume", "Weight"}
var labelKeys = []string{"Name", "Category", "Manufacturer", "Desc", "Dimensions", "ImgURL1", "ImgURL2", "ImgURL3", "ImgURL4", "ImgURL5", "SpecsURL", "ModelURL", "Width", "Height", "Depth", "Volume", "Weight"}
var selectKeys = []string{"Manufacturer", "Category", "LengthUnit", "VolumeUnit", "WeightUnit"}

func newProductView(b *backend.Backend) *productView {
	p := &productView{}

	categories, _ := b.Metadata.Categories.Get()
	b.Metadata.Categories.AddListener(binding.NewDataListener(func() {
		categories, _ := b.Metadata.Categories.Get()
		p.selects["Category"].SetOptions(categories)
	}))

	manufacturers, _ := b.Metadata.MfrNameList.Get()
	b.Metadata.MfrNameList.AddListener(binding.NewDataListener(func() {
		manufacturers, _ := b.Metadata.MfrNameList.Get()
		manufacturers = append([]string{lang.L("None")}, manufacturers...)
		p.selects["Manufacturer"].SetOptions(manufacturers)
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
		r := regexp.MustCompile(`MDL-\d+$`)
		mdl := r.FindString(uid)
		if mdl != "" {
			mdl = strings.TrimPrefix(mdl, "MDL-")
			ModelID, err := strconv.Atoi(mdl)
			if err != nil {
				log.Printf("strconv.Atoi(%s) error: %s", mdl, err)
				panic(err)
			}
			p.LoadModel(b, backend.ModelID(ModelID))
		} else {
			r := regexp.MustCompile(`MFR-\d+([^FR]*)$`)
			mfr := r.FindString(uid)
			log.Printf("uid = %s, mfr = %s", uid, mfr)
			mfr = strings.TrimPrefix(mfr, "MFR-")
			MfrID, err := strconv.Atoi(mfr)
			if err != nil {
				log.Printf("strconv.Atoi(%s) error: %s", mfr, err)
				panic(err)
			}
			p.LoadMfr(backend.MfrID(MfrID))
		}
	}
	// t.OnUnselected = func(uid widget.TreeNodeID) {}

	p.entry = make(bridge.Entries)
	p.label = make(bridge.Labels)
	p.selects = make(bridge.Selects)

	for _, key := range entryKeys {
		p.entry[key] = midget.NewEntry()
	}

	labelStrings := []string{
		lang.L("Name"),
		lang.X("item.form.label.category", "item.form.label.category"),
		lang.L("Manufacturer"),
		lang.X("metadata.product.form.description", "metadata.product.form.description"),
		lang.L("Dimensions"),
		lang.L("Image URL"),
		lang.L("Image URL"),
		lang.L("Image URL"),
		lang.L("Image URL"),
		lang.L("Image URL"),
		lang.L("Specs URL"),
		lang.L("Model URL"),
		lang.L("Width"),
		lang.L("Height"),
		lang.L("Depth"),
		lang.L("Volume"),
		lang.L("Weight"),
	}
	for i, key := range labelKeys {
		p.label[key] = ttw.NewLabel(labelStrings[i])
	}

	selectOptions := [][]string{manufacturers, categories, lengthUnits, volumeUnits, weightUnits}
	for i, key := range selectKeys {
		p.selects[key] = ttw.NewSelect(selectOptions[i], func(s string) {})
	}

	p.entry["Desc"].MultiLine = true
	p.entry["Desc"].SetMinRowsVisible(5)
	p.entry["Desc"].Wrapping = fyne.TextWrapWord

	dimBox := container.NewGridWithRows(1,
		container.NewBorder(nil, nil, p.label["Width"], nil, p.entry["Width"]),
		container.NewBorder(nil, nil, p.label["Height"], nil, p.entry["Height"]),
		container.NewBorder(nil, nil, p.label["Depth"], nil, p.entry["Depth"]),
		p.selects["LengthUnit"],
	)

	massBox := container.NewGridWithRows(1,
		container.NewBorder(nil, nil, p.label["Volume"], p.selects["VolumeUnit"], p.entry["Volume"]),
		container.NewBorder(nil, nil, p.label["Weight"], p.selects["WeightUnit"], p.entry["Weight"]),
	)

	f := container.New(layout.NewFormLayout(),
		p.label["Name"], p.entry["Name"],
		p.label["Manufacturer"], p.selects["Manufacturer"],
		p.label["Category"], p.selects["Category"],
		p.label["Desc"], p.entry["Desc"],
		p.label["ImgURL1"], p.entry["ImgURL1"],
		p.label["SpecsURL"], p.entry["SpecsURL"],
		p.label["ModelURL"], p.entry["ModelURL"],
		p.label["Dimensions"], dimBox,
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
	p.container = container.NewHSplit(t, f)
	p.container.SetOffset(0.25)
	p.Clear()
	p.Disable()
	p.Hide()
	return p
}

func (pv *productView) Clear() {
	pv.entry.Clear()
	pv.selects.Clear()
	pv.label["Name"].SetText(lang.L("Name"))
}
func (pv *productView) Disable() {
	pv.entry.Disable()
	pv.selects.Disable()
}
func (pv *productView) Enable() {
	pv.entry.Enable()
	pv.selects.Enable()
}
func (pv *productView) Hide() {
	pv.entry.Hide()
	pv.label.Hide()
	pv.selects.Hide()
}
func (pv *productView) LoadMfr(id backend.MfrID) {
	pv.Unbind()
	pv.Clear()
	pv.Hide()

	pv.entry["Name"].Bind(id.Manufacturer().Name)
	pv.entry["Name"].Enable()
	pv.entry["Name"].Show()

	pv.label["Name"].SetText(lang.L("Manufacturer"))
	pv.label["Name"].Show()
}
func (pv *productView) LoadModel(b *backend.Backend, id backend.ModelID) {
	pv.Hide()
	pv.Unbind()
	pv.Clear()

	pv.entry["Name"].Bind(id.Model().Name)
	pv.entry["Desc"].Bind(id.Model().Desc)
	pv.entry["ImgURL1"].Bind(id.Model().ImgURL1)
	pv.entry["ImgURL2"].Bind(id.Model().ImgURL2)
	pv.entry["ImgURL3"].Bind(id.Model().ImgURL3)
	pv.entry["ImgURL4"].Bind(id.Model().ImgURL4)
	pv.entry["ImgURL5"].Bind(id.Model().ImgURL5)
	pv.entry["SpecsURL"].Bind(id.Model().SpecsURL)
	pv.entry["ModelURL"].Bind(id.Model().ModelURL)
	pv.entry["Width"].Bind(id.Model().Width)
	pv.entry["Height"].Bind(id.Model().Height)
	pv.entry["Depth"].Bind(id.Model().Depth)
	pv.entry["Volume"].Bind(id.Model().Volume)
	pv.entry["Weight"].Bind(id.Model().Weight)

	pv.label["Name"].SetText(lang.L("Product"))

	pv.selects["Category"].Bind(id.Model().Category)
	pv.selects["Manufacturer"].Bind(id.Model().Manufacturer)
	pv.selects["LengthUnit"].Bind(id.Model().LengthUnit)
	pv.selects["VolumeUnit"].Bind(id.Model().VolumeUnit)
	pv.selects["WeightUnit"].Bind(id.Model().WeightUnit)

	if mfrid, _ := id.MfrID(); mfrid != 0 {
		n, _ := mfrid.Name()
		pv.selects["Manufacturer"].SetSelected(n)
	}

	cat, _ := id.Model().Category.Get()
	pv.selects["Category"].SetSelectedIndex(b.Metadata.GetListItemIDFor(cat))

	pv.entry.Enable()
	pv.selects.Enable()

	pv.entry.Show()
	pv.label.Show()
	pv.selects.Show()
}
func (pv *productView) Unbind() {
	pv.entry.Unbind()
	pv.selects.Unbind()
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
