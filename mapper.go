package sheet

// Mapper is a struct that maps the column index to a field in the struct.
type Mapper[T any] struct {
	Setters map[int]func(*T, any)
	Getters map[int]func(*T) any
}

// AddFieldSetter adds a mapping between a column index and a field in the struct.
// example:
//
//	mapper := &Mapper[SpreadsheetRow]{}
//	mapper.AddFieldSetter("Column", 0, func(row *SpreadsheetRow, value any) {
//	  row.Column = value.(string)
//	})
func (m *Mapper[T]) AddFieldSetter(field string, index int, setter func(*T, any)) {
	if m.Setters == nil {
		m.Setters = make(map[int]func(*T, any))
	}
	m.Setters[index] = setter
}

func (m *Mapper[T]) AddFieldGetter(field string, index int, getter func(*T) any) {
	if m.Getters == nil {
		m.Getters = make(map[int]func(*T) any)
	}
	m.Getters[index] = getter
}
