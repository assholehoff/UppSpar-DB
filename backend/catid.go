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

/* Categories and other metadata */
var (
	_ sql.Scanner   = (*CatID)(nil)
	_ driver.Valuer = (*CatID)(nil)
	_ fmt.Stringer  = (*CatID)(nil)
)

type CatID int

/* String implements fmt.Stringer. */
func (id CatID) String() string {
	return fmt.Sprintf("%d", id)
}

/* Value implements driver.Valuer. */
func (id CatID) Value() (driver.Value, error) {
	return int64(id), nil
}

/* Scan implements sql.Scanner. */
func (id *CatID) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		/* Note: THIS happens when the SQL value is NULL! */
		*id = 0
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		*id = CatID(src.(int))
	case "int8":
		*id = CatID(src.(int8))
	case "int16":
		*id = CatID(src.(int16))
	case "int32":
		*id = CatID(src.(int32))
	case "int64":
		*id = CatID(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				*id = CatID(math.MaxInt32)
				return ErrLossyConversion
			}
		}
	case "uint":
		*id = CatID(src.(uint))
	case "uint8":
		*id = CatID(src.(uint8))
	case "uint16":
		*id = CatID(src.(uint16))
	case "uint32":
		*id = CatID(src.(uint32))
	case "uint64":
		*id = CatID(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				*id = CatID(math.MaxUint32)
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			*id = CatID(math.MaxInt64)
			return ErrLossyConversion
		}
	default:
		log.Printf("ItemID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return ErrInvalidType
	}
	return nil
}

func (id CatID) Category() *Category {
	return getCategory(be, id)
}

/* Returns Name from SQL query */
func (id CatID) Name() (val string, err error) {
	var s sql.NullString
	query := `SELECT Name FROM Category WHERE CatID = @0`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		return val, fmt.Errorf("CatID(%d).Name() error: %w", id, err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&s)
	if err != nil {
		return val, fmt.Errorf("CatID(%d).Name() error: %w", id, err)
	}
	if s.Valid {
		val = s.String
	} else {
		err = ErrSQLNullValue
	}
	return
}

func (id CatID) SetName() error {
	s, err := id.Category().Name.Get()
	if err != nil {
		return fmt.Errorf("CatID.SetName() error: %w", err)
	}
	query := `UPDATE Category SET Name = @0 WHERE CatID = @1 AND Name <> @2`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("CatID.SetName() error: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(s, id, s)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("CatID.SetName() error: %w", err)
	}
	return err
}

/* Get the pointer to Category from map or make one and return it */
func getCategory(b *Backend, id CatID) *Category {
	if c := b.Metadata.categoryData[id]; c == nil {
		c = newCategory(b, id)
		b.Metadata.categoryData[id] = c
	}
	return b.Metadata.categoryData[id]
}
