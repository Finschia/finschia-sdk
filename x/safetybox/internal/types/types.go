package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SafetyBox struct {
	ID                   string         `json:"id"`
	Owner                sdk.AccAddress `json:"owner"`
	Address              sdk.AccAddress `json:"address"`
	Denoms               []string       `json:"denoms"`
	TotalAllocation      sdk.Coins      `json:"total_allocation"`
	CumulativeAllocation sdk.Coins      `json:"cumulative_allocation"`
	TotalIssuance        sdk.Coins      `json:"total_issuance"`
}

func NewSafetyBox(owner sdk.AccAddress, safetyBoxId string, address sdk.AccAddress, denoms []string) SafetyBox {
	return SafetyBox{
		safetyBoxId,
		owner,
		address,
		denoms,
		sdk.Coins{},
		sdk.Coins{},
		sdk.Coins{},
	}
}

func (sb SafetyBox) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, sb))
}

type MsgSafetyBoxRoleResponse struct {
	HasRole bool `json:"has_role"`
}

func (sbpr MsgSafetyBoxRoleResponse) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, sbpr))
}

func (sbpr MsgSafetyBoxRoleResponse) Bytes() []byte {
	return codec.MustMarshalJSONIndent(ModuleCdc, sbpr)
}
