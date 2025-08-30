package backend

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Metadata struct {
	b            *Backend
	categoryData map[CatID]*Category
	mfrData      map[MfrID]*Manufacturer
	modelData    map[ModelID]*Model

	catSelection   binding.UntypedList
	mfrSelection   binding.UntypedList
	modelSelection binding.UntypedList

	CatIDList   binding.UntypedList
	CatIDTree   binding.UntypedTree
	Categories  binding.StringList
	MfrIDList   binding.UntypedList
	MfrNameList binding.StringList
	ModelIDList binding.UntypedList
	ProductTree binding.StringTree

	UnitIDList       binding.UntypedList
	ItemStatusIDList binding.UntypedList
}

func NewMetadata(b *Backend) *Metadata {
	return &Metadata{
		b:              b,
		categoryData:   make(map[CatID]*Category),
		mfrData:        make(map[MfrID]*Manufacturer),
		modelData:      make(map[ModelID]*Model),
		catSelection:   binding.NewUntypedList(),
		mfrSelection:   binding.NewUntypedList(),
		modelSelection: binding.NewUntypedList(),

		Categories:       binding.NewStringList(),
		CatIDList:        binding.NewUntypedList(),
		CatIDTree:        binding.NewUntypedTree(),
		MfrIDList:        binding.NewUntypedList(),
		MfrNameList:      binding.NewStringList(),
		ModelIDList:      binding.NewUntypedList(),
		ProductTree:      binding.NewStringTree(),
		UnitIDList:       binding.NewUntypedList(),
		ItemStatusIDList: binding.NewUntypedList(),
	}
}

