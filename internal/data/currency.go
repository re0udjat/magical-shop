package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Define an error that UnmarshalJSON() method can return if we're unable to parse or convert
// the JSON string successfully
var ErrInvalidCurrencyFormat = errors.New("invalid currency format")

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

func (c *Currency) UnmarshalJSON(jsonValue []byte) error {
	// The incoming JSON value must be a string in the format "<price> coins"
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidCurrencyFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	// Sanity check the parts of the string to make sure it was in the expected format
	if len(parts) != 2 || parts[1] != "coins" {
		return ErrInvalidCurrencyFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return ErrInvalidCurrencyFormat
	}

	// Convert the int64 to a Currency type and assign this to the receiver
	*c = Currency(i)

	return nil
}
