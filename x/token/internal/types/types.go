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

func (collections Collections) String() string {
	var collectionStrings []string
	for _, t := range collections {
		collectionStrings = append(collectionStrings, t.String())
	}
	return strings.Join(collectionStrings, ",")
}
