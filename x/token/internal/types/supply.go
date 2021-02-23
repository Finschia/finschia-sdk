package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
)

type Supply interface {
	GetTotalSupply() sdk.Int
	SetTotalSupply(total sdk.Int) Supply
	GetTotalBurn() sdk.Int
	GetTotalMint() sdk.Int
	GetContractID() string

	Inflate(amount sdk.Int) Supply
	Deflate(amount sdk.Int) Supply

	String() string
}

type BaseSupply struct {
	ContractID  string  `json:"contract_id"`
	TotalSupply sdk.Int `json:"total_supply"`
	TotalMint   sdk.Int `json:"total_mint"`
	TotalBurn   sdk.Int `json:"total_burn"`
}

func NewSupply(contractID string, total sdk.Int) Supply {
	return BaseSupply{contractID, total, total, sdk.ZeroInt()}
}

func DefaultSupply(contractID string) Supply {
	return NewSupply(contractID, sdk.ZeroInt())
}

func (supply BaseSupply) SetTotalSupply(total sdk.Int) Supply {
	supply.TotalSupply = total
	supply.TotalMint = total
	supply.TotalBurn = sdk.ZeroInt()
	return supply
}

func (supply BaseSupply) GetContractID() string {
	return supply.ContractID
}

func (supply BaseSupply) GetTotalSupply() sdk.Int {
	return supply.TotalSupply
}

func (supply BaseSupply) GetTotalMint() sdk.Int {
	return supply.TotalMint
}

func (supply BaseSupply) GetTotalBurn() sdk.Int {
	return supply.TotalBurn
}

func (supply BaseSupply) Inflate(amount sdk.Int) Supply {
	supply.TotalSupply = supply.TotalSupply.Add(amount)
	supply.TotalMint = supply.TotalMint.Add(amount)
	supply.checkInvariant()
	return supply
}

func (supply BaseSupply) Deflate(amount sdk.Int) Supply {
	supply.TotalSupply = supply.TotalSupply.Sub(amount)
	supply.TotalBurn = supply.TotalBurn.Add(amount)
	supply.checkInvariant()
	return supply
}

func (supply BaseSupply) String() string {
	b, err := json.Marshal(supply)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// panic if totalSupply != totalMint - totalBurn
func (supply BaseSupply) checkInvariant() {
	if !supply.TotalSupply.Equal(supply.TotalMint.Sub(supply.TotalBurn)) {
		panic(fmt.Sprintf(
			"Token [%v]'s total supply [%v] does not match with total mint [%v] - total burn [%v]",
			supply.GetContractID(),
			supply.TotalSupply,
			supply.TotalMint,
			supply.TotalBurn,
		))
	}
}
