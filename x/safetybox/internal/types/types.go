package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SafetyBox struct {
	ID                   string         `json:"id"`
	Owner                sdk.AccAddress `json:"owner"`
	Address              sdk.AccAddress `json:"address"`
	TotalAllocation      sdk.Int        `json:"total_allocation"`
	CumulativeAllocation sdk.Int        `json:"cumulative_allocation"`
	TotalIssuance        sdk.Int        `json:"total_issuance"`
	ContractID           string         `json:"contract_id"`
}

func NewSafetyBox(owner sdk.AccAddress, safetyBoxID string, address sdk.AccAddress, contractID string) SafetyBox {
	return SafetyBox{
		safetyBoxID,
		owner,
		address,
		sdk.ZeroInt(),
		sdk.ZeroInt(),
		sdk.ZeroInt(),
		contractID,
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
