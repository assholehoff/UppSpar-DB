package backend

import (
	"log"

	"github.com/xuri/excelize/v2"
)

/* Export all items with ItemStatusAvailable */
func (m *Items) ExportExcel(p string) {
	f := excelize.NewFile()
	defer f.Close()
	f.Path = p

	f.SetSheetName("Sheet1", "Data")
	sw, err := f.NewStreamWriter("Data")
	if err != nil {
		log.Printf("Items.ExportExcel(%s) query error: %s", p, err)
	}

	/* Write R1 headers */
	if err := sw.SetRow("A1",
		[]any{
			excelize.Cell{Value: "Artikelnummer*"},                 // ItemId           *Obligatoriskt fält*
			excelize.Cell{Value: "Produktbenämning*"},              // Name             *Obligatoriskt fält*
			excelize.Cell{Value: "Pris*"},                          // Price            *Obligatoriskt fält*
			excelize.Cell{Value: "Valuta*"},                        // Currency         *Obligatoriskt fält*
			excelize.Cell{Value: "Antal enheter i pris*"},          // QuantityInPrice
			excelize.Cell{Value: "Säljenhet*"},                     // Unit             *Obligatoriskt fält*
			excelize.Cell{Value: "Beställs i multiplar av*"},       // OrderMultiple
			excelize.Cell{Value: "Minsta beställningskvantitet*"},  // MinOrder
			excelize.Cell{Value: "Momssats*"},                      // Vat              *Obligatoriskt fält*
			excelize.Cell{Value: "Antal dagar för leverans"},       // Eta
			excelize.Cell{Value: "Leveransbeskr. (ers. dagar)"},    // EtaText
			excelize.Cell{Value: "Bassortiment (\"tumme upp\")*"},  // Priority         *Obligatoriskt fält*
			excelize.Cell{Value: "Saldo"},                          // Stock
			excelize.Cell{Value: "Sökord"},                         // SearchWords
			excelize.Cell{Value: "Webblänk till bild"},             // ImgURL1
			excelize.Cell{Value: "Webblänk till bild 2"},           // ImgURL2
			excelize.Cell{Value: "Webblänk till bild 3"},           // ImgURL3
			excelize.Cell{Value: "Webblänk till bild 4"},           // ImgURL4
			excelize.Cell{Value: "Webblänk till bild 5"},           // ImgURL5
			excelize.Cell{Value: "Webblänk till produktblad"},      // SpecsURL
			excelize.Cell{Value: "UNSPSC (00.00.00.00)"},           // UNSPSC
			excelize.Cell{Value: "Utförligare beskrivning"},        // LongDesc
			excelize.Cell{Value: "Tillverkare"},                    // Manufacturer
			excelize.Cell{Value: "Tillverkarens artnr."},           // MfrItemId
			excelize.Cell{Value: "Globalt ID"},                     // GlobId
			excelize.Cell{Value: "Kvalificerare globalt ID"},       // GlobIdType
			excelize.Cell{Value: "Ersätter artikelnummer"},         // ReplacesItem
			excelize.Cell{Value: "Tillhör produkt"},                // SubItemOf
			excelize.Cell{Value: "Tilläggsfrågor"},                 // Questions
			excelize.Cell{Value: "Förpackas*"},                     // PackagingCode
			excelize.Cell{Value: "Presentation*"},                  // PresentationCode
			excelize.Cell{Value: "Automatisk leveranskvittens"},    // DeliveryAutoSign
			excelize.Cell{Value: "Visa i alternativ best."},        // DeliveryOption
			excelize.Cell{Value: "Jämförelsepris"},                 // ComparePrice
			excelize.Cell{Value: "Enhetstyp i jämförelsepris"},     // CompareUnit
			excelize.Cell{Value: "Antal enheter i jfrpris"},        // CompareQuantityInPrice
			excelize.Cell{Value: "Prisinformation"},                // PriceInfo
			excelize.Cell{Value: "Extra beskrivningsfält"},         // AddDesc
			excelize.Cell{Value: "Flöde*"},                         // ProcFlow
			excelize.Cell{Value: "Inre enhetstyp"},                 // InnerUnit
			excelize.Cell{Value: "Antal inre enheter i säljenhet"}, // QuantityInUnit
			excelize.Cell{Value: "Riskbeskrivning"},                // RiskClassification
			excelize.Cell{Value: "Kommentar till beställare"},      // Comment
			excelize.Cell{Value: "Miljömärkning"},                  // EnvClassification
			excelize.Cell{Value: "Formulär"},                       // FormId
			excelize.Cell{Value: "Artikeltyp "},                    // Article
			excelize.Cell{Value: "Bifoga filer"},                   // Attachments
			excelize.Cell{Value: "Produktgrupp"},                   // ItemGroup
		}); err != nil {
		log.Printf("Items.ExportExcel(%s) error: %s", p, err)
		return
	}

	/* Fetch ItemIds for all items set to be exported */
	var ids []ItemID
	query := `SELECT ItemID FROM Item WHERE ItemStatusID = 1`
	rows, err := m.db.Query(query)
	if err != nil {
		log.Printf("Items.ExportExcel(%s) error: %s", p, err)
	}

	for rows.Next() {
		var id NullInt
		rows.Scan(&id)
		if id.Valid {
			ids = append(ids, ItemID(id.Int))
		}
	}

	/* Iterate over items, add each one as a row */
	for i, id := range ids {
		row := make([]any, 48)
		row[0] = id.String()
		row[1] = valueOrVoid(id, "Name")            // *Obligatoriskt fält*
		row[2] = valueOrVoid(id, "Price")           // *Obligatoriskt fält*
		row[3] = valueOrVoid(id, "Currency")        // *Obligatoriskt fält* SEK
		row[4] = valueOrVoid(id, "QuantityInPrice") //
		row[5] = valueOrVoid(id, "Unit")            // *Obligatoriskt fält*
		row[6] = valueOrVoid(id, "OrderMultiple")   //
		row[7] = valueOrVoid(id, "MinOrder")        //
		row[8] = valueOrVoid(id, "Vat")             // *Obligatoriskt fält*
		row[9] = valueOrVoid(id, "Eta")             //
		row[10] = valueOrVoid(id, "EtaText")        //
		row[11] = valueOrVoid(id, "Priority")       // *Obligatoriskt fält* [Y|N]
		row[12] = valueOrVoid(id, "Stock")
		row[13] = valueOrVoid(id, "SearchWords") // TODO compile from search words list
		row[14] = valueOrVoid(id, "ImgURL1")
		row[15] = valueOrVoid(id, "ImgURL2")
		row[16] = valueOrVoid(id, "ImgURL3")
		row[17] = valueOrVoid(id, "ImgURL4")
		row[18] = valueOrVoid(id, "ImgURL5")
		row[19] = valueOrVoid(id, "SpecsURL")
		row[20] = valueOrVoid(id, "UNSPSC")
		row[21] = valueOrVoid(id, "LongDesc")
		row[22] = valueOrVoid(id, "Manufacturer")
		row[23] = valueOrVoid(id, "MfrItemId")
		row[24] = valueOrVoid(id, "GlobId")
		row[25] = valueOrVoid(id, "GlobIdType")
		row[26] = valueOrVoid(id, "ReplacesItem")
		row[27] = valueOrVoid(id, "SubItemOf") // TODO pull from separate table
		row[28] = valueOrVoid(id, "Questions")
		row[29] = valueOrVoid(id, "PackagingCode")
		row[30] = valueOrVoid(id, "PresentationCode")
		row[31] = valueOrVoid(id, "DeliveryAutoSign")
		row[32] = valueOrVoid(id, "DeliveryOption")
		row[33] = valueOrVoid(id, "ComparePrice")
		row[34] = valueOrVoid(id, "CompareUnit")
		row[35] = valueOrVoid(id, "CompareQuantityInPrice")
		row[36] = valueOrVoid(id, "PriceInfo")
		row[37] = valueOrVoid(id, "AddDesc")
		row[38] = valueOrVoid(id, "ProcFlow")
		row[39] = valueOrVoid(id, "InnerUnit")
		row[40] = valueOrVoid(id, "QuantityInUnit")
		row[41] = valueOrVoid(id, "RiskClassification")
		row[42] = valueOrVoid(id, "Comment")
		row[43] = valueOrVoid(id, "EnvClassification")
		row[44] = valueOrVoid(id, "FormId")
		row[45] = valueOrVoid(id, "Article")
		row[46] = valueOrVoid(id, "Attachments")
		row[47] = valueOrVoid(id, "ItemGroup")

		cell, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			log.Printf("Items.ExportExcel(%s) error: %s", p, err)
		}
		if err := sw.SetRow(cell, row); err != nil {
			log.Printf("Items.ExportExcel(%s) error: %s", p, err)
		}
	}

	/* Flush stream */
	if err := sw.Flush(); err != nil {
		log.Printf("Items.ExportExcel(%s) error: %s", p, err)
	}
	/* Save file */
	if err := f.SaveAs(f.Path); err != nil {
		log.Printf("Items.ExportExcel(%s) error: %s", p, err)
	}
}

/* If value equals the zero value, return an empty any */
func valueOrVoid(id ItemID, key string) (val any) {
	switch key {
	case "Name":
		val, _ = id.Name()
	case "Price":
		val, _ = id.Price()
	case "Currency":
		val, _ = id.Currency()
	case "Unit":
		val, _ = id.Unit()
	case "Vat":
		val, _ = id.Vat()
	case "Priority":
		val = "Y"
	case "Stock":
		val, _ = id.Stock()
	case "ImgURL1":
		val, _ = id.ImgURL1()
	case "ImgURL2":
		val, _ = id.ImgURL2()
	case "ImgURL3":
		val, _ = id.ImgURL3()
	case "ImgURL4":
		val, _ = id.ImgURL4()
	case "ImgURL5":
		val, _ = id.ImgURL5()
	case "SpecsURL":
		val, _ = id.SpecsURL()
	case "LongDesc":
		val, _ = id.LongDesc()
	case "Manufacturer":
		val, _ = id.Manufacturer()
	case "AddDesc":
		val, _ = id.AddDesc()
	default:
		return
	}
	return
}
