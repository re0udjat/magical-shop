package data

import (
	"fmt"
	"strconv"
)

type Currency int64

// MarshalJSON() method satisfies the json.Marshaler interface
// This allows us to customize how the Currency type is encoded to JSON (e.g. "100 coins" instead of 100)
func (c Currency) MarshalJSON() ([]byte, error) {
	// Format the currency value as a string with "coins" suffix
	jsonValue := fmt.Sprintf("%d coins", c)

	// Quote the string to make it a valid JSON string
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
