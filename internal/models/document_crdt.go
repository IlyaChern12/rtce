package models

// Char — элемент CRDT для символа документа
type Char struct {
	ID     string `json:"id"`     // уникальный идентификатор символа
	Value  string `json:"value"`  // символ
	PrevID string `json:"prevId"` // ID предыдущего символа для упорядочивания
}

// NewDocumentCRDT создаёт новый документ с пустым CRDT
func NewDocumentCRDT(id string) *Document {
	return &Document{
		ID:     id,
		CRDT:   []Char{},
		Version: 0,
	}
}

// Insert добавляет символ в CRDT
func (d *Document) Insert(c Char) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.CRDT = append(d.CRDT, c)
	d.Version++
}

// Delete удаляет символ по ID из CRDT
func (d *Document) Delete(charID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, ch := range d.CRDT {
		if ch.ID == charID {
			d.CRDT = append(d.CRDT[:i], d.CRDT[i+1:]...)
			d.Version++
			break
		}
	}
}