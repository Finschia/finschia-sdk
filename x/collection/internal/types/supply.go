package types

import (
	"encoding/json"
	"fmt"
)

type Supply interface {
	GetSymbol() string
	GetTotal() Coins
	SetTotal(total Coins) Supply

	Inflate(amount Coins) Supply
	Deflate(amount Coins) Supply

	String() string
	ValidateBasic() error
}

var _ Supply = (*BaseSupply)(nil)

type BaseSupply struct {
	Symbol string `json:"symbol"`
	Total  Coins  `json:"total"`
}

func (supply BaseSupply) GetSymbol() string {
	return supply.Symbol
}

func (supply BaseSupply) SetTotal(total Coins) Supply {
	supply.Total = total
	return supply
}

func (supply BaseSupply) GetTotal() Coins {
	return supply.Total
}

func NewSupply(symbol string, total Coins) Supply {
	return BaseSupply{Symbol: symbol, Total: total}
}

func DefaultSupply(symbol string) Supply {
	return NewSupply(symbol, NewCoins())
}

func (supply BaseSupply) Inflate(amount Coins) Supply {
	supply.Total = supply.Total.Add(amount...)
	return supply
}

func (supply BaseSupply) Deflate(amount Coins) Supply {
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

func (supply BaseSupply) ValidateBasic() error {
	if !supply.Total.IsValid() {
		return fmt.Errorf("invalid total supply: %s", supply.Total.String())
	}
	return nil
}
