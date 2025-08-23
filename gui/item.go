package gui

import (
	"UppSpar/backend"
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
)

type itemView struct {
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
				fyne.Do(func() {
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
		iv.formView.Bind(i)
	}
	iv.list.OnUnselected = func(id widget.ListItemID) {
		i, err := b.Items.GetItemIDFor(id)
		if err != nil {
			log.Println(err)
			iv.formView.Clear()
			return
		}
		b.Items.UnselectItem(i)
		iv.formView.Clear()
	}
	searchEntry := xwidget.NewCompletionEntry([]string{})
	searchEntry.Bind(b.Items.SearchString)
	searchEntry.OnChanged = func(s string) {
		if len(s) < 3 {
			searchEntry.HideCompletion()
			return
		}
		results, err := b.Items.SearchResultUniqueCompletions.Get()
		if err != nil {
			panic(err)
		}
		searchEntry.SetOptions(results)
		searchEntry.ShowCompletion()
	}
	toolbarLeft := container.NewGridWithRows(1,
		layout.NewSpacer(),
		widget.NewSelect([]string{
			lang.X("item.form.name.text", "item.form.name.text"),
			lang.X("item.form.manufacturer.text", "item.form.manufacturer.text"),
		}, func(s string) {
			if s == b.Items.SearchKey().String() {
				return
			}
			switch s {
			case lang.X("item.form.name.text", "item.form.name.text"):
				b.Items.SetSearchKey(backend.SearchKeyName)
			case lang.X("item.form.manufacturer.text", "item.form.manufacturer.text"):
				b.Items.SetSearchKey(backend.SearchKeyManufacturer)
			}
			b.Items.Search()
		}),
		widget.NewSelect([]string{
			lang.X("form.select.search.beginswith", "form.select.search.beginswith"),
			lang.X("form.select.search.endswith", "form.select.search.endswith"),
			lang.X("form.select.search.contains", "form.select.search.contains"),
			lang.X("form.select.search.equals", "form.select.search.equals"),
			// lang.X("form.select.search.regexp", "form.select.search.regexp"),
		}, func(s string) {
			switch s {
			case lang.X("form.select.search.beginswith", "form.select.search.beginswith"):
				b.Items.SetSearchConfig(backend.BeginsWith)
			case lang.X("form.select.search.endswith", "form.select.search.endswith"):
				b.Items.SetSearchConfig(backend.EndsWith)
			case lang.X("form.select.search.equals", "form.select.search.equals"):
				b.Items.SetSearchConfig(backend.Equals)
			case lang.X("form.select.search.contains", "form.select.search.contains"):
				b.Items.SetSearchConfig(backend.Contains)
			// case lang.X("form.select.search.regexp", "form.select.search.regexp"):
			// 	b.Items.SetSearchConfig(backend.RegExp)
			default:
				b.Items.SetSearchConfig(backend.Contains)
			}
			b.Items.Search()
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
				b.Items.SetSortKey(backend.SearchKeyItemID)
			case lang.X("form.select.sortby.name", "form.select.sortby.name"):
				b.Items.SetSortKey(backend.SearchKeyName)
			case lang.X("form.select.sortby.manufacturer", "form.select.sortby.manufacturer"):
				b.Items.SetSortKey(backend.SearchKeyManufacturer)
			case lang.X("form.select.sortby.datecreated", "form.select.sortby.datecreated"):
				b.Items.SetSortKey(backend.SearchKeyDateCreated)
			case lang.X("form.select.sortby.datemodified", "form.select.sortby.datemodified"):
				b.Items.SetSortKey(backend.SearchKeyDateModified)
			}
			b.Items.Search()
		}),
		widget.NewSelect([]string{
			lang.X("form.select.sortorder.ascending", "form.select.sortorder.ascending"),
			lang.X("form.select.sortorder.descending", "form.select.sortorder.descending"),
		}, func(s string) {
			if s == lang.X("form.select.sortorder.ascending", "form.select.sortorder.ascending") {
				b.Items.SetSortOrder(backend.SortAscending)
			} else {
				b.Items.SetSortOrder(backend.SortDescending)
			}
			b.Items.Search()
		}))
	iv.toolbar = container.NewBorder(nil, nil, toolbarLeft, toolbarRight, searchEntry)
	toolbarLeft.Objects[1].(*widget.Select).SetSelectedIndex(0)  // search key
	toolbarLeft.Objects[2].(*widget.Select).SetSelectedIndex(2)  // search type
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
	name         *widget.Entry
	price        *widget.Entry
	vat          *widget.Entry
	imgurl1      *widget.Entry
	imgurl2      *widget.Entry
	imgurl3      *widget.Entry
	imgurl4      *widget.Entry
	imgurl5      *widget.Entry
	specsurl     *widget.Entry
	longdesc     *widget.Entry
	manufacturer *widget.Entry
	model        *widget.Entry
	modelurl     *widget.Entry
	notes        *widget.Entry
	width        *widget.Entry
	height       *widget.Entry
	depth        *widget.Entry
	volume       *widget.Entry
	weight       *widget.Entry
}

