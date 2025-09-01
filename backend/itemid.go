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
	"time"

	"fyne.io/fyne/v2/lang"
)

type ItemID int

var (
	_ sql.Scanner   = (*ItemID)(nil)
	_ driver.Valuer = (*ItemID)(nil)
	_ fmt.Stringer  = (*ItemID)(nil)
	_ NumID         = (ItemID)(0)
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
		log.Printf("ItemID(%d).Scan(%v) type error: %s", id, src, reflect.TypeOf(src).Name())
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

/* Returns a tree-friendly identifying string */
func (id ItemID) TString() string {
	return fmt.Sprintf("ITEM-%d", id)
}
func (id ItemID) Int() int {
	return int(id)
}
func (id ItemID) TypeName() string {
	return "ItemID"
}
func (id ItemID) Item() *Item {
	return b.Items.GetItem(id)
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
func (id ItemID) MfrID() (MfrID, error) {
	mid, err := id.getInt("MfrID")
	return MfrID(mid), err
}
func (id ItemID) Manufacturer() (string, error) {
	mid, err := id.MfrID()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	return mid.Name()
}
func (id ItemID) ModelID() (ModelID, error) {
	mid, err := id.getInt("ModelID")
	return ModelID(mid), err
}
func (id ItemID) Model() (string, error) {
	mid, err := id.ModelID()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	return mid.Name()
}
func (id ItemID) ModelURL() (string, error) {
	return id.getString("ModelURL")
}
func (id ItemID) Notes() (string, error) {
	return id.getString("Notes")
}
func (id ItemID) Width() (float64, error) {
	if mi, _ := id.ModelID(); mi == 0 {
		log.Printf("ModelID is 0, shortcircuiting...")
		return id.getFloat("Width")
	}
	w, _ := id.getFloat("Width")
	h, _ := id.getFloat("Height")
	d, _ := id.getFloat("Depth")
	if w == 0 && h == 0 && d == 0 {
		v, err := id.Item().ModelID.Width()
		log.Printf("item width is 0 but ModelID is %d, show model width %f and err %s", id.Item().ModelID, v, err)
		return id.Item().ModelID.Width()
	}
	log.Printf("one or more dimensions is not 0, returning width field")
	return id.getFloat("Width")
}
func (id ItemID) Height() (float64, error) {
	if mi, _ := id.ModelID(); mi == 0 {
		return id.getFloat("Height")
	}
	w, _ := id.getFloat("Width")
	h, _ := id.getFloat("Height")
	d, _ := id.getFloat("Depth")
	if w == 0 && h == 0 && d == 0 {
		return id.Item().ModelID.Height()
	}
	return id.getFloat("Height")
}
func (id ItemID) Depth() (float64, error) {
	if mi, _ := id.ModelID(); mi == 0 {
		return id.getFloat("Depth")
	}
	w, _ := id.getFloat("Width")
	h, _ := id.getFloat("Height")
	d, _ := id.getFloat("Depth")
	if w == 0 && h == 0 && d == 0 {
		return id.Item().ModelID.Depth()
	}
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
	stat, err := id.getInt("ItemStatusID")
	if err != nil {
		log.Printf("ItemID(%d).ItemStatus() error: %s", id, err)
		return ""
	}
	return ItemStatusID(stat).LString()
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
	b, err := getValue[sql.NullBool]("Item", id, key)
	if b.Valid {
		val = b.Bool
	} else {
		log.Printf("ItemID(%d).getBool(%s) error: %s", id, key, err)
	}
	return
}
func (id ItemID) getFloat(key string) (val float64, err error) {
	f, err := getValue[sql.NullFloat64]("Item", id, key)
	if f.Valid {
		val = f.Float64
	} else {
		log.Printf("ItemID(%d).getFloat(%s) error: %s", id, key, err)
	}
	return
}
func (id ItemID) getInt(key string) (val int, err error) {
	i, err := getValue[sql.NullInt64]("Item", id, key)
	val = int(i.Int64)
	if !i.Valid {
		log.Printf("ItemID(%d).getInt(%s) error: %s", id, key, err)
	}
	return
}
func (id ItemID) getString(key string) (val string, err error) {
	s, err := getValue[sql.NullString]("Item", id, key)
	if s.Valid {
		val = s.String
	} else {
		log.Printf("ItemID(%d).getInt(%s) error: %s", id, key, err)
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
	var longDesc string
	n := false
	addStringToLine := func(s string, e error) {
		if e != nil && !errors.Is(e, sql.ErrNoRows) {
			log.Printf("ItemID(%d).CompileLongDesc().addStringToLine(%s) error: %s", id, s, e)
		}
		if s != "" {
			longDesc += fmt.Sprintf("%s ", s)
			n = true
		}
	}
	addNewlines := func(i int) {
		if n {
			for range i {
				longDesc += "\n"
			}
			n = false
		}
	}
	/* Manufacturer and model */
	addStringToLine(id.Manufacturer())
	addStringToLine(id.Model())
	addNewlines(1)
	addStringToLine(id.Name())
	addNewlines(2)
	addStringToLine(id.Notes())

	id.Item().LongDesc.Set(longDesc)
	return id.SetLongDesc()
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
	s = strings.TrimSpace(s)
	n, err := CatIDFor(s)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("ItemID(%d).SetCategory(%s) error: %s", id, s, err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows
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
func (id ItemID) SetMfrID(val MfrID) error {
	key := "MfrID"
	return id.setInt(key, int(val))
}
func (id ItemID) SetManufacturer() error {
	key := "Manufacturer"
	s, err := id.Item().Manufacturer.Get()
	if err != nil {
		return fmt.Errorf("ItemID.SetManufacturer() error: %w", err)
	}
	if len(s) == 0 {
		id.setString(key, s)
		return nil
	}
	n, err := MfrIDFor(s)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("ItemID(%d).SetManufacturer(%s) error: %s", id, s, err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		// no such manufacturer exists, set the name field instead
		id.setString(key, s)
		return nil
	}
	id.Item().MfrID = n
	return id.SetMfrID(n)
}
func (id ItemID) SetModelID(val ModelID) error {
	key := "ModelID"
	err := id.setInt(key, int(val))
	if err != nil {
		log.Printf("updating fields with model data...")
		id.Item().CatID = val.Model().CatID
		// Update fields with data from ModelID
		// General description
		if d, _ := val.Desc(); d != "" {
			id.Item().ModelDesc.Set(d)
		}
		// Measurements
		touched := false
		if w, _ := val.Width(); w != 0 {
			id.Item().widthFloat.Set(w)
			touched = true
		}
		if h, _ := val.Height(); h != 0 {
			id.Item().heightFloat.Set(h)
			touched = true
		}
		if d, _ := val.Depth(); d != 0 {
			id.Item().depthFloat.Set(d)
			touched = true
		}
		if touched {
			iu, _ := id.LengthUnitID()
			mu, _ := val.LengthUnitID()
			if iu != mu {
				id.SetLengthUnitID(mu)
			}
		}
		// Volume
		if v, _ := val.Volume(); v != 0 {
			id.Item().volumeFloat.Set(v)
			u, _ := val.VolumeUnitID()
			id.SetVolumeUnitID(u)
		}
		// Weight
		if w, _ := val.Weight(); w != 0 {
			id.Item().weightFloat.Set(w)
			u, _ := val.WeightUnitID()
			id.SetWeightUnitID(u)
		}
	}
	return err
}
func (id ItemID) SetModelName() error {
	key := "ModelName"
	s, err := id.Item().ModelName.Get()
	if err != nil {
		return fmt.Errorf("ItemID(%d).SetModel(%s) error: %s", id, s, err)
	}
	if len(s) == 0 {
		id.setString(key, s)
		id.SetModelID(0)
		return nil
	}
	mfr, _ := id.MfrID()
	n, err := ModelIDFor(mfr, s)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("ItemID(%d).SetModel(%s) error: %s", id, s, err)
	}
	if errors.Is(err, sql.ErrNoRows) || n == 0 {
		// no such model exists, set the name field
		id.setString(key, s)
		return nil
	}
	// log.Printf("ModelIDFor(s) returned something: %d, error: %s", n, err)
	// set ModelID to returned ID
	id.Item().ModelID = n
	return id.SetModelID(n)
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
func (id ItemID) updateDateModified() {
	dm, err := id.DateModified()
	if err != nil {
		log.Println(err)
	}
	id.Item().DateModified.Set(dm.Format(time.DateTime))
}

func (id ItemID) setBool(key string, val bool) error {
	err := setValue("Item", id, key, val)
	id.updateDateModified()
	return err
}
func (id ItemID) setFloat(key string, val float64) error {
	err := setValue("Item", id, key, val)
	id.updateDateModified()
	return err
}
func (id ItemID) setInt(key string, val int) error {
	err := setValue("Item", id, key, val)
	id.updateDateModified()
	return err
}
func (id ItemID) setString(key string, val string) error {
	err := setValue("Item", id, key, val)
	id.updateDateModified()
	return err
}
