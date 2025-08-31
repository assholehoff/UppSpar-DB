package backend

import (
	"database/sql"
	"errors"

	"fyne.io/fyne/v2/data/binding"
)

type Manufacturer struct {
	binding.DataItem
	MfrID MfrID
	Name  binding.String
}

func newMfr(id MfrID) *Manufacturer {
	mfr := &Manufacturer{
		MfrID: id,
		Name:  binding.NewString(),
	}

	name, err := mfr.MfrID.Name()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	mfr.Name.Set(name)
	mfr.Name.AddListener(binding.NewDataListener(func() { mfr.MfrID.SetName(); b.Metadata.GetMfrIDs(); b.Metadata.GetProductTree() }))
	return mfr
}
