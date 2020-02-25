package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Account interface {
	GetAddress() sdk.AccAddress
	GetSymbol() string
	GetCoins() Coins
	SetCoins(Coins) Account
	String() string
}

type BaseAccount struct {
	Address sdk.AccAddress `json:"address"`
	Symbol  string         `json:"symbol"`
	Coins   Coins          `json:"tokens"`
}

func NewBaseAccountWithAddress(symbol string, addr sdk.AccAddress) *BaseAccount {
	return &BaseAccount{
		Symbol:  symbol,
		Address: addr,
		Coins:   NewCoins(),
	}
}

func (acc BaseAccount) String() string {
	b, err := json.Marshal(acc)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (acc BaseAccount) GetSymbol() string {
	return acc.Symbol
}

func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

func (acc BaseAccount) GetCoins() Coins {
	return acc.Coins
}

func (acc BaseAccount) SetCoins(coins Coins) Account {
	acc.Coins = coins
	return acc
}
