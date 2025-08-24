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
	_ sql.Scanner   = (*MfrID)(nil)
	_ driver.Valuer = (*MfrID)(nil)
	_ fmt.Stringer  = (*MfrID)(nil)
)

type MfrID int

/* String implements fmt.Stringer. */
func (id MfrID) String() string {
	return fmt.Sprintf("%d", id)
}

/* Value implements driver.Valuer. */
func (id MfrID) Value() (driver.Value, error) {
	return int64(id), nil
}

/* Scan implements sql.Scanner. */
func (id *MfrID) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		/* Note: THIS happens when the SQL value is NULL! */
		*id = 0
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		*id = MfrID(src.(int))
	case "int8":
		*id = MfrID(src.(int8))
	case "int16":
		*id = MfrID(src.(int16))
	case "int32":
		*id = MfrID(src.(int32))
	case "int64":
		*id = MfrID(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				*id = MfrID(math.MaxInt32)
				return ErrLossyConversion
			}
		}
	case "uint":
		*id = MfrID(src.(uint))
	case "uint8":
		*id = MfrID(src.(uint8))
	case "uint16":
		*id = MfrID(src.(uint16))
	case "uint32":
		*id = MfrID(src.(uint32))
	case "uint64":
		*id = MfrID(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				*id = MfrID(math.MaxUint32)
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			*id = MfrID(math.MaxInt64)
			return ErrLossyConversion
		}
	default:
		log.Printf("ItemID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return ErrInvalidType
	}
	return nil
}

/* Returns Name from SQL query */
func (id MfrID) Name() (val string, err error) {
	var s sql.NullString
	query := `SELECT Name FROM Manufacturer WHERE MfrID = @0`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		return val, fmt.Errorf("MfrID(%d).Name() error: %w", id, err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&s)
	if err != nil {
		return val, fmt.Errorf("MfrID(%d).Name() error: %w", id, err)
	}
	if s.Valid {
		val = s.String
	} else {
		err = ErrSQLNullValue
	}
	return
}

func (id MfrID) SetName() error {
	s, err := id.Manufacturer().Name.Get()
	if err != nil {
		return fmt.Errorf("MfrID.SetName() error: %w", err)
	}
	query := `UPDATE Manufacturer SET Name = @0 WHERE MfrID = @1 AND Name <> @2`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("MfrID.SetName() error: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(s, id, s)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("MfrID.SetName() error: %w", err)
	}
	return err
}

func (id MfrID) Manufacturer() *Manufacturer {
	return getManufacturer(be, id)
}

/* Get the pointer to Manufacturer from map or make one and return it */
func getManufacturer(b *Backend, id MfrID) *Manufacturer {
	if mfr := b.Metadata.mfrData[id]; mfr == nil {
		mfr = newMfr(b, id)
		b.Metadata.mfrData[id] = mfr
	}
	return b.Metadata.mfrData[id]
}
