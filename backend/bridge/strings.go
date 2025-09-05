package bridge

import (
	"time"

	"fyne.io/fyne/v2/lang"
)

var (
	CategoryFormCheckKeys = []string{
		"Price",
		"Dimensions",
		"Mass",
		"Volume",
		"Condition",
	}
	CategoryFormEntryKeys = []string{
		"Name",
	}
	CategoryFormLabelKeys = []string{
		"Name",
		"Price",
		"Dimensions",
		"Mass",
		"Volume",
		"Condition",
	}
	CategoryFormRadioKeys  = []string{}
	CategoryFormSelectKeys = []string{
		"Parent",
	}
	CategoryFormValuesKeys = []string{}

	ItemFormCheckKeys = []string{
		"New",
		"Working",
	}
	ItemFormEntryKeys = []string{
		"Name",
		"Price",
		"Vat",
		"ImgURL1",
		"ImgURL2",
		"ImgURL3",
		"ImgURL4",
		"ImgURL5",
		"SpecsURL",
		"LongDesc",
		"Manufacturer",
		"ModelDesc",
		"ModelURL",
		"Notes",
		"Width",
		"Height",
		"Depth",
		"Volume",
		"Weight",
	}
	ItemFormLabelKeys = []string{
		"ItemID",
		"Name",
		"Category",
		"Currency",
		"Price",
		"Vat",
		"ImgURL1",
		"ImgURL2",
		"ImgURL3",
		"ImgURL4",
		"ImgURL5",
		"SpecsURL",
		"AddDesc",
		"LongDesc",
		"Manufacturer",
		"ModelName",
		"ModelDesc",
		"ModelURL",
		"Notes",
		"Dimensions",
		"Width",
		"Height",
		"Depth",
		"Volume",
		"Weight",
		"Status",
		"DateCreated",
		"DateModified",
		"Condition",
		"Functionality",
	}
	ItemFormRadioKeys = []string{
		"Tested",
	}
	ItemFormSelectKeys = []string{
		"Category",
		"Manufacturer",
		"ModelName",
		"LengthUnit",
		"VolumeUnit",
		"WeightUnit",
		"Status",
	}
	ItemFormValueKeys = []string{
		"ItemID",
		"AddDesc",
		"LongDesc",
		"DateCreated",
		"DateModified",
	}

	ManufacturerFormCheckKeys  = []string{}
	ManufacturerFormEntryKeys  = []string{"Name"}
	ManufacturerFormLabelKeys  = []string{"Name"}
	ManufacturerFormRadioKeys  = []string{}
	ManufacturerFormSelectKeys = []string{}
	ManufacturerFormValuesKeys = []string{}

	ModelFormCheckKeys = []string{}
	ModelFormEntryKeys = []string{
		"Name",
		"Desc",
		"ImgURL1",
		"ImgURL2",
		"ImgURL3",
		"ImgURL4",
		"ImgURL5",
		"SpecsURL",
		"Manufacturer",
		"ModelURL",
		"Width",
		"Height",
		"Depth",
		"Volume",
		"Weight",
	}
	ModelFormLabelKeys = []string{
		"Name",
		"Desc",
		"ImgURL1",
		"ImgURL2",
		"ImgURL3",
		"ImgURL4",
		"ImgURL5",
		"SpecsURL",
		"Manufacturer",
		"ModelURL",
		"Dimensions",
		"Width",
		"Height",
		"Depth",
		"Volume",
		"Weight",
	}
	ModelFormRadioKeys  = []string{}
	ModelFormSelectKeys = []string{
		"Category",
		"Manufacturer",
		"LengthUnit",
		"VolumeUnit",
		"WeightUnit",
	}
	ModelFormValuesKeys = []string{}

	SearchBarCheckKeys  = []string{}
	SearchBarEntryKeys  = []string{}
	SearchBarLabelKeys  = []string{}
	SearchBarRadioKeys  = []string{}
	SearchBarSelectKeys = []string{}
	SearchBarValuesKeys = []string{}

	ItemFormLabelStrings = make(map[string]string)
	ItemFormValueStrings = make(map[string]string)

	ProductFormLabelStrings = make(map[string]string)
	ProductFormValueStrings = make(map[string]string)
)

