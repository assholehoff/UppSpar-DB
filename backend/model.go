package backend

import (
	"database/sql"
	"log"

	"fyne.io/fyne/v2/data/binding"
)

type Model struct {
	binding.DataItem
	db      *sql.DB
	ModelID ModelID
	Name    binding.String
}

func newModel(b *Backend, id ModelID) *Model {
	mdl := &Model{
		db:      b.db,
		ModelID: id,
		Name:    binding.NewString(),
	}

	name, err := mdl.ModelID.Name()
	if err != nil {
		log.Println(err)
	}
	mdl.Name.Set(name)
	mdl.Name.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetName() }))
	return mdl
}
