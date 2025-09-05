package gui

import (
	"UppSpar/backend"
	"UppSpar/backend/bridge"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

// TODO redo with maps just like metadata

type items struct {
	container *fyne.Container
	form      *bridge.Form
	list      *bridge.List
	search    *bridge.Tools
}

func newItems(a *App) *items {
	w := a.window
	b := a.backend

	v := &items{
		form:   bridge.NewItemForm(b, w),
		list:   bridge.NewList(b, w),
		search: bridge.NewSearchBar(b, w),
	}

	b.Items.ItemIDSelection.AddListener(binding.NewDataListener(func() {
		ids, err := b.Items.ItemIDSelection.Get()
		if err != nil {
			panic(err)
		}
		if len(ids) < 1 {
			return
		}
		ItemID := ids[0].(backend.ItemID)
		v.form.LoadItem(b, ItemID)
	}))

	split := container.NewHSplit(v.list.Container, v.form.Container)
	split.SetOffset(0.2)

	v.container = container.NewBorder(v.search.Container, nil, nil, nil, split)

	// searchEntry := xwidget.NewCompletionEntry([]string{})
	// searchEntry.Bind(b.Items.Search.Term)
	// searchEntry.OnChanged = func(s string) {
	// 	if len(s) < 3 {
	// 		searchEntry.HideCompletion()
	// 		return
	// 	}
	// 	results, err := b.Items.Search.Completions.Get()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	searchEntry.SetOptions(results)
	// 	searchEntry.ShowCompletion()
	// }
	// searchKeys := make(bridge.Checks)
	// searchKeys["Name"] = ttw.NewCheckWithData("Art", b.Items.Search.Scope["Name"])
	// searchKeys["Manufacturer"] = ttw.NewCheckWithData("Mfr", b.Items.Search.Scope["Manufacturer"])
	// searchKeys["ModelName"] = ttw.NewCheckWithData("Mdl", b.Items.Search.Scope["ModelName"])
	// searchKeys["ModelDesc"] = ttw.NewCheckWithData("Dsc", b.Items.Search.Scope["ModelDesc"])
	// toolbarLeft := container.NewGridWithRows(1,
	// 	// layout.NewSpacer(),
	// 	container.NewHBox(
	// 		searchKeys["Name"],
	// 		searchKeys["Manufacturer"],
	// 		searchKeys["ModelName"],
	// 		searchKeys["ModelDesc"],
	// 	),
	// 	widget.NewSelect([]string{
	// 		lang.X("form.select.search.beginswith", "form.select.search.beginswith"),
	// 		lang.X("form.select.search.endswith", "form.select.search.endswith"),
	// 		lang.X("form.select.search.contains", "form.select.search.contains"),
	// 		lang.X("form.select.search.equals", "form.select.search.equals"),
	// 		// lang.X("form.select.search.regexp", "form.select.search.regexp"),
	// 	}, func(s string) {
	// 		switch s {
	// 		case lang.X("form.select.search.beginswith", "form.select.search.beginswith"):
	// 			b.Items.Search.Match = backend.MatchBeginsWith
	// 		case lang.X("form.select.search.endswith", "form.select.search.endswith"):
	// 			b.Items.Search.Match = backend.MatchEndsWith
	// 		case lang.X("form.select.search.equals", "form.select.search.equals"):
	// 			b.Items.Search.Match = backend.MatchEquals
	// 		case lang.X("form.select.search.contains", "form.select.search.contains"):
	// 			b.Items.Search.Match = backend.MatchContains
	// 		default:
	// 			b.Items.Search.Match = backend.MatchContains
	// 		}
	// 		b.Items.GetItemIDs()
	// 	}),
	// )
	// toolbarRight := container.NewGridWithRows(1,
	// 	midget.NewLabel(
	// 		lang.X("form.select.sortby.text", "form.select.sortby.text"),
	// 		"",
	// 		"",
	// 	),
	// 	widget.NewSelect([]string{
	// 		lang.X("form.select.sortby.itemid", "form.select.sortby.itemid"),
	// 		lang.X("form.select.sortby.name", "form.select.sortby.name"),
	// 		lang.X("form.select.sortby.manufacturer", "form.select.sortby.manufacturer"),
	// 		lang.X("form.select.sortby.datecreated", "form.select.sortby.datecreated"),
	// 		lang.X("form.select.sortby.datemodified", "form.select.sortby.datemodified"),
	// 	}, func(s string) {
	// 		switch s {
	// 		case lang.X("form.select.sortby.itemid", "form.select.sortby.itemid"):
	// 			b.Items.Search.SortBy = backend.SearchKeyItemID
	// 		case lang.X("form.select.sortby.name", "form.select.sortby.name"):
	// 			b.Items.Search.SortBy = backend.SearchKeyName
	// 		case lang.X("form.select.sortby.manufacturer", "form.select.sortby.manufacturer"):
	// 			b.Items.Search.SortBy = backend.SearchKeyManufacturer
	// 		case lang.X("form.select.sortby.datecreated", "form.select.sortby.datecreated"):
	// 			b.Items.Search.SortBy = backend.SearchKeyDateCreated
	// 		case lang.X("form.select.sortby.datemodified", "form.select.sortby.datemodified"):
	// 			b.Items.Search.SortBy = backend.SearchKeyDateModified
	// 		}
	// 		b.Items.GetItemIDs()
	// 	}),
	// 	widget.NewSelect([]string{
	// 		lang.X("form.select.sortorder.ascending", "form.select.sortorder.ascending"),
	// 		lang.X("form.select.sortorder.descending", "form.select.sortorder.descending"),
	// 	}, func(s string) {
	// 		if s == lang.X("form.select.sortorder.ascending", "form.select.sortorder.ascending") {
	// 			b.Items.Search.Order = backend.SortAscending
	// 		} else {
	// 			b.Items.Search.Order = backend.SortDescending
	// 		}
	// 		b.Items.GetItemIDs()
	// 	}))
	// filterWidth := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Width")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Width))
	// filterHeight := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Height")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Height))
	// filterDepth := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Depth")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Depth))
	// filterVolume := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Volume")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Volume))
	// filterWeight := container.NewBorder(nil, nil, widget.NewLabel(lang.L("Weight")), nil, widget.NewEntryWithData(a.backend.Items.Filter.Weight))
	// filterBar := container.NewGridWithRows(1,
	// 	filterWidth, filterHeight, filterDepth, filterVolume, filterWeight,
	// )
	// searchBar := container.NewBorder(nil, nil, toolbarLeft, toolbarRight, searchEntry)
	// v.search = container.NewBorder(searchBar, filterBar, nil, nil)
	// toolbarLeft.Objects[1].(*widget.Select).SetSelectedIndex(2)  // search type
	// toolbarRight.Objects[1].(*widget.Select).SetSelectedIndex(0) // sort by
	// toolbarRight.Objects[2].(*widget.Select).SetSelectedIndex(0) // sort order
	// imgView := container.NewBorder(nil, nil, nil, nil)
	// listView := container.NewBorder(nil, nil, nil, nil, v.list)
	// formView := container.NewBorder(nil, nil, nil, imgView, v.formView.container)
	// split := container.NewHSplit(listView, formView)
	// split.SetOffset(0.25)
	// toolbar := container.NewBorder(nil, nil, v.listTools, nil, v.search)
	// statusbar := container.NewGridWithRows(1)
	// v.container = container.NewBorder(toolbar, statusbar, nil, nil, split)

	return v
}

// func newFormView(b *backend.Backend) *formObjects {
// 	var categories []string
// 	fetchCategories := func() []string {
// 		cats, _ := b.Metadata.Categories.Get()
// 		return cats
// 	}
// 	categories = fetchCategories()

// 	lengthUnits := []string{"mm", "cm", "dm", "m"}
// 	volumeUnits := []string{"ml", "cl", "dl", "l"}
// 	weightUnits := []string{"g", "hg", "kg"}

// 	itemStatus := func() []string {
// 		var ss []string
// 		stats, _ := b.Metadata.ItemStatusIDList.Get()
// 		for _, stat := range stats {
// 			ss = append(ss, stat.(backend.ItemStatusID).LString())
// 		}
// 		return ss
// 	}()

// 	v := &formObjects{
// 		entries: &formEntries{
// 			Name:         midget.NewEntry(),
// 			Price:        midget.NewEntry(),
// 			Vat:          midget.NewEntry(),
// 			ImgURL1:      midget.NewEntry(),
// 			ImgURL2:      midget.NewEntry(),
// 			ImgURL3:      midget.NewEntry(),
// 			ImgURL4:      midget.NewEntry(),
// 			ImgURL5:      midget.NewEntry(),
// 			SpecsURL:     midget.NewEntry(),
// 			LongDesc:     midget.NewEntry(),
// 			Manufacturer: widget.NewSelectEntry([]string{}),
// 			ModelName:    widget.NewSelectEntry([]string{}),
// 			ModelDesc:    midget.NewEntry(),
// 			ModelURL:     midget.NewEntry(),
// 			Notes:        midget.NewEntry(),
// 			Width:        midget.NewEntry(),
// 			Height:       midget.NewEntry(),
// 			Depth:        midget.NewEntry(),
// 			Volume:       midget.NewEntry(),
// 			Weight:       midget.NewEntry(),
// 		},
// 		labels: &formLabels{
// 			ItemID:       widget.NewLabel(lang.X("item.form.label.itemid", "item.form.label.itemid")),
// 			Name:         widget.NewLabel(lang.X("item.form.label.name", "item.form.label.name")),
// 			Category:     widget.NewLabel(lang.X("item.form.label.category", "item.form.label.category")),
// 			Currency:     widget.NewLabel("SEK"),
// 			Price:        widget.NewLabel(lang.X("item.form.label.price", "item.form.label.price")),
// 			Vat:          widget.NewLabel(lang.X("item.form.label.vat", "item.form.label.vat")),
// 			ImgURL1:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
// 			ImgURL2:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
// 			ImgURL3:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
// 			ImgURL4:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
// 			ImgURL5:      widget.NewLabel(lang.X("item.form.label.imgurl", "item.form.label.imgurl")),
// 			SpecsURL:     widget.NewLabel(lang.X("item.form.label.specsurl", "item.form.label.specsurl")),
// 			AddDesc:      widget.NewLabel(lang.X("item.form.label.adddesc", "item.form.label.adddesc")),
// 			LongDesc:     widget.NewLabel(lang.X("item.form.label.longdesc", "item.form.label.longdesc")),
// 			Manufacturer: widget.NewLabel(lang.X("item.form.label.manufacturer", "item.form.label.manufacturer")),
// 			ModelName:    widget.NewLabel(lang.X("item.form.label.modelname", "item.form.label.modelname")),
// 			ModelDesc:    widget.NewLabel(lang.X("item.form.label.modeldesc", "item.form.label.modeldesc")),
// 			ModelURL:     widget.NewLabel(lang.X("item.form.label.modelurl", "item.form.label.modelurl")),
// 			Notes:        widget.NewLabel(lang.X("item.form.label.notes", "item.form.label.notes")),
// 			Dimensions:   widget.NewLabel(lang.X("item.form.label.dimensions", "item.form.label.dimensions")),
// 			Width:        widget.NewLabel(lang.X("item.form.label.width", "item.form.label.width")),
// 			Height:       widget.NewLabel(lang.X("item.form.label.height", "item.form.label.height")),
// 			Depth:        widget.NewLabel(lang.X("item.form.label.depth", "item.form.label.depth")),
// 			Volume:       widget.NewLabel(lang.X("item.form.label.volume", "item.form.label.volume")),
// 			Weight:       widget.NewLabel(lang.X("item.form.label.weight", "item.form.label.weight")),
// 			Status:       widget.NewLabel(lang.X("item.form.label.status", "item.form.label.status")),
// 			DateCreated:  widget.NewLabel(lang.X("item.form.label.datecreated", "item.form.label.datecreated")),
// 			DateModified: widget.NewLabel(lang.X("item.form.label.datemodified", "item.form.label.datemodified")),
// 		},
// 		values: &formDataLabels{
// 			ItemID:       widget.NewLabel(lang.X("item.form.data.itemid", "item.form.data.itemid")),
// 			AddDesc:      widget.NewLabel(lang.X("item.form.data.adddesc", "item.form.data.adddesc")),
// 			LongDesc:     widget.NewLabel(lang.X("item.form.data.longdesc", "item.form.data.longdesc")),
// 			DateCreated:  widget.NewLabel(lang.X("item.form.data.datecreated", "item.form.data.datecreated")),
// 			DateModified: widget.NewLabel(lang.X("item.form.data.datemodified", "item.form.data.datemodified")),
// 		},
// 		selects: &formSelects{
// 			Category:   widget.NewSelect(categories, func(s string) {}),
// 			LengthUnit: widget.NewSelect(lengthUnits, func(s string) {}),
// 			VolumeUnit: widget.NewSelect(volumeUnits, func(s string) {}),
// 			WeightUnit: widget.NewSelect(weightUnits, func(s string) {}),
// 			Status:     widget.NewSelect(itemStatus, func(s string) {}),
// 		},
// 	}
// 	b.Metadata.CatIDList.AddListener(binding.NewDataListener(func() {
// 		categories = fetchCategories()
// 		v.selects.Category.Options = categories
// 		v.selects.Category.Refresh()
// 	}))
// 	v.entries.ModelDesc.MultiLine = true
// 	v.entries.ModelDesc.SetMinRowsVisible(5)
// 	v.entries.ModelDesc.Wrapping = fyne.TextWrapWord
// 	v.entries.LongDesc.MultiLine = true
// 	v.entries.LongDesc.SetMinRowsVisible(5)
// 	v.entries.LongDesc.Wrapping = fyne.TextWrapWord
// 	v.entries.Notes.MultiLine = true
// 	v.entries.Notes.SetMinRowsVisible(5)
// 	v.entries.Notes.Wrapping = fyne.TextWrapWord

// 	v.values.LongDesc.Wrapping = fyne.TextWrapWord
// 	v.values.AddDesc.Wrapping = fyne.TextWrapWord

// 	b.Metadata.MfrNameList.AddListener(binding.NewDataListener(func() {
// 		manufacturers, _ := b.Metadata.MfrNameList.Get()
// 		v.entries.Manufacturer.SetOptions(manufacturers)
// 	}))

// 	idbox := container.NewBorder(nil, nil, v.values.ItemID, nil, container.NewHBox(v.selects.Status))
// 	spacebox := container.NewGridWithRows(1,
// 		container.NewBorder(nil, nil, v.labels.Width, nil, v.entries.Width),
// 		container.NewBorder(nil, nil, v.labels.Height, nil, v.entries.Height),
// 		container.NewBorder(nil, nil, v.labels.Depth, nil, v.entries.Depth),
// 		v.selects.LengthUnit,
// 	)
// 	massbox := container.NewGridWithRows(1,
// 		container.NewBorder(nil, nil, v.labels.Volume, v.selects.VolumeUnit, v.entries.Volume),
// 		container.NewBorder(nil, nil, v.labels.Weight, v.selects.WeightUnit, v.entries.Weight),
// 	)
// 	v.container = container.NewVScroll(container.New(layout.NewFormLayout(),
// 		layout.NewSpacer(), container.NewHBox(v.labels.DateCreated, v.values.DateCreated),
// 		layout.NewSpacer(), container.NewHBox(v.labels.DateModified, v.values.DateModified),
// 		v.labels.ItemID, idbox,
// 		v.labels.Name, v.entries.Name,
// 		v.labels.Category, v.selects.Category,
// 		v.labels.Manufacturer, v.entries.Manufacturer,
// 		v.labels.ModelName, v.entries.ModelName,
// 		v.labels.ModelDesc, v.entries.ModelDesc,
// 		v.labels.ModelURL, v.entries.ModelURL,
// 		v.labels.Dimensions, spacebox,
// 		layout.NewSpacer(), massbox,
// 		v.labels.Price, container.NewBorder(nil, nil, nil, v.labels.Currency, v.entries.Price),
// 		v.labels.ImgURL1, v.entries.ImgURL1,
// 		v.labels.Notes, v.entries.Notes,
// 		layout.NewSpacer(), widget.NewLabel(" "),
// 		layout.NewSpacer(), widget.NewRichTextFromMarkdown(`### `+lang.L("Preview")),
// 		v.labels.LongDesc, v.values.LongDesc,
// 		v.labels.AddDesc, v.values.AddDesc,
// 	))
// 	v.Clear()
// 	return v
// }

// func (v formObjects) Bind(b *backend.Backend, id backend.ItemID) {
// 	v.Clear()

// 	v.values.ItemID.Bind(id.Item().ItemIDString)
// 	v.values.DateCreated.Bind(id.Item().DateCreated)
// 	v.values.DateModified.Bind(id.Item().DateModified)
// 	v.values.AddDesc.Bind(id.Item().AddDesc)
// 	v.values.LongDesc.Bind(id.Item().LongDesc)

// 	v.entries.Name.Bind(id.Item().Name)
// 	v.entries.Price.Bind(id.Item().PriceString)
// 	v.entries.Vat.Bind(id.Item().VatString)
// 	v.entries.ImgURL1.Bind(id.Item().ImgURL1)
// 	v.entries.SpecsURL.Bind(id.Item().SpecsURL)
// 	v.entries.Manufacturer.Bind(id.Item().Manufacturer)
// 	v.entries.ModelName.Bind(id.Item().ModelName)
// 	v.entries.ModelDesc.Bind(id.Item().ModelDesc)
// 	v.entries.ModelURL.Bind(id.Item().ModelURL)
// 	v.entries.Notes.Bind(id.Item().Notes)
// 	v.entries.Width.Bind(id.Item().WidthString)
// 	v.entries.Height.Bind(id.Item().HeightString)
// 	v.entries.Depth.Bind(id.Item().DepthString)
// 	v.entries.Volume.Bind(id.Item().VolumeString)
// 	v.entries.Weight.Bind(id.Item().WeightString)

// 	v.selects.Category.Bind(id.Item().Category)
// 	v.selects.LengthUnit.Bind(id.Item().LengthUnit)
// 	v.selects.VolumeUnit.Bind(id.Item().VolumeUnit)
// 	v.selects.WeightUnit.Bind(id.Item().WeightUnit)
// 	v.selects.Status.Bind(id.Item().ItemStatus)

// 	id.Item().Manufacturer.AddListener(binding.NewDataListener(func() {
// 		models := func() []string {
// 			b.Metadata.GetModelIDs(id.Item().MfrID)
// 			var names []string
// 			ids := id.Item().MfrID.Children()
// 			for _, id := range ids {
// 				name, _ := id.Name()
// 				names = append(names, name)
// 			}
// 			return names
// 		}()
// 		v.entries.ModelName.SetOptions(models)
// 	}))

// 	/* This step is needed because child categories have spaces prepended to them in the select list */
// 	cat, _ := id.Item().Category.Get()
// 	v.selects.Category.SetSelectedIndex(b.Metadata.GetListItemIDForCategory(cat))

// 	id.Item().CatID.Category().Config["ShowPrice"].AddListener(binding.NewDataListener(func() {
// 		if cid, _ := id.CatID(); cid == backend.CatID(1) {
// 			// TODO fixa asap
// 			v.hideStatus()
// 			// v.hideCategory()
// 			v.hideImgURL()
// 			v.hideSpecsURL()
// 			v.hideMfrModel()
// 			v.hideLength()
// 			v.hideVolume()
// 			v.hideWeight()
// 			v.showPrice()
// 			v.hidePreviewAddDesc()
// 			return
// 		} else {
// 			v.showStatus()
// 			v.showCategory()
// 			v.showImgURL()
// 			v.showSpecsURL()
// 			v.showMfrModel()
// 			v.showLength()
// 			v.showWeight()
// 			v.showVolume()
// 			v.showPreviewAddDesc()
// 		}
// 		p, _ := id.Item().CatID.Category().Config["ShowPrice"].Get()
// 		if p {
// 			v.showPrice()
// 		} else {
// 			v.hidePrice()
// 		}
// 	}))
// 	id.Item().CatID.Category().Config["ShowLength"].AddListener(binding.NewDataListener(func() {
// 		p, _ := id.Item().CatID.Category().Config["ShowLength"].Get()
// 		if p {
// 			v.showLength()
// 		} else {
// 			v.hideLength()
// 		}
// 	}))
// 	id.Item().CatID.Category().Config["ShowVolume"].AddListener(binding.NewDataListener(func() {
// 		p, _ := id.Item().CatID.Category().Config["ShowVolume"].Get()
// 		if p {
// 			v.showVolume()
// 		} else {
// 			v.hideVolume()
// 		}
// 	}))
// 	id.Item().CatID.Category().Config["ShowWeight"].AddListener(binding.NewDataListener(func() {
// 		p, _ := id.Item().CatID.Category().Config["ShowWeight"].Get()
// 		if p {
// 			v.showWeight()
// 		} else {
// 			v.hideWeight()
// 		}
// 	}))

// 	v.Enable()
// }
