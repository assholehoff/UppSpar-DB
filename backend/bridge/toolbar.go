package bridge

import (
	"UppSpar/backend"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Tools struct {
	Container *fyne.Container
	Check     Checks
	Entry     Entries
	Label     Labels
	Radio     Radios
	Select    Selects
}

func NewTools(b *backend.Backend, w fyne.Window) *Tools {
	t := &Tools{}
	return t
}

func NewSearchBar(b *backend.Backend, w fyne.Window) *Tools {
	c := container.NewBorder(nil, nil, nil, nil)
	t := &Tools{
		Container: c,
		Check:     make(Checks),
		Entry:     make(Entries),
		Label:     make(Labels),
		Radio:     make(Radios),
		Select:    make(Selects),
	}
	return t
}
