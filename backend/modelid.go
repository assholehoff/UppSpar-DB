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

var _ NumID = (ModelID)(0)

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
		log.Printf("ModelID.Scan(%v) error: unknown type %s", src, reflect.TypeOf(src).Name())
		return ErrInvalidType
	}
	return nil
}

func (id ModelID) TypeName() string {
	return "ModelID"
}

func (id ModelID) Name() (string, error) {
	return id.getString("Name")
}
func (id ModelID) CatID() (CatID, error) {
	cid, err := id.getInt("CatID")
	return CatID(cid), err
}
func (id ModelID) MfrID() (MfrID, error) {
	val, err := id.getInt("MfrID")
	return MfrID(val), err
}
func (id ModelID) Desc() (string, error) {
	return id.getString("Desc")
}
func (id ModelID) ImgURL1() (string, error) {
	return id.getString("ImgURL1")
}
func (id ModelID) ImgURL2() (string, error) {
	return id.getString("ImgURL2")
}
func (id ModelID) ImgURL3() (string, error) {
	return id.getString("ImgURL3")
}
func (id ModelID) ImgURL4() (string, error) {
	return id.getString("ImgURL4")
}
func (id ModelID) ImgURL5() (string, error) {
	return id.getString("ImgURL5")
}
func (id ModelID) SpecsURL() (string, error) {
	return id.getString("SpecsURL")
}
func (id ModelID) ModelURL() (string, error) {
	return id.getString("ModelURL")
}
func (id ModelID) Width() (float64, error) {
	return id.getFloat("Width")
}
func (id ModelID) Height() (float64, error) {
	return id.getFloat("Height")
}
func (id ModelID) Depth() (float64, error) {
	return id.getFloat("Depth")
}
func (id ModelID) Weight() (float64, error) {
	return id.getFloat("Weight")
}
func (id ModelID) LengthUnit() (string, error) {
	val, err := id.getInt("LengthUnitID")
	return UnitID(val).String(), err
}
func (id ModelID) VolumeUnit() (string, error) {
	val, err := id.getInt("VolumeUnitID")
	return UnitID(val).String(), err
}
func (id ModelID) WeightUnit() (string, error) {
	val, err := id.getInt("WeightUnitID")
	return UnitID(val).String(), err
}
func (id ModelID) LengthUnitID() (UnitID, error) {
	uid, err := id.getInt("LengthUnitID")
	return UnitID(uid), err
}
func (id ModelID) VolumeUnitID() (UnitID, error) {
	uid, err := id.getInt("VolumeUnitID")
	return UnitID(uid), err
}
func (id ModelID) WeightUnitID() (UnitID, error) {
	uid, err := id.getInt("WeightUnitID")
	return UnitID(uid), err
}

func (id ModelID) getBool(key string) (val bool, err error) {
	b, err := getValue[sql.NullBool]("Model", id, key)
	if b.Valid {
		val = b.Bool
	} else {
		log.Printf("getBool(%s) b is invalid (NULL), err is %v", key, err)
		err = ErrSQLNullValue
	}
	return
}
func (id ModelID) getFloat(key string) (val float64, err error) {
	f, err := getValue[sql.NullFloat64]("Model", id, key)
	if f.Valid {
		val = f.Float64
	} else {
		log.Printf("getFloat(%s) %s is invalid (NULL), err is %v", key, key, err)
		err = ErrSQLNullValue
	}
	return
}
func (id ModelID) getInt(key string) (val int, err error) {
	i, err := getValue[sql.NullInt64]("Model", id, key)
	val = int(i.Int64)
	if !i.Valid {
		log.Printf("getInt(%s) %s is invalid (NULL), err is %v", key, key, err)
		err = ErrSQLNullValue
	}
	return
}
func (id ModelID) getString(key string) (val string, err error) {
	s, err := getValue[sql.NullString]("Model", id, key)
	if s.Valid {
		val = s.String
	} else {
		log.Printf("getInt(%s) %s is invalid (NULL), err is %v", key, key, err)
		err = ErrSQLNullValue
	}
	return
}