type formLabels struct {
	itemidtitle       *widget.Label
	itemid            *widget.Label
	name              *widget.Label
	category          *widget.Label
	price             *widget.Label
	vat               *widget.Label
	imgurl1           *widget.Label
	imgurl2           *widget.Label
	imgurl3           *widget.Label
	imgurl4           *widget.Label
	imgurl5           *widget.Label
	specsurl          *widget.Label
	adddesctitle      *widget.Label
	adddesc           *widget.Label
	longdesc          *widget.Label
	manufacturer      *widget.Label
	model             *widget.Label
	modelurl          *widget.Label
	notes             *widget.Label
	dimensions        *widget.Label
	width             *widget.Label
	height            *widget.Label
	depth             *widget.Label
	volume            *widget.Label
	weight            *widget.Label
	status            *widget.Label
	datecreatedtitle  *widget.Label
	datemodifiedtitle *widget.Label
	datecreated       *widget.Label
	datemodified      *widget.Label
}

type formSelects struct {
	category   *widget.Select
	lengthunit *widget.Select
	volumeunit *widget.Select
	weightunit *widget.Select
	status     *widget.Select
}

type formView struct {
	container *fyne.Container
	entries   *formEntries
	selects   *formSelects
	labels    *formLabels
}

func newFormView(b *backend.Backend) *formView {
	var categories []string
	fetchCategories := func() []string {
		var cats []string
		ids, err := b.Metadata.CatIDList.Get()
		if err != nil {
			log.Println(err)
		}
		for _, id := range ids {
			cat, err := id.(backend.CatID).Name()
			if err != nil {
				log.Println(err)
			}
			cats = append(cats, cat)
		}
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
			name:         widget.NewEntry(),
			price:        widget.NewEntry(),
			vat:          widget.NewEntry(),
			imgurl1:      widget.NewEntry(),
			imgurl2:      widget.NewEntry(),
			imgurl3:      widget.NewEntry(),
			imgurl4:      widget.NewEntry(),
			imgurl5:      widget.NewEntry(),
			specsurl:     widget.NewEntry(),
			longdesc:     widget.NewEntry(),
			manufacturer: widget.NewEntry(),
			model:        widget.NewEntry(),
			modelurl:     widget.NewEntry(),
			notes:        widget.NewEntry(),
			width:        widget.NewEntry(),
			height:       widget.NewEntry(),
			depth:        widget.NewEntry(),
			volume:       widget.NewEntry(),
			weight:       widget.NewEntry(),
		},
		labels: &formLabels{
			itemidtitle:       widget.NewLabel(lang.X("item.form.itemid.text", "item.form.itemid.text")),
			itemid:            widget.NewLabel("0000000000"),
			name:              widget.NewLabel(lang.X("item.form.name.text", "item.form.name.text")),
			category:          widget.NewLabel(lang.X("item.form.category.text", "item.form.category.text")),
			price:             widget.NewLabel(lang.X("item.form.price.text", "item.form.price.text")),
			vat:               widget.NewLabel(lang.X("item.form.vat.text", "item.form.vat.text")),
			imgurl1:           widget.NewLabel(lang.X("item.form.imgurl.text", "item.form.imgurl.text")),
			imgurl2:           widget.NewLabel(lang.X("item.form.imgurl.text", "item.form.imgurl.text")),
			imgurl3:           widget.NewLabel(lang.X("item.form.imgurl.text", "item.form.imgurl.text")),
			imgurl4:           widget.NewLabel(lang.X("item.form.imgurl.text", "item.form.imgurl.text")),
			imgurl5:           widget.NewLabel(lang.X("item.form.imgurl.text", "item.form.imgurl.text")),
			specsurl:          widget.NewLabel(lang.X("item.form.specsurl.text", "item.form.specsurl.text")),
			adddesctitle:      widget.NewLabel(lang.X("item.form.adddesc.text", "item.form.adddesc.text")),
			adddesc:           widget.NewLabel(lang.X("item.form.adddesc.text", "item.form.adddesc.text")),
			longdesc:          widget.NewLabel(lang.X("item.form.longdesc.text", "item.form.longdesc.text")),
			manufacturer:      widget.NewLabel(lang.X("item.form.manufacturer.text", "item.form.manufacturer.text")),
			model:             widget.NewLabel(lang.X("item.form.model.text", "item.form.model.text")),
			modelurl:          widget.NewLabel(lang.X("item.form.modelurl.text", "item.form.modelurl.text")),
			notes:             widget.NewLabel(lang.X("item.form.notes.text", "item.form.notes.text")),
			dimensions:        widget.NewLabel(lang.X("item.form.dimensions.text", "item.form.dimensions.text")),
			width:             widget.NewLabel(lang.X("item.form.width.text", "item.form.width.text")),
			height:            widget.NewLabel(lang.X("item.form.height.text", "item.form.height.text")),
			depth:             widget.NewLabel(lang.X("item.form.depth.text", "item.form.depth.text")),
			volume:            widget.NewLabel(lang.X("item.form.volume.text", "item.form.volume.text")),
			weight:            widget.NewLabel(lang.X("item.form.weight.text", "item.form.weight.text")),
			status:            widget.NewLabel(lang.X("item.form.status.text", "item.form.status.text")),
			datecreatedtitle:  widget.NewLabel(lang.X("item.form.datecreated.text", "item.form.datecreated.text")),
			datemodifiedtitle: widget.NewLabel(lang.X("item.form.datemodified.text", "item.form.datemodified.text")),
			datecreated:       widget.NewLabel(time.DateTime),
			datemodified:      widget.NewLabel(time.DateTime),
		},
		selects: &formSelects{
			category:   widget.NewSelect(categories, func(s string) {}),
			lengthunit: widget.NewSelect(lengthUnits, func(s string) {}),
			volumeunit: widget.NewSelect(volumeUnits, func(s string) {}),
			weightunit: widget.NewSelect(weightUnits, func(s string) {}),
			status:     widget.NewSelect(itemStatus, func(s string) {}),
		},
	}
	b.Metadata.CatIDList.AddListener(binding.NewDataListener(func() {
		categories = fetchCategories()
		v.selects.category.Options = categories
		v.selects.category.Refresh()
	}))
	// v.entries.adddesc.MultiLine = true
	// v.entries.adddesc.SetMinRowsVisible(2)
	// v.entries.adddesc.Wrapping = fyne.TextWrapWord
	v.entries.longdesc.MultiLine = true
	v.entries.longdesc.SetMinRowsVisible(5)
	v.entries.longdesc.Wrapping = fyne.TextWrapWord
	v.entries.notes.MultiLine = true
	v.entries.notes.SetMinRowsVisible(5)
	v.entries.notes.Wrapping = fyne.TextWrapWord
	spacebox := container.NewGridWithRows(1,
		container.NewBorder(nil, nil, v.labels.width, nil, v.entries.width),
		container.NewBorder(nil, nil, v.labels.height, nil, v.entries.height),
		container.NewBorder(nil, nil, v.labels.depth, nil, v.entries.depth),
		v.selects.lengthunit,
	)
	massbox := container.NewGridWithRows(1,
		container.NewBorder(nil, nil, v.labels.volume, v.selects.volumeunit, v.entries.volume),
		container.NewBorder(nil, nil, v.labels.weight, v.selects.weightunit, v.entries.weight),
	)
	v.container = container.New(layout.NewFormLayout(),
		layout.NewSpacer(), container.NewHBox(v.labels.datecreatedtitle, v.labels.datecreated),
		layout.NewSpacer(), container.NewHBox(v.labels.datemodifiedtitle, v.labels.datemodified),
		v.labels.itemidtitle, v.labels.itemid,
		layout.NewSpacer(), v.selects.status,
		v.labels.name, v.entries.name,
		v.labels.category, v.selects.category,
		v.labels.manufacturer, v.entries.manufacturer,
		v.labels.model, v.entries.model,
		v.labels.modelurl, v.entries.modelurl,
		// v.labels.longdesc, v.entries.longdesc,
		v.labels.dimensions, spacebox,
		layout.NewSpacer(), massbox,
		v.labels.price, container.NewBorder(nil, nil, nil, midget.NewLabel("SEK", "", ""), v.entries.price),
		// v.labels.adddesc, v.entries.adddesc,
		v.labels.imgurl1, v.entries.imgurl1,
		v.labels.notes, v.entries.notes,
		layout.NewSpacer(), widget.NewLabel(" "),
		layout.NewSpacer(), widget.NewRichTextFromMarkdown(`### `+lang.L("Preview")),
		v.labels.longdesc, layout.NewSpacer(),
		v.labels.adddesctitle, v.labels.adddesc,
	)
	v.Clear()
	return v
}

