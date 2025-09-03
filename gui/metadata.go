package gui

import (
	"UppSpar/backend"
	"UppSpar/backend/bridge"
	"fmt"
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
			return container.NewHBox(&widget.Label{
				Text:      "Template branch Manufacturer name",
				TextStyle: fyne.TextStyle{Bold: true},
			}, widget.NewLabel("(00000)"))
		}
		return widget.NewLabel("Template leaf Product name")
	}
	updateItem := func(di binding.DataItem, branch bool, co fyne.CanvasObject) {
		v, err := di.(binding.Untyped).Get()
		if err != nil {
			panic(err)
		}
		if branch {
			MfrID := v.(backend.MfrID)
			co.(*fyne.Container).Objects[0].(*widget.Label).Bind(MfrID.Manufacturer().Name)
			co.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("(%d)", len(MfrID.Children())))
			if len(MfrID.Children()) > 0 {
				co.(*fyne.Container).Objects[1].(*widget.Label).Show()
			} else {
				co.(*fyne.Container).Objects[1].(*widget.Label).Hide()
			}
		} else {
			if v.(backend.NumID).TypeName() == "MfrID" {
				id := v.(backend.MfrID)
				co.(*widget.Label).Bind(id.Manufacturer().Name)
			} else {
				id := v.(backend.ModelID)
				co.(*widget.Label).Bind(id.Model().Name)
				co.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			}
		}
	}

	tree := widget.NewTreeWithData(b.Metadata.ProductTree, createItem, updateItem)
	tree.OnSelected = func(uid widget.TreeNodeID) {
		tree.OpenBranch(uid)
		r := regexp.MustCompile(`MDL-\d+$`)
		mdl := r.FindString(uid)
		if mdl != "" {
			mdl = strings.TrimPrefix(mdl, "MDL-")
			ModelID, err := strconv.Atoi(mdl)
			if err != nil {
				log.Printf("strconv.Atoi(%s) error: %s", mdl, err)
				panic(err)
			}
			b.Metadata.SelectProduct(backend.ModelID(ModelID))
			p.LoadModel(b, backend.ModelID(ModelID))
		} else {
			r := regexp.MustCompile(`MFR-\d+([^FR]*)$`)
			mfr := r.FindString(uid)
			mfr = strings.TrimPrefix(mfr, "MFR-")
			MfrID, err := strconv.Atoi(mfr)
			if err != nil {
				log.Printf("strconv.Atoi(%s) error: %s", mfr, err)
				panic(err)
			}
			b.Metadata.SelectProduct(backend.MfrID(MfrID))
			p.LoadMfr(backend.MfrID(MfrID))
		}
	}
	tree.OnUnselected = func(uid widget.TreeNodeID) { b.Metadata.ClearProdSelection() }

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
		lang.L("Image URL") + " 1",
		lang.L("Image URL") + " 2",
		lang.L("Image URL") + " 3",
		lang.L("Image URL") + " 4",
		lang.L("Image URL") + " 5",
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
		p.label["ImgURL2"], p.entry["ImgURL2"],
		p.label["ImgURL3"], p.entry["ImgURL3"],
		p.label["ImgURL4"], p.entry["ImgURL4"],
		p.label["ImgURL5"], p.entry["ImgURL5"],
		p.label["SpecsURL"], p.entry["SpecsURL"],
		p.label["ModelURL"], p.entry["ModelURL"],
		p.label["Dimensions"], dimBox,
		layout.NewSpacer(), massBox,
	)

	t := container.NewBorder(
		container.NewHBox(
			widget.NewButton(lang.L("New Manufacturer"), func() {
				id, err := b.Metadata.CreateNewManufacturer()
				if err != nil {
					return
				}
				p.LoadMfr(id)
			}),
			widget.NewButton(lang.L("New Product"), func() {
				id, err := b.Metadata.CreateNewProduct()
				if err != nil {
					return
				}
				p.LoadModel(b, id)
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
	pv.selects["Category"].SetSelectedIndex(b.Metadata.GetListItemIDForCategory(cat))

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

type categoryView struct {
	container *container.Split
	check     bridge.Checks
	entry     bridge.Entries
	label     bridge.Labels
	selects   bridge.Selects
	tree      *widget.Tree
	toolbar   *widget.Toolbar
}

func newCategoryView(b *backend.Backend) *categoryView {
	var cv *categoryView

	createTreeItem := func(branch bool) fyne.CanvasObject {
		if branch {
			return container.NewHBox(
				widget.NewLabel("Template branch category name"),
				widget.NewLabel("(0000)"),
			)
		}
		return widget.NewLabel("Template leaf category name")
	}
	updateTreeItem := func(di binding.DataItem, branch bool, co fyne.CanvasObject) {
		v, err := di.(binding.Untyped).Get()
		if err != nil {
			log.Println(err)
			return
		}
		CatID := v.(backend.CatID)
		if branch {
			co.(*fyne.Container).Objects[0].(*widget.Label).Bind(CatID.Category().Name)
			co.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("(%d)", len(CatID.Children())))
		} else {
			co.(*widget.Label).Bind(CatID.Category().Name)
		}
	}
	cv = &categoryView{
		tree: widget.NewTreeWithData(b.Metadata.CatIDTree, createTreeItem, updateTreeItem),
	}
	cv.tree.OnSelected = func(uid widget.TreeNodeID) {
		cv.tree.OpenBranch(uid)
		b.Metadata.SelectCategory(b.Metadata.GetCatIDForTreeItem(uid))
		cv.Load(b, b.Metadata.GetCatIDForTreeItem(uid))
	}
	cv.tree.OnUnselected = func(uid widget.TreeNodeID) {
		b.Metadata.UnselectCategory(b.Metadata.GetCatIDForTreeItem(uid))
		cv.Unload()
	}

	cv.check = make(bridge.Checks)
	cv.entry = make(bridge.Entries)
	cv.label = make(bridge.Labels)
	cv.selects = make(bridge.Selects)

	cv.check["Price"] = ttw.NewCheck(lang.X("metadata.category.check.price", "metadata.category.check.price"), func(b bool) {})
	cv.check["Length"] = ttw.NewCheck(lang.X("metadata.category.check.length", "metadata.category.check.length"), func(b bool) {})
	cv.check["Volume"] = ttw.NewCheck(lang.X("metadata.category.check.volume", "metadata.category.check.volume"), func(b bool) {})
	cv.check["Weight"] = ttw.NewCheck(lang.X("metadata.category.check.weight", "metadata.category.check.weight"), func(b bool) {})

	cv.entry["Name"] = midget.NewEntry()

	cv.label["Name"] = ttw.NewLabel(lang.X("metadata.form.name", "metadata.form.name"))
	cv.label["Parent"] = ttw.NewLabel(lang.X("metadata.form.parent", "metadata.form.parent"))

	cv.selects["Parent"] = ttw.NewSelect([]string{}, func(s string) {})

	form := container.New(layout.NewFormLayout(),
		cv.label["Parent"], cv.selects["Parent"],
		cv.label["Name"], cv.entry["Name"],
		layout.NewSpacer(), cv.check["Price"],
		layout.NewSpacer(), container.NewHBox(
			cv.check["Length"],
			cv.check["Volume"],
			cv.check["Weight"],
		),
	)

	cv.toolbar = widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			cv.toolbar.Items[0].(*widget.ToolbarAction).Disable()
			go func() {
				b.Metadata.CreateNewCategory()
				fyne.Do(func() {
					time.Sleep(300 * time.Millisecond)
					cv.toolbar.Items[0].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
			cv.toolbar.Items[1].(*widget.ToolbarAction).Disable()
			go func() {
				b.Metadata.DeleteCategory()
				fyne.Do(func() {
					time.Sleep(300 * time.Millisecond)
					cv.toolbar.Items[1].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			cv.toolbar.Items[2].(*widget.ToolbarAction).Disable()
			go func() {
				b.Metadata.CopyCategory()
				fyne.Do(func() {
					time.Sleep(300 * time.Millisecond)
					cv.toolbar.Items[2].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
	)

	list := container.NewBorder(cv.toolbar, nil, nil, nil, cv.tree)
	cv.container = container.NewHSplit(list, form)
	cv.container.SetOffset(0.25)
	cv.Clear()
	cv.Disable()
	return cv
}

func (c *categoryView) Clear() {
	c.entry.Clear()
	c.selects.Clear()
}
func (c *categoryView) Disable() {
	c.entry.Disable()
	c.selects.Disable()
}
func (c *categoryView) Enable() {
	c.entry.Enable()
	c.selects.Enable()
}
func (c *categoryView) Load(b *backend.Backend, id backend.CatID) {
	b.Metadata.Categories.AddListener(binding.NewDataListener(func() {
		categories, _ := b.Metadata.Categories.Get()
		self, _ := id.Name()
		remove(categories, self)
		c.selects["Parent"].SetOptions(append([]string{lang.L("None")}, categories...))
	}))

	c.entry["Name"].Bind(id.Category().Name)
	c.check["Price"].Bind(id.Category().Config["ShowPrice"])
	c.check["Length"].Bind(id.Category().Config["ShowLength"])
	c.check["Volume"].Bind(id.Category().Config["ShowVolume"])
	c.check["Weight"].Bind(id.Category().Config["ShowWeight"])

	c.selects["Parent"].Bind(id.Category().Parent)

	p, _ := id.ParentID()
	if p == 0 {
		c.selects["Parent"].SetSelectedIndex(0)
	}

	id.Category().Name.AddListener(binding.NewDataListener(func() { b.Metadata.UpdateCatList() }))
	id.Category().Parent.AddListener(binding.NewDataListener(func() { b.Metadata.UpdateCatList() }))
}
func (c *categoryView) Unload() {}
