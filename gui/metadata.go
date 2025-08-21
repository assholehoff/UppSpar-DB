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

type metadataChecks struct {
	price *ttw.Check
}

type metadataEntries struct {
	name *xwidget.CompletionEntry
}

type metadataLabels struct {
	name *midget.Label
}

type metadataSelects struct {
}

type metadataFormView struct {
	container *fyne.Container
	checks    *struct{}
	entries   *metadataEntries
	labels    *metadataLabels
	selects   *struct{}
}

type metadataView struct {
	container *container.Split
	form      *metadataFormView
	list      *widget.List
	toolbar   *widget.Toolbar
}

func newMetadataView(b *backend.Backend) *metadataView {
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
		CatID := v.(backend.CatID)
		co.(*widget.Label).Bind(CatID.Category().Name)
	}
	m := &metadataView{
		list: widget.NewListWithData(b.Metadata.CatIDList, createItem, updateItem),
	}
	m.form = &metadataFormView{
		entries: &metadataEntries{
			name: widget.NewEntry(),
		},
		labels: &metadataLabels{
			name: midget.NewLabel(lang.X("metadata.form.name", "metadata.form.name"), "", ""),
		},
	}
	m.form.container = container.New(layout.NewFormLayout(),
		widget.NewRichTextFromMarkdown(
			`## `+lang.X("metadata.subtitle.categories", "metadata.subtitle.categories"),
		), layout.NewSpacer(),
		m.form.labels.name, m.form.entries.name,
	)
	m.list.OnSelected = func(id widget.ListItemID) {
		b.Metadata.SelectCategory(b.Metadata.GetCatIDFor(id))
		m.form.Bind(b, b.Metadata.GetCatIDFor(id))
	}
	m.list.OnUnselected = func(id widget.ListItemID) {
		b.Metadata.UnselectCategory(b.Metadata.GetCatIDFor(id))
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
	listView := container.NewBorder(m.toolbar, nil, nil, nil, m.list)
	m.container = container.NewHSplit(listView, m.form.container)
	m.container.SetOffset(0.25)
	return m
}

func (m *metadataFormView) Bind(b *backend.Backend, id backend.CatID) {
	m.entries.name.Bind(id.Category().Name)
	id.Category().Name.AddListener(binding.NewDataListener(func() { b.Metadata.UpdateCatList() }))
}