func initItemStringMaps() {
	ItemFormLabelStrings["ItemID"] = lang.X("item.form.label.itemid", "item.form.label.itemid")
	ItemFormLabelStrings["Name"] = lang.X("item.form.label.name", "item.form.label.name")
	ItemFormLabelStrings["Category"] = lang.X("item.form.label.category", "item.form.label.category")
	ItemFormLabelStrings["Currency"] = "SEK"
	ItemFormLabelStrings["Price"] = lang.X("item.form.label.price", "item.form.label.price")
	ItemFormLabelStrings["Vat"] = lang.X("item.form.label.vat", "item.form.label.vat")
	ItemFormLabelStrings["ImgURL1"] = lang.X("item.form.label.imgurl", "item.form.label.imgurl")
	ItemFormLabelStrings["ImgURL2"] = lang.X("item.form.label.imgurl", "item.form.label.imgurl")
	ItemFormLabelStrings["ImgURL3"] = lang.X("item.form.label.imgurl", "item.form.label.imgurl")
	ItemFormLabelStrings["ImgURL4"] = lang.X("item.form.label.imgurl", "item.form.label.imgurl")
	ItemFormLabelStrings["ImgURL5"] = lang.X("item.form.label.imgurl", "item.form.label.imgurl")
	ItemFormLabelStrings["SpecsURL"] = lang.X("item.form.label.specsurl", "item.form.label.specsurl")
	ItemFormLabelStrings["AddDesc"] = lang.X("item.form.label.adddesc", "item.form.label.adddesc")
	ItemFormLabelStrings["LongDesc"] = lang.X("item.form.label.longdesc", "item.form.label.longdesc")
	ItemFormLabelStrings["Manufacturer"] = lang.X("item.form.label.manufacturer", "item.form.label.manufacturer")
	ItemFormLabelStrings["ModelName"] = lang.X("item.form.label.modelname", "item.form.label.modelname")
	ItemFormLabelStrings["ModelDesc"] = lang.X("item.form.label.modeldesc", "item.form.label.modeldesc")
	ItemFormLabelStrings["ModelURL"] = lang.X("item.form.label.modelurl", "item.form.label.modelurl")
	ItemFormLabelStrings["Notes"] = lang.X("item.form.label.notes", "item.form.label.notes")
	ItemFormLabelStrings["Dimensions"] = lang.X("item.form.label.dimensions", "item.form.label.dimensions")
	ItemFormLabelStrings["Width"] = lang.X("item.form.label.width", "item.form.label.width")
	ItemFormLabelStrings["Height"] = lang.X("item.form.label.height", "item.form.label.height")
	ItemFormLabelStrings["Depth"] = lang.X("item.form.label.depth", "item.form.label.depth")
	ItemFormLabelStrings["Volume"] = lang.X("item.form.label.volume", "item.form.label.volume")
	ItemFormLabelStrings["Weight"] = lang.X("item.form.label.weight", "item.form.label.weight")
	ItemFormLabelStrings["Status"] = lang.X("item.form.label.status", "item.form.label.status")
	ItemFormLabelStrings["DateCreated"] = lang.X("item.form.label.datecreated", "item.form.label.datecreated")
	ItemFormLabelStrings["DateModified"] = lang.X("item.form.label.datemodified", "item.form.label.datemodified")

	ItemFormValueStrings["ItemID"] = "0000000"
	ItemFormValueStrings["DateCreated"] = time.DateTime
	ItemFormValueStrings["DateModified"] = time.DateTime
	ItemFormValueStrings["AddDesc"] = lang.X("item.form.label.adddesc", "item.form.label.adddesc")
	ItemFormValueStrings["LongDesc"] = lang.X("item.form.label.longdesc", "item.form.label.longdesc")
}

func initProductStringMaps() {
	ProductFormLabelStrings["Name"] = lang.L("Name")
	ProductFormLabelStrings["Category"] = lang.X("item.form.label.category", "item.form.label.category")
	ProductFormLabelStrings["Manufacturer"] = lang.L("Manufacturer")
	ProductFormLabelStrings["ModelDesc"] = lang.X("metadata.product.form.description", "metadata.product.form.description")
	ProductFormLabelStrings["Dimensions"] = lang.L("Dimensions")
	ProductFormLabelStrings["ImgURL1"] = lang.L("Image URL") + " 1"
	ProductFormLabelStrings["ImgURL2"] = lang.L("Image URL") + " 2"
	ProductFormLabelStrings["ImgURL3"] = lang.L("Image URL") + " 3"
	ProductFormLabelStrings["ImgURL4"] = lang.L("Image URL") + " 4"
	ProductFormLabelStrings["ImgURL5"] = lang.L("Image URL") + " 5"
	ProductFormLabelStrings["SpecsURL"] = lang.L("Specs URL")
	ProductFormLabelStrings["ModelURL"] = lang.L("Model URL")
	ProductFormLabelStrings["Width"] = lang.L("Width")
	ProductFormLabelStrings["Height"] = lang.L("Height")
	ProductFormLabelStrings["Depth"] = lang.L("Depth")
	ProductFormLabelStrings["Volume"] = lang.L("Volume")
	ProductFormLabelStrings["Weight"] = lang.L("Weight")
}

/* Return a slice with the strings from all the slices without repeating any string */
func Combine(s ...[]string) []string {
	var list []string
	keys := make(map[string]bool)
	for _, slc := range s {
		for _, str := range slc {
			if _, val := keys[str]; !val {
				keys[str] = true
				list = append(list, str)
			}
		}
	}
	return list
}

/* Return a slice without the strings found in r */
func Exclude(s []string, r []string) []string {
	var list []string
	return list
}
