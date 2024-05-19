package sheetutils

import (
	"reflect"
	"strconv"
	"strings"
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
		switch fieldType.Kind() {
		case reflect.String:
			strValue, ok := value.(string)
			if ok {
				field.SetString(strValue)
			}
		case reflect.Int:
			intValue, ok := value.(int)
			if ok {
				field.SetInt(int64(intValue))
			}
		case reflect.Float64:
			floatValue, ok := value.(float64)
			if ok {
				field.SetFloat(floatValue)
			}
		case reflect.Bool:
			boolValue, ok := value.(bool)
			if ok {
				field.SetBool(boolValue)
			}
		}
	}
}
