package types

import (
	"encoding/json"

	sdk "github.com/line/lbm-sdk/types"
)

type TokenID string

type Account interface {
	GetAddress() sdk.AccAddress
	GetContractID() string
	GetBalance() sdk.Int
	SetBalance(amount sdk.Int) Account
	String() string
}

type BaseAccount struct {
	ContractID string         `json:"contract_id"`
	Address    sdk.AccAddress `json:"address"`
	Amount     sdk.Int        `json:"amount"`
}

func NewBaseAccountWithAddress(contractID string, addr sdk.AccAddress) *BaseAccount {
	return &BaseAccount{
		ContractID: contractID,
		Address:    addr,
		Amount:     sdk.ZeroInt(),
	}
}

func (acc BaseAccount) GetContractID() string {
	return acc.ContractID
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
