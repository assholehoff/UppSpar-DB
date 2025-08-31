package backend

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"math"
	"reflect"
	"runtime"
)

var (
	_ sql.Scanner   = (*UnitID)(nil)
	_ driver.Valuer = (*UnitID)(nil)
	_ fmt.Stringer  = (*UnitID)(nil)
)

type UnitID int

/* Scan implements sql.Scanner. */
func (id *UnitID) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		/* Note: THIS happens when the SQL value is NULL! */
		*id = 0
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		*id = UnitID(src.(int))
	case "int8":
		*id = UnitID(src.(int8))
	case "int16":
		*id = UnitID(src.(int16))
	case "int32":
		*id = UnitID(src.(int32))
	case "int64":
		*id = UnitID(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				*id = UnitID(math.MaxInt32)
				return ErrLossyConversion
			}
		}
	case "uint":
		*id = UnitID(src.(uint))
	case "uint8":
		*id = UnitID(src.(uint8))
	case "uint16":
		*id = UnitID(src.(uint16))
	case "uint32":
		*id = UnitID(src.(uint32))
	case "uint64":
		*id = UnitID(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				*id = UnitID(math.MaxUint32)
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			*id = UnitID(math.MaxInt64)
			return ErrLossyConversion
		}
	default:
		log.Printf("UnitID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return ErrInvalidType
	}
	return nil
}

// Value implements driver.Valuer.
func (id *UnitID) Value() (driver.Value, error) {
	return int64(*id), nil
}

// String implements fmt.Stringer.
func (id UnitID) String() string {
	switch id {
	case millimeter:
		return "mm"
	case centimeter:
		return "cm"
	case decimeter:
		return "dm"
	case meter:
		return "m"
	case gram:
		return "g"
	case hectogram:
		return "hg"
	case kilogram:
		return "kg"
	case milliliter:
		return "ml"
	case centiliter:
		return "cl"
	case deciliter:
		return "dl"
	case liter:
		return "l"
	default:
		return "m"
	}
}

/* Returns Name from SQL query */
func (id UnitID) Name() (val string, err error) {
	var s sql.NullString
	query := `SELECT Text FROM Metric WHERE UnitID = @0`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return val, fmt.Errorf("UnitID(%d).Name() error: %w", id, err)
	}
	defer stmt.Close()
	// log.Println(strings.Replace(query, "@0", fmt.Sprintf("%d", id), 1))
	err = stmt.QueryRow(id).Scan(&s)
	if err != nil {
		return val, fmt.Errorf("UnitID(%d).Name() error: %w", id, err)
	}
	if s.Valid {
		val = s.String
	} else {
		err = ErrSQLNullValue
	}
	return
}

const (
	millimeter UnitID = iota + 1
	centimeter
	decimeter
	meter
	gram
	hectogram
	kilogram
	milliliter
	centiliter
	deciliter
	liter
)
