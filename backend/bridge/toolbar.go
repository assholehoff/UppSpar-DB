package bridge

import (
	"UppSpar/backend"

	"fyne.io/fyne/v2"
)

type Tools struct {
	Check  Checks
	Entry  Entries
	Label  Labels
	Radio  Radios
	Select Selects
}

func NewTools(b *backend.Backend, w fyne.Window) *Tools {
	t := &Tools{}
	return t
}

func NewSearchBar(b *backend.Backend, w fyne.Window) *Tools {
	t := &Tools{}
	return t
}
