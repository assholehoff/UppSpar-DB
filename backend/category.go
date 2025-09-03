package backend

import (
	"database/sql"
	"errors"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
)

type Category struct {
	binding.DataItem
	CatID  CatID
	Name   binding.String
	Parent binding.String
	Config map[string]binding.Bool
}

func newCategory(id CatID) *Category {
	c := &Category{
		CatID:  id,
		Name:   binding.NewString(),
		Parent: binding.NewString(),
		Config: make(map[string]binding.Bool),
	}

	c.getNameStrings()
	c.makeConfigMap()
	c.Name.AddListener(binding.NewDataListener(func() { c.CatID.SetName() }))
	c.Parent.AddListener(binding.NewDataListener(func() { c.CatID.SetParent() }))
	return c
}
func (c *Category) Bindings() map[string]binding.String {
	m := make(map[string]binding.String)
	m["Name"] = c.Name
	m["Parent"] = c.Parent
	return m
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
	b.db.QueryRow(query, c.CatID).Scan(&Name, &ParentID)

	c.Name.Set(Name.String)

	if ParentID == 0 {
		c.Parent.Set(lang.L("None"))
		return
	}
	n, _ := ParentID.Name()
	c.Parent.Set(n)
}
