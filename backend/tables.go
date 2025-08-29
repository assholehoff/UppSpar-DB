package backend

import (
	"log"
	"slices"
)

// TODO validation and repair functions

/* Table initialisation, validation and repair */
func (b *Backend) listTables() []string {
	var name string
	var tables []string
	stmt, err := b.db.Prepare(`SELECT name FROM sqlite_master WHERE type='table'`)
	if err != nil {
		log.Println("listTables() panic!")
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println("listTables() panic!")
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&name)
		tables = append(tables, name)
	}
	return tables
}
func (b *Backend) createTables() {
	j := b.Journal
	tables := b.listTables()
	touched := false
	if !slices.Contains(tables, "Config") {
		log.Printf("!slices.Contains(tables \"Config\")")
		b.db.Exec(`CREATE TABLE Config(
ConfigKey TEXT PRIMARY KEY,
ConfigVal TEXT)`)
		b.db.Exec(`INSERT INTO Config (ConfigKey, ConfigVal)
VALUES ("ItemIDWidth", "7")`)
		touched = true
	}

	if !slices.Contains(tables, "Item") {
		log.Printf("!slices.Contains(tables \"Item\")")
		b.db.Exec(`CREATE TABLE Item(
-- Proceedo defined column names --
ItemID                  INTEGER PRIMARY KEY AUTOINCREMENT, 
Name                    TEXT DEFAULT 'Nytt föremål', 
Price                   REAL DEFAULT 0, 
Currency                TEXT DEFAULT 'SEK', 
QuantityInPrice         REAL DEFAULT 1, 
Unit                    TEXT DEFAULT 'st', 
OrderMultiple           REAL DEFAULT 0, 
MinOrder                REAL DEFAULT 0, 
Vat                     REAL DEFAULT 0, 
Eta                     INT DEFAULT 0, 
EtaText                 TEXT DEFAULT '', 
Priority                BOOL DEFAULT true, 
Stock                   REAL DEFAULT 1, 
ImgURL1                 TEXT DEFAULT '', 
ImgURL2                 TEXT DEFAULT '', 
ImgURL3                 TEXT DEFAULT '', 
ImgURL4                 TEXT DEFAULT '', 
ImgURL5                 TEXT DEFAULT '', 
SpecsURL                TEXT DEFAULT '', 
UNSPSC                  TEXT DEFAULT '', 
LongDesc                TEXT DEFAULT '', 
Manufacturer            TEXT DEFAULT '', 
MfrItemId               TEXT DEFAULT '', 
GlobId                  TEXT DEFAULT '', 
GlobIdType              TEXT DEFAULT '', 
ReplacesItem            INT DEFAULT 0, 
Questions               TEXT DEFAULT '', 
PackagingCode           BOOL DEFAULT false, 
PresentationCode        BOOL DEFAULT false, 
DeliveryAutoSign        BOOL DEFAULT false, 
DeliveryOption          BOOL DEFAULT false, 
ComparePrice            REAL DEFAULT 0, 
CompareUnit             TEXT DEFAULT '', 
CompareQuantityInPrice  REAL DEFAULT 0, 
PriceInfo               TEXT DEFAULT '', 
AddDesc                 TEXT DEFAULT '', 
ProcFlow                TEXT DEFAULT '', 
InnerUnit               TEXT DEFAULT '', 
QuantityInUnit          REAL DEFAULT 0, 
RiskClassification      TEXT DEFAULT '', 
Comment                 TEXT DEFAULT '', 
EnvClassification       TEXT DEFAULT '', 
FormId                  TEXT DEFAULT '', 
Article                 TEXT DEFAULT '', 
Attachments             BOOL DEFAULT false, 
ItemGroup               TEXT DEFAULT '', 
-- Custom defined fields --
MfrID                   INT DEFAULT 0, 
ModelID                 INT DEFAULT 0, 
ModelName               TEXT DEFAULT '',
ModelDesc               TEXT DEFAULT '', 
ModelURL                TEXT DEFAULT '',
Notes                   TEXT DEFAULT '', 
Width                   REAL DEFAULT 0, 
Height                  REAL DEFAULT 0, 
Depth                   REAL DEFAULT 0, 
Volume                  REAL DEFAULT 0, 
Weight                  REAL DEFAULT 0, 
LengthUnitID            INT DEFAULT 2, 
VolumeUnitID            INT DEFAULT 11, 
WeightUnitID            INT DEFAULT 7, 
CatID                   INT DEFAULT 1, 
GroupID                 INT DEFAULT 0, 
StorageID               INT DEFAULT 0, 
ItemStatusID            INT DEFAULT 1, 
ItemConditionID         INT DEFAULT 0, 
DateCreated             TEXT DEFAULT(datetime('now', 'subsec')), 
DateModified            TEXT DEFAULT(datetime('now', 'subsec')), 
FOREIGN KEY(MfrID) REFERENCES Manufacturer(MfrID), 
FOREIGN KEY(ModelID) REFERENCES Model(ModelID), 
FOREIGN KEY(LengthUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(VolumeUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(WeightUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(CatID) REFERENCES Category(CatID), 
FOREIGN KEY(GroupID) REFERENCES Item_Group(GroupID), 
FOREIGN KEY(ItemStatusID) REFERENCES ItemStatus(ItemStatusID),  
FOREIGN KEY(StorageID) REFERENCES Storage(StorageID))`)
		b.db.Exec(`CREATE TRIGGER UpdateDateModified
AFTER UPDATE ON Item FOR EACH ROW
BEGIN
UPDATE Item SET DateModified = datetime('now', 'subsec') WHERE ItemID = old.ItemID;
END`)
		touched = true
	}
	if !slices.Contains(tables, "Temp_Item") {
		log.Printf("!slices.Contains(tables \"Temp_Item\")")
		b.db.Exec(`CREATE TABLE Temp_Item(
-- Proceedo defined column names --
ItemID                  INT, 
Name                    TEXT, 
Price                   REAL, 
Currency                TEXT, 
QuantityInPrice         REAL, 
Unit                    TEXT, 
OrderMultiple           REAL, 
MinOrder                REAL, 
Vat                     REAL, 
Eta                     INT, 
EtaText                 TEXT, 
Priority                BOOL, 
Stock                   REAL, 
ImgURL1                 TEXT, 
ImgURL2                 TEXT, 
ImgURL3                 TEXT, 
ImgURL4                 TEXT, 
ImgURL5                 TEXT, 
SpecsURL                TEXT, 
UNSPSC                  TEXT, 
LongDesc                TEXT, 
Manufacturer            TEXT, 
MfrItemId               TEXT, 
GlobId                  TEXT, 
GlobIdType              TEXT, 
ReplacesItem            INT, 
Questions               TEXT, 
PackagingCode           BOOL, 
PresentationCode        BOOL, 
DeliveryAutoSign        BOOL, 
DeliveryOption          BOOL, 
ComparePrice            REAL, 
CompareUnit             TEXT, 
CompareQuantityInPrice  REAL, 
PriceInfo               TEXT, 
AddDesc                 TEXT, 
ProcFlow                TEXT, 
InnerUnit               TEXT, 
QuantityInUnit          REAL, 
RiskClassification      TEXT, 
Comment                 TEXT, 
EnvClassification       TEXT, 
FormId                  TEXT, 
Article                 TEXT, 
Attachments             BOOL, 
ItemGroup               TEXT, 
-- Custom defined fields --
MfrID                   INT, 
ModelID                 INT, 
ModelName               TEXT,
ModelDesc               TEXT, 
ModelURL                TEXT,
Notes                   TEXT, 
Width                   REAL, 
Height                  REAL, 
Depth                   REAL, 
Volume                  REAL, 
Weight                  REAL, 
LengthUnitID            INT, 
VolumeUnitID            INT, 
WeightUnitID            INT, 
CatID                   INT, 
GroupID                 INT, 
StorageID               INT, 
ItemStatusID            INT, 
ItemConditionID         INT, 
DateCreated             TEXT, 
DateModified            TEXT, 
FOREIGN KEY(MfrID) REFERENCES Manufacturer(MfrID), 
FOREIGN KEY(ModelID) REFERENCES Model(ModelID), 
FOREIGN KEY(LengthUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(VolumeUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(WeightUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(CatID) REFERENCES Category(CatID), 
FOREIGN KEY(GroupID) REFERENCES Item_Group(GroupID), 
FOREIGN KEY(ItemStatusID) REFERENCES ItemStatus(ItemStatusID),  
FOREIGN KEY(StorageID) REFERENCES Storage(StorageID))`)
		b.db.Exec(`CREATE TRIGGER Temp_UpdateDateModified
AFTER UPDATE ON Temp_Item FOR EACH ROW
BEGIN
UPDATE Temp_Item SET DateModified = datetime('now', 'subsec') WHERE ItemID = old.ItemID;
END`)
		touched = true
	}

	if !slices.Contains(tables, "Item_Condition") {
		log.Printf("!slices.Contains(tables \"Item_Condition\")")
		b.db.Exec(`CREATE TABLE Item_Condition(
ItemID INT, 
Rate INT, 
Comment TEXT, 
FOREIGN KEY(ItemID) REFERENCES Item(ItemID) ON DELETE CASCADE)`)
		touched = true
	}
	if !slices.Contains(tables, "Item_Group") {
		log.Printf("!slices.Contains(tables \"Item_Group\")")
		b.db.Exec(`CREATE TABLE Item_Group(
GroupID INTEGER PRIMARY KEY AUTOINCREMENT,
ParentID INT DEFAULT 0,
Name TEXT DEFAULT '',
Deleted BOOL DEFAULT false)`)
		touched = true
	}
	if !slices.Contains(tables, "Item_Function") {
		log.Printf("!slices.Contains(tables \"Item_Function\")")
		b.db.Exec(`CREATE TABLE Item_Function(
ItemID INT, 
FuncID INT, 
IsTested BOOL, 
IsWorking BOOL, 
Comment TEXT, 
FOREIGN KEY(ItemID) REFERENCES Item(ItemID) ON DELETE CASCADE, 
FOREIGN KEY(FuncID) REFERENCES Function_Data(FuncID))`)
		touched = true
	}
	if !slices.Contains(tables, "Function_Data") {
		log.Printf("!slices.Contains(tables \"Function_Data\")")
		b.db.Exec(`CREATE TABLE Function_Data(
FuncID INTEGER PRIMARY KEY AUTOINCREMENT, 
Name TEXT)`)
		touched = true
	}
	if !slices.Contains(tables, "ItemStatus") {
		log.Printf("!slices.Contains(tables \"ItemStatus\")")
		b.db.Exec(`CREATE TABLE ItemStatus(
ItemStatusID INTEGER PRIMARY KEY AUTOINCREMENT, 
Name TEXT)`)
		b.db.Exec(`INSERT INTO ItemStatus (Name) 
VALUES ("available"), ("sold"), ("archived"), ("deleted")`)
		touched = true
	}
	if !slices.Contains(tables, "Manufacturer") {
		log.Printf("!slices.Contains(tables \"Manufacturer\")")
		b.db.Exec(`CREATE TABLE Manufacturer(
MfrID INTEGER PRIMARY KEY AUTOINCREMENT, 
Name TEXT DEFAULT 'Ny tillverkare',
Deleted BOOL DEFAULT false)`)
		b.db.Exec(`INSERT INTO Manufacturer (Name) 
VALUES ("UppSpar"), ("IKEA"), ("Kinnarps")`)
		touched = true
	}
	if !slices.Contains(tables, "Model") {
		log.Printf("!slices.Contains(tables \"Model\")")
		b.db.Exec(`CREATE TABLE Model(
ModelID      INTEGER PRIMARY KEY AUTOINCREMENT, 
Name         TEXT DEFAULT 'Ny modell', 
Manufacturer TEXT DEFAULT '',
MfrID        INT DEFAULT 0, 
Desc         TEXT DEFAULT '', 
ImgURL1      TEXT DEFAULT '', 
ImgURL2      TEXT DEFAULT '', 
ImgURL3      TEXT DEFAULT '', 
ImgURL4      TEXT DEFAULT '', 
ImgURL5      TEXT DEFAULT '', 
SpecsURL     TEXT DEFAULT '', 
ModelURL     TEXT DEFAULT '', 
Width        REAL DEFAULT 0, 
Height       REAL DEFAULT 0, 
Depth        REAL DEFAULT 0, 
Volume       REAL DEFAULT 0, 
Weight       REAL DEFAULT 0, 
LengthUnitID INT DEFAULT 2, 
VolumeUnitID INT DEFAULT 11, 
WeightUnitID INT DEFAULT 7, 
CatID        INT DEFAULT 1, 
Deleted      BOOL DEFAULT false, 
FOREIGN KEY(MfrID) REFERENCES Manufacturer(MfrID)
FOREIGN KEY(LengthUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(VolumeUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(WeightUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(CatID) REFERENCES Category(CatID))`)
		touched = true
	}

	if !slices.Contains(tables, "Category") {
		log.Printf("!slices.Contains(tables \"Category\")")
		b.createCategoryTable()
		touched = true
	}
	if !slices.Contains(tables, "Category_Config") {
		log.Printf("!slices.Contains(tables \"Category_Config\")")
		b.db.Exec(`CREATE TABLE Category_Config(
CatID INT, 
ConfigKey TEXT, 
ConfigVal BOOL, 
FOREIGN KEY(CatID) REFERENCES Category(CatID) ON DELETE CASCADE)`)
		touched = true
	}
	if !slices.Contains(tables, "Category_Data") {
		log.Printf("!slices.Contains(tables \"Category_Data\")")
		b.db.Exec(`CREATE TABLE Category_Data(
CatID INT, 
DataKey TEXT, 
DataVal TEXT, 
FOREIGN KEY(CatID) REFERENCES Category(CatID) ON DELETE CASCADE)`)
		touched = true
	}

	if !slices.Contains(tables, "Image") {
		b.db.Exec(`CREATE TABLE Image(
ImgID INTEGER PRIMARY KEY AUTOINCREMENT, 
ImgData BLOB, 
ImgThumb BLOB, 
ImgURL TEXT DEFAULT '', 
Deleted BOOL DEFAULT false),`)
		touched = true
	}

	if !slices.Contains(tables, "Metric") {
		log.Printf("!slices.Contains(tables \"Metric\")")
		b.createMetricTable()
		touched = true
	}

	if !slices.Contains(tables, "SearchWords_Association") {
		log.Printf("!slices.Contains(tables \"SearchWords_Association\")")
		b.db.Exec(`CREATE TABLE SearchWords_Association(
ItemID INT, 
WordID INT, 
FOREIGN KEY(ItemID) REFERENCES SearchWords_Vocabulary(WordID) ON DELETE CASCADE, 
FOREIGN KEY(WordID) REFERENCES Item(ItemID))`)
		touched = true
	}
	if !slices.Contains(tables, "SearchWords_Vocabulary") {
		log.Printf("!slices.Contains(tables \"SearchWords_Vocabulary\")")
		b.db.Exec(`CREATE TABLE SearchWords_Vocabulary(
WordID INTEGER PRIMARY KEY, 
WordString TEXT)`)
		touched = true
	}

	if !slices.Contains(tables, "Storage") {
		log.Printf("!slices.Contains(tables \"Storage\")")
		b.db.Exec(`CREATE TABLE Storage(
StorageID INTEGER PRIMARY KEY, 
Place TEXT, 
Comment TEXT)`)
		touched = true
	}

	if !slices.Contains(tables, "WishList") {
		log.Printf("!slices.Contains(tables \"WishList\")")
		b.db.Exec(`CREATE TABLE WishList(
WishID INTEGER PRIMARY KEY, 
ContactID TEXT, 
WishItemID TEXT, 
Stock INT 
DateCreated TEXT DEFAULT(datetime('now', 'subsec')), 
DateModified TEXT DEFAULT(datetime('now', 'subsec')), 
DateExpires TEXT, 
FOREIGN KEY(ContactID) REFERENCES WishList_Contact(ContactID), 
FOREIGN KEY(WishItemID) REFERENCES WishList_Item(WishItemID))`)
		touched = true
	}
	if !slices.Contains(tables, "WishList_Contact") {
		log.Printf("!slices.Contains(tables \"WishList_Contact\")")
		b.db.Exec(`CREATE TABLE WishList_Contact(
ContactID INTEGER PRIMARY KEY, 
FirstName TEXT DEFAULT 'Förnamn', 
LastName TEXT DEFAULT 'Efternamn', 
Email TEXT, 
Phone TEXT, 
Comment TEXT, 
DateCreated TEXT DEFAULT(datetime('now', 'subsec')), 
DateModified TEXT DEFAULT(datetime('now', 'subsec')), 
DateExpires TEXT)`)
		touched = true
	}
	if !slices.Contains(tables, "WishList_Item") {
		log.Printf("!slices.Contains(tables \"WishList_Item\")")
		b.db.Exec(`CREATE TABLE WishList_Item(
WishItemID INTEGER PRIMARY KEY AUTOINCREMENT, 
Name TEXT DEFAULT 'Nytt föremål', 
Comment TEXT, 
Width REAL, 
Height REAL, 
Depth REAL, 
Weight REAL, 
LengthUnitID INT, 
WeightUnitID INT, 
CatID INT, 
DateCreated TEXT DEFAULT(datetime('now', 'subsec')), 
DateModified TEXT DEFAULT(datetime('now', 'subsec')), 
DateExpires TEXT, 
FOREIGN KEY(LengthUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(WeightUnitID) REFERENCES Metric(UnitID), 
FOREIGN KEY(CatID) REFERENCES Category(CatID))`)
		touched = true
	}
	if !slices.Contains(tables, "WishList_Item_Function") {
		log.Printf("!slices.Contains(tables \"WishList_Item_Function\")")
		b.db.Exec(`CREATE TABLE WishList_Item_Function(
WishItemID INT, 
FuncID INT, 
Comment TEXT, 
FOREIGN KEY(WishItemID) REFERENCES WishList_Item(WishItemID) ON DELETE CASCADE,  
FOREIGN KEY(FuncID) REFERENCES Function_Data(FuncID))`)
		touched = true
	}

	if touched {
		log.Printf("touched!")
		j.NewMessage("Skapade nya tabeller för föremål i databasen.")
		j.Refresh()
	}
}

