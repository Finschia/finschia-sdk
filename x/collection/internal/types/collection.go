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
	SetName(name string) Collection
	GetBaseImgURI() string
	SetBaseImgURI(baseImgURI string) Collection
	String() string
}
type BaseCollection struct {
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	BaseImgURI string `json:"base_img_uri"`
}

func NewCollection(symbol, name, baseImgURI string) Collection {
	return &BaseCollection{
		Symbol:     symbol,
		Name:       name,
		BaseImgURI: baseImgURI,
	}
}

func (c BaseCollection) GetSymbol() string { return c.Symbol }
func (c BaseCollection) GetName() string   { return c.Name }
func (c *BaseCollection) SetName(name string) Collection {
	c.Name = name
	return c
}

func (c BaseCollection) GetBaseImgURI() string { return c.BaseImgURI }
func (c *BaseCollection) SetBaseImgURI(baseImgURI string) Collection {
	c.BaseImgURI = baseImgURI
	return c
}

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
