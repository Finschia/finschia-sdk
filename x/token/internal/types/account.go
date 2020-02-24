package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenID string

type Account interface {
	GetAddress() sdk.AccAddress
	GetSymbol() string
	GetBalance() sdk.Int
	SetBalance(amount sdk.Int) Account
	String() string
}

type BaseAccount struct {
	Symbol  string         `json:"symbol"`
	Address sdk.AccAddress `json:"address"`
	Amount  sdk.Int        `json:"amount"`
}

func NewBaseAccountWithAddress(symbol string, addr sdk.AccAddress) *BaseAccount {
	return &BaseAccount{
		Symbol:  symbol,
		Address: addr,
		Amount:  sdk.ZeroInt(),
	}
}

func (acc BaseAccount) GetSymbol() string {
	return acc.Symbol
}

func (acc BaseAccount) String() string {
	b, err := json.Marshal(acc)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

func (acc BaseAccount) GetBalance() sdk.Int {
	return acc.Amount
}

func (acc BaseAccount) SetBalance(amount sdk.Int) Account {
	acc.Amount = amount
	return acc
}
