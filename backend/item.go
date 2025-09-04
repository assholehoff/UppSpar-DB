package backend

import (
	"UppSpar/backend/journal"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

type SearchTermMatch int

const (
	MatchBeginsWith SearchTermMatch = iota
	MatchEndsWith
	MatchContains
	MatchEquals
	RegExp
)

type SearchKey int

const (
	SearchKeyName SearchKey = iota
	SearchKeyDesc
	SearchKeyManufacturer
	SearchKeyModel
	SearchKeyItemID
	SearchKeyDateCreated
	SearchKeyDateModified
)

func (k SearchKey) String() string {
	switch k {
	case SearchKeyItemID:
		return "ItemID"
	case SearchKeyName:
		return "Name"
	case SearchKeyDesc:
		return "Descr"
	case SearchKeyManufacturer:
		return "Manufacturer"
	case SearchKeyModel:
		return "Model"
	case SearchKeyDateCreated:
		return "DateCreated"
	case SearchKeyDateModified:
		return "DateModified"
	default:
		return ""
	}
}

type SortOrder int

const (
	SortAscending SortOrder = iota
	SortDescending
)

func (o SortOrder) String() string {
	if o == SortAscending {
		return "ASC"
	}
	return "DESC"
}

var (
	_ sql.Scanner   = (*ItemStatusID)(nil)
	_ driver.Valuer = (*ItemStatusID)(nil)
	_ fmt.Stringer  = (*ItemStatusID)(nil)
)

type ItemStatusID int

/* Returns a localized string */
func (id ItemStatusID) LString() string {
	switch id {
	case ItemStatusAvailable:
		return lang.X("itemstatus.available", "itemstatus.available")
	case ItemStatusSold:
		return lang.X("itemstatus.sold", "itemstatus.sold")
	case ItemStatusReserved:
		return lang.X("itemstatus.reserved", "itemstatus.reserved")
	case ItemStatusArchived:
		return lang.X("itemstatus.archived", "itemstatus.archived")
	case ItemStatusDeleted:
		return lang.X("itemstatus.deleted", "itemstatus.deleted")
	default:
		return lang.X("itemstatus.available", "itemstatus.available")
	}
}

/* String implements fmt.Stringer. */
func (id ItemStatusID) String() string {
	return fmt.Sprintf("%d", id)
}

/* Value implements driver.Valuer. */
func (id ItemStatusID) Value() (driver.Value, error) {
	return int64(id), nil
}

/* Scan implements sql.Scanner. */
func (id *ItemStatusID) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		*id = 0
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		*id = ItemStatusID(src.(int))
	case "int8":
		*id = ItemStatusID(src.(int8))
	case "int16":
		*id = ItemStatusID(src.(int16))
	case "int32":
		*id = ItemStatusID(src.(int32))
	case "int64":
		*id = ItemStatusID(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				*id = ItemStatusID(math.MaxInt32)
				return ErrLossyConversion
			}
		}
	case "uint":
		*id = ItemStatusID(src.(uint))
	case "uint8":
		*id = ItemStatusID(src.(uint8))
	case "uint16":
		*id = ItemStatusID(src.(uint16))
	case "uint32":
		*id = ItemStatusID(src.(uint32))
	case "uint64":
		*id = ItemStatusID(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				*id = ItemStatusID(math.MaxUint32)
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			*id = ItemStatusID(math.MaxInt64)
			return ErrLossyConversion
		}
	default:
		log.Printf("ItemStatusID(%d).Scan(%v) error: invalid type %s", id, src, reflect.TypeOf(src).Name())
		return ErrInvalidType
	}
	return nil
}

const (
	ItemStatusAvailable ItemStatusID = iota + 1
	ItemStatusSold
	ItemStatusReserved
	ItemStatusArchived
	ItemStatusDeleted
)

type Items struct {
	j    *journal.Journal
	data map[ItemID]*Item

	ItemIDList      binding.UntypedList
	ItemIDSelection binding.UntypedList
	Filter          *Filter
	Search          *Search
}

func NewItems() *Items {
	m := &Items{
		j:               b.Journal,
		data:            make(map[ItemID]*Item),
		ItemIDList:      binding.NewUntypedList(),
		ItemIDSelection: binding.NewUntypedList(),
		Filter:          newFilter(),
		Search:          newSearch(),
	}
	return m
}

