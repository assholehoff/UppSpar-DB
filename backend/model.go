package backend

import (
	"database/sql"
	"log"

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
		db:      b.db,
		ModelID: id,
		Name:    binding.NewString(),
	}

	name, err := mdl.ModelID.Name()
	if err != nil {
		log.Println(err)
	}
	mdl.Name.Set(name)
	mdl.Name.AddListener(binding.NewDataListener(func() { mdl.ModelID.SetName() }))
	return mdl
}
