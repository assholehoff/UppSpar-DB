package backend

import (
	"database/sql"
	"errors"

	"fyne.io/fyne/v2/data/binding"
)

/* Note: "Model" should probably be called "Product" */
type Model struct {
	binding.DataItem
	db           *sql.DB
	ModelID      ModelID
	Name         binding.String
	MfrID        MfrID
	Manufacturer binding.String
	CatID        CatID
	Category     binding.String
	Desc         binding.String
	ImgURL1      binding.String
	ImgURL2      binding.String
	ImgURL3      binding.String
	ImgURL4      binding.String
	ImgURL5      binding.String
	SpecsURL     binding.String
	ModelURL     binding.String
	Width        binding.String
	Height       binding.String
	Depth        binding.String
	Volume       binding.String
	Weight       binding.String
	widthFloat   binding.Float
	heightFloat  binding.Float
	depthFloat   binding.Float
	volumeFloat  binding.Float
	weightFloat  binding.Float
	LengthUnit   binding.String
	VolumeUnit   binding.String
	WeightUnit   binding.String
	branch       bool
	touched      bool
}

func newModel(b *Backend, id ModelID) *Model {
	mdl := &Model{
		db:           b.db,
		ModelID:      id,
		Name:         binding.NewString(),
		Manufacturer: binding.NewString(),
		Category:     binding.NewString(),
		Desc:         binding.NewString(),
		ImgURL1:      binding.NewString(),
		ImgURL2:      binding.NewString(),
		ImgURL3:      binding.NewString(),
		ImgURL4:      binding.NewString(),
		ImgURL5:      binding.NewString(),
		SpecsURL:     binding.NewString(),
		ModelURL:     binding.NewString(),
		Width:        binding.NewString(),
		Height:       binding.NewString(),
		Depth:        binding.NewString(),
		Volume:       binding.NewString(),
		Weight:       binding.NewString(),
		widthFloat:   binding.NewFloat(),
		heightFloat:  binding.NewFloat(),
		depthFloat:   binding.NewFloat(),
		volumeFloat:  binding.NewFloat(),
		weightFloat:  binding.NewFloat(),
		LengthUnit:   binding.NewString(),
		VolumeUnit:   binding.NewString(),
		WeightUnit:   binding.NewString(),
	}

	var Name, Desc, ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL, ModelURL, Manufacturer sql.NullString
	var Width, Height, Depth, Volume, Weight sql.NullFloat64
	var CatID CatID
	var MfrID MfrID
	var LengthUnitID, VolumeUnitID, WeightUnitID UnitID

	query := `SELECT Name, Manufacturer, MfrID, Desc, ImgURL1, ImgURL2, ImgURL3, ImgURL4, ImgURL5, SpecsURL, ModelURL, 
Width, Height, Depth, Volume, Weight, LengthUnitID, VolumeUnitID, WeightUnitID, CatID
FROM Model WHERE ModelID = @0`
	err := be.db.QueryRow(query, mdl.ModelID).Scan(
		&Name, &Manufacturer, &MfrID, &Desc, &ImgURL1, &ImgURL2, &ImgURL3, &ImgURL4, &ImgURL5, &SpecsURL, &ModelURL,
		&Width, &Height, &Depth, &Volume, &Weight, &LengthUnitID, &VolumeUnitID, &WeightUnitID, &CatID,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}

	Category, err := CatID.Name()
	manufacturer := Manufacturer.String
	LengthUnit, err := LengthUnitID.Name()
	VolumeUnit, err := VolumeUnitID.Name()
	WeightUnit, err := WeightUnitID.Name()

	if MfrID != 0 {
		if n, _ := MfrID.Name(); n != manufacturer {
			manufacturer = n
		}
	}

	mdl.Name.Set(Name.String)
	mdl.Category.Set(Category)
	mdl.CatID = CatID
	mdl.MfrID = MfrID
	mdl.Manufacturer.Set(manufacturer)
	mdl.Desc.Set(Desc.String)
	mdl.ImgURL1.Set(ImgURL1.String)
	mdl.ImgURL2.Set(ImgURL2.String)
	mdl.ImgURL3.Set(ImgURL3.String)
	mdl.ImgURL4.Set(ImgURL4.String)
	mdl.ImgURL5.Set(ImgURL5.String)
	mdl.SpecsURL.Set(SpecsURL.String)
	mdl.ModelURL.Set(ModelURL.String)
	mdl.widthFloat.Set(Width.Float64)
	mdl.heightFloat.Set(Height.Float64)
	mdl.depthFloat.Set(Depth.Float64)
	mdl.volumeFloat.Set(Volume.Float64)
	mdl.weightFloat.Set(Weight.Float64)
	mdl.LengthUnit.Set(LengthUnit)
	mdl.VolumeUnit.Set(VolumeUnit)
	mdl.WeightUnit.Set(WeightUnit)

	mdl.Width = binding.FloatToStringWithFormat(mdl.widthFloat, "%.2f")
	mdl.Height = binding.FloatToStringWithFormat(mdl.heightFloat, "%.2f")
	mdl.Depth = binding.FloatToStringWithFormat(mdl.depthFloat, "%.2f")
	mdl.Volume = binding.FloatToStringWithFormat(mdl.volumeFloat, "%.2f")
	mdl.Weight = binding.FloatToStringWithFormat(mdl.weightFloat, "%.2f")

	mdl.Name.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetName(); b.Metadata.GetProductTree() }))
	mdl.Manufacturer.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetManufacturer(); b.Metadata.GetProductTree() }))
	mdl.Category.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetCategory() }))
	mdl.Desc.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetDesc() }))
	mdl.ImgURL1.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetImgURL1() }))
	mdl.ImgURL2.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetImgURL2() }))
	mdl.ImgURL3.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetImgURL3() }))
	mdl.ImgURL4.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetImgURL4() }))
	mdl.ImgURL5.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetImgURL5() }))
	mdl.SpecsURL.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetSpecsURL() }))
	mdl.ModelURL.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetModelURL() }))
	mdl.widthFloat.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetWidth() }))
	mdl.heightFloat.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetHeight() }))
	mdl.depthFloat.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetDepth() }))
	mdl.volumeFloat.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetVolume() }))
	mdl.weightFloat.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetWeight() }))
	mdl.LengthUnit.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetLengthUnit() }))
	mdl.VolumeUnit.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetVolumeUnit() }))
	mdl.WeightUnit.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetWeightUnit() }))
	return mdl
}
