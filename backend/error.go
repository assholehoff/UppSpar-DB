package backend

import (
	"errors"
)

/* Custom Errors */

var (
	ErrNotFound         = errors.New("not found")
	ErrIndexOutOfBounds = errors.New("index out of bounds")
	ErrInvalidType      = errors.New("invalid type")
	ErrInvalidValue     = errors.New("invalid value")
	ErrLossyConversion  = errors.New("lossy conversion")
	ErrSQLNullValue     = errors.New("value is null")
)
