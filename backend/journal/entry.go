package journal

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"runtime"
	"time"
)

var (
	_ driver.Valuer = (*EntryID)(nil)
	_ driver.Valuer = (*Event)(nil)
	_ driver.Valuer = (*Level)(nil)

	_ sql.Scanner = (*EntryID)(nil)
	_ sql.Scanner = (*Event)(nil)
	_ sql.Scanner = (*Level)(nil)

	_ fmt.Stringer = (*EntryID)(nil)
	_ fmt.Stringer = (*Entry)(nil)
	_ fmt.Stringer = (*Event)(nil)
	_ fmt.Stringer = (*Level)(nil)
)

type EntryID int

// String implements fmt.Stringer.
func (e *EntryID) String() string {
	return fmt.Sprintf("%d", *e)
}

// Scan implements sql.Scanner.
func (e *EntryID) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		*e = 0
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		*e = EntryID(src.(int))
	case "int8":
		*e = EntryID(src.(int8))
	case "int16":
		*e = EntryID(src.(int16))
	case "int32":
		*e = EntryID(src.(int32))
	case "int64":
		*e = EntryID(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				*e = EntryID(math.MaxInt32)
				return ErrLossyConversion
			}
		}
	case "uint":
		*e = EntryID(src.(uint))
	case "uint8":
		*e = EntryID(src.(uint8))
	case "uint16":
		*e = EntryID(src.(uint16))
	case "uint32":
		*e = EntryID(src.(uint32))
	case "uint64":
		*e = EntryID(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				*e = EntryID(math.MaxUint32)
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			*e = EntryID(math.MaxInt64)
			return ErrLossyConversion
		}
	default:
		log.Printf("ItemID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return ErrInvalidType
	}
	return nil
}

// Value implements driver.Valuer.
func (e *EntryID) Value() (driver.Value, error) {
	return int64(*e), nil
}

func (e *EntryID) Int() int {
	return int(*e)
}

type Event uint8

// String implements fmt.Stringer.
func (e *Event) String() string {
	switch *e {
	case Log:
		return "Log"
	case Add:
		return "Add"
	case Copy:
		return "Copy"
	case Edit:
		return "Edit"
	case Delete:
		return "Delete"
	case SQL:
		return "SQL"
	default:
		return ""
	}
}

