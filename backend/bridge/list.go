package bridge

import (
	"UppSpar/backend"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type List struct {
	Container *fyne.Container
	list      *widget.List
	tree      *widget.Tree
	toolbar   *widget.Toolbar
}

func NewList(b *backend.Backend, w fyne.Window) *List {
	var list *widget.List
	var tree *widget.Tree
	var toolbar *widget.Toolbar

	list = widget.NewListWithData(
		b.Items.ItemIDList,
		func() fyne.CanvasObject {
			return container.NewVBox(
				container.NewHBox(
					&widget.Label{
						Text:     "00000000",
						SizeName: theme.SizeNameCaptionText,
					},
					&widget.Label{
						Text:     "Category template text",
						SizeName: theme.SizeNameCaptionText,
					},
				),
				&widget.Label{
					Text:      "Item name template",
					TextStyle: fyne.TextStyle{Bold: true},
				},
			)
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			val, err := di.(binding.Untyped).Get()
			if err != nil {
				panic(err)
			}
			ItemID := val.(backend.ItemID)
			co.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).Bind(ItemID.Item().ItemIDString)
			co.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).Bind(ItemID.Item().Category)
			co.(*fyne.Container).Objects[1].(*widget.Label).Bind(ItemID.Item().Name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		ItemID, err := b.Items.GetItemIDFor(id)
		if err != nil {
			panic(err)
		}
		b.Items.SelectItem(ItemID)
	}
	list.OnUnselected = func(id widget.ListItemID) {
		ItemID, err := b.Items.GetItemIDFor(id)
		if err != nil {
			panic(err)
		}
		b.Items.UnselectItem(ItemID)
	}

	tree = widget.NewTree(
		func(tni widget.TreeNodeID) []widget.TreeNodeID {
			return []widget.TreeNodeID{}
		},
		func(tni widget.TreeNodeID) bool {
			return false
		},
		func(b bool) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(tni widget.TreeNodeID, b bool, co fyne.CanvasObject) {},
	)

	toolbar = widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			toolbar.Items[0].(*widget.ToolbarAction).Disable()
			go func() {
				id, _ := b.Items.CreateNewItem()
				index, _ := b.Items.GetListItemIDFor(id)
				fyne.Do(func() {
					list.Select(index)
					time.Sleep(100 * time.Millisecond)
					toolbar.Items[0].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
			toolbar.Items[1].(*widget.ToolbarAction).Disable()
			go func() {
				items, err := b.Items.ItemIDSelection.Get()
				if err != nil {
					panic(err)
				}
				fyne.Do(func() {
					list.UnselectAll()
					for _, item := range items {
						b.Items.DeleteItem(item.(backend.ItemID))
					}
					time.Sleep(100 * time.Millisecond)
					toolbar.Items[1].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			toolbar.Items[2].(*widget.ToolbarAction).Disable()
			go func() {
				items, err := b.Items.ItemIDSelection.Get()
				if err != nil {
					panic(err)
				}
				fyne.Do(func() {
					list.UnselectAll()
					for _, item := range items {
						id, err := b.Items.CopyItem(item.(backend.ItemID))
						if err != nil {
							panic(err)
						}
						index, err := b.Items.GetListItemIDFor(id)
						if err != nil {
							panic(err)
						}
						list.Select(index)
					}
					time.Sleep(100 * time.Millisecond)
					toolbar.Items[2].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			toolbar.Items[3].(*widget.ToolbarAction).Disable()
			go func() {
				fyne.Do(func() {
					NewSaveFileDialog(b, w).Show()
					time.Sleep(100 * time.Millisecond)
					toolbar.Items[3].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
	)

	statbar := widget.NewLabel("List/Tree statusbar")

	c := container.NewBorder(toolbar, statbar, nil, nil, tree)

	return &List{
		Container: c,
		list:      list,
		tree:      tree,
		toolbar:   toolbar,
	}
}
