package bridge

import (
	"UppSpar/backend"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	midget "github.com/assholehoff/fyne-midget"
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
			co := midget.NewLabel("Template item name", "00000000", "")
			co.SetTop()
			return co
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			val, err := di.(binding.Untyped).Get()
			if err != nil {
				panic(err)
			}
			ItemID := val.(backend.ItemID)
			subtext := binding.NewString()

			ItemID.Item().Category.AddListener(binding.NewDataListener(func() {
				id, _ := ItemID.Item().ItemIDString.Get()
				cat, _ := ItemID.Item().Category.Get()
				subtext.Set(fmt.Sprintf("%s : %s", id, strings.TrimSpace(strings.ToUpper(cat))))
			}))

			co.(*midget.Label).BindText(ItemID.Item().Name)
			co.(*midget.Label).BindSubtext(subtext)
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
					NewExportExcelDialog(b, w).Show()
					time.Sleep(100 * time.Millisecond)
					toolbar.Items[3].(*widget.ToolbarAction).Enable()
				})
			}()
		}),
	)

	statbar := widget.NewLabel("List/Tree statusbar")

	c := container.NewBorder(toolbar, statbar, nil, nil, list)

	return &List{
		Container: c,
		list:      list,
		tree:      tree,
		toolbar:   toolbar,
	}
}
