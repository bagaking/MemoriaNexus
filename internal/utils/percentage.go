package utils

import (
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/khicago/irr"
)

// Percentage is a custom type to handle JSON serialization of percentages as uint8.
type Percentage uint8

// Raw returns the raw uint8 value of the percentage.
func (p *Percentage) Raw() uint8 {
	return uint8(*p)
}

// AsFloat returns the percentage as a float64 value (e.g., 50 becomes 50.0).
func (p *Percentage) AsFloat() float64 {
	return float64(*p)
}

// MarshalJSON serializes the Percentage as a string to avoid precision loss in JavaScript.
func (p *Percentage) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.Itoa(int(*p)) + `"`), nil
}

// UnmarshalJSON supports parsing the Percentage from a number or a string in JSON.
func (p *Percentage) UnmarshalJSON(b []byte) error {
	var percentVal uint8
	var strValue string

	// Attempt to unmarshal as a string first to accommodate both quoted and non-quoted JSON numbers.
	if err := sonic.Unmarshal(b, &strValue); err == nil {
		tempVal, err := strconv.ParseUint(strValue, 10, 8)
		if err != nil {
			return irr.Wrap(err, "parse string val failed")
		}
		percentVal = uint8(tempVal)
	} else {
		// If unmarshaling as a string fails, it means the bytes represent an actual JSON number.
		tempVal, err := strconv.ParseUint(string(b), 10, 8)
		if err != nil {
			return irr.Wrap(err, "parse uint val failed")
		}
		percentVal = uint8(tempVal)
	}

	*p = Percentage(percentVal)
	return nil
}

// ParsePercentageFromString supports parsing from a string value directly.
func ParsePercentageFromString(value string) (Percentage, error) {
	parsedValue, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return Percentage(uint8(parsedValue)), nil
}
