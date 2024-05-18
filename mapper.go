package sheet

// Mapper is a struct that maps the column index to a field in the struct.
type Mapper[T any] struct {
	Columns map[int]func(*T, any)
	Fields  map[string]int
}

// AddFieldMapping adds a mapping between a column index and a field in the struct.
// example:
//
//	mapper := &Mapper[SpreadsheetRow]{}
//	mapper.AddFieldMapping("Column", 0, func(row *SpreadsheetRow, value any) {
//	  row.Column = value.(string)
//	})
func (m *Mapper[T]) AddFieldMapping(field string, index int, setter func(*T, any)) {
	if m.Columns == nil {
		m.Columns = make(map[int]func(*T, any))
	}
	if m.Fields == nil {
		m.Fields = make(map[string]int)
	}
	m.Columns[index] = setter
	m.Fields[field] = index
}
