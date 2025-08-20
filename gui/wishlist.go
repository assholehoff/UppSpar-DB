package gui

import (
	"UppSpar/backend"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type wishlistView struct {
	container *fyne.Container
}

func newWishlistView(b *backend.Backend) *wishlistView {
	return &wishlistView{
		container: container.NewBorder(nil, nil, nil, nil),
	}
}
