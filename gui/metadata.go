package gui

import (
	"UppSpar/backend"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	xwidget "fyne.io/x/fyne/widget"
	midget "github.com/assholehoff/fyne-midget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type metadataView struct {
	category     *categoryView
	manufacturer *manufacturerView
	model        *modelView
	tabs         *container.AppTabs
}

func newMetadataView(b *backend.Backend) *metadataView {
	categoryView := newCategoryView(b)
	manufacturerView := newManufacturerView(b)
	modelView := newModelView(b)

	return &metadataView{
		category:     categoryView,
		manufacturer: manufacturerView,
		model:        modelView,
		tabs:         newMetadataTabs(categoryView, manufacturerView, modelView),
	}
}

func newMetadataTabs(c *categoryView, mfr *manufacturerView, mdl *modelView) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem(lang.L("Manufacturers"), mfr.container),
		container.NewTabItem(lang.L("Models"), mdl.container),
		container.NewTabItem(lang.L("Categories"), c.container),
	)
	return tabs
}

type manufacturerView struct {
	container *container.Split
}

func newManufacturerView(b *backend.Backend) *manufacturerView {
	createItem := func() fyne.CanvasObject {
		l := widget.NewLabel("Template category name")
		return l
	}
	updateItem := func(di binding.DataItem, co fyne.CanvasObject) {
		v, err := di.(binding.Untyped).Get()
		if err != nil {
			log.Println(err)
			return
		}
		MfrID := v.(backend.MfrID)
		co.(*widget.Label).Bind(MfrID.Manufacturer().Name)
	}
	l := widget.NewListWithData(b.Metadata.MfrIDList, createItem, updateItem)
	f := container.New(layout.NewFormLayout())
	s := container.NewHSplit(l, f)
	s.SetOffset(0.25)

	// TODO convert to tree structure
	createTreeItem := func(branch bool) fyne.CanvasObject {
		if branch {
		}
		return widget.NewLabel("")
	}
	updateTreeItem := func(di binding.DataItem, branch bool, co fyne.CanvasObject) {}
	tree := binding.NewUntypedTree()
	t := widget.NewTreeWithData(tree, createTreeItem, updateTreeItem)
	t.CloseAllBranches()

	return &manufacturerView{container: s}
}

type modelView struct {
	container *container.Split
}

func newModelView(b *backend.Backend) *modelView {
	createItem := func(branch bool) fyne.CanvasObject {
		// if branch {
		// 	return widget.NewLabel("Template product name")
		// }
		return widget.NewLabel("Template category name")
	}
	updateItem := func(di binding.DataItem, branch bool, co fyne.CanvasObject) {
		co.(*widget.Label).Bind(di.(binding.String))
	}
	t := widget.NewTreeWithData(b.Metadata.ProductTree, createItem, updateItem)
	f := container.New(layout.NewFormLayout())
	s := container.NewHSplit(t, f)
	s.SetOffset(0.25)
	return &modelView{container: s}
}

type categoryChecks struct {
	price  *ttw.Check
	length *ttw.Check
	volume *ttw.Check
	weight *ttw.Check
}

type categoryEntries struct {
	name *xwidget.CompletionEntry
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
			name: xwidget.NewCompletionEntry([]string{}),
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
	m.entries.name.Bind(id.Category().Name)
	m.checks.price.Bind(id.Category().Config["ShowPrice"])
	m.checks.length.Bind(id.Category().Config["ShowLength"])
	m.checks.volume.Bind(id.Category().Config["ShowVolume"])
	m.checks.weight.Bind(id.Category().Config["ShowWeight"])
	p, _ := id.ParentID()
	if p == 0 {
		m.selects.parent.SetSelected(lang.L("None"))
	}
	id.Category().Name.AddListener(binding.NewDataListener(func() { b.Metadata.UpdateCatList() }))
}
