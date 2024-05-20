package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khgame/ranger_iam/pkg/authcli"

	"github.com/bytedance/sonic"
	"github.com/khicago/irr"
)

// UInt64 is a custom type to handle JSON serialization of 64-bit integers.
type UInt64 uint64

func (val UInt64) Raw() uint64 {
	return uint64(val)
}

// MarshalJSON serializes the UInt64 as a string to avoid precision loss in JavaScript.
func (val UInt64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(val), 10) + `"`), nil
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

func GetUIDFromGinCtx(c *gin.Context) (UInt64, bool) {
	id, exist := authcli.GetUIDFromGinCtx(c)
	return UInt64(id), exist
}
