package backend

import (
	"database/sql"
	"errors"

	"fyne.io/fyne/v2/data/binding"
)

type Category struct {
	binding.DataItem
	db      *sql.DB
	CatID   CatID
	Name    binding.String
	Config  map[string]bool   // which fields to display in form
	Data    map[string]string // what text to put in spreadsheet
	branch  bool
	touched bool
}

func newCategory(b *Backend, id CatID) *Category {
	c := &Category{
		db:     b.db,
		CatID:  id,
		Config: make(map[string]bool),
		Data:   make(map[string]string),
	}

	c.getAllFields()
	c.Name.AddListener(binding.NewDataListener(func() { c.CatID.SetName() }))
	return c
}

func (c *Category) getAllFields() {
	var Name sql.NullString
	query := `SELECT Name FROM Category WHERE CatID = @0`
	c.db.QueryRow(query, c.CatID).Scan(&Name)

	c.Name = binding.NewString()
	c.Name.Set(Name.String)

	query = `SELECT ConfigKey, ConfigVal FROM Category_Config WHERE CatID = @0`
	rows, err := c.db.Query(query, c.CatID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	for rows.Next() {
		var key sql.NullString
		var val sql.NullBool
		rows.Scan(&key, &val)
		c.Config[key.String] = val.Bool
	}

	query = `SELECT DataKey, DataVal FROM Category_Data WHERE CatID = @0`
	rows, err = c.db.Query(query, c.CatID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	for rows.Next() {
		var key, val sql.NullString
		rows.Scan(&key, &val)
		c.Data[key.String] = val.String
	}
}
