package types

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
)

type Collection struct {
	Symbol string `json:"symbol"`
}

func NewCollection(symbol string) Collection {
	return Collection{
		Symbol: symbol,
	}
}

func (c Collection) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, c))
}

type Collections []Collection

func (collections *Collections) String() string {
	var collectionStrings = make([]string, len(*collections))
	for idx, t := range *collections {
		collectionStrings[idx] = t.String()
	}
	return strings.Join(collectionStrings, ",")
}
