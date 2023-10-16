package testutil

import (
	"encoding/json"
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
)

func MustJSONMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return b
}

type FmtStringer interface {
	string | sdk.AccAddress | sdk.Int
}

func W[T FmtStringer](input T) []byte {
	return []byte(fmt.Sprintf("\"%s\"", input))
}
