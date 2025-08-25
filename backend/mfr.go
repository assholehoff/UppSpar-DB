package backend

import (
	"database/sql"
	"log"

	"fyne.io/fyne/v2/data/binding"
)

type Manufacturer struct {
	binding.DataItem
	db      *sql.DB
	MfrID   MfrID
	Name    binding.String
	branch  bool
	touched bool
}

func newMfr(b *Backend, id MfrID) *Manufacturer {
	mfr := &Manufacturer{
		db:    b.db,
		MfrID: id,
		Name:  binding.NewString(),
	}

	name, err := mfr.MfrID.Name()
	if err != nil {
		log.Println(err)
	}
	mfr.Name.Set(name)
	mfr.Name.AddListener(binding.NewDataListener(func() { mfr.MfrID.SetName() }))
	return mfr
}
