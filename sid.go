package sid

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

const (
	// base is the base of the ID.
	base = 10
	// bitSize is the size of the ID in bits.
	bitSize = 64
	// zero is the zero value of the ID in int64.
	zero = 0
)

var (
	// reDigit is a regular expression to check if a string is a digit.
	reDigit = regexp.MustCompile(`^[0-9]+$`)
	// Zero is the zero value of the ID type.
	Zero = New(zero)
)

// ID represents a unique identifier.
// It is a wrapper around an int64.
type ID int64

// New creates a new ID from an int64.
func New(v int64) ID {
	return NewFromInt64(v)
}

// NewFromInt64 creates a new ID from an int64.
func NewFromInt64(v int64) ID {
	return ID(v)
}

// NewFromUint64 creates a new ID from a uint64.
func NewFromUint64(v uint64) (ID, error) {
	if v > math.MaxInt64 {
		return Zero, fmt.Errorf("sid: uint64 value %d overflows int64", v)
	}

	return NewFromInt64(int64(v)), nil
}

// NewFromString creates a new ID from a string.
func NewFromString(s string) (ID, error) {
	if !reDigit.MatchString(s) {
		return Zero, fmt.Errorf("sid: can't convert %s to sid", s)
	}

	i, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		return Zero, fmt.Errorf("sid: can't convert %s to int64", s)
	}

	return New(i), nil
}

// IsZero returns true if the ID is the zero value.
func (id ID) IsZero() bool {
	return id == Zero
}

// MarshalText implements the encoding.TextMarshaler interface.
func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (id *ID) UnmarshalText(b []byte) error {
	text := bytes.TrimSpace(b)
	if len(text) == 0 {
		*id = Zero
		return nil
	}

	s := string(text)

	i, err := NewFromString(s)
	if err != nil {
		return err
	}

	*id = i
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(id.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *ID) UnmarshalJSON(b []byte) error {
	val := bytes.TrimSpace(b)

	if len(val) == 0 {
		*id = Zero
		return nil
	}

	if isString(b) {
		s, err := strconv.Unquote(string(b))
		if err != nil {
			return err
		}

		if s == "" {
			*id = Zero
			return nil
		}

		return id.UnmarshalText([]byte(s))
	}

	// if it's not a string, it should be a number
	var n int64

	if err := json.Unmarshal(b, &n); err != nil {
		return fmt.Errorf("sid: can't convert %s to int64", string(b))
	}

	*id = New(n)

	return nil
}

// Scan implements the sql.Scanner interface.
func (id *ID) Scan(value any) error {
	if value == nil || value == (*ID)(nil) {
		*id = Zero

		return nil
	}

	switch val := value.(type) {
	case string:
		return id.UnmarshalText([]byte(val))
	case []byte:
		return id.UnmarshalJSON(val)
	case int64:
		*id = NewFromInt64(val)
		return nil
	case uint64:
		if val > math.MaxInt64 {
			return fmt.Errorf("sid: uint64 value %d overflows int64", val)
		}
		*id = NewFromInt64(int64(val))
		return nil
	case ID:
		*id = val
		return nil
	case *ID:
		*id = *val
		return nil
	default:
		return fmt.Errorf("sid: can't convert %T to ID", value)
	}
}

// Compare compares the numbers represented by id and v and returns:
// -1 if id < v
// 0 if id == v
// +1 if id > v
func (id ID) Compare(v ID) int {
	return id.Cmp(v)
}

// Cmp compares the numbers represented by id and v and returns:
// -1 if id < v
// 0 if id == v
// +1 if id > v
func (id ID) Cmp(v ID) int {
	if id < v {
		return -1
	}

	if id > v {
		return 1
	}

	return 0
}

// Equal returns true if id is equal to v.
func (id ID) Equal(v ID) bool {
	return id.Cmp(v) == 0
}

// GreaterThan returns true if id is greater than v.
func (id ID) GreaterThan(v ID) bool {
	return id.Cmp(v) == 1
}

// GreaterThanOrEqual returns true if id is greater than or equal to v.
func (id ID) GreaterThanOrEqual(v ID) bool {
	cmp := id.Cmp(v)

	return cmp == 1 || cmp == 0
}

// LessThan returns true if id is less than i.
func (id ID) LessThan(v ID) bool {
	return id.Cmp(v) == -1
}

// LessThanOrEqual returns true if id is less than or equal to v.
func (id ID) LessThanOrEqual(v ID) bool {
	cmp := id.Cmp(v)

	return cmp == -1 || cmp == 0
}

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	return id.Int64(), nil
}

// Int64 returns the ID as a int64.
func (id ID) Int64() int64 {
	return int64(id)
}

// String returns the ID as a string.
func (id ID) String() string {
	return strconv.FormatInt(id.Int64(), base)
}

// IDs implements sort.Interface for converting a slice of ID.
type IDs []ID

func (x IDs) Len() int {
	return len(x)
}

func (x IDs) Less(i, j int) bool {
	return x[i].LessThan(x[j])
}

func (x IDs) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// isString checks if the byte slice is a string
func isString(b []byte) bool {
	return len(b) >= 2 && b[0] == '"' && b[len(b)-1] == '"'
}