// Scan implements sql.Scanner.
func (e *Event) Scan(src any) error {
	InvalidEntryIDError := errors.New("invalid Event")
	InvalidEntryIDZeroError := errors.New("Event must not be zero")
	InvalidTypeError := errors.New("invalid type")
	IntegerOverflowError := errors.New("integer overflow")
	if !reflect.ValueOf(src).IsValid() {
		*e = 0
		return InvalidEntryIDError
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		if src.(int) > math.MaxUint8 {
			*e = Event(math.MaxUint8)
			return IntegerOverflowError
		} else {
			if src.(int) < 1 {
				*e = 0
			} else {
				*e = Event(src.(int))
			}
		}
	case "int8":
		if src.(int8) < 1 {
			*e = 0
		} else {
			*e = Event(src.(int8))
		}
	case "int16":
		if src.(int16) > math.MaxUint8 {
			*e = Event(math.MaxUint8)
			return IntegerOverflowError
		} else {
			if src.(int16) < 1 {
				*e = 0
			} else {
				*e = Event(src.(int16))
			}
		}
	case "int32":
		if src.(int32) > math.MaxUint8 {
			*e = Event(math.MaxUint8)
			return IntegerOverflowError
		} else {
			if src.(int32) < 1 {
				*e = 0
			} else {
				*e = Event(src.(int32))
			}
		}
	case "int64":
		if src.(int64) > math.MaxUint8 {
			*e = Event(math.MaxUint8)
			return IntegerOverflowError
		} else {
			if src.(int64) < 1 {
				*e = 0
			} else {
				*e = Event(src.(int64))
			}
		}
	case "uint":
		if src.(uint) > math.MaxUint8 {
			*e = Event(math.MaxUint8)
			return IntegerOverflowError
		} else {
			*e = Event(src.(uint))
		}
	case "uint8":
		*e = Event(src.(uint8))
	case "uint16":
		if src.(uint16) > math.MaxUint8 {
			*e = Event(math.MaxUint8)
			return IntegerOverflowError
		} else {
			*e = Event(src.(uint16))
		}
	case "uint32":
		if src.(uint32) > math.MaxUint8 {
			*e = Event(math.MaxUint8)
			return IntegerOverflowError
		} else {
			*e = Event(src.(uint32))
		}
	case "uint64":
		if src.(uint64) > math.MaxUint8 {
			*e = Event(math.MaxUint8)
			return IntegerOverflowError
		} else {
			*e = Event(src.(uint64))
		}
	default:
		log.Printf("ItemID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return InvalidTypeError
	}
	if *e == Event(0) {
		return InvalidEntryIDZeroError
	}
	return nil
}

// Value implements driver.Valuer.
func (e *Event) Value() (driver.Value, error) {
	return int64(*e), nil
}

func (e *Event) Int() int {
	return int(*e)
}

const (
	Log Event = iota + 1
	Add
	Copy
	Edit
	Delete
	SQL
)

type Level uint8

// String implements fmt.Stringer.
func (l *Level) String() string {
	switch *l {
	case Message:
		return "Message"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	default:
		return ""
	}
}

// Scan implements sql.Scanner.
func (l *Level) Scan(src any) error {
	InvalidLevelError := errors.New("invalid Level")
	InvalidLevelZeroError := errors.New("Level must not be zero")
	InvalidTypeError := errors.New("invalid type")
	IntegerOverflowError := errors.New("integer overflow")
	if !reflect.ValueOf(src).IsValid() {
		*l = 0
		return InvalidLevelError
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		if src.(int) > math.MaxUint8 {
			*l = Level(math.MaxUint8)
			return IntegerOverflowError
		} else {
			if src.(int) < 1 {
				*l = 0
			} else {
				*l = Level(src.(int))
			}
		}
	case "int8":
		if src.(int8) < 1 {
			*l = 0
		} else {
			*l = Level(src.(int8))
		}
	case "int16":
		if src.(int16) > math.MaxUint8 {
			*l = Level(math.MaxUint8)
			return IntegerOverflowError
		} else {
			if src.(int16) < 1 {
				*l = 0
			} else {
				*l = Level(src.(int16))
			}
		}
	case "int32":
		if src.(int32) > math.MaxUint8 {
			*l = Level(math.MaxUint8)
			return IntegerOverflowError
		} else {
			if src.(int32) < 1 {
				*l = 0
			} else {
				*l = Level(src.(int32))
			}
		}
	case "int64":
		if src.(int64) > math.MaxUint8 {
			*l = Level(math.MaxUint8)
			return IntegerOverflowError
		} else {
			if src.(int64) < 1 {
				*l = 0
			} else {
				*l = Level(src.(int64))
			}
		}
	case "uint":
		if src.(uint) > math.MaxUint8 {
			*l = Level(math.MaxUint8)
			return IntegerOverflowError
		} else {
			*l = Level(src.(uint))
		}
	case "uint8":
		*l = Level(src.(uint8))
	case "uint16":
		if src.(uint16) > math.MaxUint8 {
			*l = Level(math.MaxUint8)
			return IntegerOverflowError
		} else {
			*l = Level(src.(uint16))
		}
	case "uint32":
		if src.(uint32) > math.MaxUint8 {
			*l = Level(math.MaxUint8)
			return IntegerOverflowError
		} else {
			*l = Level(src.(uint32))
		}
	case "uint64":
		if src.(uint64) > math.MaxUint8 {
			*l = Level(math.MaxUint8)
			return IntegerOverflowError
		} else {
			*l = Level(src.(uint64))
		}
	default:
		log.Printf("ItemID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return InvalidTypeError
	}
	if *l == Level(0) {
		return InvalidLevelZeroError
	}
	return nil
}

// Value implements driver.Valuer.
func (l *Level) Value() (driver.Value, error) {
	return int64(*l), nil
}

func (l *Level) Int() int {
	return int(*l)
}

const (
	Message Level = iota + 1
	Warning
	Error
)

type Entry struct {
	EntryID EntryID
	Level   Level
	Event   Event
	Message string
	Time    time.Time
}

// String implements fmt.Stringer.
func (e *Entry) String() string {
	return fmt.Sprintf("%s %s", e.Time.Format(time.DateTime), e.Message)
}

func (e *Entry) Strings() (id, level, event, msg, tme string) {
	return e.EntryID.String(),
		e.Level.String(),
		e.Event.String(),
		e.Message,
		e.Time.Format(time.DateTime)
}