func (v formView) Bind(id backend.ItemID) {
	v.Clear()
	v.Enable()

	v.labels.itemid.Bind(id.Item().ItemIDString)
	v.labels.status.Bind(id.Item().ItemStatus)
	v.labels.datecreated.Bind(id.Item().DateCreated)
	v.labels.datemodified.Bind(id.Item().DateModified)
	v.labels.adddesc.Bind(id.Item().AddDesc)

	v.entries.name.Bind(id.Item().Name)
	v.entries.price.Bind(id.Item().PriceString)
	v.entries.vat.Bind(id.Item().VatString)
	v.entries.imgurl1.Bind(id.Item().ImgURL1)
	v.entries.specsurl.Bind(id.Item().SpecsURL)
	// v.entries.adddesc.Bind(id.Item().AddDesc)
	v.entries.longdesc.Bind(id.Item().LongDesc)
	v.entries.manufacturer.Bind(id.Item().Manufacturer)
	v.entries.model.Bind(id.Item().Model)
	v.entries.modelurl.Bind(id.Item().ModelURL)
	v.entries.notes.Bind(id.Item().Notes)
	v.entries.width.Bind(id.Item().WidthString)
	v.entries.height.Bind(id.Item().HeightString)
	v.entries.depth.Bind(id.Item().DepthString)
	v.entries.volume.Bind(id.Item().VolumeString)
	v.entries.weight.Bind(id.Item().WeightString)

	v.selects.category.Bind(id.Item().Category)
	v.selects.lengthunit.Bind(id.Item().LengthUnit)
	v.selects.volumeunit.Bind(id.Item().VolumeUnit)
	v.selects.weightunit.Bind(id.Item().WeightUnit)
	v.selects.status.Bind(id.Item().ItemStatus)
}
func (v formView) Clear() {
	v.entries.name.Unbind()
	v.entries.price.Unbind()
	v.entries.vat.Unbind()
	v.entries.imgurl1.Unbind()
	v.entries.specsurl.Unbind()
	// v.entries.adddesc.Unbind()
	v.entries.longdesc.Unbind()
	v.entries.manufacturer.Unbind()
	v.entries.model.Unbind()
	v.entries.modelurl.Unbind()
	v.entries.notes.Unbind()
	v.entries.width.Unbind()
	v.entries.height.Unbind()
	v.entries.depth.Unbind()
	v.entries.volume.Unbind()
	v.entries.weight.Unbind()

	v.labels.itemid.Unbind()
	v.labels.status.Unbind()
	v.labels.datecreated.Unbind()
	v.labels.datemodified.Unbind()
	v.labels.adddesc.Unbind()

	v.selects.category.Unbind()
	v.selects.lengthunit.Unbind()
	v.selects.volumeunit.Unbind()
	v.selects.weightunit.Unbind()
	v.selects.status.Unbind()

	v.entries.name.SetText("")
	v.entries.price.SetText("")
	v.entries.vat.SetText("")
	v.entries.imgurl1.SetText("")
	v.entries.specsurl.SetText("")
	// v.entries.adddesc.SetText("")
	v.entries.longdesc.SetText("")
	v.entries.manufacturer.SetText("")
	v.entries.model.SetText("")
	v.entries.modelurl.SetText("")
	v.entries.notes.SetText("")
	v.entries.volume.SetText("")
	v.entries.width.SetText("")
	v.entries.height.SetText("")
	v.entries.depth.SetText("")
	v.entries.volume.SetText("")
	v.entries.weight.SetText("")

	v.labels.itemid.SetText("0000000000")
	v.labels.status.SetText("")
	v.labels.datecreated.SetText("")
	v.labels.datemodified.SetText("")
	v.labels.adddesc.SetText("")

	v.selects.category.ClearSelected()
	v.selects.lengthunit.ClearSelected()
	v.selects.volumeunit.ClearSelected()
	v.selects.weightunit.ClearSelected()
	v.selects.status.ClearSelected()

	v.Disable()
}

