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

func W(input string) []byte {
	return []byte(fmt.Sprintf("\"%s\"", input))
}
