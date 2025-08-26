package backend

import (
	"database/sql"
	"errors"

	"fyne.io/fyne/v2/data/binding"
)

type Category struct {
	binding.DataItem
	db     *sql.DB
	CatID  CatID
	Name   binding.String
	Parent binding.String
	Config map[string]binding.Bool
}

func newCategory(b *Backend, id CatID) *Category {
	c := &Category{
		db:     b.db,
		CatID:  id,
		Name:   binding.NewString(),
		Parent: binding.NewString(),
		Config: make(map[string]binding.Bool),
	}

	c.getNameStrings()
	c.makeConfigMap()
	c.Name.AddListener(binding.NewDataListener(func() { c.CatID.SetName() }))
	return c
}
func (c *Category) makeConfigMap() {
	addConfig := func(key string) {
		c.Config[key] = binding.NewBool()
		b, err := c.CatID.getConfig(key)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		}
		c.Config[key].Set(b)
		c.Config[key].AddListener(binding.NewDataListener(func() {
			ba, _ := c.Config[key].Get()
			if bb, _ := c.CatID.getConfig(key); ba != bb {
				c.CatID.setConfig(key, ba)
			}
		}))
	}
	addConfig("ShowPrice")
	addConfig("ShowLength")
	addConfig("ShowVolume")
	addConfig("ShowWeight")
}
func (c *Category) getNameStrings() {
	var Name sql.NullString
	var ParentID CatID
	query := `SELECT Name, ParentID FROM Category WHERE CatID = @0`
	c.db.QueryRow(query, c.CatID).Scan(&Name, &ParentID)

	c.Name.Set(Name.String)

	n, _ := ParentID.Name()
	c.Parent.Set(n)
}
