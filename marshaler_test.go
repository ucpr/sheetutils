package sheet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshaler_Marshal(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Column1 string
		Column2 int
	}

	patterns := []struct {
		name  string
		setup func() (Marshaler[testStruct], []testStruct)
		want  [][]any
		err   error
	}{
		{
			name: "success",
			setup: func() (Marshaler[testStruct], []testStruct) {
				m := NewMarshaler[testStruct]()
				m.AddFieldGetter("Column1", 0, func(row *testStruct) any {
					return row.Column1
				})
				m.AddFieldGetter("Column2", 1, func(row *testStruct) any {
					return row.Column2
				})

				data := []testStruct{
					{Column1: "value1", Column2: 1},
					{Column1: "value2", Column2: 2},
				}
				return m, data
			},
			want: [][]any{
				{"value1", 1},
				{"value2", 2},
			},
			err: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mapper, data := tt.setup()
			got, err := mapper.Marshal(data)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Column1 string `sheet:"column1,0"`
		Column2 int    `sheet:"column2,1"`
	}

	data := []testStruct{
		{Column1: "value1", Column2: 1},
		{Column1: "value2", Column2: 2},
	}

	got, err := Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, [][]any{
		{"value1", 1},
		{"value2", 2},
	}, got)
}
