package testutil

import (
	"encoding/json"
	"fmt"
)

func MustJSONMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return b
}

// W wraps input with double quotes if it is a string or fmt.Stringer.
func W(input any) []byte {
	switch input.(type) {
	case string, fmt.Stringer:
		return []byte(fmt.Sprintf("\"%s\"", input))
	default:
		panic("unsupported type")
	}
}
