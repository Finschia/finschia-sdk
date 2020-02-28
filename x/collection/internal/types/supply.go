package types

import (
	"encoding/json"
	"fmt"
)

type Supply interface {
	GetContractID() string
	GetTotalSupply() Coins
	SetTotalSupply(total Coins) Supply
	GetTotalMint() Coins
	GetTotalBurn() Coins

	Inflate(amount Coins) Supply
	Deflate(amount Coins) Supply

	String() string
	ValidateBasic() error
}

var _ Supply = (*BaseSupply)(nil)

type BaseSupply struct {
	ContractID  string `json:"contract_id"`
	TotalSupply Coins  `json:"total_supply"`
	TotalMint   Coins  `json:"total_mint"`
	TotalBurn   Coins  `json:"total_burn"`
}

func (supply BaseSupply) GetContractID() string {
	return supply.ContractID
}

func (supply BaseSupply) SetTotalSupply(total Coins) Supply {
	supply.TotalSupply = total
	supply.TotalMint = total
	supply.TotalBurn = NewCoins()
	return supply
}

func (supply BaseSupply) GetTotalSupply() Coins {
	return supply.TotalSupply
}

func (supply BaseSupply) GetTotalMint() Coins {
	return supply.TotalMint
}

func (supply BaseSupply) GetTotalBurn() Coins {
	return supply.TotalBurn
}

func NewSupply(contractID string, total Coins) Supply {
	return BaseSupply{ContractID: contractID, TotalSupply: total, TotalMint: total, TotalBurn: NewCoins()}
}

func DefaultSupply(contractID string) Supply {
	return NewSupply(contractID, NewCoins())
}

func (supply BaseSupply) Inflate(amount Coins) Supply {
	supply.TotalSupply = supply.TotalSupply.Add(amount...)
	supply.TotalMint = supply.TotalMint.Add(amount...)
	supply.checkInvariant()
	return supply
}

func (supply BaseSupply) Deflate(amount Coins) Supply {
	supply.TotalSupply = supply.TotalSupply.Sub(amount)
	supply.TotalBurn = supply.TotalBurn.Add(amount...)
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

func (supply BaseSupply) ValidateBasic() error {
	if !supply.TotalSupply.IsValid() {
		return fmt.Errorf("invalid total supply: %s", supply.TotalSupply.String())
	}
	if !supply.TotalMint.IsValid() {
		return fmt.Errorf("invalid total mint: %s", supply.TotalMint.String())
	}
	if !supply.TotalBurn.IsValid() {
		return fmt.Errorf("invalid total burn: %s", supply.TotalBurn.String())
	}
	return nil
}

// panic if totalSupply != totalMint - totalBurn
func (supply BaseSupply) checkInvariant() {
	if !supply.TotalSupply.IsEqual(supply.TotalMint.Sub(supply.TotalBurn)) {
		panic(fmt.Sprintf(
			"Collection [%v]'s total supply [%v] does not match with total mint [%v] - total burn [%v]",
			supply.GetContractID(),
			supply.TotalSupply,
			supply.TotalMint,
			supply.TotalBurn,
		))
	}
}
