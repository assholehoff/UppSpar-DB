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

/* Returns a tree-friendly identifying string */
func (id CatID) TString() string {
	return fmt.Sprintf("CAT-%d", id)
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

func (id CatID) TypeName() string {
	return "CatID"
}
func (id CatID) Branch() bool {
	var CatID CatID
	query := `SELECT CatID FROM Category WHERE ParentID = @0 LIMIT 1`
	err := be.db.QueryRow(query, id).Scan(&CatID)
	return !errors.Is(err, sql.ErrNoRows)
}
func (id CatID) ParentID() (CatID, error) {
	val, err := id.getInt("ParentID")
	return CatID(val), err
}
func (id CatID) Children() []CatID {
	var children []CatID
	query := `SELECT CatID FROM Category WHERE ParentID = @0 ORDER BY Name ASC`
	rows, err := be.db.Query(query, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	for rows.Next() {
		var child CatID
		rows.Scan(&child)
		children = append(children, child)
	}
	return children
}

/* Returns Name from SQL query */
func (id CatID) Name() (val string, err error) {
	return id.getString("Name")
}

func (id CatID) SetName() error {
	key := "Name"
	val, err := id.Category().Name.Get()
	if err != nil {
		return fmt.Errorf("MfrID.SetName() error: %w", err)
	}
	return id.setString(key, val)
}

func (id CatID) Category() *Category {
	return getCategory(be, id)
}
func (id CatID) ShowPrice() bool {
	// TODO fix this lazy s..t
	if id == CatID(1) {
		return true
	}
	return false
}

/* Get the pointer to Category from map or make one and return it */
func getCategory(b *Backend, id CatID) *Category {
	if c := b.Metadata.categoryData[id]; c == nil {
		c = newCategory(b, id)
		b.Metadata.categoryData[id] = c
	}
	return b.Metadata.categoryData[id]
}
func (id CatID) getBool(key string) (val bool, err error) {
	b, err := getValue[sql.NullBool]("Category", id, key)
	if b.Valid {
		val = b.Bool
	} else {
		log.Printf("CatID.getBool(%s) b is invalid (NULL), err is %v", key, err)
	}
	return
}
func (id CatID) getFloat(key string) (val float64, err error) {
	f, err := getValue[sql.NullFloat64]("Category", id, key)
	if f.Valid {
		val = f.Float64
	} else {
		log.Printf("CatID.getFloat(%s) %s is invalid (NULL), err is %v", key, key, err)
	}
	return
}
func (id CatID) getInt(key string) (val int, err error) {
	i, err := getValue[sql.NullInt64]("Category", id, key)
	val = int(i.Int64)
	if !i.Valid {
		log.Printf("CatID.getInt(%s) %s is invalid (NULL), err is %v", key, key, err)
	}
	return
}
func (id CatID) getString(key string) (val string, err error) {
	s, err := getValue[sql.NullString]("Category", id, key)
	if s.Valid {
		val = s.String
	} else {
		log.Printf("CatID.getInt(%s) %s is invalid (NULL), err is %v", key, key, err)
	}
	return
}
func (id CatID) setBool(key string, val bool) error {
	err := setValue("Category", id, key, val)
	return err
}
func (id CatID) setFloat(key string, val float64) error {
	err := setValue("Category", id, key, val)
	return err
}
func (id CatID) setInt(key string, val int) error {
	err := setValue("Category", id, key, val)
	return err
}
func (id CatID) setString(key string, val string) error {
	err := setValue("Category", id, key, val)
	return err
}
