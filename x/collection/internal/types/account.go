package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Account interface {
	GetAddress() sdk.AccAddress
	GetContractID() string
	GetCoins() Coins
	SetCoins(Coins) Account
	String() string
}

type BaseAccount struct {
	Address    sdk.AccAddress `json:"address"`
	ContractID string         `json:"contract_id"`
	Coins      Coins          `json:"tokens"`
}

func NewBaseAccountWithAddress(contractID string, addr sdk.AccAddress) *BaseAccount {
	return &BaseAccount{
		ContractID: contractID,
		Address:    addr,
		Coins:      NewCoins(),
	}
}

func (acc BaseAccount) String() string {
	b, err := json.Marshal(acc)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (acc BaseAccount) GetContractID() string {
	return acc.ContractID
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

type Accounts []Account
