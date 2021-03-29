package types

import (
	"encoding/json"
)

type Findable interface {
	IDAtIndex(index int) string
	Len() int
}
type Collection interface {
	GetContractID() string
	GetName() string
	SetName(name string)
	GetBaseImgURI() string
	SetBaseImgURI(baseImgURI string)
	GetMeta() string
	SetMeta(meta string)
	String() string
}
type BaseCollection struct {
	ContractID string `json:"contract_id"`
	Name       string `json:"name"`
	Meta       string `json:"meta"`
	BaseImgURI string `json:"base_img_uri"`
}

func NewCollection(contractID, name, meta, baseImgURI string) Collection {
	return &BaseCollection{
		ContractID: contractID,
		Name:       name,
		Meta:       meta,
		BaseImgURI: baseImgURI,
	}
}

func (c BaseCollection) GetContractID() string { return c.ContractID }
func (c BaseCollection) GetName() string       { return c.Name }
func (c *BaseCollection) SetName(name string) {
	c.Name = name
}

func (c BaseCollection) GetMeta() string { return c.Meta }
func (c *BaseCollection) SetMeta(meta string) {
	c.Meta = meta
}

func (c BaseCollection) GetBaseImgURI() string { return c.BaseImgURI }
func (c *BaseCollection) SetBaseImgURI(baseImgURI string) {
	c.BaseImgURI = baseImgURI
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
