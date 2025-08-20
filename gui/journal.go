package gui

import (
	"UppSpar/backend"
	"UppSpar/backend/bridge"
	"UppSpar/backend/journal"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type journalView struct {
	container *fyne.Container
	entryList *widget.List
	toolbar   *fyne.Container
}

func newJournalView(b *backend.Backend) *journalView {
	j := b.Journal
	msg := "ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ0123456789abcdefghijklmnopqrstuvwxyzåäö0123456789"
	createItem := func() fyne.CanvasObject {
		return bridge.NewJournalEntry(&journal.Entry{
			Event:   journal.Log,
			Level:   journal.Message,
			Message: msg,
			Time:    time.Now(),
		})
	}
	updateItem := func(di binding.DataItem, co fyne.CanvasObject) {
		val, err := di.(binding.Untyped).Get()
		if err != nil {
			panic(err)
		}
		id := val.(journal.EntryID)
		entry := j.GetEntry(id)
		co.(*bridge.JournalEntry).Bind(entry)
		co.(*bridge.JournalEntry).Format(backend.ItemIDWidth())
	}
	jv := &journalView{
		entryList: widget.NewListWithData(j.List, createItem, updateItem),
		toolbar:   container.NewHBox(),
	}
	list := jv.entryList
	scrollToNewEntry := binding.NewDataListener(func() {
		if j.Autoscroll() {
			EntryID := j.LastEntryId()
			ListID, present := j.EntryListItemID(EntryID)
			if present {
				list.ScrollTo(ListID)
			}
		}
	})
	j.List.AddListener(scrollToNewEntry)
	sortAscending := func() {
		j.SortAscending()
		list.Refresh()
	}
	sortDescending := func() {
		j.SortDescending()
		list.Refresh()
	}
	jv.toolbar.Add(widget.NewToolbar(
		widget.NewToolbarAction(theme.MoveUpIcon(), sortAscending),
		widget.NewToolbarAction(theme.MoveDownIcon(), sortDescending),
	))
	jv.toolbar.Add(ttw.NewCheck(lang.L("scroll to new entries"), func(b bool) {
		if b {
			j.SetAutoscroll()
		} else {
			j.UnsetAutoscroll()
		}
	}))
	jv.toolbar.Objects[1].(*ttw.Check).SetChecked(true)
	jv.container = container.NewBorder(jv.toolbar, nil, nil, nil, jv.entryList)
	return jv
}
