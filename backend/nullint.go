package backend

import (
	"database/sql"
	"database/sql/driver"
	"math"
	"reflect"
	"runtime"
)

type NullInt struct {
	Int   int
	Valid bool
}

var _ sql.Scanner = (*NullInt)(nil)
var _ driver.Valuer = (*NullInt)(nil)

func (n *NullInt) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		n.Int = 0
		n.Valid = false
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		n.Int = src.(int)
	case "int8":
		n.Int = int(src.(int8))
	case "int16":
		n.Int = int(src.(int16))
	case "int32":
		n.Int = int(src.(int32))
	case "int64":
		n.Int = int(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				return ErrLossyConversion
			}
		}
	case "uint":
		n.Int = int(src.(uint))
	case "uint8":
		n.Int = int(src.(uint8))
	case "uint16":
		n.Int = int(src.(uint16))
	case "uint32":
		n.Int = int(src.(uint32))
	case "uint64":
		n.Int = int(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			return ErrLossyConversion
		}
	default:
		return ErrInvalidType
	}
	n.Valid = true
	return nil
}
func (n *NullInt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Int), nil
}
