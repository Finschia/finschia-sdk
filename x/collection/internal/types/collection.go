package types

import (
	"encoding/json"
)

type Findable interface {
	IDAtIndex(index int) string
	Len() int
}
type Collection interface {
	GetSymbol() string
	GetName() string
	String() string
}
type BaseCollection struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func NewCollection(symbol, name string) Collection {
	return &BaseCollection{
		Symbol: symbol,
		Name:   name,
	}
}

func (c BaseCollection) GetSymbol() string { return c.Symbol }
func (c BaseCollection) GetName() string   { return c.Name }

func (c BaseCollection) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type Collections []Collection

func (collections Collections) String() string {
	b, err := json.Marshal(collections)
	if err != nil {
		panic(err)
	}
	return string(b)
}
