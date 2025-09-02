package backend

import "fyne.io/fyne/v2/data/binding"

/* Wishlist and function to apply wishes to items in inventory */

type Wishlist binding.UntypedTree

func NewWishlist() Wishlist {
	return binding.NewUntypedTree()
}
