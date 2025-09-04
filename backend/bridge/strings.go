package bridge

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
		"Dimensions",
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
)

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
