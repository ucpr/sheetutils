package sheetutils

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Unmarshaler is an interface that unmarshals the data into a slice of the given type.
type Unmarshaler[T any] interface {
	Unmarshal(data [][]any) ([]T, error)
}

// UnmarshalerImpl is an implementation of Unmarshaler interface.
type UnmarshalerImpl[T any] struct {
	mapper *Mapper[T]
}

// NewUnmarshaler creates a new instance of UnmarshalerImpl.
func NewUnmarshaler[T any]() *UnmarshalerImpl[T] {
	return &UnmarshalerImpl[T]{mapper: &Mapper[T]{}}
}

// AddFieldMapping adds a mapping between a column index and a field in the struct.
func (m *UnmarshalerImpl[T]) AddFieldMapping(field string, index int, setter func(*T, any)) {
	m.mapper.AddFieldSetter(field, index, setter)
}

// Unmarshal the data based on the mapping added by AddFieldMapping method.
func (m *UnmarshalerImpl[T]) Unmarshal(data [][]any) ([]T, error) {
	result := make([]T, 0, len(data))

	for _, row := range data {
		var instance T
		for idx, setter := range m.mapper.Setters {
			if idx < len(row) {
				setter(&instance, row[idx])
			}
		}
		result = append(result, instance)
	}

	return result, nil
}

// Unmarshal unmarshals the data into a slice of the given type.
func Unmarshal[T any](data [][]any) ([]T, error) {
	var tp T
	t := reflect.TypeOf(tp)

	um := NewUnmarshaler[T]()
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
				continue
			}

			fieldName := field.Name
			um.mapper.AddFieldSetter(fieldName, index, func(ptr *T, value any) {
				v := reflect.ValueOf(ptr).Elem().FieldByName(fieldName)
				setFieldValue(v, value, field.Type)
			})
		}
	}

	return um.Unmarshal(data)
}

// setFieldValue sets the value of the field based on the type of the field.
// TODO: Add support for more types.
func setFieldValue(field reflect.Value, value any, fieldType reflect.Type) {
	if field.IsValid() && field.CanSet() {
		v, ok := value.(string)
		if !ok {
			return
		}
		switch fieldType.Kind() {
		case reflect.String:
			field.SetString(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil || field.OverflowInt(n) {
				return
			}
			field.SetInt(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(v, 64)
			if err != nil || field.OverflowFloat(n) {
				return
			}
			field.SetFloat(n)
		case reflect.Bool:
			field.SetBool(v == "TRUE")
		case reflect.Struct:
			// TODO: Add support for more types.
			if fieldType == reflect.TypeOf(time.Time{}) {
				t, err := time.Parse(time.RFC3339, v) // TODO: Add support for more time formats.
				if err == nil {
					field.Set(reflect.ValueOf(t))
				}
			}
		}
	}
}
