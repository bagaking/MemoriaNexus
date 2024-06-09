package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/khicago/irr"
)

// UInt64 is a custom type to handle JSON serialization of 64-bit integers.
type UInt64 uint64

func (val UInt64) Raw() uint64 {
	return uint64(val)
}

var (
	Zero                  = UInt64(0)
	_    json.Marshaler   = &Zero
	_    json.Unmarshaler = &Zero
)

// Scan implements the sql.Scanner interface.
func (val *UInt64) Scan(value any) error {
	if value == nil {
		*val = Zero
		return nil
	}

	switch v := value.(type) {
	case int64:
		*val = UInt64(v)
	case uint64:
		*val = UInt64(v)
	case []byte:
		var parsed uint64
		if err := sonic.Unmarshal(v, &parsed); err != nil {
			return err
		}
		*val = UInt64(parsed)
	case string:
		var parsed uint64
		if err := sonic.Unmarshal([]byte(v), &parsed); err != nil {
			return err
		}
		*val = UInt64(parsed)
	default:
		return errors.New("unsupported Scan source")
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (val UInt64) Value() (driver.Value, error) {
	return int64(val), nil
}

// MarshalJSON serializes the UInt64 as a string to avoid precision loss in JavaScript.
func (val *UInt64) MarshalJSON() ([]byte, error) {
	bytes, err := sonic.Marshal(uint64(*val))
	if err != nil {
		return nil, err
	}
	data := make([]byte, 0, len(bytes)+2)
	data = append(data, '"')
	data = append(data, bytes...)
	data = append(data, '"')
	return data, nil
}

// UnmarshalJSON supports parsing the UInt64 from a number or a string in JSON.
func (val *UInt64) UnmarshalJSON(b []byte) error {
	// Unmarshalling the bytes as a string takes care of both quoted and non-quoted numbers.
	var idVal uint64
	var strValue string

	// Attempt to unmarshal as a string first to accommodate both quoted and non-quoted JSON numbers.
	if err := sonic.Unmarshal(b, &strValue); err == nil {
		idVal, err = strconv.ParseUint(strValue, 10, 64)
		if err != nil {
			return irr.Wrap(err, "parse string val failed")
		}
	} else {
		// If unmarshaling as a string fails, it means the bytes represent an actual JSON number.
		idVal, err = strconv.ParseUint(string(b), 10, 64)
		if err != nil {
			return irr.Wrap(err, "parse uint val failed")
		}
	}

	*val = UInt64(idVal)
	return nil
}

// ParseIDFromString Support parsing from a string value directly.
func ParseIDFromString(value string) (UInt64, error) {
	parsedValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return UInt64(parsedValue), nil
}
