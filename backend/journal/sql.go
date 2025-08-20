package journal

import (
	"database/sql"
	"reflect"
	"slices"
	"time"
)

func (j *Journal) createMainTable() {
	j.db.Exec(`CREATE TABLE Journal(
EntryID INTEGER PRIMARY KEY AUTOINCREMENT,
Time TEXT DEFAULT(datetime('now', 'subsec')),
LevelID INT DEFAULT 1,
EventID INT DEFAULT 1,
Message TEXT,
FOREIGN KEY(LevelID) REFERENCES Journal_EntryLevel(LevelID),
FOREIGN KEY(EventID) REFERENCES Journal_EntryEvent(EventID))`)
}
func (j *Journal) createEntryLevelTable() {
	j.db.Exec(`CREATE TABLE Journal_EntryLevel(
LevelID INTEGER PRIMARY KEY AUTOINCREMENT,
Name TEXT UNIQUE)`)
	j.db.Exec(`INSERT INTO Journal_EntryLevel 
VALUES ("Message"), ("Warning"), ("Error")`)
}
func (j *Journal) createEntryEventTable() {
	j.db.Exec(`CREATE TABLE Journal_EntryEvent(
EventID INTEGER PRIMARY KEY AUTOINCREMENT,
Name TEXT UNIQUE)`)
	j.db.Exec(`INSERT INTO Journal_EntryEvent 
VALUES ("Log"), ("Add"), ("Copy"), ("Edit"), ("Delete"), ("SQL")`)
}
func (j *Journal) createTables() {
	tables := j.listTables()
	touched := false
	if !slices.Contains(tables, "Journal") {
		j.createMainTable()
		touched = true
	}
	if !slices.Contains(tables, "Journal_EntryLevel") {
		j.createEntryLevelTable()
		touched = true
	}
	if !slices.Contains(tables, "Journal_EntryEvent") {
		j.createEntryEventTable()
		touched = true
	}
	if len(tables) == 0 {
		j.db.Exec(`INSERT INTO Journal (EventID, Message)
VALUES ("6", "Databasen tom. Skapar nya tabeller.")`)
	}
	if touched {
		j.db.Exec(`INSERT INTO Journal (EventID, Message)
VALUES ("6", "Skapade nya tabeller för Journalen i databasen.")`)
	}
}

func (j *Journal) listTables() []string {
	var name string
	var tables []string
	stmt, err := j.db.Prepare(`SELECT name FROM sqlite_master WHERE type='table'`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&name)
		tables = append(tables, name)
	}
	return tables
}
func (j *Journal) verifyTables() bool { panic("unimplemented") }
func (j *Journal) repairTables() bool { panic("unimplemented") }

func (j *Journal) getAllEntryIDs() []any {
	var id EntryID
	var ids []any
	stmt, err := j.db.Prepare(`SELECT EntryID FROM Journal`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids
}

func (j *Journal) getRecentEntryIds() []any {
	var id EntryID
	var ids []any
	stmt, err := j.db.Prepare(`SELECT * FROM 
(SELECT EntryId FROM Journal ORDER BY EntryId DESC LIMIT @0) 
ORDER BY EntryId ` + j.config.sorting.String())
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(j.config.limit)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids
}

func (j *Journal) getEntry(id EntryID) *Entry {
	if reflect.ValueOf(j.entries[id]).IsZero() {
		j.entries[id] = j.getEntryFromSQL(id)
	}
	return j.entries[id]
}

func (j *Journal) getEntryFromSQL(id EntryID) *Entry {
	var tim sql.NullString
	entry := &Entry{}
	stmt, err := j.db.Prepare(`SELECT * FROM Journal WHERE EntryID = @0`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&entry.EntryID, &tim, &entry.Level, &entry.Event, &entry.Message)
	if err != nil {
		panic(err)
	}
	stockholm, err := time.LoadLocation("Europe/Stockholm")
	utc, err := time.Parse(subsec, tim.String)
	entry.Time = utc.In(stockholm)
	if err != nil {
		panic(err)
	}
	return entry
}

func (j *Journal) newEntry(level Level, event Event, message string) EntryID {
	stmt, err := j.db.Prepare(`INSERT INTO Journal (LevelID, EventID, Message)
VALUES (@0, @1, @2)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(level, event, message)
	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return EntryID(id)
}