func (m *Metadata) CreateNewCategory() (id CatID, err error) {
	query := `INSERT INTO Category DEFAULT VALUES`
	res, err := m.b.db.Exec(query)
	if err != nil {
		err = fmt.Errorf("Metadata.CreateNewCategory() error: %w", err)
		return
	}
	i, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("Metadata.CreateNewCategory() error: %w", err)
		return
	}
	id = CatID(i)
	m.UpdateCatList()
	return
}
func (m *Metadata) CreateNewManufacturer() (id MfrID, err error) {
	query := `INSERT INTO Manufacturer DEFAULT VALUES`
	res, err := m.b.db.Exec(query)
	if err != nil {
		err = fmt.Errorf("Metadata.CreateNewManufacturer() error: %w", err)
		return
	}
	i, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("Metadata.CreateNewManufacturer() error: %w", err)
		return
	}
	id = MfrID(i)
	m.GetMfrIDs()
	m.GetProductTree()
	return
}
func (m *Metadata) CreateNewProduct() (id ModelID, err error) {
	query := `INSERT INTO Model DEFAULT VALUES`
	res, err := m.b.db.Exec(query)
	if err != nil {
		err = fmt.Errorf("Metadata.CreateNewModel() error: %w", err)
		return
	}
	i, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("Metadata.CreateNewModel() error: %w", err)
		return
	}
	id = ModelID(i)
	m.GetModelIDs()
	m.GetProductTree()
	return
}
func (m *Metadata) CopyCategory() error {
	// TODO consider looping through selection slice
	sid, err := m.catSelection.GetValue(0)
	if err != nil {
		return fmt.Errorf("Metadata.CopyCategory() error: %w", err)
	}
	selectedCatID := sid.(CatID)
	query := `INSERT INTO Category (PrentID, Name) SELECT ParentID, Name FROM Category WHERE CatID = @0`
	stmt, err := m.b.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Metadata.CopyCategory() error: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(selectedCatID)
	if err != nil {
		return fmt.Errorf("Metadata.CopyCategory() error: %w", err)
	}
	m.UnselectCategory(selectedCatID)
	m.UpdateCatList()
	return err
}
func (m *Metadata) CopyProduct(id ModelID) (newID ModelID, err error) {
	return
}
func (m *Metadata) DeleteCategory() error {
	// TODO consider looping through selection slice
	sid, err := m.catSelection.GetValue(0)
	if err != nil {
		return fmt.Errorf("Metadata.DeleteCategory() error: %w", err)
	}
	selectedCatID := sid.(CatID)
	query := `DELETE FROM Category WHERE CatID = @0`
	stmt, err := m.b.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Metadata.DeleteCategory() error: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(selectedCatID)
	if err != nil {
		return fmt.Errorf("Metadata.DeleteCategory() error: %w", err)
	}
	err = m.UnselectCategory(selectedCatID)
	if err != nil {
		return fmt.Errorf("Metadata.DeleteCategory() error: %w", err)
	}
	err = m.CatIDList.Remove(selectedCatID)
	if err != nil {
		return fmt.Errorf("Metadata.DeleteCategory() error: %w", err)
	}
	return err
}
func (m *Metadata) DeleteManufacturer(id MfrID) error {
	return nil
}
func (m *Metadata) DeleteProduct(id ModelID) error {
	return nil
}
func (m *Metadata) GetCatIDForListItem(index widget.ListItemID) CatID {
	id, err := m.CatIDList.GetValue(index)
	if err != nil {
		log.Printf("Metadata.GetCatIDFor(%d) error: %s", index, err)
		panic(err)
	}
	return id.(CatID)
}
func (m *Metadata) GetListItemIDFor(s string) widget.ListItemID {
	cats, err := m.Categories.Get()
	if err != nil {
		log.Printf("Metadata.GetListItemIDFor(%s) error: %s", s, err)
		panic(err)
	}
	index := slices.IndexFunc(cats, func(n string) bool {
		return strings.TrimSpace(n) == strings.TrimSpace(s)
	})
	return index
}
func (m *Metadata) GetProductIDFor(index widget.TreeNodeID) string {
	id, err := m.ProductTree.GetValue(index)
	if err != nil {
		log.Printf("Metadata.GetProductIDFor(%s) error: %s", index, err)
	}
	return id
}
func (m *Metadata) GetCatIDForTreeItem(index widget.TreeNodeID) CatID {
	id, err := m.CatIDTree.GetValue(index)
	if err != nil {
		log.Printf("Metadata.GetCatIDForTreeItem(%s) error: %s", index, err)
		panic(err)
	}
	return id.(CatID)
}
func (m *Metadata) SelectCategory(id CatID) error {
	return m.catSelection.Append(id)
}
func (m *Metadata) UnselectCategory(id CatID) error {
	return m.catSelection.Remove(id)
}
func (m *Metadata) ClearSelection() error {
	return m.catSelection.Set([]any{})
}
func (m *Metadata) UpdateCatList() error {
	// TODO fix this
	m.GetCatIDTree()
	m.getCatIDList()
	return nil
}

func (m *Metadata) GetProductTree() error {
	// TODO break this thing up
	// TODO also write an *excluding* function (get all except this one)
	query := `SELECT MfrID FROM Manufacturer ORDER BY Name ASC`
	rows, err := m.b.db.Query(query)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		}
	}
	defer rows.Close()
	m.ProductTree.Set(make(map[string][]string), make(map[string]string))
	for rows.Next() {
		var MfrID MfrID
		err := rows.Scan(&MfrID)
		if err != nil {
			panic(err)
		}
		name, _ := MfrID.Name()
		m.ProductTree.Append("", MfrID.TString(), name)
		if MfrID.Branch() {
			query := `SELECT ModelID FROM Model WHERE MfrID = @0 ORDER BY Name ASC`
			rows, err := m.b.db.Query(query, MfrID)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					panic(err)
				}
			}
			defer rows.Close()
			for rows.Next() {
				var ModelID ModelID
				err := rows.Scan(&ModelID)
				if err != nil {
					panic(err)
				}
				name, _ := ModelID.Name()
				m.ProductTree.Append(MfrID.TString(), MfrID.TString()+ModelID.TString(), name)
			}
		}
	}
	rows.Close()
	query = `SELECT ModelID FROM Model WHERE MfrID = 0 ORDER BY NAME ASC`
	rows, err = m.b.db.Query(query)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	for rows.Next() {
		var ModelID ModelID
		err = rows.Scan(&ModelID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		}
		name, _ := ModelID.Name()
		m.ProductTree.Append("", ModelID.TString(), name)
	}
	return nil
}

