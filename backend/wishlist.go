package backend

/* Wishlist and function to apply wishes to items in inventory */

type Wishlist struct{}

func NewWishlist(b *Backend) *Wishlist {
	return &Wishlist{}
}
