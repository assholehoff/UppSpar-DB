package backend

import (
	"UppSpar/backend/journal"
	"database/sql"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

/* Storing and applying settings */

type Setting struct {
	key   string
	value binding.String
}

func newSetting(key string) *Setting {
	return &Setting{key: key, value: binding.NewString()}
}
func (s *Setting) get() {
	query := `SELECT ConfigVal FROM Config WHERE ConfigKey = @0`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	var t sql.NullString
	err = stmt.QueryRow(s.key).Scan(&t)
	s.value.Set(t.String)
}
func (s *Setting) set() error {
	val, err := s.value.Get()
	query := `UPDATE Config SET ConfigVal = @0 WHERE ConfigKey = @1 AND ConfigVal <> @2`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(val, s.key, val)
	return err
}

type Settings struct {
	j                 *journal.Journal
	m                 map[string]*Setting
	ItemIDWidth       binding.Int
	ResumeLastSession binding.Bool
}

func NewSettings() *Settings {
	s := &Settings{
		j:                 b.Journal,
		m:                 make(map[string]*Setting),
		ItemIDWidth:       binding.NewInt(),
		ResumeLastSession: binding.NewBool(),
	}
	s.initItemIDWidth()
	s.ResumeLastSession.Set(fyne.CurrentApp().Preferences().BoolWithFallback("resume", false))
	s.ResumeLastSession.AddListener(binding.NewDataListener(func() {
		b, err := s.ResumeLastSession.Get()
		if err != nil {
			log.Println(err)
			return
		}
		fyne.CurrentApp().Preferences().SetBool("resume", b)
	}))
	return s
}

func (s *Settings) getSetting(key string) *Setting {
	// TODO debug this thing
	t := s.m[key]
	if t == nil {
		t = newSetting(key)
		t.get()
	}
	return t
}
func (s *Settings) initItemIDWidth() {
	key := "ItemIDWidth"
	s.m[key] = newSetting(key)
	s.m[key].get()
	s.ItemIDWidth = binding.StringToInt(s.m[key].value)
	s.m[key].value.AddListener(binding.NewDataListener(func() {
		s.m[key].set()
	}))
}
