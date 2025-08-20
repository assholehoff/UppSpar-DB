package journal

import (
	"database/sql"
	"slices"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const (
	minute = "2006-01-02 15:04"
	subsec = "2006-01-02 15:04:05.999"
)

var (
	db *sql.DB
)

type Order int

const (
	Ascending Order = iota + 1
	Descending
)

func (o Order) String() string {
	if o == Ascending {
		return "ASC"
	}
	return "DESC"
}

type config struct {
	limit   int
	scroll  bool
	sorting Order
	time    string
}

type Journal struct {
	db      *sql.DB
	config  config
	entries map[EntryID]*Entry

	List binding.UntypedList
}

func NewJournal(dc *sql.DB) *Journal {
	db = dc
	j := &Journal{
		db: dc,
		config: config{
			limit:   100,
			sorting: Ascending,
			time:    minute,
		},
		entries: make(map[EntryID]*Entry),
		List:    binding.NewUntypedList(),
	}
	// 1. connect and set up/verify tables/schema
	j.createTables()
	// 2. get `limit` last entries from database
	j.List.Set(j.getRecentEntryIds())
	return j
}

func (j *Journal) NewEntry(level Level, event Event, message string) {
	id := j.newEntry(level, event, message)
	entry := j.getEntry(id)
	j.addEntry(id, entry)
}

func (j *Journal) NewMessage(message string) {
	id := j.newEntry(Message, Log, message)
	entry := j.getEntry(id)
	j.addEntry(id, entry)
}

func (j *Journal) NewWarning(message string) {
	id := j.newEntry(Warning, Log, message)
	entry := j.getEntry(id)
	j.addEntry(id, entry)
}

func (j *Journal) NewError(message string) {
	id := j.newEntry(Error, Log, message)
	entry := j.getEntry(id)
	j.addEntry(id, entry)
}

func (j *Journal) GetEntry(id EntryID) *Entry {
	return j.getEntry(id)
}

func (j *Journal) GetAllEntryIDs() []any {
	return j.getAllEntryIDs()
}

func (j *Journal) addEntry(id EntryID, entry *Entry) {
	j.entries[id] = entry
	if j.config.sorting == Ascending {
		j.List.Append(id)
	} else {
		j.List.Prepend(id)
	}
}

func (j *Journal) EntryListItemID(id EntryID) (widget.ListItemID, bool) {
	ids, _ := j.List.Get()
	listId := slices.IndexFunc(ids, func(n any) bool {
		return n.(EntryID) == id
	})
	if listId == -1 {
		return listId, false
	}
	return listId, true
}

func (j *Journal) Refresh() {
	// TODO: fix way to reload currently loaded entries
	j.List.Set(j.getRecentEntryIds())
}

func (j *Journal) LastEntryId() EntryID {
	id, err := j.List.GetValue(j.List.Length() - 1)
	if err != nil {
		panic(err)
	}
	return id.(EntryID)
}

func (j *Journal) SetEntryLimit(n int) {
	j.config.limit = n
}

func (j *Journal) SetTimeLayout(s string) {
	j.config.time = s
}

func (j *Journal) Autoscroll() bool {
	return j.config.scroll
}
func (j *Journal) SetAutoscroll() {
	j.config.scroll = true
}
func (j *Journal) UnsetAutoscroll() {
	j.config.scroll = false
}
func (j *Journal) SortAscending() {
	j.config.sorting = Ascending
	j.Refresh()
}
func (j *Journal) SortDescending() {
	j.config.sorting = Descending
	j.Refresh()
}
