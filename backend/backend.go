package backend

import (
	"UppSpar/backend/journal"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var be *Backend
var subsec = "2006-01-02 15:04:05.999"

type Backend struct {
	db       *sql.DB
	Items    *Items
	Journal  *journal.Journal
	Metadata *Metadata
	Settings *Settings
	Wishlist *Wishlist
}

func NewBackend(file string) (*Backend, error) {
	sqldb, err := sql.Open("sqlite3", file)

	be = &Backend{
		db: sqldb,
	}

	be.Journal = journal.NewJournal(be.db)
	be.createTables()
	be.Items = NewItems(be)
	be.Metadata = NewMetadata(be)
	be.Settings = NewSettings(be)
	be.Wishlist = NewWishlist(be)

	be.Items.GetAllItemIDs()

	be.Metadata.getAllCatIDs()
	be.Metadata.getAllUnitIDs()
	be.Metadata.getAllItemStatusIDs()

	return be, err
}

func (b *Backend) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

func ItemIDWidth() int {
	return be.Settings.getItemIDWidth()
}
func CatIDFor(s string) (CatID, error) {
	return be.Metadata.findCatIDFor(s)
}
func UnitIDFor(s string) (UnitID, error) {
	var i NullInt
	var id UnitID

	query := `SELECT UnitID FROM Metric WHERE Text = @0`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		return id, fmt.Errorf("UnitIDFor(%s) error: %w", s, err)
	}
	defer stmt.Close()
	stmt.QueryRow(s).Scan(&id)

	if !i.Valid {
		return id, ErrNotFound
	}

	id = UnitID(i.Int)
	return id, nil
}
func ItemStatusIDFor(s string) (ItemStatusID, error) {
	var i NullInt
	var is ItemStatusID

	query := `SELECT ItemStatusID FROM ItemStatus WHERE Name = @0`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		return is, fmt.Errorf("ItemStatusFor(%s) error: %w", s, err)
	}
	defer stmt.Close()
	stmt.QueryRow(s).Scan(&i)

	if !i.Valid {
		return is, ErrNotFound
	}

	is = ItemStatusID(i.Int)
	return is, nil
}