func (b *Backend) createCategoryTable() {
	b.db.Exec(`CREATE TABLE Category(
CatID INTEGER PRIMARY KEY AUTOINCREMENT, 
ParentID INT DEFAULT 0,
Name TEXT DEFAULT 'Ny kategori')`)
	b.db.Exec(`INSERT INTO Category (Name, ParentID) 
VALUES  ("Administration", 0), 
        ("Hushåll", 0), 
        ("Kontor", 0), 
        ("Tjänster", 0), 
        ("Övrigt", 0),
        ("Badrum", 2), 
        ("Belysning", 5), 
        ("Bord", 2), 
        ("Dekor", 5), 
        ("Elektronik", 3),  
        ("Förvaring", 3), 
        ("Husgeråd", 2), 
        ("Hylla", 3), 
        ("Kök & vitvaror", 2), 
        ("Textilier & mattor", 5), 
        ("Skrivbord", 3), 
        ("Skåp", 3), 
        ("Soffor & fåtöljer", 2), 
        ("Stolar", 3), 
        ("Tvätt & städ", 5), 
        ("Sängar & madrasser", 2)`)
}

func (b *Backend) createMetricTable() {
	b.db.Exec(`CREATE TABLE Metric(
UnitID INTEGER PRIMARY KEY, 
Text TEXT)`)
	b.db.Exec(`INSERT INTO Metric (Text) VALUES ("mm"), ("cm"), ("dm"), ("m"), ("g"), ("hg"), ("kg"), ("ml"), ("cl"), ("dl"), ("l")`)
}
