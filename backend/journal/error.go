package journal

import "errors"

var (
	ErrInvalidEntryID  = errors.New("invalid EntryID")
	ErrInvalidType     = errors.New("invalid type")
	ErrLossyConversion = errors.New("lossy conversion")
)
