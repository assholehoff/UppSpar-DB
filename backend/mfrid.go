package backend

import (
	"database/sql"
	"database/sql/driver"
	"errors"
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

/* Returns a tree-friendly identifying string */
func (id MfrID) TString() string {
	return fmt.Sprintf("MFR-%d", id)
}
func (id MfrID) Branch() bool {
	var MfrID MfrID
	query := `SELECT Manufacturer.MfrID FROM Manufacturer LEFT JOIN Model WHERE Model.MfrID = @0 LIMIT 1`
	err := b.db.QueryRow(query, id).Scan(&MfrID)
	return !errors.Is(err, sql.ErrNoRows)
}
func (id MfrID) TypeName() string {
	return "MfrID"
}
func (id MfrID) Children() (children []ModelID) {
	query := `SELECT ModelID FROM Model WHERE MfrID = @0 ORDER BY Name ASC`
	rows, err := b.db.Query(query, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var child ModelID
		rows.Scan(&child)
		children = append(children, child)
	}
	return
}

/* Returns Name from SQL query */
func (id MfrID) Name() (val string, err error) {
	return id.getString("Name")
}

func (id MfrID) SetName() error {
	key := "Name"
	val, err := id.Manufacturer().Name.Get()
	if err != nil {
		return fmt.Errorf("MfrID.SetName() error: %w", err)
	}
	return id.setString(key, val)
}

func (id MfrID) Manufacturer() *Manufacturer {
	return getManufacturer(id)
}

/* Get the pointer to Manufacturer from map or make one and return it */
func getManufacturer(id MfrID) *Manufacturer {
	if mfr := b.Metadata.mfrData[id]; mfr == nil {
		mfr = newMfr(id)
		b.Metadata.mfrData[id] = mfr
	}
	return b.Metadata.mfrData[id]
}

func (id MfrID) getBool(key string) (val bool, err error) {
	b, err := getValue[sql.NullBool]("Manufacturer", id, key)
	if b.Valid && err == nil {
		val = b.Bool
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("MfrID(%d).getBool(%s) error: %s", id, key, err)
	}
	return
}
func (id MfrID) getFloat(key string) (val float64, err error) {
	f, err := getValue[sql.NullFloat64]("Manufacturer", id, key)
	if f.Valid && err == nil {
		val = f.Float64
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("MfrID(%d).getFloat(%s) error: %s", id, key, err)
	}
	return
}
func (id MfrID) getInt(key string) (val int, err error) {
	i, err := getValue[sql.NullInt64]("Manufacturer", id, key)
	if i.Valid && err == nil {
		val = int(i.Int64)
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("MfrID(%d).getInt(%s) error: %s", id, key, err)
	}
	return
}
func (id MfrID) getString(key string) (val string, err error) {
	s, err := getValue[sql.NullString]("Manufacturer", id, key)
	if s.Valid && err == nil {
		val = s.String
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("MfrID(%d).getString(%s) error: %s", id, key, err)
	}
	return
}
func (id MfrID) setBool(key string, val bool) error {
	err := setValue("Manufacturer", id, key, val)
	return err
}
func (id MfrID) setFloat(key string, val float64) error {
	err := setValue("Manufacturer", id, key, val)
	return err
}
func (id MfrID) setInt(key string, val int) error {
	err := setValue("Manufacturer", id, key, val)
	return err
}
func (id MfrID) setString(key string, val string) error {
	err := setValue("Manufacturer", id, key, val)
	return err
}