func (v formView) Disable() {
	v.entries.name.Disable()
	v.entries.price.Disable()
	v.entries.vat.Disable()
	v.entries.imgurl1.Disable()
	v.entries.specsurl.Disable()
	// v.entries.adddesc.Disable()
	v.entries.longdesc.Disable()
	v.entries.manufacturer.Disable()
	v.entries.model.Disable()
	v.entries.modelurl.Disable()
	v.entries.notes.Disable()
	v.entries.width.Disable()
	v.entries.height.Disable()
	v.entries.depth.Disable()
	v.entries.volume.Disable()
	v.entries.weight.Disable()

	v.selects.category.Disable()
	v.selects.lengthunit.Disable()
	v.selects.volumeunit.Disable()
	v.selects.weightunit.Disable()
	v.selects.status.Disable()
}
func (v formView) Enable() {
	v.entries.name.Enable()
	v.entries.price.Enable()
	v.entries.vat.Enable()
	v.entries.imgurl1.Enable()
	v.entries.specsurl.Enable()
	// v.entries.adddesc.Enable()
	v.entries.longdesc.Enable()
	v.entries.manufacturer.Enable()
	v.entries.model.Enable()
	v.entries.modelurl.Enable()
	v.entries.notes.Enable()
	v.entries.width.Enable()
	v.entries.height.Enable()
	v.entries.depth.Enable()
	v.entries.volume.Enable()
	v.entries.weight.Enable()

	v.selects.category.Enable()
	v.selects.lengthunit.Enable()
	v.selects.volumeunit.Enable()
	v.selects.weightunit.Enable()
	v.selects.status.Enable()
}
