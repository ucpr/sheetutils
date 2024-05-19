package sheetutils

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshaler_Unmarshal(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Column1 string
		Column2 int
	}

	patterns := []struct {
		name  string
		setup func() (Unmarshaler[testStruct], [][]any)
		want  []testStruct
		err   error
	}{
		{
			name: "success",
			setup: func() (Unmarshaler[testStruct], [][]any) {
				um := NewUnmarshaler[testStruct]()
				um.AddFieldMapping("Column1", 0, func(row *testStruct, value any) {
					row.Column1 = value.(string)
				})
				um.AddFieldMapping("Column2", 1, func(row *testStruct, value any) {
					v, ok := value.(string)
					if !ok {
						return
					}
					i, _ := strconv.ParseInt(v, 10, 64)
					row.Column2 = int(i)
				})

				data := [][]any{
					{"value1", "1"},
					{"value2", "2"},
				}

				return um, data
			},
			want: []testStruct{
				{Column1: "value1", Column2: 1},
				{Column1: "value2", Column2: 2},
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mapper, data := tt.setup()
			got, err := mapper.Unmarshal(data)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Column1 string    `sheet:"column1,0"`
		Column2 int       `sheet:"column2,1"`
		Column3 float64   `sheet:"column3,2"`
		Column4 bool      `sheet:"column4,3"`
		Column5 time.Time `sheet:"column5,4"`
	}

	data := [][]any{
		{"value1", "1", "1.1", "TRUE", "2024-05-13T09:07:44+09:00"},
		{"value2", "2", "2.2", "FALSE", "2024-05-13T09:07:44+09:00"},
	}

	tm, err := time.Parse(time.RFC3339, "2024-05-13T09:07:44+09:00")
	assert.NoError(t, err)

	got, err := Unmarshal[testStruct](data)
	assert.NoError(t, err)
	assert.Equal(t, []testStruct{
		{Column1: "value1", Column2: 1, Column3: 1.1, Column4: true, Column5: tm},
		{Column1: "value2", Column2: 2, Column3: 2.2, Column4: false, Column5: tm},
	}, got)
}
