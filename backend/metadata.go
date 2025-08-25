package backend

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Metadata struct {
	b            *Backend
	categoryData map[CatID]*Category
	mfrData      map[MfrID]*Manufacturer
	modelData    map[ModelID]*Model

	catSelection binding.UntypedList
	CatIDList    binding.UntypedList
	CatIDTree    binding.UntypedTree

	mfrSelection binding.UntypedList
	MfrIDList    binding.UntypedList
	MfrIDTree    binding.UntypedTree

	// TODO Merge Mfr and Model to a binding.UntypedTree

	modelSelection binding.UntypedList
	ModelIDList    binding.UntypedList

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

		CatIDList:        binding.NewUntypedList(),
		CatIDTree:        binding.NewUntypedTree(),
		MfrIDList:        binding.NewUntypedList(),
		ModelIDList:      binding.NewUntypedList(),
		UnitIDList:       binding.NewUntypedList(),
		ItemStatusIDList: binding.NewUntypedList(),
	}
}

func (m *Metadata) CreateNewCategory() error {
	query := `INSERT INTO Category DEFAULT VALUES`
	stmt, err := m.b.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Metadata.CreateNewCategory() error: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("Metadata.CreateNewCategory() error: %w", err)
	}
	m.UpdateCatList()
	return err
}
func (m *Metadata) CopyCategory() error {
	// TODO consider looping through selection slice
	sid, err := m.catSelection.GetValue(0)
	if err != nil {
		return fmt.Errorf("Metadata.CopyCategory() error: %w", err)
	}
	log.Printf("copying %d from selection slice", sid)
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
	res, err := stmt.Exec(selectedCatID)
	if err != nil {
		return fmt.Errorf("Metadata.DeleteCategory() error: %w", err)
	}
	raf, _ := res.RowsAffected()
	log.Printf("%d rows affected", raf)
	m.UnselectCategory(selectedCatID)
	m.CatIDList.Remove(selectedCatID)
	return err
}
func (m *Metadata) GetCatIDForListItem(index widget.ListItemID) CatID {
	id, err := m.CatIDList.GetValue(index)
	if err != nil {
		log.Println("Metadata.GetCatIDFor(index widget.ListItemID) panic!")
		panic(err)
	}
	return id.(CatID)
}
func (m *Metadata) GetCatIDForTreeItem(index widget.TreeNodeID) CatID {
	id, err := m.CatIDTree.GetValue(index)
	if err != nil {
		log.Println("Metadata.GetCatIDForTreeItem(index widget.TreeNodeID) panic!")
		panic(err)
	}
	return id.(CatID)
}
func (m *Metadata) SelectCategory(id CatID) error {
	log.Printf("adding %d to selection slice", id)
	return m.catSelection.Append(id)
}
func (m *Metadata) UnselectCategory(id CatID) error {
	log.Printf("removing %d from selection slice", id)
	return m.catSelection.Remove(id)
}
func (m *Metadata) ClearSelection() error {
	return m.catSelection.Set([]any{})
}
func (m *Metadata) UpdateCatList() error {
	// TODO
	return m.getAllCatIDs()
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
	if err != nil {
		return id, fmt.Errorf("findCatIDFor error: %w", err)
	}

	if !i.Valid {
		log.Printf("findCatIDFor(s) i: %v", i)
		return id, ErrNotFound
	}

	id = CatID(i.Int)
	return id, nil
}
func (m *Metadata) getAllCatIDs() error {
	query := `SELECT CatID FROM Category ORDER BY Name`
	rows, err := m.b.db.Query(query)
	if err != nil {
		log.Println(err)
	}
	m.CatIDList.Set([]any{})
	for rows.Next() {
		var CatID CatID

		rows.Scan(&CatID)

		m.CatIDList.Append(CatID)
	}
	return err
}
func (m *Metadata) getCatIDTree() error {
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
func (m *Metadata) getAllMfrIDs() {
	query := `SELECT MfrID FROM Manufacturer`
	rows, err := m.b.db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	m.MfrIDList.Set([]any{})
	for rows.Next() {
		var MfrID MfrID
		rows.Scan(&MfrID)
		m.MfrIDList.Append(MfrID)
	}
}
func (m *Metadata) getAllModelIDs() {
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
		m.MfrIDList.Append(ModelID)
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
