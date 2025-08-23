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
	"strings"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

type SearchType int

const (
	BeginsWith SearchType = iota
	EndsWith
	Contains
	Equals
	RegExp
)

type SearchKey int

const (
	SearchKeyName SearchKey = iota
	SearchKeyManufacturer
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
	case SearchKeyManufacturer:
		return "Manufacturer"
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
func (is ItemStatusID) LString() string {
	switch is {
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
func (is ItemStatusID) String() string {
	return fmt.Sprintf("%d", is)
}

/* Value implements driver.Valuer. */
func (is ItemStatusID) Value() (driver.Value, error) {
	return int64(is), nil
}

/* Scan implements sql.Scanner. */
func (is *ItemStatusID) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		*is = 0
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		*is = ItemStatusID(src.(int))
	case "int8":
		*is = ItemStatusID(src.(int8))
	case "int16":
		*is = ItemStatusID(src.(int16))
	case "int32":
		*is = ItemStatusID(src.(int32))
	case "int64":
		*is = ItemStatusID(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				*is = ItemStatusID(math.MaxInt32)
				return ErrLossyConversion
			}
		}
	case "uint":
		*is = ItemStatusID(src.(uint))
	case "uint8":
		*is = ItemStatusID(src.(uint8))
	case "uint16":
		*is = ItemStatusID(src.(uint16))
	case "uint32":
		*is = ItemStatusID(src.(uint32))
	case "uint64":
		*is = ItemStatusID(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				*is = ItemStatusID(math.MaxUint32)
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			*is = ItemStatusID(math.MaxInt64)
			return ErrLossyConversion
		}
	default:
		log.Printf("ItemStatus.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
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

/* Item and Bindings */

type ItemID int

var (
	_ sql.Scanner   = (*ItemID)(nil)
	_ driver.Valuer = (*ItemID)(nil)
	_ fmt.Stringer  = (*ItemID)(nil)
)

func (id *ItemID) Scan(src any) error {
	if !reflect.ValueOf(src).IsValid() {
		*id = 0
		return nil
	}
	switch reflect.TypeOf(src).Name() {
	case "int":
		*id = ItemID(src.(int))
	case "int8":
		*id = ItemID(src.(int8))
	case "int16":
		*id = ItemID(src.(int16))
	case "int32":
		*id = ItemID(src.(int32))
	case "int64":
		*id = ItemID(src.(int64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(int64) > math.MaxInt32 {
				*id = ItemID(math.MaxInt32)
				return ErrLossyConversion
			}
		}
	case "uint":
		*id = ItemID(src.(uint))
	case "uint8":
		*id = ItemID(src.(uint8))
	case "uint16":
		*id = ItemID(src.(uint16))
	case "uint32":
		*id = ItemID(src.(uint32))
	case "uint64":
		*id = ItemID(src.(uint64))
		if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
			if src.(uint64) > math.MaxUint32 {
				*id = ItemID(math.MaxUint32)
				return ErrLossyConversion
			}
		}
		if src.(uint64) > math.MaxInt64 {
			*id = ItemID(math.MaxInt64)
			return ErrLossyConversion
		}
	default:
		log.Printf("ItemID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return ErrInvalidType
	}
	return nil
}
func (id ItemID) Value() (driver.Value, error) {
	return int64(id), nil
}
func (id ItemID) String() string {
	return fmt.Sprintf("%0*d", ItemIDWidth(), id)
}
func (id ItemID) Int() int {
	return int(id)
}
func (id ItemID) Item() *Item {
	return getItem(be, id)
}

/* For the Journal */
func (id ItemID) Touch() {
	id.Item().touched = true
}
func (id ItemID) Touched() bool {
	return id.Item().touched
}
func (id ItemID) Untouch() {
	id.Item().touched = false
}

/* Returning data */

func (id ItemID) Name() (string, error) {
	return id.getString("Name")
}
func (id ItemID) CatID() (CatID, error) {
	cid, err := id.getInt("CatID")
	return CatID(cid), err
}
func (id ItemID) Category() (string, error) {
	return id.Item().CatID.Name()
}
func (id ItemID) Price() (float64, error) {
	return id.getFloat("Price")
}
func (id ItemID) Currency() (string, error) {
	return id.getString("Currency")
}
func (id ItemID) Unit() (string, error) {
	return id.getString("Unit")
}
func (id ItemID) Vat() (float64, error) {
	return id.getFloat("Vat")
}
func (id ItemID) Priority() (bool, error) {
	return id.getBool("Priority")
}
func (id ItemID) Stock() (float64, error) {
	return id.getFloat("Stock")
}
func (id ItemID) ImgURL1() (string, error) {
	return id.getString("ImgURL1")
}
func (id ItemID) ImgURL2() (string, error) {
	return id.getString("ImgURL2")
}
func (id ItemID) ImgURL3() (string, error) {
	return id.getString("ImgURL3")
}
func (id ItemID) ImgURL4() (string, error) {
	return id.getString("ImgURL4")
}
func (id ItemID) ImgURL5() (string, error) {
	return id.getString("ImgURL5")
}
func (id ItemID) SpecsURL() (string, error) {
	return id.getString("SpecsURL")
}
func (id ItemID) AddDesc() (string, error) {
	return id.getString("AddDesc")
}
func (id ItemID) LongDesc() (string, error) {
	return id.getString("LongDesc")
}
func (id ItemID) Manufacturer() (string, error) {
	return id.getString("Manufacturer")
}
func (id ItemID) Model() (string, error) {
	return id.getString("Model")
}
func (id ItemID) ModelURL() (string, error) {
	return id.getString("ModelURL")
}
func (id ItemID) Notes() (string, error) {
	return id.getString("Notes")
}
func (id ItemID) Width() (float64, error) {
	return id.getFloat("Width")
}
func (id ItemID) Height() (float64, error) {
	return id.getFloat("Height")
}
func (id ItemID) Depth() (float64, error) {
	return id.getFloat("Depth")
}
func (id ItemID) Volume() (float64, error) {
	return id.getFloat("Volume")
}
func (id ItemID) Weight() (float64, error) {
	return id.getFloat("Weight")
}
func (id ItemID) LengthUnit() (string, error) {
	uid, err := id.getInt("LengthUnitID")
	return UnitID(uid).String(), err
}
func (id ItemID) WeightUnit() (string, error) {
	uid, err := id.getInt("WeightUnitID")
	return UnitID(uid).String(), err
}
func (id ItemID) VolumeUnit() (string, error) {
	uid, err := id.getInt("VolumeUnitID")
	return UnitID(uid).String(), err
}
func (id ItemID) LengthUnitID() (UnitID, error) {
	uid, err := id.getInt("LengthUnitID")
	return UnitID(uid), err
}
func (id ItemID) VolumeUnitID() (UnitID, error) {
	uid, err := id.getInt("VolumeUnitID")
	return UnitID(uid), err
}
func (id ItemID) WeightUnitID() (UnitID, error) {
	uid, err := id.getInt("WeightUnitID")
	return UnitID(uid), err
}
func (id ItemID) ItemStatus() string {
	is, err := id.getInt("ItemStatusID")
	if err != nil {
		log.Println(err)
		return ""
	}
	return ItemStatusID(is).LString()
}
func (id ItemID) ItemStatusID() (ItemStatusID, error) {
	is, err := id.getInt("ItemStatusID")
	return ItemStatusID(is), err
}
func (id ItemID) DateCreated() (t time.Time, err error) {
	ts, err := id.getString("DateCreated")
	utc, err := time.Parse(subsec, ts)
	stockholm, err := time.LoadLocation("Europe/Stockholm")
	t = utc.In(stockholm)
	return
}
func (id ItemID) DateModified() (t time.Time, err error) {
	ts, err := id.getString("DateModified")
	utc, err := time.Parse(subsec, ts)
	stockholm, err := time.LoadLocation("Europe/Stockholm")
	t = utc.In(stockholm)
	return
}

func (id ItemID) getBool(key string) (val bool, err error) {
	b, err := getItemIDValue[sql.NullBool](&id, key)
	if b.Valid {
		val = b.Bool
	} else {
		log.Printf("getBool(%s) b is invalid (NULL), err is %v", key, err)
		err = ErrSQLNullValue
	}
	return
}
func (id ItemID) getFloat(key string) (val float64, err error) {
	f, err := getItemIDValue[sql.NullFloat64](&id, key)
	if f.Valid {
		val = f.Float64
	} else {
		log.Printf("getFloat(%s) %s is invalid (NULL), err is %v", key, key, err)
		err = ErrSQLNullValue
	}
	return
}
func (id ItemID) getInt(key string) (val int, err error) {
	i, err := getItemIDValue[sql.NullInt64](&id, key)
	val = int(i.Int64)
	if !i.Valid {
		log.Printf("getInt(%s) %s is invalid (NULL), err is %v", key, key, err)
		err = ErrSQLNullValue
	}
	return
}
func (id ItemID) getString(key string) (val string, err error) {
	s, err := getItemIDValue[sql.NullString](&id, key)
	if s.Valid {
		val = s.String
	} else {
		log.Printf("getInt(%s) %s is invalid (NULL), err is %v", key, key, err)
		err = ErrSQLNullValue
	}
	return
}
func getItemIDValue[T sql.NullBool | sql.NullFloat64 | NullInt | sql.NullInt64 | sql.NullString](id *ItemID, key string) (val T, err error) {
	query := `SELECT ` + key + ` FROM Item WHERE ItemID = @1`
	stmt, err := be.db.Prepare(query)
	if err != nil {
		log.Printf("getItemIDValue(%d, %s) panic!", id, key)
		panic(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&val)
	if err != nil {
		log.Printf("getItemIDValue(%d, %s) panic!", id, key)
		panic(err)
	}
	return
}

/* Compiling data from multiple cells */

func (id ItemID) CompileAddDesc() error {
	var addDesc, u string
	var w, h, d, v float64
	w, _ = id.Width()
	h, _ = id.Height()
	d, _ = id.Depth()
	u, _ = id.LengthUnit()
	if w > 0 || h > 0 || d > 0 {
		addDesc += fmt.Sprintf("Mått: %.0fx%.0fx%.0f %s\n", w, h, d, u)
	}
	v, _ = id.Volume()
	u, _ = id.VolumeUnit()
	if v > 0 {
		addDesc += fmt.Sprintf("Volym: %.2f %s\n", v, u)
	}
	w, _ = id.Weight()
	u, _ = id.WeightUnit()
	if w > 0 {
		addDesc += fmt.Sprintf("Vikt: %.2f %s\n", w, u)
	}
	n, _ := id.Notes()
	if n != "" {
		addDesc += fmt.Sprintf("Anmärkningar: %s\n", n)
	}
	id.Item().AddDesc.Set(addDesc)
	return id.SetAddDesc()
}

func (id ItemID) CompileLongDesc() error {

}

/* Updating data */

func (id ItemID) SetName() error {
	key := "Name"
	val, err := id.Item().Name.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetName() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetCatID(val CatID) error {
	key := "CatID"
	return id.setInt(key, int(val))
}
func (id ItemID) SetCategory() error {
	s, err := id.Item().Category.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetCategory() error: %w", err)
	}
	n, err := CatIDFor(s)
	if err != nil {
		return fmt.Errorf("ItemID.SetCategory() error: %w", err)
	}
	id.Item().CatID = n
	return id.SetCatID(n)
}
func (id ItemID) SetPrice() error {
	key := "Price"
	val, err := id.Item().priceFloat.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetPrice() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ItemID) SetCurrency() error {
	key := "Currency"
	val, err := id.Item().Currency.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetCurrency() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetUnit() error {
	key := "Unit"
	val, err := id.Item().Unit.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetUnit() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetVat() error {
	key := "Vat"
	val, err := id.Item().vatFloat.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetVat() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ItemID) SetPriority() error {
	key := "Priority"
	val, err := id.Item().Priority.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetPriority() error: %w", err)
	}
	return id.setBool(key, val)
}
func (id ItemID) SetStock() error {
	key := "Vat"
	val, err := id.Item().stockFloat.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetStock() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ItemID) SetImgURL1() error {
	key := "ImgURL1"
	val, err := id.Item().ImgURL1.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetImgURL1() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetImgURL2() error {
	key := "ImgURL2"
	val, err := id.Item().ImgURL2.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetImgURL2() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetImgURL3() error {
	key := "ImgURL3"
	val, err := id.Item().ImgURL3.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetImgURL3() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetImgURL4() error {
	key := "ImgURL4"
	val, err := id.Item().ImgURL4.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetImgURL4() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetImgURL5() error {
	key := "ImgURL5"
	val, err := id.Item().ImgURL5.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetImgURL5() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetSpecsURL() error {
	key := "SpecsURL"
	val, err := id.Item().SpecsURL.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetSpecsURL() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetAddDesc() error {
	key := "AddDesc"
	val, err := id.Item().AddDesc.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetAddDesc() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetLongDesc() error {
	key := "LongDesc"
	val, err := id.Item().LongDesc.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetLongDesc() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetManufacturer() error {
	key := "Manufacturer"
	val, err := id.Item().Manufacturer.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetManufacturer() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetModel() error {
	key := "Model"
	val, err := id.Item().Model.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetModel() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetModelURL() error {
	key := "ModelURL"
	val, err := id.Item().ModelURL.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetModelURL() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetNotes() error {
	key := "Notes"
	val, err := id.Item().Notes.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetNotes() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ItemID) SetWidth() error {
	key := "Width"
	val, err := id.Item().widthFloat.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetWidth() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ItemID) SetHeight() error {
	key := "Height"
	val, err := id.Item().heightFloat.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetHeight() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ItemID) SetDepth() error {
	key := "Depth"
	val, err := id.Item().depthFloat.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetDepth() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ItemID) SetVolume() error {
	key := "Volume"
	val, err := id.Item().volumeFloat.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetVolume() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ItemID) SetWeight() error {
	key := "Weight"
	val, err := id.Item().weightFloat.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetWeight() error: %w", err)
	}
	return id.setFloat(key, val)
}

func (id ItemID) SetLengthUnit() error {
	str, err := id.Item().LengthUnit.Get()
	if err != nil {
		return fmt.Errorf("SetLengthUnit error: %w", err)
	}
	switch str {
	case "mm":
		return id.SetLengthUnitID(millimeter)
	case "cm":
		return id.SetLengthUnitID(centimeter)
	case "dm":
		return id.SetLengthUnitID(decimeter)
	case "m":
		return id.SetLengthUnitID(meter)
	default:
		return fmt.Errorf("invalid length UnitID")
	}
}
func (id ItemID) SetLengthUnitID(l UnitID) error {
	key := "LengthUnitID"
	val := int(l)
	return id.setInt(key, val)
}
func (id ItemID) SetVolumeUnit() error {
	str, err := id.Item().VolumeUnit.Get()
	if err != nil {
		return fmt.Errorf("SetVolumeUnit error: %w", err)
	}
	switch str {
	case "ml":
		return id.SetVolumeUnitID(milliliter)
	case "cl":
		return id.SetVolumeUnitID(centiliter)
	case "dl":
		return id.SetVolumeUnitID(deciliter)
	case "l":
		return id.SetVolumeUnitID(liter)
	default:
		return fmt.Errorf("invalid volume UnitID")
	}
}
func (id ItemID) SetVolumeUnitID(v UnitID) error {
	key := "VolumeUnitID"
	val := int(v)
	return id.setInt(key, val)
}

func (id ItemID) SetWeightUnit() error {
	str, err := id.Item().WeightUnit.Get()
	if err != nil {
		return fmt.Errorf("SetWeightUnit error: %w", err)
	}
	switch str {
	case "g":
		return id.SetWeightUnitID(gram)
	case "hg":
		return id.SetWeightUnitID(hectogram)
	case "kg":
		return id.SetWeightUnitID(kilogram)
	default:
		return fmt.Errorf("invalid weight UnitID")
	}
}
func (id ItemID) SetWeightUnitID(w UnitID) error {
	key := "WeightUnitID"
	val := int(w)
	return id.setInt(key, val)
}

func (id ItemID) SetItemStatus() error {
	str, err := id.Item().ItemStatus.Get()
	if err != nil {
		return fmt.Errorf("SetItemStatus error: %w", err)
	}
	switch str {
	case lang.X("itemstatus.available", "itemstatus.available"):
		return id.SetItemStatusID(ItemStatusAvailable)
	case lang.X("itemstatus.archived", "itemstatus.archived"):
		return id.SetItemStatusID(ItemStatusArchived)
	case lang.X("itemstatus.deleted", "itemstatus.deleted"):
		return id.SetItemStatusID(ItemStatusDeleted)
	case lang.X("itemstatus.reserved", "itemstatus.reserved"):
		return id.SetItemStatusID(ItemStatusReserved)
	case lang.X("itemstatus.sold", "itemstatus.sold"):
		return id.SetItemStatusID(ItemStatusSold)
	default:
		return id.SetItemStatusID(ItemStatusAvailable)
	}
}

/* Set ItemStatusID based on contents of ItemID.Item().ItemStatus string */
func (id ItemID) SetItemStatusID(t ItemStatusID) error {
	key := "ItemStatusID"
	val := int(t)
	return id.setInt(key, val)
}

func (id ItemID) setBool(key string, val bool) error {
	return setItemIDValue(id, key, val)
}
func (id ItemID) setFloat(key string, val float64) error {
	return setItemIDValue(id, key, val)
}
func (id ItemID) setInt(key string, val int) error {
	return setItemIDValue(id, key, val)
}
func (id ItemID) setString(key string, val string) error {
	return setItemIDValue(id, key, val)
}
func setItemIDValue[T bool | float64 | int | string](id ItemID, key string, val T) (err error) {
	query := `UPDATE Item SET ` + key + ` = @1 WHERE ItemID = @2 AND ` + key + ` <> @3`
	log.Printf("UPDATE Item SET %s = %v WHERE ItemID = %d AND %s <> %v", key, val, id, key, val)
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
	dm, _ := id.DateModified()
	id.Item().DateModified.Set(dm.Format(time.DateTime))
	return
}

/* Get the pointer to Item from map or make one and return it */
func getItem(b *Backend, id ItemID) *Item {
	if t := b.Items.data[id]; t == nil {
		t = newItem(b, id)
		b.Items.data[id] = t
	}
	return b.Items.data[id]
}

type Items struct {
	db   *sql.DB
	j    *journal.Journal
	data map[ItemID]*Item

	ItemIDList                    binding.UntypedList
	ItemIDSelection               binding.UntypedList
	SearchResultUniqueCompletions binding.StringList
	SearchString                  binding.String
	searchType                    SearchType
	searchKey                     SearchKey
	sortKey                       SearchKey
	sortOrder                     SortOrder
}

func NewItems(b *Backend) *Items {
	m := &Items{
		db:                            b.db,
		j:                             b.Journal,
		data:                          make(map[ItemID]*Item),
		ItemIDList:                    binding.NewUntypedList(),
		ItemIDSelection:               binding.NewUntypedList(),
		SearchResultUniqueCompletions: binding.NewStringList(),
		SearchString:                  binding.NewString(),
		searchKey:                     SearchKeyName,
		sortKey:                       SearchKeyItemID,
		sortOrder:                     SortAscending,
	}
	m.SearchString.AddListener(binding.NewDataListener(func() { m.Search() }))
	return m
}
func (m *Items) GetAllItemIDs() {
	// TODO redo this to fetch all according to current selection/search config, then call after any mod to list
	query := `SELECT ItemID FROM Item WHERE ItemStatusID <> @0`
	stmt, err := m.db.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(ItemStatusDeleted)
	if err != nil {
		panic(err)
	}
	var id ItemID
	m.ItemIDList.Set([]any{})
	for rows.Next() {
		rows.Scan(&id)
		m.ItemIDList.Append(id)
	}
}
func (m *Items) GetItem(id ItemID) *Item {
	return getItem(be, id)
}
func (m *Items) GetItemIDFor(index widget.ListItemID) (ItemID, error) {
	id, err := m.ItemIDList.GetValue(index)
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
	stmt, err := m.db.Prepare(query)
	if err != nil {
		log.Printf("Items.CreateNewItem() error: %v", err)
		return i, fmt.Errorf("Items.CreateNewItem() error: %w", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec()
	if err != nil {
		log.Printf("Items.CreateNewItem() error: %v", err)
		return i, fmt.Errorf("Items.CreateNewItem() error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Items.CreateNewItem() error: %v", err)
		return i, fmt.Errorf("Items.CreateNewItem() error: %w", err)
	}
	log.Printf("Items.CreateNewItem() result.LastInsertId() = %d", id)
	i = ItemID(id)
	// TODO this should not just append, but put in the correct position according to current config
	m.Search()
	return i, err
}
func (m *Items) CopyItem(id ItemID) (ItemID, error) {
	query := `INSERT INTO Item (Name, Price, Currency, QuantityInPrice, Unit, OrderMultiple, MinOrder, Vat, Eta, EtaText, 
Priority, Stock, ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL, UNSPSC, LongDesc, Manufacturer, 
MfrItemId, GlobId, GlobIdType, ReplacesItem, Questions, PackagingCode, PresentationCode, 
DeliveryAutoSign, DeliveryOption, ComparePrice, CompareUnit, CompareQuantityInPrice, PriceInfo, 
AddDesc, ProcFlow, InnerUnit, QuantityInUnit, RiskClassification, Comment, EnvClassification, 
FormId, Article, Attachments, ItemGroup, Width, Height, Depth, Weight, 
LengthUnitID, WeightUnitID, CatID, StorageId, ItemStatusID)
SELECT Name, Price, Currency, QuantityInPrice, Unit, OrderMultiple, MinOrder, Vat, Eta, EtaText, 
Priority, Stock, ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL, UNSPSC, LongDesc, Manufacturer, 
MfrItemId, GlobId, GlobIdType, ReplacesItem, Questions, PackagingCode, PresentationCode, 
DeliveryAutoSign, DeliveryOption, ComparePrice, CompareUnit, CompareQuantityInPrice, PriceInfo, 
AddDesc, ProcFlow, InnerUnit, QuantityInUnit, RiskClassification, Comment, EnvClassification, 
FormId, Article, Attachments, ItemGroup, Width, Height, Depth, Weight, 
LengthUnitID, WeightUnitID, CatID, StorageId, ItemStatusID
FROM Item WHERE ItemID = @0`
	stmt, err := m.db.Prepare(query)
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
	m.Search()
	return newid, err
}
func (m *Items) DeleteItem(id ItemID) error {
	query := `UPDATE Item SET ItemStatusID = @0 WHERE ItemID = @1`
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("DeleteItem error: %w", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(ItemStatusDeleted, id)
	if err != nil {
		return fmt.Errorf("DeleteItem error: %w", err)
	}
	raf, _ := res.RowsAffected()
	log.Printf("%d rows affected", raf)
	m.ItemIDList.Remove(id)
	delete(m.data, id)
	return err
}
func (m *Items) SelectItem(id ItemID) error {
	return m.ItemIDSelection.Append(id)
}
func (m *Items) UnselectItem(id ItemID) error {
	return m.ItemIDSelection.Remove(id)
}
func (m *Items) ClearSelection() error {
	return m.ItemIDSelection.Set([]any{})
}
func (m *Items) Search() {
	searchString, err := m.SearchString.Get()
	var query string
	searchKey := m.searchKey.String()
	switch m.searchType {
	case BeginsWith:
		searchString = fmt.Sprintf("%s%%", searchString)
		query = `SELECT ItemID FROM Item WHERE ` + searchKey + ` LIKE @0 AND ItemStatusID <> @1 ORDER BY ` + m.sortKey.String() + ` ` + m.sortOrder.String()
	case EndsWith:
		searchString = fmt.Sprintf("%%%s", searchString)
		query = `SELECT ItemID FROM Item WHERE ` + searchKey + ` LIKE @0 AND ItemStatusID <> @1 ORDER BY ` + m.sortKey.String() + ` ` + m.sortOrder.String()
	case Contains:
		searchString = fmt.Sprintf("%%%s%%", searchString)
		query = `SELECT ItemID FROM Item WHERE ` + searchKey + ` LIKE @0 AND ItemStatusID <> @1 ORDER BY ` + m.sortKey.String() + ` ` + m.sortOrder.String()
	// case RegExp:
	// 	query = `SELECT ItemID FROM Item WHERE ` + searchKey + ` REGEXP @0 AND ItemStatusID <> @1 ORDER BY ` + m.sortKey.String() + ` ` + m.sortOrder.String()
	default:
		// Equals
		query = `SELECT ItemID FROM Item WHERE ` + searchKey + ` LIKE @0 AND ItemStatusID <> @1 ORDER BY ` + m.sortKey.String() + ` ` + m.sortOrder.String()
	}
	stmt, err := m.db.Prepare(query)
	if err != nil {
		log.Println(fmt.Errorf("search: prepare statement failed: %w", err))
		return
	}
	defer stmt.Close()

	clearQuery := strings.Replace(query, "@0", searchString, 1)
	clearQuery = strings.Replace(clearQuery, "@1", fmt.Sprintf("%d", ItemStatusDeleted), 1)
	log.Println(clearQuery)

	rows, err := stmt.Query(searchString, ItemStatusDeleted)
	if err != nil {
		log.Println(fmt.Errorf("search: statement query failed: %w", err))
		return
	}
	m.ItemIDList.Set([]any{})
	uniqueResults := make(map[string]bool)
	m.SearchResultUniqueCompletions.Set([]string{})
	for rows.Next() {
		var hit string
		var id ItemID
		rows.Scan(&id)
		m.ItemIDList.Append(id)
		if m.searchKey == SearchKeyName {
			hit, _ = id.Name()
		}
		if m.searchKey == SearchKeyManufacturer {
			hit, _ = id.Manufacturer()
		}
		if !uniqueResults[hit] {
			uniqueResults[hit] = true
			m.SearchResultUniqueCompletions.Append(hit)
		}
	}
}
func (m *Items) SetSearchConfig(c SearchType) error {
	m.searchType = c
	return nil
}
func (m *Items) SetSearchKey(k SearchKey) error {
	m.searchKey = k
	return nil
}
func (m *Items) SetSortKey(k SearchKey) error {
	m.sortKey = k
	return nil
}
func (m *Items) SetSortOrder(o SortOrder) error {
	m.sortOrder = o
	return nil
}
func (m *Items) SearchKey() SearchKey {
	return m.searchKey
}
func (m *Items) SortKey() SearchKey {
	return m.sortKey
}
func (m *Items) SortOrder() SortOrder {
	return m.sortOrder
}

type Item struct {
	binding.DataItem
	db           *sql.DB
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
	Manufacturer binding.String
	Model        binding.String
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
	dateCreated  time.Time
	dateModified time.Time
	DateCreated  binding.String
	DateModified binding.String
	Condition    map[string]any // TODO fix this
	Properties   map[string]any // TODO fix this
	touched      bool
}

func newItem(b *Backend, id ItemID) *Item {
	t := &Item{
		db:     b.db,
		ItemID: id,
		CatID:  CatID(0),
	}

	t.getAllFields()

	t.PriceString = binding.FloatToStringWithFormat(t.priceFloat, "%.2f")
	t.VatString = binding.FloatToStringWithFormat(t.vatFloat, "%.2f")
	t.StockString = binding.FloatToStringWithFormat(t.stockFloat, "%.0f")
	t.WidthString = binding.FloatToStringWithFormat(t.widthFloat, "%.0f")
	t.HeightString = binding.FloatToStringWithFormat(t.heightFloat, "%.0f")
	t.DepthString = binding.FloatToStringWithFormat(t.depthFloat, "%.0f")
	t.VolumeString = binding.FloatToStringWithFormat(t.volumeFloat, "%.2f")
	t.WeightString = binding.FloatToStringWithFormat(t.weightFloat, "%.2f")

	t.Name.AddListener(binding.NewDataListener(func() { t.ItemID.SetName(); b.Items.Search() }))
	t.Category.AddListener(binding.NewDataListener(func() { t.ItemID.SetCategory() }))
	t.priceFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetPrice() }))
	t.Currency.AddListener(binding.NewDataListener(func() { t.ItemID.SetCurrency() }))
	t.Unit.AddListener(binding.NewDataListener(func() { t.ItemID.SetUnit() }))
	t.vatFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetVat() }))
	t.Priority.AddListener(binding.NewDataListener(func() { t.ItemID.SetPriority() }))
	t.stockFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetStock() }))
	// TODO implement SearchWords
	// t.SearchWords.AddListener(binding.NewDataListener(func() { panic("SearchWords unimplemented") }))
	t.ImgURL1.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL1() }))
	t.ImgURL2.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL2() }))
	t.ImgURL3.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL3() }))
	t.ImgURL4.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL4() }))
	t.ImgURL5.AddListener(binding.NewDataListener(func() { t.ItemID.SetImgURL5() }))
	t.SpecsURL.AddListener(binding.NewDataListener(func() { t.ItemID.SetSpecsURL() }))
	t.AddDesc.AddListener(binding.NewDataListener(func() { t.ItemID.SetAddDesc() }))
	t.LongDesc.AddListener(binding.NewDataListener(func() { t.ItemID.SetLongDesc() }))
	t.Manufacturer.AddListener(binding.NewDataListener(func() { t.ItemID.SetManufacturer() }))
	t.Model.AddListener(binding.NewDataListener(func() { t.ItemID.SetModel() }))
	t.ModelURL.AddListener(binding.NewDataListener(func() { t.ItemID.SetModelURL() }))
	t.Notes.AddListener(binding.NewDataListener(func() { t.ItemID.SetNotes() }))
	t.widthFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetWidth() }))
	t.heightFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetHeight() }))
	t.depthFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetDepth() }))
	t.volumeFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetVolume() }))
	t.weightFloat.AddListener(binding.NewDataListener(func() { t.ItemID.SetWeight() }))
	t.LengthUnit.AddListener(binding.NewDataListener(func() { t.ItemID.SetLengthUnit() }))
	t.VolumeUnit.AddListener(binding.NewDataListener(func() { t.ItemID.SetVolumeUnit() }))
	t.WeightUnit.AddListener(binding.NewDataListener(func() { t.ItemID.SetWeightUnit() }))
	t.ItemStatus.AddListener(binding.NewDataListener(func() { t.ItemID.SetItemStatus() }))
	return t
}

