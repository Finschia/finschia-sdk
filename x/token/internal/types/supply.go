package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Supply interface {
	GetTotal() sdk.Int
	SetTotal(total sdk.Int) Supply
	GetSymbol() string

	Inflate(amount sdk.Int) Supply
	Deflate(amount sdk.Int) Supply

	String() string
}

type BaseSupply struct {
	Symbol string  `json:"symbol"`
	Total  sdk.Int `json:"total"`
}

func NewSupply(symbol string, total sdk.Int) Supply {
	return BaseSupply{symbol, total}
}

func DefaultSupply(symbol string) Supply {
	return NewSupply(symbol, sdk.ZeroInt())
}

func (supply BaseSupply) SetTotal(total sdk.Int) Supply {
	supply.Total = total
	return supply
}

func (supply BaseSupply) GetSymbol() string {
	return supply.Symbol
}

func (supply BaseSupply) GetTotal() sdk.Int {
	return supply.Total
}

func (supply BaseSupply) Inflate(amount sdk.Int) Supply {
	supply.Total = supply.Total.Add(amount)
	return supply
}

func (supply BaseSupply) Deflate(amount sdk.Int) Supply {
	supply.Total = supply.Total.Sub(amount)
	return supply
}

func (supply BaseSupply) String() string {
	b, err := json.Marshal(supply)
	if err != nil {
		panic(err)
	}
	return string(b)
}