func (m *Items) GetItem(id ItemID) *Item {
	if t := m.data[id]; t == nil {
		t = newItem(id)
		m.data[id] = t
	}
	return m.data[id]
}
func (m *Items) GetItemIDFor(index widget.ListItemID) (ItemID, error) {
	id, err := m.ItemIDList.GetValue(index)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return 0, err
	}
	return id.(ItemID), err
}
func (m *Items) GetListItemIDFor(id ItemID) (widget.ListItemID, error) {
	ids, err := m.ItemIDList.Get()
	if err != nil {
		panic(err)
	}
	index := slices.IndexFunc(ids, func(n any) bool { return n == id })
	if index == -1 {
		return index, ErrNotFound
	}
	return widget.ListItemID(index), nil
}
func (m *Items) CreateNewItem() (ItemID, error) {
	var i ItemID
	query := `INSERT INTO Item DEFAULT VALUES`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return i, fmt.Errorf("Items.CreateNewItem() error: %w", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec()
	if err != nil {
		return i, fmt.Errorf("Items.CreateNewItem() error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return i, fmt.Errorf("Items.CreateNewItem() error: %w", err)
	}
	i = ItemID(id)
	m.GetItemIDs()
	return i, err
}
func (m *Items) CopyItem(id ItemID) (ItemID, error) {
	query := `INSERT INTO Item (Name, Price, Currency, QuantityInPrice, Unit, 
OrderMultiple, MinOrder, Vat, Eta, EtaText, Priority, Stock, 
ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL, 
UNSPSC, LongDesc, Manufacturer, MfrItemId, GlobId, GlobIdType, 
ReplacesItem, Questions, PackagingCode, PresentationCode, 
DeliveryAutoSign, DeliveryOption, ComparePrice, CompareUnit, 
CompareQuantityInPrice, PriceInfo, AddDesc, ProcFlow, InnerUnit, 
QuantityInUnit, RiskClassification, Comment, EnvClassification, 
FormId, Article, Attachments, ItemGroup, 
MfrID, ModelID, Notes, Width, Height, Depth, Volume, Weight, 
LengthUnitID, VolumeUnitID, WeightUnitID, CatID, GroupID, 
StorageID, ItemStatusID, ItemConditionID)
SELECT Name, Price, Currency, QuantityInPrice, Unit, 
OrderMultiple, MinOrder, Vat, Eta, EtaText, Priority, Stock, 
ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL, 
UNSPSC, LongDesc, Manufacturer, MfrItemId, GlobId, GlobIdType, 
ReplacesItem, Questions, PackagingCode, PresentationCode, 
DeliveryAutoSign, DeliveryOption, ComparePrice, CompareUnit, 
CompareQuantityInPrice, PriceInfo, AddDesc, ProcFlow, InnerUnit, 
QuantityInUnit, RiskClassification, Comment, EnvClassification, 
FormId, Article, Attachments, ItemGroup, 
MfrID, ModelID, Notes, Width, Height, Depth, Volume, Weight, 
LengthUnitID, VolumeUnitID, WeightUnitID, CatID, GroupID, 
StorageID, ItemStatusID, ItemConditionID
FROM Item WHERE ItemID = @0`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return id, fmt.Errorf("CopyItem error: %w", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		return id, fmt.Errorf("CopyItem error: %w", err)
	}
	lid, _ := res.LastInsertId()
	newid := ItemID(lid)
	m.GetItemIDs()
	return newid, err
}
func (m *Items) DeleteItem(id ItemID) error {
	query := `UPDATE Item SET ItemStatusID = @0 WHERE ItemID = @1`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("DeleteItem error: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(ItemStatusDeleted, id)
	if err != nil {
		return fmt.Errorf("DeleteItem error: %w", err)
	}
	m.ItemIDList.Remove(id)
	delete(m.data, id)
	return err
}
func (m *Items) SelectItem(id ItemID) error {
	id.Item().FetchAllFields()
	return m.ItemIDSelection.Append(id)
}
func (m *Items) UnselectItem(id ItemID) error {
	return m.ItemIDSelection.Remove(id)
}
func (m *Items) ClearSelection() error {
	return m.ItemIDSelection.Set([]any{})
}
func (m *Items) GetItemIDs() {
	m.ItemIDList.Set([]any{})
	query := `SELECT ItemID FROM Item WHERE ItemID <> 0 `
	e := m.Search.complex()
	query, term := e.addSearchStrings(query)
	f := m.Filter.complex()
	query = f.addFilterStrings(query)
	query += fmt.Sprintf("AND ItemStatusID <> %d ", ItemStatusDeleted) // TODO update this
	query += "ORDER BY " + e.sortby.String() + " " + e.order.String()

	// log.Println(query)

	var err error
	var rows *sql.Rows
	switch len(e.scope) {
	case 1:
		rows, err = b.db.Query(query, term)
	case 2:
		rows, err = b.db.Query(query, term, term)
	case 3:
		rows, err = b.db.Query(query, term, term, term)
	case 4:
		rows, err = b.db.Query(query, term, term, term, term)
	default:
		rows, err = b.db.Query(query)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	defer rows.Close()

	m.ItemIDList.Set([]any{})
	m.Search.Completions.Set([]string{})
	uniqueResults := make(map[string]bool)
	for rows.Next() {
		var hit string
		var id ItemID
		rows.Scan(&id)
		m.ItemIDList.Append(id)
		if e.scope["Name"] {
			hit, _ = id.Name()
			if !uniqueResults[hit] {
				uniqueResults[hit] = true
				m.Search.Completions.Append(hit)
			}
		}
		if e.scope["Manufacturer"] {
			hit, _ = id.Manufacturer()
			if !uniqueResults[hit] {
				uniqueResults[hit] = true
				m.Search.Completions.Append(hit)
			}
		}
		// if e.scope["ModelName"] {
		// 	hit, _ = id.ModelName()
		// 	if !uniqueResults[hit] {
		// 		uniqueResults[hit] = true
		// 		m.Search.Completions.Append(hit)
		// 	}
		// }
	}
}

type Search struct {
	Completions binding.StringList
	Term        binding.String
	Scope       map[string]binding.Bool
	Match       SearchTermMatch
	SortBy      SearchKey
	Order       SortOrder
}

func newSearch() *Search {
	s := &Search{
		Completions: binding.NewStringList(),
		Term:        binding.NewString(),
		Scope:       make(map[string]binding.Bool),
		Match:       MatchContains,
		SortBy:      SearchKeyItemID,
		Order:       SortAscending,
	}
	columns := []string{"Name", "Manufacturer", "ModelName", "ModelDesc"}
	for _, key := range columns {
		s.Scope[key] = binding.NewBool()
		s.Scope[key].AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	}
	s.Scope["Name"].Set(true)
	s.Term.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	return s
}

func (e *Search) complex() searchComplex {
	c := searchComplex{
		scope: make(map[string]bool),
	}
	c.term, _ = e.Term.Get()
	c.scope["Name"], _ = e.Scope["Name"].Get()
	c.scope["Manufacturer"], _ = e.Scope["Manufacturer"].Get()
	c.scope["ModelName"], _ = e.Scope["ModelName"].Get()
	c.scope["ModelDesc"], _ = e.Scope["ModelDesc"].Get()
	c.match = e.Match
	c.sortby = e.SortBy
	c.order = e.Order
	return c
}

type searchComplex struct {
	term   string
	scope  map[string]bool
	match  SearchTermMatch
	sortby SearchKey
	order  SortOrder
}

func (e searchComplex) addSearchStrings(query string) (string, string) {
	var term string
	if e.term == "" {
		return query, term
	}
	columns := []string{"Name", "Manufacturer", "ModelName", "ModelDesc"}
	var keys []string
	for _, column := range columns {
		if e.scope[column] {
			keys = append(keys, column)
		}
	}
	if len(keys) < 1 {
		return query, term
	}
	switch e.match {
	case MatchBeginsWith:
		term = fmt.Sprintf("%s%%", e.term)
		if len(keys) > 1 {
			query += "AND ("
			for i, key := range keys {
				if i > 0 {
					query += " OR "
				}
				query += fmt.Sprintf("%s LIKE ?", key)
			}
			query += ") "
		} else {
			query += fmt.Sprintf("AND %s LIKE ? ", keys[0])
		}
	case MatchEndsWith:
		term = fmt.Sprintf("%%%s", e.term)
		if len(keys) > 1 {
			query += "AND ("
			for i, key := range keys {
				if i > 0 {
					query += " OR "
				}
				query += fmt.Sprintf("%s LIKE ?", key)
			}
			query += ") "
		} else {
			query += fmt.Sprintf("AND %s LIKE ? ", keys[0])
		}
	case MatchContains:
		term = fmt.Sprintf("%%%s%%", e.term)
		if len(keys) > 1 {
			query += "AND ("
			for i, key := range keys {
				if i > 0 {
					query += " OR "
				}
				query += fmt.Sprintf("%s LIKE ?", key)
			}
			query += ") "
		} else {
			query += fmt.Sprintf("AND %s LIKE ? ", keys[0])
		}
	default:
		// MatchEquals
		term = fmt.Sprintf("%s", e.term)
		if len(keys) > 1 {
			query += "AND ("
			for i, key := range keys {
				if i > 0 {
					query += " OR "
				}
				query += fmt.Sprintf("%s LIKE ?", key)
			}
			query += ") "
		} else {
			query += fmt.Sprintf("AND %s LIKE ? ", keys[0])
		}
	}
	return query, term
}

type Filter struct {
	Category             binding.String
	Manufacturer         binding.String
	Model                binding.String
	Width, Height, Depth binding.String
	Volume, Weight       binding.String
}

func newFilter() *Filter {
	f := &Filter{
		Category:     binding.NewString(),
		Manufacturer: binding.NewString(),
		Model:        binding.NewString(),
		Width:        binding.NewString(),
		Height:       binding.NewString(),
		Depth:        binding.NewString(),
		Volume:       binding.NewString(),
		Weight:       binding.NewString(),
	}
	f.Category.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	f.Manufacturer.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	f.Model.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	f.Width.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	f.Height.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	f.Depth.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	f.Volume.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	f.Weight.AddListener(binding.NewDataListener(func() { b.Items.GetItemIDs() }))
	return f
}

func (f Filter) complex() filterComplex {
	c := filterComplex{}
	if s, _ := f.Category.Get(); s != "" {
		c.CatID, _ = CatIDFor(s)
	}
	if s, _ := f.Manufacturer.Get(); s != "" {
		if id, err := MfrIDFor(s); id != 0 && err == nil {
			c.MfrID = id
		} else {
			c.Manufacturer = s
		}
	}
	if s, _ := f.Model.Get(); s != "" {
		if c.MfrID != 0 {
			if id, err := ModelIDFor(c.MfrID, s); id != 0 && err == nil {
				c.ModelID = id
			} else {
				c.Model = s
			}
		}
		if s, _ := f.Width.Get(); s != "" {
			if strings.Contains(s, "-") {
				t := strings.Split(s, "-")
				min, err := strconv.ParseFloat(t[0], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinWidth = min
				}
				max, err := strconv.ParseFloat(t[1], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MaxWidth = max
				}
			} else {
				w, err := strconv.ParseFloat(s, 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinWidth = w
					c.MaxWidth = w
				}
			}
		}
		if s, _ := f.Height.Get(); s != "" {
			if strings.Contains(s, "-") {
				t := strings.Split(s, "-")
				min, err := strconv.ParseFloat(t[0], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinHeight = min
				}
				max, err := strconv.ParseFloat(t[1], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MaxHeight = max
				}
			} else {
				w, err := strconv.ParseFloat(s, 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinHeight = w
					c.MaxHeight = w
				}
			}
		}
		if s, _ := f.Depth.Get(); s != "" {
			if strings.Contains(s, "-") {
				t := strings.Split(s, "-")
				min, err := strconv.ParseFloat(t[0], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinDepth = min
				}
				max, err := strconv.ParseFloat(t[1], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MaxDepth = max
				}
			} else {
				w, err := strconv.ParseFloat(s, 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinDepth = w
					c.MaxDepth = w
				}
			}
		}
		if s, _ := f.Volume.Get(); s != "" {
			if strings.Contains(s, "-") {
				t := strings.Split(s, "-")
				min, err := strconv.ParseFloat(t[0], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinVolume = min
				}
				max, err := strconv.ParseFloat(t[1], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MaxVolume = max
				}
			} else {
				w, err := strconv.ParseFloat(s, 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinVolume = w
					c.MaxVolume = w
				}
			}
		}
		if s, _ := f.Weight.Get(); s != "" {
			if strings.Contains(s, "-") {
				t := strings.Split(s, "-")
				min, err := strconv.ParseFloat(t[0], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinWeight = min
				}
				max, err := strconv.ParseFloat(t[1], 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MaxWeight = max
				}
			} else {
				w, err := strconv.ParseFloat(s, 64)
				if err != nil {
					log.Println(err)
				} else {
					c.MinWeight = w
					c.MaxWeight = w
				}
			}
		}
	}
	return c
}

type filterComplex struct {
	CatID                CatID
	MfrID                MfrID
	ModelID              ModelID
	Manufacturer         string
	Model                string
	MinWidth, MaxWidth   float64
	MinHeight, MaxHeight float64
	MinDepth, MaxDepth   float64
	MinVolume, MaxVolume float64
	MinWeight, MaxWeight float64
}

func (f filterComplex) addFilterStrings(query string) string {
	if f.CatID != 0 {
		query += fmt.Sprintf("AND CatID = %d ", f.CatID)
	}
	if f.MfrID != 0 {
		query += fmt.Sprintf("AND MfrID = %d ", f.MfrID)
	} else if f.Manufacturer != "" {
		query += fmt.Sprintf("AND Model = '%s' ", f.Model)
	}
	if f.ModelID != 0 {
		query += fmt.Sprintf("AND ModelID = %d ", f.ModelID)
	} else if f.Model != "" {
		query += fmt.Sprintf("AND Manufacturer = '%s' ", f.Manufacturer)
	}
	if f.MinWidth != 0 || f.MaxWidth != 0 {
		query += fmt.Sprintf("AND Width BETWEEN %f AND %f ", f.MinWidth, f.MaxWidth)
	}
	if f.MinHeight != 0 || f.MaxHeight != 0 {
		query += fmt.Sprintf("AND Height BETWEEN %f AND %f ", f.MinHeight, f.MaxHeight)
	}
	if f.MinDepth != 0 || f.MaxDepth != 0 {
		query += fmt.Sprintf("AND Depth BETWEEN %f AND %f ", f.MinDepth, f.MaxDepth)
	}
	if f.MinVolume != 0 || f.MaxVolume != 0 {
		query += fmt.Sprintf("AND Volume BETWEEN %f AND %f ", f.MinVolume, f.MaxVolume)
	}
	if f.MinWeight != 0 || f.MaxWeight != 0 {
		query += fmt.Sprintf("AND Weight BETWEEN %f AND %f ", f.MinWeight, f.MaxWeight)
	}
	return query
}

type Item struct {
	binding.DataItem
	ItemID       ItemID
	ItemIDString binding.String
	Name         binding.String
	CatID        CatID
	Category     binding.String
	priceFloat   binding.Float
	PriceString  binding.String
	Currency     binding.String
	Unit         binding.String
	vatFloat     binding.Float
	VatString    binding.String
	Priority     binding.Bool
	stockFloat   binding.Float
	StockString  binding.String
	SearchWords  binding.StringList // TODO implement this
	ImgURL1      binding.String
	ImgURL2      binding.String
	ImgURL3      binding.String
	ImgURL4      binding.String
	ImgURL5      binding.String
	SpecsURL     binding.String
	AddDesc      binding.String
	LongDesc     binding.String
	MfrID        MfrID
	Manufacturer binding.String
	ModelID      ModelID
	ModelName    binding.String
	ModelDesc    binding.String
	ModelURL     binding.String
	Notes        binding.String
	widthFloat   binding.Float
	heightFloat  binding.Float
	depthFloat   binding.Float
	volumeFloat  binding.Float
	weightFloat  binding.Float
	WidthString  binding.String
	HeightString binding.String
	DepthString  binding.String
	VolumeString binding.String
	WeightString binding.String
	LengthUnit   binding.String
	VolumeUnit   binding.String
	WeightUnit   binding.String
	ItemStatus   binding.String
	DateCreated  binding.String
	DateModified binding.String
}

func newItem(id ItemID) *Item {
	t := &Item{
		ItemID:  id,
		CatID:   CatID(0),
		MfrID:   MfrID(0),
		ModelID: ModelID(0),
	}

	t.ItemIDString = binding.NewString()
	t.Name = binding.NewString()
	t.Category = binding.NewString()
	t.priceFloat = binding.NewFloat()
	t.Currency = binding.NewString()
	t.Unit = binding.NewString()
	t.vatFloat = binding.NewFloat()
	t.Priority = binding.NewBool()
	t.stockFloat = binding.NewFloat()
	t.ImgURL1 = binding.NewString()
	t.ImgURL2 = binding.NewString()
	t.ImgURL3 = binding.NewString()
	t.ImgURL4 = binding.NewString()
	t.ImgURL5 = binding.NewString()
	t.SpecsURL = binding.NewString()
	t.AddDesc = binding.NewString()
	t.LongDesc = binding.NewString()
	t.Manufacturer = binding.NewString()
	t.ModelName = binding.NewString()
	t.ModelDesc = binding.NewString()
	t.ModelURL = binding.NewString()
	t.Notes = binding.NewString()
	t.widthFloat = binding.NewFloat()
	t.heightFloat = binding.NewFloat()
	t.depthFloat = binding.NewFloat()
	t.weightFloat = binding.NewFloat()
	t.volumeFloat = binding.NewFloat()
	t.LengthUnit = binding.NewString()
	t.VolumeUnit = binding.NewString()
	t.WeightUnit = binding.NewString()
	t.ItemStatus = binding.NewString()

	t.DateCreated = binding.NewString()
	t.DateModified = binding.NewString()

	t.FetchAllFields()

	t.PriceString = binding.FloatToStringWithFormat(t.priceFloat, "%.2f")
	t.VatString = binding.FloatToStringWithFormat(t.vatFloat, "%.2f")
	t.StockString = binding.FloatToStringWithFormat(t.stockFloat, "%.0f")
	t.WidthString = binding.FloatToStringWithFormat(t.widthFloat, "%.0f")
	t.HeightString = binding.FloatToStringWithFormat(t.heightFloat, "%.0f")
	t.DepthString = binding.FloatToStringWithFormat(t.depthFloat, "%.0f")
	t.VolumeString = binding.FloatToStringWithFormat(t.volumeFloat, "%.2f")
	t.WeightString = binding.FloatToStringWithFormat(t.weightFloat, "%.2f")

	t.Name.AddListener(binding.NewDataListener(func() { t.ItemID.SetName(); b.Items.GetItemIDs(); t.ItemID.CompileLongDesc() }))
	t.Category.AddListener(binding.NewDataListener(func() { t.ItemID.SetCategory(); t.ItemID.CompileLongDesc() }))
	t.priceFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetPrice(); t.ItemID.CompileLongDesc() }))
	t.Currency.AddListener(binding.NewDataListener(func() { t.ItemID.SetCurrency(); t.ItemID.CompileLongDesc() }))
	t.Unit.AddListener(binding.NewDataListener(func() { t.ItemID.SetUnit(); t.ItemID.CompileLongDesc() }))
	t.vatFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetVat(); t.ItemID.CompileLongDesc() }))
	t.Priority.AddListener(binding.NewDataListener(func() { t.ItemID.SetPriority(); t.ItemID.CompileLongDesc() }))
	t.stockFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetStock(); t.ItemID.CompileLongDesc() }))
	t.ImgURL1.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL1(); t.ItemID.CompileLongDesc() }))
	t.ImgURL2.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL2(); t.ItemID.CompileLongDesc() }))
	t.ImgURL3.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL3(); t.ItemID.CompileLongDesc() }))
	t.ImgURL4.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL4(); t.ItemID.CompileLongDesc() }))
	t.ImgURL5.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL5(); t.ItemID.CompileLongDesc() }))
	t.SpecsURL.AddListener(binding.NewDataListener(func() { t.ItemID.SetSpecsURL(); t.ItemID.CompileLongDesc() }))
	t.AddDesc.AddListener(binding.NewDataListener(func() { t.ItemID.SetAddDesc(); t.ItemID.CompileLongDesc() }))
	t.LongDesc.AddListener(binding.NewDataListener(func() { t.ItemID.SetLongDesc(); t.ItemID.CompileLongDesc() }))
	t.Manufacturer.AddListener(binding.NewDataListener(func() { t.ItemID.SetManufacturer(); t.ItemID.CompileLongDesc() }))
	t.ModelName.AddListener(binding.NewDataListener(func() { t.ItemID.SetModelName(); t.FetchAllFields(); t.ItemID.CompileLongDesc() }))
	t.ModelURL.AddListener(binding.NewDataListener(func() { t.ItemID.SetModelURL(); t.ItemID.CompileLongDesc() }))
	t.Notes.AddListener(binding.NewDataListener(func() { t.ItemID.SetNotes(); t.ItemID.CompileLongDesc() }))
	t.widthFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetWidth(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))
	t.heightFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetHeight(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))
	t.depthFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetDepth(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))
	t.volumeFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetVolume(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))
	t.weightFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetWeight(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))
	t.LengthUnit.AddListener(binding.NewDataListener(func() { t.ItemID.SetLengthUnit(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))
	t.VolumeUnit.AddListener(binding.NewDataListener(func() { t.ItemID.SetVolumeUnit(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))
	t.WeightUnit.AddListener(binding.NewDataListener(func() { t.ItemID.SetWeightUnit(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))
	t.ItemStatus.AddListener(binding.NewDataListener(func() { t.ItemID.SetItemStatus(); t.ItemID.CompileAddDesc(); t.ItemID.CompileLongDesc() }))

	// TODO implement SearchWords

	return t
}

func (t *Item) Bindings() map[string]binding.String {
	m := make(map[string]binding.String)
	m["ItemIDString"] = t.ItemIDString
	m["Name"] = t.Name
	m["Category"] = t.Category
	m["PriceString"] = t.PriceString
	m["Currency"] = t.Currency
	m["Unit"] = t.Unit
	m["VatString"] = t.VatString
	m["StockString"] = t.StockString
	m["ImgURL1"] = t.ImgURL1
	m["ImgURL2"] = t.ImgURL2
	m["ImgURL3"] = t.ImgURL3
	m["ImgURL4"] = t.ImgURL4
	m["ImgURL5"] = t.ImgURL5
	m["SpecsURL"] = t.SpecsURL
	m["AddDesc"] = t.AddDesc
	m["LongDesc"] = t.LongDesc
	m["Manufacturer"] = t.Manufacturer
	m["ModelName"] = t.ModelName
	m["ModelDesc"] = t.ModelDesc
	m["ModelURL"] = t.ModelURL
	m["Notes"] = t.Notes
	m["WidthString"] = t.WidthString
	m["HeightString"] = t.HeightString
	m["DepthString"] = t.DepthString
	m["VolumeString"] = t.VolumeString
	m["WeightString"] = t.WeightString
	m["LengthUnit"] = t.LengthUnit
	m["VolumeUnit"] = t.VolumeUnit
	m["WeightUnit"] = t.WeightUnit
	m["ItemStatus"] = t.ItemStatus
	m["DateCreated"] = t.DateCreated
	m["DateModified"] = t.DateModified
	return m
}

func (t *Item) FetchAllFields() error {
	var Name, Currency, Unit, ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL sql.NullString
	var AddDesc, LongDesc, Manufacturer, ModelName, ModelDesc, ModelURL, Notes, DateCreated, DateModified sql.NullString
	var Price, QuantityInPrice, Vat, Stock, Width, Height, Depth, Volume, Weight sql.NullFloat64
	var Priority sql.NullBool
	var CatID CatID
	var MfrID MfrID
	var ModelID ModelID
	var LengthUnitID, VolumeUnitID, WeightUnitID UnitID
	var ItemStatusID ItemStatusID

	query := `SELECT 
Name, CatID, Price, Currency, QuantityInPrice, Unit, Vat, 
Priority, Stock, ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL, 
AddDesc, LongDesc, Manufacturer, MfrID, ModelID, ModelName, ModelDesc, ModelURL, Notes, 
Width, Height, Depth, Volume, Weight, 
LengthUnitID, VolumeUnitID, WeightUnitID, 
ItemStatusID, DateCreated, DateModified 
FROM Item WHERE ItemID = @0`
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("get all item fields error: %w", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.ItemID).Scan(
		&Name, &CatID, &Price, &Currency, &QuantityInPrice, &Unit, &Vat,
		&Priority, &Stock, &ImgURL1, &ImgURL2, &ImgURL3, &ImgURL4, &ImgURL5, &SpecsURL,
		&AddDesc, &LongDesc, &Manufacturer, &MfrID, &ModelID, &ModelName, &ModelDesc, &ModelURL, &Notes,
		&Width, &Height, &Depth, &Volume, &Weight,
		&LengthUnitID, &VolumeUnitID, &WeightUnitID,
		&ItemStatusID, &DateCreated, &DateModified,
	)
	if err != nil {
		return fmt.Errorf("get all item fields error: %w", err)
	}

	var category, manufacturer, model, modelUrl, modelDesc string

	category, err = CatID.Name()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get all item fields error: %w", err)
	}

	manufacturer = Manufacturer.String
	model = ModelName.String
	modelUrl = ModelURL.String
	modelDesc = ModelDesc.String

	width := Width.Float64
	height := Height.Float64
	depth := Depth.Float64
	volume := Volume.Float64
	weight := Weight.Float64

	LengthUnitString, err := LengthUnitID.Name()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get all item fields error: %w", err)
	}
	VolumeUnitString, err := VolumeUnitID.Name()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get all item fields error: %w", err)
	}
	WeightUnitString, err := WeightUnitID.Name()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get all item fields error: %w", err)
	}
	ItemStatusString := ItemStatusID.LString()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get all item fields error: %w", err)
	}

	if MfrID != 0 {
		if n, _ := MfrID.Name(); manufacturer == "" && n != manufacturer {
			// log.Printf("MfrID is set to %d", MfrID)
			manufacturer = n
		}
	} else {
		// log.Printf("MfrID is unset")
	}

	if ModelID != 0 {
		// log.Printf("ModelID is set to %d", ModelID)
		if n, _ := ModelID.Name(); model == "" && n != model {
			model = n
		}
		modelCatID, err := ModelID.CatID()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get all item fields error: %w", err)
		}
		if modelCategory, _ := modelCatID.Name(); category == "" && modelCategory != category {
			category = modelCategory
		}
		modelMfrID, err := ModelID.MfrID()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get all item fields error: %w", err)
		}
		if t.MfrID == 0 && modelMfrID != t.MfrID {
			t.MfrID = modelMfrID
		}
		if n, _ := modelMfrID.Name(); manufacturer == "" && n != manufacturer {
			manufacturer = n
		}
		if u, _ := ModelID.ModelURL(); modelUrl == "" && u != modelUrl {
			modelUrl = u
		}
		if d, _ := ModelID.Desc(); modelDesc == "" && d != modelDesc {
			modelDesc = d
		}
		if weight == 0 && height == 0 && depth == 0 {
			if w, _ := ModelID.Weight(); w != weight {
				weight = w
			}
			if h, _ := ModelID.Height(); h != height {
				height = h
			}
			if d, _ := ModelID.Depth(); d != depth {
				depth = d
			}
			if u, _ := ModelID.LengthUnitID(); u != LengthUnitID {
				LengthUnitID = u
			}
		}
		if volume == 0 {
			if v, _ := ModelID.Volume(); v != volume {
				volume = v
			}
			if u, _ := ModelID.VolumeUnitID(); u != VolumeUnitID {
				VolumeUnitID = u
			}
		}
		if weight == 0 {
			if w, _ := ModelID.Weight(); w != weight {
				weight = w
			}
			if u, _ := ModelID.WeightUnitID(); u != WeightUnitID {
				WeightUnitID = u
			}
		}
	} else {
		// log.Printf("ModelID is unset")
	}

	t.ItemIDString.Set(t.ItemID.String())
	t.Name.Set(Name.String)
	t.Category.Set(category)
	t.priceFloat.Set(Price.Float64)
	t.Currency.Set(Currency.String)
	t.Unit.Set(Unit.String)
	t.vatFloat.Set(Vat.Float64)
	t.Priority.Set(Priority.Bool)
	t.stockFloat.Set(Stock.Float64)
	t.ImgURL1.Set(ImgURL1.String)
	t.ImgURL2.Set(ImgURL2.String)
	t.ImgURL3.Set(ImgURL3.String)
	t.ImgURL4.Set(ImgURL4.String)
	t.ImgURL5.Set(ImgURL5.String)
	t.SpecsURL.Set(SpecsURL.String)
	t.AddDesc.Set(AddDesc.String)
	t.LongDesc.Set(LongDesc.String)
	t.Manufacturer.Set(manufacturer)
	t.ModelName.Set(model)
	t.ModelDesc.Set(modelDesc)
	t.ModelURL.Set(modelUrl)
	t.Notes.Set(Notes.String)
	t.widthFloat.Set(width)
	t.heightFloat.Set(height)
	t.depthFloat.Set(depth)
	t.volumeFloat.Set(volume)
	t.weightFloat.Set(weight)
	t.LengthUnit.Set(LengthUnitString)
	t.VolumeUnit.Set(VolumeUnitString)
	t.WeightUnit.Set(WeightUnitString)
	t.ItemStatus.Set(ItemStatusString)

	var created, modified time.Time
	stockholm, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		log.Printf("time.LoadLocation error: %s", err)
	}
	if DateCreated.Valid {
		utc, err := time.Parse(subsec, DateCreated.String)
		created = utc.In(stockholm)
		if err != nil {
			t.DateCreated.Set(fmt.Sprintf("error parsing DateCreated: %v", err))
		} else {
			t.DateCreated.Set(created.Format(time.DateTime))
		}
	}
	if DateModified.Valid {
		utc, err := time.Parse(subsec, DateCreated.String)
		modified = utc.In(stockholm)
		if err != nil {
			t.DateModified.Set(fmt.Sprintf("error parsing DateModified: %v", err))
		} else {
			t.DateModified.Set(modified.Format(time.DateTime))
		}
	}
	return nil
}
