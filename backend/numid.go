package backend

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

/* A NumID is an int indentifying a row in a specific SQL table */
type NumID interface {
	/* Returns the Name column for the row with the NumID from the relevant SQL table */
	Name() (string, error)
	/* Returns the name of the struct, such as "ItemID", "CatID" or "UnitID" */
	TypeName() string
	/* Returns a string with the underlying int */
	String() string
	/* Returns a tree-friendly identifier string */
	TString() string

	getBool(key string) (val bool, err error)
	getFloat(key string) (val float64, err error)
	getInt(key string) (val int, err error)
	getString(key string) (val string, err error)
	setBool(key string, val bool) error
	setFloat(key string, val float64) error
	setInt(key string, val int) error
	setString(key string, val string) error
}

/* Get value 'val' for column 'key' for row 'id' from table 't' */
func getValue[T sql.NullBool | sql.NullFloat64 | NullInt | sql.NullInt64 | sql.NullString](t string, id NumID, key string) (val T, err error) {
	query := `SELECT ` + key + ` FROM ` + t + ` WHERE ` + id.TypeName() + ` = @0`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		log.Println(strings.Replace(query, "@0", id.String(), 1))
		log.Printf("getValue[%s](%s, %d, %v) panic!", id.TypeName(), t, id, key)
		panic(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&val)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		log.Println(strings.Replace(query, "@0", id.String(), 1))
		log.Printf("getValue[%s](%s, %d, %v) panic!", id.TypeName(), t, id, key)
		panic(err)
	}
	return
}

/* Set value for column 'key' for row 'id' in table 't' to 'val' */
func setValue[T bool | float64 | int | string](t string, id NumID, key string, val T) (err error) {
	query := `UPDATE ` + t + ` SET ` + key + ` = @1 WHERE ` + id.TypeName() + ` = @2 AND ` + key + ` <> @3`
	// log.Printf("UPDATE %s SET %s = %v WHERE %s = %d AND %s <> %v", t, key, val, id.TypeName(), id, key, val)
	stmt, err := be.db.Prepare(query)
	if err != nil {
		log.Printf("setItemIDValue(%d, %s, %v) panic!", id, key, val)
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(val, id, val)
	if err != nil {
		log.Printf("setItemIDValue(%d, %s, %v) panic!", id, key, val)
		panic(err)
	}
	return
}
