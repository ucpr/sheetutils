package sheetutils

import (
	"reflect"
	"strconv"
	"strings"
)

type Marshaler[T any] interface {
	Marshal(data []T) ([][]any, error)
}

type MarshalerImpl[T any] struct {
	mapper *Mapper[T]
}

func NewMarshaler[T any]() *MarshalerImpl[T] {
	return &MarshalerImpl[T]{mapper: &Mapper[T]{}}
}

func (m *MarshalerImpl[T]) AddFieldGetter(field string, index int, getter func(*T) any) {
	m.mapper.AddFieldGetter(field, index, getter)
}

func (m *MarshalerImpl[T]) Marshal(data []T) ([][]any, error) {
	result := make([][]any, 0, len(data))

	for _, instance := range data {
		row := make([]any, 0, len(m.mapper.Getters))
		for _, getter := range m.mapper.Getters {
			row = append(row, getter(&instance))
		}
		result = append(result, row)
	}

	return result, nil
}

func Marshal[T any](data []T) ([][]any, error) {
	var tp T
	t := reflect.TypeOf(tp)

	m := NewMarshaler[T]()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(sheetTag)
		if tag != "" {
			tagParts := strings.Split(tag, sheetTagSeparator)
			if len(tagParts) != 2 {
				continue
			}
			index, err := strconv.Atoi(tagParts[1])
			if err != nil {
				return nil, err
			}
			m.AddFieldGetter(field.Name, index, func(row *T) any {
				return reflect.ValueOf(row).Elem().Field(i).Interface()
			})
		}
	}

	return m.Marshal(data)
}