func (id ModelID) SetName() error {
	key := "Name"
	val, err := id.Model().Name.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetName() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetCatID(val CatID) error {
	key := "CatID"
	return id.setInt(key, int(val))
}
func (id ModelID) SetCategory() error {
	s, err := id.Model().Category.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetCategory() error: %w", err)
	}
	n, err := CatIDFor(s)
	if err != nil {
		return fmt.Errorf("ModelID.SetCategory() error: %w", err)
	}
	id.Model().CatID = n
	return id.SetCatID(n)
}
func (id ModelID) SetMfrID(val MfrID) error {
	key := "MfrID"
	return id.setInt(key, int(val))
}
func (id ModelID) SetManufacturer() error {
	s, err := id.Model().Manufacturer.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetManufacturer() error: %w", err)
	}
	n, err := MfrIDFor(s)
	if err != nil {
		return fmt.Errorf("ModelID.SetManufacturer() error: %w", err)
	}
	id.Model().MfrID = n
	return id.SetMfrID(n)
}
func (id ModelID) SetDesc() error {
	key := "Desc"
	val, err := id.Model().ImgURL1.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetDesc() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetImgURL1() error {
	key := "ImgURL1"
	val, err := id.Model().ImgURL1.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetImgURL1() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetImgURL2() error {
	key := "ImgURL2"
	val, err := id.Model().ImgURL2.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetImgURL2() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetImgURL3() error {
	key := "ImgURL3"
	val, err := id.Model().ImgURL3.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetImgURL3() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetImgURL4() error {
	key := "ImgURL4"
	val, err := id.Model().ImgURL4.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetImgURL4() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetImgURL5() error {
	key := "ImgURL5"
	val, err := id.Model().ImgURL5.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetImgURL5() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetSpecsURL() error {
	key := "SpecsURL"
	val, err := id.Model().SpecsURL.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetSpecsURL() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetModelURL() error {
	key := "ModelURL"
	val, err := id.Model().ModelURL.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetModelURL() error: %w", err)
	}
	return id.setString(key, val)
}
func (id ModelID) SetWidth() error {
	key := "Width"
	val, err := id.Model().widthFloat.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetWidth() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ModelID) SetHeight() error {
	key := "Height"
	val, err := id.Model().heightFloat.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetHeight() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ModelID) SetDepth() error {
	key := "Depth"
	val, err := id.Model().depthFloat.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetDepth() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ModelID) SetVolume() error {
	key := "Volume"
	val, err := id.Model().volumeFloat.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetVolume() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ModelID) SetWeight() error {
	key := "Weight"
	val, err := id.Model().weightFloat.Get()
	if err != nil {
		return fmt.Errorf("ModelID.SetWeight() error: %w", err)
	}
	return id.setFloat(key, val)
}
func (id ModelID) SetLengthUnit() error {
	str, err := id.Model().LengthUnit.Get()
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
func (id ModelID) SetLengthUnitID(l UnitID) error {
	key := "LengthUnitID"
	val := int(l)
	return id.setInt(key, val)
}
func (id ModelID) SetVolumeUnit() error {
	str, err := id.Model().VolumeUnit.Get()
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
func (id ModelID) SetVolumeUnitID(v UnitID) error {
	key := "VolumeUnitID"
	val := int(v)
	return id.setInt(key, val)
}
func (id ModelID) SetWeightUnit() error {
	str, err := id.Model().WeightUnit.Get()
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
func (id ModelID) SetWeightUnitID(w UnitID) error {
	key := "WeightUnitID"
	val := int(w)
	return id.setInt(key, val)
}

func (id ModelID) setBool(key string, val bool) error {
	err := setValue("Model", id, key, val)
	return err
}
func (id ModelID) setFloat(key string, val float64) error {
	err := setValue("Model", id, key, val)
	return err
}
func (id ModelID) setInt(key string, val int) error {
	err := setValue("Model", id, key, val)
	return err
}
func (id ModelID) setString(key string, val string) error {
	err := setValue("Model", id, key, val)
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
