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

type ModelID int

/* String implements fmt.Stringer. */
func (id ModelID) String() string {
	return fmt.Sprintf("%d", id)
}

/* Value implements driver.Valuer. */
func (id ModelID) Value() (driver.Value, error) {
	return int64(id), nil
}

/* Scan implements sql.Scanner. */
func (id *ModelID) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		/* Note: THIS happens when the SQL value is NULL! */
		*id = 0
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		*id = ModelID(src.(int))
	case "int8":
		*id = ModelID(src.(int8))
	case "int16":
		*id = ModelID(src.(int16))
	case "int32":
		*id = ModelID(src.(int32))
	case "int64":
		*id = ModelID(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				*id = ModelID(math.MaxInt32)
				return ErrLossyConversion
			}
		}
	case "uint":
		*id = ModelID(src.(uint))
	case "uint8":
		*id = ModelID(src.(uint8))
	case "uint16":
		*id = ModelID(src.(uint16))
	case "uint32":
		*id = ModelID(src.(uint32))
	case "uint64":
		*id = ModelID(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				*id = ModelID(math.MaxUint32)
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			*id = ModelID(math.MaxInt64)
			return ErrLossyConversion
		}
	default:
		log.Printf("ItemID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return ErrInvalidType
	}
	return nil
}

/* Returns Name from SQL query */
func (id ModelID) Name() (val string, err error) {
	var s sql.NullString
	query := `SELECT Name FROM Model WHERE ModelID = @0`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		return val, fmt.Errorf("ModelID(%d).Name() error: %w", id, err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&s)
	if err != nil {
		return val, fmt.Errorf("ModelID(%d).Name() error: %w", id, err)
	}
	if s.Valid {
		val = s.String
	} else {
		err = ErrSQLNullValue
	}
	return
}

func (id ModelID) SetName() error {
	s, err := id.Model().Name.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetName() error: %w", err)
	}
	query := `UPDATE Model SET Name = @0 WHERE ModelID = @1 AND Name <> @2`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("ModelID.SetName() error: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(s, id, s)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("ModelID.SetName() error: %w", err)
	}
	return err
}

func (id ModelID) Model() *Model {
	return getModel(be, id)
}

/* Get the pointer to Model from map or make one and return it */
func getModel(b *Backend, id ModelID) *Model {
	if mdl := b.Metadata.modelData[id]; mdl == nil {
		mdl = newModel(b, id)
		b.Metadata.modelData[id] = mdl
	}
	return b.Metadata.modelData[id]
}