func (m *Metadata) findCatIDFor(s string) (CatID, error) {
	var i NullInt
	var id CatID

	query := `SELECT CatID FROM Category WHERE Name = @0`
	stmt, err := m.b.db.Prepare(query)
	if err != nil {
		return id, fmt.Errorf("findCatIDFor error: %w", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(s).Scan(&i)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return id, fmt.Errorf("findCatIDFor error: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return id, err
	}

	if !i.Valid {
		log.Printf("findCatIDFor(%s); i.Valid: %t, i.Int: %v", s, i.Valid, i.Int)
		return id, err
	}

	id = CatID(i.Int)
	return id, nil
}
func (m *Metadata) appendCatIDAndChildren(id CatID, spc string) {
	n, _ := id.Name()
	m.CatIDList.Append(id)
	m.Categories.Append(spc + n)
	if !id.Branch() {
		return
	}
	spc += "  "
	for _, child := range id.Children() {
		m.appendCatIDAndChildren(child, spc)
	}
}
func (m *Metadata) getCatIDList() error {
	query := `SELECT CatID FROM Category WHERE ParentID = 0 ORDER BY Name ASC`
	rows, err := m.b.db.Query(query)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	defer rows.Close()
	m.Categories.Set([]string{})
	m.CatIDList.Set([]any{})
	if errors.Is(err, sql.ErrNoRows) {
		return err
	}
	for rows.Next() {
		var CatID CatID
		rows.Scan(&CatID)
		spc := ""
		m.appendCatIDAndChildren(CatID, spc)
	}
	return err
}
func (m *Metadata) GetCatIDTree() error {
	query := `SELECT CatID, ParentID FROM Category ORDER BY Name`
	rows, err := m.b.db.Query(query)
	if err != nil {
		log.Println(err)
		return err
	}
	m.CatIDTree.Set(make(map[string][]string), make(map[string]any))
	for rows.Next() {
		var CatID, ParentID CatID
		rows.Scan(&CatID, &ParentID)
		ps := ""
		if ParentID.String() != "0" {
			ps = ParentID.String()
		}
		m.CatIDTree.Append(ps, CatID.String(), CatID)
	}
	return err
}
func (m *Metadata) GetMfrIDs() {
	query := `SELECT MfrID FROM Manufacturer ORDER BY Name ASC`
	rows, err := m.b.db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	m.MfrIDList.Set([]any{})
	m.MfrNameList.Set([]string{})
	for rows.Next() {
		var MfrID MfrID
		rows.Scan(&MfrID)
		m.MfrIDList.Append(MfrID)
		n, _ := MfrID.Name()
		m.MfrNameList.Append(n)
	}
}
func (m *Metadata) GetModelIDs() {
	query := `SELECT ModelID FROM Model`
	rows, err := m.b.db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	m.ModelIDList.Set([]any{})
	for rows.Next() {
		var ModelID ModelID
		rows.Scan(&ModelID)
		m.ModelIDList.Append(ModelID)
	}
}
func (m *Metadata) getAllUnitIDs() {
	query := `SELECT UnitID FROM Metric`
	rows, err := m.b.db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	m.UnitIDList.Set([]any{})
	for rows.Next() {
		var UnitID UnitID

		rows.Scan(&UnitID)

		m.UnitIDList.Append(UnitID)
	}
}
func (m *Metadata) getAllItemStatusIDs() {
	query := `SELECT ItemStatusID FROM ItemStatus`
	stmt, err := m.b.db.Prepare(query)
	if err != nil {
		log.Printf("Metadata.getAllItemStatusIDs() error: %v", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("Metadata.getAllItemStatusIDs() error: %v", err)
	}
	m.ItemStatusIDList.Set([]any{})
	for rows.Next() {
		var ItemStatusID ItemStatusID

		rows.Scan(&ItemStatusID)

		m.ItemStatusIDList.Append(ItemStatusID)
	}
}
