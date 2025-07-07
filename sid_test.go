package sid

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initTestIDs() IDs {
	return IDs{
		15941887676443916,
		17536210255801612,
		17676953901395212,
		18128057957149964,
		18273735630974220,
		18283195028008204,
	}
}

func TestIDs_Len(t *testing.T) {
	ids := initTestIDs()
	assert.Equal(t, 6, ids.Len())
}

func TestIDs_Less(t *testing.T) {
	ids := initTestIDs()
	assert.True(t, ids.Less(0, 1))
	assert.False(t, ids.Less(1, 0))
}

func TestIDs_Swap(t *testing.T) {
	ids := initTestIDs()
	ids.Swap(0, 1)
	assert.Equal(t, IDs{
		17536210255801612,
		15941887676443916,
		17676953901395212,
		18128057957149964,
		18273735630974220,
		18283195028008204,
	}, ids)
}

func TestID_Cmp(t *testing.T) {
	assert.Equal(t, 0, ID(18283195028008204).Cmp(ID(18283195028008204)))
	assert.Equal(t, -1, ID(15941887676443916).Cmp(ID(17676953901395212)))
	assert.Equal(t, 1, ID(18128057957149964).Cmp(ID(17536210255801612)))
}

func TestID_Compare(t *testing.T) {
	assert.Equal(t, 0, ID(18283195028008204).Compare(ID(18283195028008204)))
	assert.Equal(t, -1, ID(15941887676443916).Compare(ID(17676953901395212)))
	assert.Equal(t, 1, ID(18128057957149964).Compare(ID(17536210255801612)))
}

func TestID_Equal(t *testing.T) {
	assert.True(t, ID(18283195028008204).Equal(ID(18283195028008204)))
	assert.False(t, ID(15941887676443916).Equal(ID(17676953901395212)))
	assert.False(t, ID(18128057957149964).Equal(ID(17536210255801612)))
}

func TestID_GreaterThan(t *testing.T) {
	assert.True(t, ID(18283195028008204).GreaterThan(ID(18128057957149964)))
	assert.False(t, ID(18128057957149964).GreaterThan(ID(18283195028008204)))
	assert.False(t, ID(18128057957149964).GreaterThan(ID(18128057957149964)))
}

func TestID_GreaterThanOrEqual(t *testing.T) {
	assert.True(t, ID(18283195028008204).GreaterThanOrEqual(ID(18128057957149964)))
	assert.True(t, ID(18128057957149964).GreaterThanOrEqual(ID(18128057957149964)))
	assert.False(t, ID(18128057957149964).GreaterThanOrEqual(ID(18283195028008204)))
}

func TestID_LessThan(t *testing.T) {
	assert.True(t, ID(18128057957149964).LessThan(ID(18283195028008204)))
	assert.False(t, ID(18283195028008204).LessThan(ID(18128057957149964)))
	assert.False(t, ID(18128057957149964).LessThan(ID(18128057957149964)))
}

func TestID_LessThanOrEqual(t *testing.T) {
	assert.True(t, ID(18128057957149964).LessThanOrEqual(ID(18283195028008204)))
	assert.True(t, ID(18128057957149964).LessThanOrEqual(ID(18128057957149964)))
	assert.False(t, ID(18283195028008204).LessThanOrEqual(ID(18128057957149964)))
}

func TestID_Value(t *testing.T) {
	value, err := ID(18283195028008204).Value()
	assert.NoError(t, err)
	assert.Equal(t, int64(18283195028008204), value)
}

func TestID_Int64(t *testing.T) {
	assert.Equal(t, int64(18283195028008204), ID(18283195028008204).Int64())
}

func TestID_String(t *testing.T) {
	assert.Equal(t, "18283195028008204", ID(18283195028008204).String())
}

func TestNewFromInt64(t *testing.T) {
	id := NewFromInt64(int64(18283195028008204))
	assert.Equal(t, ID(18283195028008204), id)
}

func TestNewFromString(t *testing.T) {
	id, err := NewFromString("18283195028008204")
	assert.NoError(t, err)
	assert.Equal(t, ID(18283195028008204), id)
}

func TestSortInterface(t *testing.T) {
	// Test IsSorted for sort.Interface
	checkSorted := make(IDs, 0, 256)
	for i := 0; i < 256; i++ {
		checkSorted = append(checkSorted, ID(i))
	}
	assert.True(t, sort.IsSorted(checkSorted))

	// Test sort for sort.Interface
	checkSorted = make(IDs, 0, 128)
	for i := 0; i < 128; i++ {
		checkSorted = append(checkSorted, ID(i))
	}
	sort.Sort(checkSorted)
	assert.True(t, sort.IsSorted(checkSorted))
}

func TestID_IsZero(t *testing.T) {
	assert.True(t, ID(0).IsZero())
	assert.False(t, ID(1).IsZero())
}

func TestNew(t *testing.T) {
	id := New(18283195028008204)
	assert.Equal(t, ID(18283195028008204), id)
	assert.False(t, id.IsZero())
}

func TestID_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected ID
		wantErr  bool
	}{
		{"Scan uint64", uint64(18283195028008204), ID(18283195028008204), false},
		{"Scan uint64 overflow", uint64(math.MaxInt64) + 1, ID(0), true},
		{"Scan string", "18283195028008204", ID(18283195028008204), false},
		{"Scan byte slice", []byte("18283195028008204"), ID(18283195028008204), false},
		{"Scan nil", nil, ID(0), false},
		{"Scan invalid string", "not a number", ID(0), true},
		{"Scan invalid byte slice", []byte("not a number"), ID(0), true},
		{"Scan invalid type", struct{}{}, ID(0), true},
		{"Scan int64", int64(18283195028008204), ID(18283195028008204), false},
		{"Scan ID", ID(18283195028008204), ID(18283195028008204), false},
		{"Scan nil pointer", (*ID)(nil), ID(0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := id.Scan(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, id, fmt.Sprintf("input=%v", tt.input))
			}
		})
	}
}

func TestID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		id       ID
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Positive ID",
			id:       ID(18283195028008204),
			expected: []byte(`"18283195028008204"`),
			wantErr:  false,
		},
		{
			name:     "Zero ID",
			id:       ID(0),
			expected: []byte(`"0"`),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.id.MarshalJSON()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, data)
		})
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ID
		wantErr  bool
	}{
		{
			name:     "Valid string input",
			input:    `"18283195028008204"`,
			expected: ID(18283195028008204),
			wantErr:  false,
		},
		{
			name:     "Valid number input",
			input:    `18283195028008204`,
			expected: ID(18283195028008204),
			wantErr:  false,
		},
		{
			name:     "Invalid input",
			input:    `"not a number"`,
			expected: ID(0),
			wantErr:  true,
		},
		{
			name:     "Empty input",
			input:    ``,
			expected: ID(0),
			wantErr:  false,
		},
		{
			name:     "Null input",
			input:    `null`,
			expected: ID(0),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID

			err := id.UnmarshalJSON([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, id)
		})
	}
}

func TestID_MarshalText(t *testing.T) {
	data, err := ID(18283195028008204).MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`18283195028008204`), data)

	data, err = ID(0).MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`0`), data)
}

func TestID_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ID
		wantErr  bool
	}{
		{
			name:     "Valid number input",
			input:    "18283195028008204",
			expected: ID(18283195028008204),
			wantErr:  false,
		},
		{
			name:     "Invalid input",
			input:    "not a number",
			expected: ID(0),
			wantErr:  true,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: ID(0),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID

			err := id.UnmarshalText([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, id)
		})
	}
}