func (t *Item) getAllFields() {
	var Name, Currency, Unit, ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL sql.NullString
	var AddDesc, LongDesc, Manufacturer, Model, ModelURL, Notes, DateCreated, DateModified sql.NullString
	var Price, QuantityInPrice, Vat, Stock, Width, Height, Depth, Volume, Weight sql.NullFloat64
	var Priority sql.NullBool
	var CatID CatID
	var LengthUnitID, VolumeUnitID, WeightUnitID UnitID
	var ItemStatusID ItemStatusID

	query := `SELECT Name, CatID, Price, Currency, QuantityInPrice, Unit, Vat, 
Priority, Stock, ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL, 
AddDesc, LongDesc, Manufacturer, Model, ModelURL, Notes, 
Width, Height, Depth, Volume, Weight, 
LengthUnitID, VolumeUnitID, WeightUnitID, 
ItemStatusID, DateCreated, DateModified 
FROM Item WHERE ItemID = @0`
	stmt, err := t.db.Prepare(query)
	if err != nil {
		log.Println(fmt.Errorf("getAllFields error: %w", err))
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.ItemID).Scan(
		&Name, &CatID, &Price, &Currency, &QuantityInPrice, &Unit, &Vat,
		&Priority, &Stock, &ImgURL1, &ImgURL2, &ImgURL3, &ImgURL4, &ImgURL5, &SpecsURL,
		&AddDesc, &LongDesc, &Manufacturer, &Model, &ModelURL, &Notes,
		&Width, &Height, &Depth, &Volume, &Weight,
		&LengthUnitID, &VolumeUnitID, &WeightUnitID,
		&ItemStatusID, &DateCreated, &DateModified,
	)
	if err != nil {
		log.Println("Item.getAllFields() error!")
		panic(err)
	}

	category, err := CatID.Name()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		} else {
			log.Println(err)
			category = ""
		}
	}
	LengthUnitString, err := LengthUnitID.Name()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		} else {
			log.Println(err)
			LengthUnitString = ""
		}
	}
	VolumeUnitString, err := VolumeUnitID.Name()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		} else {
			log.Println(err)
			VolumeUnitString = ""
		}
	}
	WeightUnitString, err := WeightUnitID.Name()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		} else {
			log.Println(err)
			WeightUnitString = ""
		}
	}
	ItemStatusString := ItemStatusID.LString()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		} else {
			log.Println(err)
			ItemStatusString = ""
		}
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
	t.Model = binding.NewString()
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
	t.Manufacturer.Set(Manufacturer.String)
	t.Model.Set(Model.String)
	t.ModelURL.Set(ModelURL.String)
	t.Notes.Set(Notes.String)
	t.widthFloat.Set(Width.Float64)
	t.heightFloat.Set(Height.Float64)
	t.depthFloat.Set(Depth.Float64)
	t.volumeFloat.Set(Volume.Float64)
	t.weightFloat.Set(Weight.Float64)
	t.LengthUnit.Set(LengthUnitString)
	t.VolumeUnit.Set(VolumeUnitString)
	t.WeightUnit.Set(WeightUnitString)
	t.ItemStatus.Set(ItemStatusString)

	var created, modified time.Time
	stockholm, err := time.LoadLocation("Europe/Stockholm")
	if DateCreated.Valid {
		utc, err := time.Parse(subsec, DateCreated.String)
		created = utc.In(stockholm)
		if err != nil {
			log.Println(fmt.Errorf("error parsing DateCreated: %w", err))
			t.DateCreated.Set(fmt.Sprintf("error parsing DateCreated: %v", err))
		} else {
			t.DateCreated.Set(created.Format(time.DateTime))
		}
	}
	if DateModified.Valid {
		utc, err := time.Parse(subsec, DateCreated.String)
		modified = utc.In(stockholm)
		if err != nil {
			log.Println(fmt.Errorf("error parsing DateCreated: %w", err))
			t.DateModified.Set(fmt.Sprintf("error parsing DateModified: %v", err))
		} else {
			t.DateModified.Set(modified.Format(time.DateTime))
		}
	}
}
