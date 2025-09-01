package backend

import (
	"UppSpar/backend/journal"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var b *Backend
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
	DB, err := sql.Open("sqlite3", file)

	b = &Backend{
		db: DB,
	}

	b.Journal = journal.NewJournal(b.db)
	b.createTables()
	b.Items = NewItems()
	b.Metadata = NewMetadata()
	b.Settings = NewSettings()
	b.Wishlist = NewWishlist()

	b.Items.GetItemIDs()

	b.Metadata.getCatIDList()
	b.Metadata.GetCatIDTree()
	b.Metadata.GetMfrIDs()
	b.Metadata.GetModelIDs()
	b.Metadata.GetProductTree()
	b.Metadata.getAllUnitIDs()
	b.Metadata.getAllItemStatusIDs()

	return b, err
}

func (backend *Backend) Close() error {
	if backend.db != nil {
		return backend.db.Close()
	}
	return nil
}

func ItemIDWidth() int {
	defaultWidth := 7
	i, err := b.Settings.ItemIDWidth.Get()
	if err != nil {
		log.Println(err)
		return defaultWidth
	}
	return i
}
func CatIDFor(s string) (CatID, error) {
	// TODO handle when database contains multiple rows with 's'
	var i NullInt
	var id CatID

	query := `SELECT CatID FROM Category WHERE Name = @0`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return id, fmt.Errorf("findCatIDFor error: %w", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(s).Scan(&i)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return id, fmt.Errorf("findCatIDFor error: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return id, err
	}

	if !i.Valid {
		log.Printf("findCatIDFor(%s); i.Valid: %t, i.Int: %v", s, i.Valid, i.Int)
		return id, err
	}

	id = CatID(i.Int)
	return id, nil
}
func MfrIDFor(s string) (MfrID, error) {
	// TODO handle when database contains multiple rows with 's'
	var i NullInt
	var id MfrID

	query := `SELECT MfrID FROM Manufacturer WHERE Name = @0`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return id, fmt.Errorf("MfrIDFor(%s) error: %w", s, err)
	}
	defer stmt.Close()
	stmt.QueryRow(s).Scan(&id)

	if !i.Valid {
		return id, err
	}

	id = MfrID(i.Int)
	// log.Printf("MfrIDFor(%s) is %d", s, id)
	return id, nil
}
func ModelIDFor(mfr MfrID, s string) (ModelID, error) {
	// TODO handle when database contains multiple rows with 's'
	var i NullInt
	var id ModelID

	query := `SELECT ModelID FROM Model WHERE MfrID = ? AND Name = ?`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	stmt.QueryRow(mfr, s).Scan(&id)

	if !i.Valid {
		return id, err
	}

	id = ModelID(i.Int)
	// log.Printf("ModelIDFor(%s) is %d", s, id)
	return id, nil
}
func UnitIDFor(s string) (UnitID, error) {
	var i NullInt
	var id UnitID

	query := `SELECT UnitID FROM Metric WHERE Text = @0`
	stmt, err := b.db.Prepare(query)
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
	stmt, err := b.db.Prepare(query)
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
