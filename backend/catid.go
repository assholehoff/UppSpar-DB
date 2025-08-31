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
	"strings"

	"fyne.io/fyne/v2/lang"
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
		log.Printf("ItemID(%d).Scan(%v) unknown type %s", id, src, reflect.TypeOf(src).Name())
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
	err := b.db.QueryRow(query, id).Scan(&CatID)
	return !errors.Is(err, sql.ErrNoRows)
}
func (id CatID) ParentID() (CatID, error) {
	val, err := id.getInt("ParentID")
	return CatID(val), err
}
func (id CatID) Parents() int {
	n := 0
	return ancestors(id, n)
}
func ancestors(p CatID, n int) int {
	if a, err := p.ParentID(); a != 0 && err == nil {
		n++
		ancestors(a, n)
	}
	return n
}
func (id CatID) Children() []CatID {
	var children []CatID
	query := `SELECT CatID FROM Category WHERE ParentID = @0 ORDER BY Name ASC`
	rows, err := b.db.Query(query, id)
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
		return fmt.Errorf("MfrID(%d).SetName(%s) error: %s", id, val, err)
	}
	return id.setString(key, val)
}
func (id CatID) SetParentID(v CatID) error {
	return id.setInt("ParentID", int(v))
}
func (id CatID) SetParent() error {
	p, err := id.Category().Parent.Get()
	if err != nil {
		panic(err)
	}
	p = strings.TrimSpace(p)
	if p == lang.L("None") || p == "" {
		return id.SetParentID(CatID(0))
	}
	pid, err := CatIDFor(p)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("no rows, pid is %d", pid)
		pid = 0
	}
	return id.SetParentID(pid)

}
func (id CatID) Category() *Category {
	return getCategory(id)
}
func (id CatID) ShowPrice() bool {
	b, err := id.getConfig("ShowPrice")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("CatID(%d).ShowPrice() error: %s", id, err)
	}
	return b
}
func (id CatID) getConfig(key string) (val bool, err error) {
	query := `SELECT ConfigVal FROM Category_Config WHERE CatID = @0 AND ConfigKey = @1`
	var nb sql.NullBool
	err = b.db.QueryRow(query, id, key).Scan(&nb)
	if errors.Is(err, sql.ErrNoRows) {
		err = id.createConfigKey(key)
	}
	if nb.Valid {
		val = nb.Bool
	} else {
		log.Printf("CatID(%d).getConfig(%s) nb is invalid!", id, key)
	}
	// log.Printf("CatID(%d).getConfig(%s) returned %v", id, key, nb.Bool)
	return
}

/* Get the pointer to Category from map or make one and return it */
func getCategory(id CatID) *Category {
	if c := b.Metadata.categoryData[id]; c == nil {
		c = newCategory(id)
		b.Metadata.categoryData[id] = c
	}
	return b.Metadata.categoryData[id]
}
func (id CatID) getBool(key string) (val bool, err error) {
	b, err := getValue[sql.NullBool]("Category", id, key)
	if b.Valid {
		val = b.Bool
	} else {
		log.Printf("CatID(%d).getBool(%s) error: %v", id, key, err)
	}
	return
}
func (id CatID) getFloat(key string) (val float64, err error) {
	f, err := getValue[sql.NullFloat64]("Category", id, key)
	if f.Valid {
		val = f.Float64
	} else {
		log.Printf("CatID(%d).getFloat(%s) error: %v", id, key, err)
	}
	return
}
func (id CatID) getInt(key string) (val int, err error) {
	i, err := getValue[sql.NullInt64]("Category", id, key)
	val = int(i.Int64)
	if !i.Valid {
		log.Printf("CatID(%d).getInt(%s) error: %v", id, key, err)
	}
	return
}
func (id CatID) getString(key string) (val string, err error) {
	s, err := getValue[sql.NullString]("Category", id, key)
	if s.Valid {
		val = s.String
	} else {
		log.Printf("CatID(%d).getString(%s) error: %v", id, key, err)
	}
	return
}
func (id CatID) createConfigKey(key string) error {
	query := `INSERT INTO Category_Config (CatID, ConfigKey, ConfigVal) VALUES (@0, @1, false)`
	_, err := b.db.Exec(query, id, key)
	// log.Printf("createConfigKey(%s) returned %v", key, err)
	return err
}
func (id CatID) setConfig(key string, val bool) error {
	query := `UPDATE Category_Config SET ConfigVal = @0 WHERE CatID = @1 AND ConfigKey = @2`
	_, err := b.db.Exec(query, val, id, key)
	// log.Printf("CatID(%d).setConfig(%s, %v) returned %v", id, key, val, err)
	return err
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
