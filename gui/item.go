package gui

import (
	"UppSpar/backend"
	"UppSpar/backend/bridge"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"

	midget "github.com/assholehoff/fyne-midget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type itemView struct {
	bound     backend.ItemID
	container *fyne.Container
	formView  *formView
	toolbar   *fyne.Container
	listView  *fyne.Container
	list      *widget.List
	listTools *widget.Toolbar
	imgView   *fyne.Container
}

func newItemView(a *App) *itemView {
	b := a.backend
	iv := &itemView{
		formView: newFormView(b),
	}
	newSaveFileDialog := func() *dialog.FileDialog {
		d := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer != nil {
				// save and export
				b.Items.ExportExcel(writer.URI().Path())
			} else {
				// close
				return
			}
		}, a.window)
		d.Resize(fyne.NewSize(900, 600))
		d.SetTitleText(lang.X("dialog.save.excel.title", "dialog.save.excel.title"))
		d.SetConfirmText(lang.L("Export"))
		d.SetDismissText(lang.L("Close"))
		d.SetFileName("UppSpar-" + time.Now().Format("20060102-150405") + ".xlsx")
		d.SetFilter(storage.NewMimeTypeFileFilter([]string{"application/excel"}))
		return d
	}
	data := b.Items.ItemIDList
	createItem := func() fyne.CanvasObject {
		l := widget.NewLabel("Item name template")
		return l
	}
	updateItem := func(di binding.DataItem, co fyne.CanvasObject) {
		val, err := di.(binding.Untyped).Get()
		if err != nil {
			panic(err)
		}
		ItemID := val.(backend.ItemID)
		co.(*widget.Label).Bind(ItemID.Item().Name)
	}
	iv.listTools = widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			iv.listTools.Items[0].(*widget.ToolbarAction).Disable()
			go func() {
				id, _ := b.Items.CreateNewItem()
				index, _ := b.Items.GetListItemIDFor(id)
				fyne.Do(func() {
					iv.list.Select(index)
					time.Sleep(100 * time.Millisecond)
					iv.listTools.Items[0].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
			iv.listTools.Items[1].(*widget.ToolbarAction).Disable()
			go func() {
				items, err := b.Items.ItemIDSelection.Get()
				if err != nil {
					log.Println(err)
					return
				}
				fyne.Do(func() {
					iv.list.UnselectAll()
					for _, item := range items {
						b.Items.DeleteItem(item.(backend.ItemID))
					}
					log.Println(b.Items.ItemIDSelection)
					time.Sleep(100 * time.Millisecond)
					iv.listTools.Items[1].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			iv.listTools.Items[2].(*widget.ToolbarAction).Disable()
			go func() {
				items, err := b.Items.ItemIDSelection.Get()
				if err != nil {
					log.Println(err)
					return
				}
				fyne.Do(func() {
					iv.list.UnselectAll()
					log.Println(b.Items.ItemIDSelection)
					var id backend.ItemID
					for _, item := range items {
						id, err = b.Items.CopyItem(item.(backend.ItemID))
						if err != nil {
							log.Println(err)
						}
					}
					index, err := b.Items.GetListItemIDFor(id)
					if err != nil {
						log.Println(err)
					}
					iv.list.Select(index)
					time.Sleep(100 * time.Millisecond)
					iv.listTools.Items[2].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.StorageIcon(), func() {
			d := newSaveFileDialog()
			d.Show()
		}),
	)

	iv.list = widget.NewListWithData(data, createItem, updateItem)
	iv.list.OnSelected = func(id widget.ListItemID) {
		i, err := b.Items.GetItemIDFor(id)
		if err != nil {
			log.Println(err)
			return
		}
		b.Items.SelectItem(i)
	}
	iv.list.OnUnselected = func(id widget.ListItemID) {
		i, err := b.Items.GetItemIDFor(id)
		if err != nil {
			log.Println(err)
			iv.formView.Clear()
			return
		}
		b.Items.UnselectItem(i)
	}
	b.Items.ItemIDSelection.AddListener(binding.NewDataListener(func() {
		ids, err := b.Items.ItemIDSelection.Get()
		if err != nil {
			log.Println(err)
			iv.formView.Clear()
			return
		}
		if len(ids) < 1 {
			iv.formView.Clear()
			return
		}
		if len(ids) == 1 {
			if iv.bound == ids[0].(backend.ItemID) {
				return
			}
			iv.formView.Clear()
			iv.formView.Bind(b, ids[0].(backend.ItemID))
			return
		}
		if len(ids) > 1 {
			iv.formView.Clear()
			// TODO bind a dummy structure to edit multiple as one
			iv.formView.Bind(b, ids[0].(backend.ItemID))
			return
		}
	}))
	searchEntry := xwidget.NewCompletionEntry([]string{})
	searchEntry.Bind(b.Items.Search.Term)
	searchEntry.OnChanged = func(s string) {
		if len(s) < 3 {
			searchEntry.HideCompletion()
			return
		}
		results, err := b.Items.Search.Completions.Get()
		if err != nil {
			panic(err)
		}
		searchEntry.SetOptions(results)
		searchEntry.ShowCompletion()
	}
	searchKeys := make(bridge.Checks)
	searchKeys["Name"] = ttw.NewCheckWithData("Art", b.Items.Search.Scope["Name"])
	searchKeys["Manufacturer"] = ttw.NewCheckWithData("Mfr", b.Items.Search.Scope["Manufacturer"])
	searchKeys["ModelName"] = ttw.NewCheckWithData("Mdl", b.Items.Search.Scope["ModelName"])
	searchKeys["ModelDesc"] = ttw.NewCheckWithData("Dsc", b.Items.Search.Scope["ModelDesc"])
	toolbarLeft := container.NewGridWithRows(1,
		// layout.NewSpacer(),
		container.NewHBox(
			searchKeys["Name"],
			searchKeys["Manufacturer"],
			searchKeys["ModelName"],
			searchKeys["ModelDesc"],
		),
		// widget.NewSelect([]string{
		// 	lang.X("item.form.label.name", "item.form.label.name"),
		// 	lang.X("item.form.label.manufacturer", "item.form.label.manufacturer"),
		// }, func(s string) {
		// 	// TODO redo this to fit map[string]bool
		// 	// b.Items.GetItemIDs()
		// }),
		widget.NewSelect([]string{
			lang.X("form.select.search.beginswith", "form.select.search.beginswith"),
			lang.X("form.select.search.endswith", "form.select.search.endswith"),
			lang.X("form.select.search.contains", "form.select.search.contains"),
			lang.X("form.select.search.equals", "form.select.search.equals"),
			// lang.X("form.select.search.regexp", "form.select.search.regexp"),
		}, func(s string) {
			switch s {
			case lang.X("form.select.search.beginswith", "form.select.search.beginswith"):
				b.Items.Search.Match = backend.MatchBeginsWith
			case lang.X("form.select.search.endswith", "form.select.search.endswith"):
				b.Items.Search.Match = backend.MatchEndsWith
			case lang.X("form.select.search.equals", "form.select.search.equals"):
				b.Items.Search.Match = backend.MatchEquals
			case lang.X("form.select.search.contains", "form.select.search.contains"):
				b.Items.Search.Match = backend.MatchContains
			default:
				b.Items.Search.Match = backend.MatchContains
			}
			b.Items.GetItemIDs()
		}),
	)
	toolbarRight := container.NewGridWithRows(1,
		midget.NewLabel(
			lang.X("form.select.sortby.text", "form.select.sortby.text"),
			"",
			"",
		),
		widget.NewSelect([]string{
			lang.X("form.select.sortby.itemid", "form.select.sortby.itemid"),
			lang.X("form.select.sortby.name", "form.select.sortby.name"),
			lang.X("form.select.sortby.manufacturer", "form.select.sortby.manufacturer"),
			lang.X("form.select.sortby.datecreated", "form.select.sortby.datecreated"),
			lang.X("form.select.sortby.datemodified", "form.select.sortby.datemodified"),
		}, func(s string) {
			switch s {
			case lang.X("form.select.sortby.itemid", "form.select.sortby.itemid"):
				b.Items.Search.SortBy = backend.SearchKeyItemID
			case lang.X("form.select.sortby.name", "form.select.sortby.name"):
				b.Items.Search.SortBy = backend.SearchKeyName
			case lang.X("form.select.sortby.manufacturer", "form.select.sortby.manufacturer"):
				b.Items.Search.SortBy = backend.SearchKeyManufacturer
			case lang.X("form.select.sortby.datecreated", "form.select.sortby.datecreated"):
				b.Items.Search.SortBy = backend.SearchKeyDateCreated
			case lang.X("form.select.sortby.datemodified", "form.select.sortby.datemodified"):
				b.Items.Search.SortBy = backend.SearchKeyDateModified
			}
			b.Items.GetItemIDs()
		}),
		widget.NewSelect([]string{
			lang.X("form.select.sortorder.ascending", "form.select.sortorder.ascending"),
			lang.X("form.select.sortorder.descending", "form.select.sortorder.descending"),
		}, func(s string) {
			if s == lang.X("form.select.sortorder.ascending", "form.select.sortorder.ascending") {
				b.Items.Search.Order = backend.SortAscending
			} else {
				b.Items.Search.Order = backend.SortDescending
			}
			b.Items.GetItemIDs()
		}))
	filterWidth := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Width")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Width))
	filterHeight := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Height")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Height))
	filterDepth := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Depth")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Depth))
	filterVolume := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Volume")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Volume))
	filterWeight := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Weight")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Weight))
	filterBar := container.NewGridWithRows(1,
		filterWidth, filterHeight, filterDepth, filterVolume, filterWeight,
	)
	searchBar := container.NewBorder(nil, nil, toolbarLeft, toolbarRight, searchEntry)
	iv.toolbar = container.NewBorder(searchBar, filterBar, nil, nil)
	// toolbarLeft.Objects[1].(*widget.Select).SetSelectedIndex(0)  // search key
	toolbarLeft.Objects[1].(*widget.Select).SetSelectedIndex(2)  // search type
	toolbarRight.Objects[1].(*widget.Select).SetSelectedIndex(0) // sort by
	toolbarRight.Objects[2].(*widget.Select).SetSelectedIndex(0) // sort order
	imgView := container.NewBorder(nil, nil, nil, nil)
	listView := container.NewBorder(nil, nil, nil, nil, iv.list)
	formView := container.NewBorder(nil, nil, nil, imgView, iv.formView.container)
	split := container.NewHSplit(listView, formView)
	split.SetOffset(0.25)
	toolbar := container.NewBorder(nil, nil, iv.listTools, nil, iv.toolbar)
	statusbar := container.NewGridWithRows(1)
	iv.container = container.NewBorder(toolbar, statusbar, nil, nil, split)

	return iv
}

type formEntries struct {
	Name         *midget.Entry
	Price        *midget.Entry
	Vat          *midget.Entry
	ImgURL1      *midget.Entry
	ImgURL2      *midget.Entry
	ImgURL3      *midget.Entry
	ImgURL4      *midget.Entry
	ImgURL5      *midget.Entry
	SpecsURL     *midget.Entry
	LongDesc     *midget.Entry
	Manufacturer *widget.SelectEntry
	ModelName    *widget.SelectEntry
	ModelDesc    *midget.Entry
	ModelURL     *midget.Entry
	Notes        *midget.Entry
	Width        *midget.Entry
	Height       *midget.Entry
	Depth        *midget.Entry
	Volume       *midget.Entry
	Weight       *midget.Entry
}

type formDataLabels struct {
	ItemID       *widget.Label
	AddDesc      *widget.Label
	LongDesc     *widget.Label
	DateCreated  *widget.Label
	DateModified *widget.Label
}

type formLabels struct {
	ItemID        *widget.Label
	Name          *widget.Label
	Category      *widget.Label
	Currency      *widget.Label
	Price         *widget.Label
	Vat           *widget.Label
	ImgURL1       *widget.Label
	ImgURL2       *widget.Label
	ImgURL3       *widget.Label
	ImgURL4       *widget.Label
	ImgURL5       *widget.Label
	SpecsURL      *widget.Label
	AddDesc       *widget.Label
	LongDesc      *widget.Label
	Manufacturer  *widget.Label
	ModelName     *widget.Label
	ModelDesc     *widget.Label
	ModelURL      *widget.Label
	Notes         *widget.Label
	Descr         *widget.Label
	Dimensions    *widget.Label
	Width         *widget.Label
	Height        *widget.Label
	Depth         *widget.Label
	Volume        *widget.Label
	Weight        *widget.Label
	Status        *widget.Label
	DateCreated   *widget.Label
	DateModified  *widget.Label
	Condition     *widget.Label
	Functionality *widget.Label
}

type formSelects struct {
	Category   *widget.Select
	LengthUnit *widget.Select
	VolumeUnit *widget.Select
	WeightUnit *widget.Select
	Status     *widget.Select
}

type formView struct {
	container *fyne.Container
	entries   *formEntries
	selects   *formSelects
	labels    *formLabels
	values    *formDataLabels
}

func newFormView(b *backend.Backend) *formView {
	var categories []string
	fetchCategories := func() []string {
		cats, _ := b.Metadata.Categories.Get()
		return cats
	}
	categories = fetchCategories()

	lengthUnits := []string{"mm", "cm", "dm", "m"}
	volumeUnits := []string{"ml", "cl", "dl", "l"}
	weightUnits := []string{"g", "hg", "kg"}

	itemStatus := func() []string {
		var ss []string
		stats, _ := b.Metadata.ItemStatusIDList.Get()
		for _, stat := range stats {
			ss = append(ss, stat.(backend.ItemStatusID).LString())
		}
		return ss
	}()

	v := &formView{
		entries: &formEntries{
			Name:         midget.NewEntry(),
			Price:        midget.NewEntry(),
			Vat:          midget.NewEntry(),
			ImgURL1:      midget.NewEntry(),
			ImgURL2:      midget.NewEntry(),
			ImgURL3:      midget.NewEntry(),
			ImgURL4:      midget.NewEntry(),
			ImgURL5:      midget.NewEntry(),
			SpecsURL:     midget.NewEntry(),
			LongDesc:     midget.NewEntry(),
			Manufacturer: widget.NewSelectEntry([]string{}),
			ModelName:    widget.NewSelectEntry([]string{}),
			ModelDesc:    midget.NewEntry(),
			ModelURL:     midget.NewEntry(),
			Notes:        midget.NewEntry(),
			Width:        midget.NewEntry(),
			Height:       midget.NewEntry(),
			Depth:        midget.NewEntry(),
			Volume:       midget.NewEntry(),
			Weight:       midget.NewEntry(),
		},
		labels: &formLabels{
			ItemID:       widget.NewLabel(lang.X("item.form.label.itemid", "item.form.label.itemid")),
			Name:         widget.NewLabel(lang.X("item.form.label.name", "item.form.label.name")),
			Category:     widget.NewLabel(lang.X("item.form.label.category", "item.form.label.category")),
			Currency:     widget.NewLabel("SEK"),
			Price:        widget.NewLabel(lang.X("item.form.label.price", "item.form.label.price")),
			Vat:          widget.NewLabel(lang.X("item.form.label.vat", "item.form.label.vat")),
			ImgURL1:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
			ImgURL2:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
			ImgURL3:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
			ImgURL4:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
			ImgURL5:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
			SpecsURL:     widget.NewLabel(lang.X("item.form.label.specsurl", "item.form.label.specsurl")),
			AddDesc:      widget.NewLabel(lang.X("item.form.label.adddesc", "item.form.label.adddesc")),
			LongDesc:     widget.NewLabel(lang.X("item.form.label.longdesc", "item.form.label.longdesc")),
			Manufacturer: widget.NewLabel(lang.X("item.form.label.manufacturer", "item.form.label.manufacturer")),
			ModelName:    widget.NewLabel(lang.X("item.form.label.modelname", "item.form.label.modelname")),
			ModelDesc:    widget.NewLabel(lang.X("item.form.label.modeldesc", "item.form.label.modeldesc")),
			ModelURL:     widget.NewLabel(lang.X("item.form.label.modelurl", "item.form.label.modelurl")),
			Notes:        widget.NewLabel(lang.X("item.form.label.notes", "item.form.label.notes")),
			Dimensions:   widget.NewLabel(lang.X("item.form.label.dimensions", "item.form.label.dimensions")),
			Width:        widget.NewLabel(lang.X("item.form.label.width", "item.form.label.width")),
			Height:       widget.NewLabel(lang.X("item.form.label.height", "item.form.label.height")),
			Depth:        widget.NewLabel(lang.X("item.form.label.depth", "item.form.label.depth")),
			Volume:       widget.NewLabel(lang.X("item.form.label.volume", "item.form.label.volume")),
			Weight:       widget.NewLabel(lang.X("item.form.label.weight", "item.form.label.weight")),
			Status:       widget.NewLabel(lang.X("item.form.label.status", "item.form.label.status")),
			DateCreated:  widget.NewLabel(lang.X("item.form.label.datecreated", "item.form.label.datecreated")),
			DateModified: widget.NewLabel(lang.X("item.form.label.datemodified", "item.form.label.datemodified")),
		},
		values: &formDataLabels{
			ItemID:       widget.NewLabel(lang.X("item.form.data.itemid", "item.form.data.itemid")),
			AddDesc:      widget.NewLabel(lang.X("item.form.data.adddesc", "item.form.data.adddesc")),
			LongDesc:     widget.NewLabel(lang.X("item.form.data.longdesc", "item.form.data.longdesc")),
			DateCreated:  widget.NewLabel(lang.X("item.form.data.datecreated", "item.form.data.datecreated")),
			DateModified: widget.NewLabel(lang.X("item.form.data.datemodified", "item.form.data.datemodified")),
		},
		selects: &formSelects{
			Category:   widget.NewSelect(categories, func(s string) {}),
			LengthUnit: widget.NewSelect(lengthUnits, func(s string) {}),
			VolumeUnit: widget.NewSelect(volumeUnits, func(s string) {}),
			WeightUnit: widget.NewSelect(weightUnits, func(s string) {}),
			Status:     widget.NewSelect(itemStatus, func(s string) {}),
		},
	}
	b.Metadata.CatIDList.AddListener(binding.NewDataListener(func() {
		categories = fetchCategories()
		v.selects.Category.Options = categories
		v.selects.Category.Refresh()
	}))
	v.entries.ModelDesc.MultiLine = true
	v.entries.ModelDesc.SetMinRowsVisible(5)
	v.entries.ModelDesc.Wrapping = fyne.TextWrapWord
	v.entries.LongDesc.MultiLine = true
	v.entries.LongDesc.SetMinRowsVisible(5)
	v.entries.LongDesc.Wrapping = fyne.TextWrapWord
	v.entries.Notes.MultiLine = true
	v.entries.Notes.SetMinRowsVisible(5)
	v.entries.Notes.Wrapping = fyne.TextWrapWord

	b.Metadata.MfrNameList.AddListener(binding.NewDataListener(func() {
		manufacturers, _ := b.Metadata.MfrNameList.Get()
		v.entries.Manufacturer.SetOptions(manufacturers)
	}))

	idbox := container.NewBorder(nil, nil, v.values.ItemID, nil, container.NewHBox(v.selects.Status))
	spacebox := container.NewGridWithRows(1,
		container.NewBorder(nil, nil, v.labels.Width, nil, v.entries.Width),
		container.NewBorder(nil, nil, v.labels.Height, nil, v.entries.Height),
		container.NewBorder(nil, nil, v.labels.Depth, nil, v.entries.Depth),
		v.selects.LengthUnit,
	)
	massbox := container.NewGridWithRows(1,
		container.NewBorder(nil, nil, v.labels.Volume, v.selects.VolumeUnit, v.entries.Volume),
		container.NewBorder(nil, nil, v.labels.Weight, v.selects.WeightUnit, v.entries.Weight),
	)
	v.container = container.New(layout.NewFormLayout(),
		layout.NewSpacer(), container.NewHBox(v.labels.DateCreated, v.values.DateCreated),
		layout.NewSpacer(), container.NewHBox(v.labels.DateModified, v.values.DateModified),
		v.labels.ItemID, idbox,
		v.labels.Name, v.entries.Name,
		v.labels.Category, v.selects.Category,
		v.labels.Manufacturer, v.entries.Manufacturer,
		v.labels.ModelName, v.entries.ModelName,
		v.labels.ModelDesc, v.entries.ModelDesc,
		v.labels.ModelURL, v.entries.ModelURL,
		v.labels.Dimensions, spacebox,
		layout.NewSpacer(), massbox,
		v.labels.Price, container.NewBorder(nil, nil, nil, v.labels.Currency, v.entries.Price),
		v.labels.ImgURL1, v.entries.ImgURL1,
		v.labels.Notes, v.entries.Notes,
		layout.NewSpacer(), widget.NewLabel(" "),
		layout.NewSpacer(), widget.NewRichTextFromMarkdown(`### `+lang.L("Preview")),
		v.labels.LongDesc, v.values.LongDesc,
		v.labels.AddDesc, v.values.AddDesc,
	)
	v.Clear()
	return v
}

func (v formView) Bind(b *backend.Backend, id backend.ItemID) {
	v.Clear()
	v.Enable()

	v.values.ItemID.Bind(id.Item().ItemIDString)
	v.values.DateCreated.Bind(id.Item().DateCreated)
	v.values.DateModified.Bind(id.Item().DateModified)
	v.values.AddDesc.Bind(id.Item().AddDesc)
	v.values.LongDesc.Bind(id.Item().LongDesc)

	v.entries.Name.Bind(id.Item().Name)
	v.entries.Price.Bind(id.Item().PriceString)
	v.entries.Vat.Bind(id.Item().VatString)
	v.entries.ImgURL1.Bind(id.Item().ImgURL1)
	v.entries.SpecsURL.Bind(id.Item().SpecsURL)
	v.entries.Manufacturer.Bind(id.Item().Manufacturer)
	v.entries.ModelName.Bind(id.Item().ModelName)
	v.entries.ModelDesc.Bind(id.Item().ModelDesc)
	v.entries.ModelURL.Bind(id.Item().ModelURL)
	v.entries.Notes.Bind(id.Item().Notes)
	v.entries.Width.Bind(id.Item().WidthString)
	v.entries.Height.Bind(id.Item().HeightString)
	v.entries.Depth.Bind(id.Item().DepthString)
	v.entries.Volume.Bind(id.Item().VolumeString)
	v.entries.Weight.Bind(id.Item().WeightString)

	v.selects.Category.Bind(id.Item().Category)
	v.selects.LengthUnit.Bind(id.Item().LengthUnit)
	v.selects.VolumeUnit.Bind(id.Item().VolumeUnit)
	v.selects.WeightUnit.Bind(id.Item().WeightUnit)
	v.selects.Status.Bind(id.Item().ItemStatus)

	id.Item().Manufacturer.AddListener(binding.NewDataListener(func() {
		models := func() []string {
			var names []string
			ids := id.Item().MfrID.Children()
			for _, id := range ids {
				name, _ := id.Name()
				names = append(names, name)
			}
			return names
		}()
		v.entries.ModelName.SetOptions(models)
	}))

	/* This step is needed because child categories have spaces prepended to them in the select list */
	cat, _ := id.Item().Category.Get()
	v.selects.Category.SetSelectedIndex(b.Metadata.GetListItemIDFor(cat))

	id.Item().CatID.Category().Config["ShowPrice"].AddListener(binding.NewDataListener(func() {
		if cid, _ := id.CatID(); cid == backend.CatID(1) {
			v.hideStatus()
			v.hideCategory()
			v.hideImgURL()
			v.hideSpecsURL()
			v.hideMfrModel()
			v.hideLength()
			v.hideVolume()
			v.hideWeight()
			v.showPrice()
			v.hidePreviewAddDesc()
			return
		} else {
			v.showStatus()
			v.showCategory()
			v.showImgURL()
			v.showSpecsURL()
			v.showMfrModel()
			v.showLength()
			v.showWeight()
			v.showVolume()
			v.showPreviewAddDesc()
		}
		p, _ := id.Item().CatID.Category().Config["ShowPrice"].Get()
		if p {
			v.showPrice()
		} else {
			v.hidePrice()
		}
	}))
	id.Item().CatID.Category().Config["ShowLength"].AddListener(binding.NewDataListener(func() {
		p, _ := id.Item().CatID.Category().Config["ShowLength"].Get()
		if p {
			v.showLength()
		} else {
			v.hideLength()
		}
	}))
	id.Item().CatID.Category().Config["ShowVolume"].AddListener(binding.NewDataListener(func() {
		p, _ := id.Item().CatID.Category().Config["ShowVolume"].Get()
		if p {
			v.showVolume()
		} else {
			v.hideVolume()
		}
	}))
	id.Item().CatID.Category().Config["ShowWeight"].AddListener(binding.NewDataListener(func() {
		p, _ := id.Item().CatID.Category().Config["ShowWeight"].Get()
		if p {
			v.showWeight()
		} else {
			v.hideWeight()
		}
	}))
}
func (v formView) Clear() {
	v.entries.Name.Unbind()
	v.entries.Price.Unbind()
	v.entries.Vat.Unbind()
	v.entries.ImgURL1.Unbind()
	v.entries.SpecsURL.Unbind()
	v.entries.Manufacturer.Unbind()
	v.entries.ModelName.Unbind()
	v.entries.ModelURL.Unbind()
	v.entries.Notes.Unbind()
	v.entries.Width.Unbind()
	v.entries.Height.Unbind()
	v.entries.Depth.Unbind()
	v.entries.Volume.Unbind()
	v.entries.Weight.Unbind()

	v.values.ItemID.Unbind()
	v.values.DateCreated.Unbind()
	v.values.DateModified.Unbind()
	v.values.AddDesc.Unbind()
	v.values.LongDesc.Unbind()

	v.selects.Category.Unbind()
	v.selects.LengthUnit.Unbind()
	v.selects.VolumeUnit.Unbind()
	v.selects.WeightUnit.Unbind()
	v.selects.Status.Unbind()

	v.entries.Name.SetText("")
	v.entries.Price.SetText("")
	v.entries.Vat.SetText("")
	v.entries.ImgURL1.SetText("")
	v.entries.SpecsURL.SetText("")
	v.entries.Manufacturer.SetText("")
	v.entries.ModelName.SetText("")
	v.entries.ModelURL.SetText("")
	v.entries.Notes.SetText("")
	v.entries.Volume.SetText("")
	v.entries.Width.SetText("")
	v.entries.Height.SetText("")
	v.entries.Depth.SetText("")
	v.entries.Volume.SetText("")
	v.entries.Weight.SetText("")

	v.values.ItemID.SetText("0000000000")
	v.values.DateCreated.SetText("")
	v.values.DateModified.SetText("")
	v.values.AddDesc.SetText("")
	v.values.LongDesc.SetText("")

	v.selects.Category.ClearSelected()
	v.selects.LengthUnit.ClearSelected()
	v.selects.VolumeUnit.ClearSelected()
	v.selects.WeightUnit.ClearSelected()
	v.selects.Status.ClearSelected()

	v.Disable()
}

func (v formView) Disable() {
	v.entries.Name.Disable()
	v.entries.Price.Disable()
	v.entries.Vat.Disable()
	v.entries.ImgURL1.Disable()
	v.entries.SpecsURL.Disable()
	v.entries.LongDesc.Disable()
	v.entries.Manufacturer.Disable()
	v.entries.ModelName.Disable()
	v.entries.ModelURL.Disable()
	v.entries.Notes.Disable()
	v.entries.Width.Disable()
	v.entries.Height.Disable()
	v.entries.Depth.Disable()
	v.entries.Volume.Disable()
	v.entries.Weight.Disable()

	v.selects.Category.Disable()
	v.selects.LengthUnit.Disable()
	v.selects.VolumeUnit.Disable()
	v.selects.WeightUnit.Disable()
	v.selects.Status.Disable()
}
func (v formView) Enable() {
	v.entries.Name.Enable()
	v.entries.Price.Enable()
	v.entries.Vat.Enable()
	v.entries.ImgURL1.Enable()
	v.entries.SpecsURL.Enable()
	v.entries.LongDesc.Enable()
	v.entries.Manufacturer.Enable()
	v.entries.ModelName.Enable()
	v.entries.ModelURL.Enable()
	v.entries.Notes.Enable()
	v.entries.Width.Enable()
	v.entries.Height.Enable()
	v.entries.Depth.Enable()
	v.entries.Volume.Enable()
	v.entries.Weight.Enable()

	v.selects.Category.Enable()
	v.selects.LengthUnit.Enable()
	v.selects.VolumeUnit.Enable()
	v.selects.WeightUnit.Enable()
	v.selects.Status.Enable()
}

func (v *formView) showPreviewAddDesc() {
	v.labels.AddDesc.Show()
	v.values.AddDesc.Show()
}
func (v *formView) hidePreviewAddDesc() {
	v.labels.AddDesc.Hide()
	v.values.AddDesc.Hide()
}
func (v *formView) showStatus() {
	v.labels.Status.Show()
	v.selects.Status.Enable()
	v.selects.Status.Show()
}
func (v *formView) hideStatus() {
	v.labels.Status.Hide()
	v.selects.Status.Disable()
	v.selects.Status.Hide()
}
func (v *formView) showCategory() {
	v.labels.Category.Show()
	v.selects.Category.Enable()
	v.selects.Category.Show()
}
func (v *formView) hideCategory() {
	v.labels.Category.Hide()
	v.selects.Category.Disable()
	v.selects.Category.Hide()
}
func (v *formView) showImgURL() {
	v.entries.ImgURL1.Enable()
	v.entries.ImgURL2.Enable()
	v.entries.ImgURL3.Enable()
	v.entries.ImgURL4.Enable()
	v.entries.ImgURL5.Enable()

	v.labels.ImgURL1.Show()
	v.labels.ImgURL2.Show()
	v.labels.ImgURL3.Show()
	v.labels.ImgURL4.Show()
	v.labels.ImgURL5.Show()

	v.entries.ImgURL1.Show()
	v.entries.ImgURL2.Show()
	v.entries.ImgURL3.Show()
	v.entries.ImgURL4.Show()
	v.entries.ImgURL5.Show()
}
func (v *formView) hideImgURL() {
	v.entries.ImgURL1.Disable()
	v.entries.ImgURL2.Disable()
	v.entries.ImgURL3.Disable()
	v.entries.ImgURL4.Disable()
	v.entries.ImgURL5.Disable()

	v.labels.ImgURL1.Hide()
	v.labels.ImgURL2.Hide()
	v.labels.ImgURL3.Hide()
	v.labels.ImgURL4.Hide()
	v.labels.ImgURL5.Hide()

	v.entries.ImgURL1.Hide()
	v.entries.ImgURL2.Hide()
	v.entries.ImgURL3.Hide()
	v.entries.ImgURL4.Hide()
	v.entries.ImgURL5.Hide()
}
func (v *formView) showSpecsURL() {
	v.entries.SpecsURL.Enable()
	v.entries.SpecsURL.Show()

	v.labels.SpecsURL.Show()
}
func (v *formView) hideSpecsURL() {
	v.entries.SpecsURL.Disable()
	v.entries.SpecsURL.Hide()

	v.labels.SpecsURL.Hide()
}
func (v *formView) showMfrModel() {
	v.entries.Manufacturer.Enable()
	v.entries.ModelName.Enable()
	v.entries.ModelURL.Enable()

	v.entries.Manufacturer.Show()
	v.entries.ModelName.Show()
	v.entries.ModelURL.Show()

	v.labels.Manufacturer.Show()
	v.labels.ModelName.Show()
	v.labels.ModelURL.Show()
}
func (v *formView) hideMfrModel() {
	v.entries.Manufacturer.Disable()
	v.entries.ModelName.Disable()
	v.entries.ModelURL.Disable()

	v.entries.Manufacturer.Hide()
	v.entries.ModelName.Hide()
	v.entries.ModelURL.Hide()

	v.labels.Manufacturer.Hide()
	v.labels.ModelName.Hide()
	v.labels.ModelURL.Hide()
}
func (v *formView) showPrice() {
	v.entries.Price.Enable()
	v.entries.Price.Show()
	v.labels.Currency.Show()
	v.labels.Price.Show()
}
func (v *formView) hidePrice() {
	v.entries.Price.Disable()
	v.entries.Price.Hide()
	v.labels.Currency.Hide()
	v.labels.Price.Hide()
}
func (v *formView) showLength() {
	v.entries.Width.Enable()
	v.entries.Height.Enable()
	v.entries.Depth.Enable()
	v.entries.Width.Show()
	v.entries.Height.Show()
	v.entries.Depth.Show()
	v.labels.Width.Show()
	v.labels.Height.Show()
	v.labels.Depth.Show()
	v.labels.Dimensions.Show()
	v.selects.LengthUnit.Enable()
	v.selects.LengthUnit.Show()
}
func (v *formView) hideLength() {
	v.entries.Width.Disable()
	v.entries.Height.Disable()
	v.entries.Depth.Disable()
	v.entries.Width.Hide()
	v.entries.Height.Hide()
	v.entries.Depth.Hide()
	v.labels.Width.Hide()
	v.labels.Height.Hide()
	v.labels.Depth.Hide()
	v.labels.Dimensions.Hide()
	v.selects.LengthUnit.Disable()
	v.selects.LengthUnit.Hide()
}
func (v *formView) showVolume() {
	v.entries.Volume.Enable()
	v.entries.Volume.Show()
	v.labels.Volume.Show()
	v.selects.VolumeUnit.Enable()
	v.selects.VolumeUnit.Show()
}
func (v *formView) hideVolume() {
	v.entries.Volume.Disable()
	v.entries.Volume.Hide()
	v.labels.Volume.Hide()
	v.selects.VolumeUnit.Disable()
	v.selects.VolumeUnit.Hide()
}
func (v *formView) showWeight() {
	v.entries.Weight.Enable()
	v.entries.Weight.Show()
	v.labels.Weight.Show()
	v.selects.WeightUnit.Enable()
	v.selects.WeightUnit.Show()
}
func (v *formView) hideWeight() {
	v.entries.Weight.Disable()
	v.entries.Weight.Hide()
	v.labels.Weight.Hide()
	v.selects.WeightUnit.Disable()
	v.selects.WeightUnit.Hide()
}
