package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EncodeHandler func(jsonMsg json.RawMessage) ([]sdk.Msg, error)
type EncodeQuerier func(ctx sdk.Context, jsonQuerier json.RawMessage) ([]byte, error)

const (
	EncodeRouterKey = "tokenencode"
)

type MsgRoute string

const (
	RIssue        = MsgRoute("issue")
	RTransfer     = MsgRoute("transfer")
	RTransferFrom = MsgRoute("transfer_from")
	RMint         = MsgRoute("mint")
	RBurn         = MsgRoute("burn")
	RBurnFrom     = MsgRoute("burn_from")
	RGrantPerm    = MsgRoute("grant_perm")
	RRevokePerm   = MsgRoute("revoke_perm")
	RModify       = MsgRoute("modify")
	RApprove      = MsgRoute("approve")
)

// WasmCustomMsg - wasm custom msg parser
type WasmCustomMsg struct {
	Route string          `json:"route"`
	Data  json.RawMessage `json:"data"`
}

type WasmCustomQuerier struct {
	Route string          `json:"route"`
	Data  json.RawMessage `json:"data"`
}

type QueryTokenWrapper struct {
	TokenParam TokenParam `json:"token_param"`
}

type TokenParam struct {
	ContractID string `json:"contract_id"`
}

type QueryBalanceWrapper struct {
	BalanceParam BalanceParam `json:"balance_param"`
}

type BalanceParam struct {
	ContractID string         `json:"contract_id"`
	Address    sdk.AccAddress `json:"address"`
}

type QueryTotalWrapper struct {
	TotalParam TotalParam `json:"total_param"`
}

type TotalParam struct {
	ContractID string `json:"contract_id"`
}

type QueryPermWrapper struct {
	PermParam PermParam `json:"perm_param"`
}

type PermParam struct {
	ContractID string         `json:"contract_id"`
	Address    sdk.AccAddress `json:"address"`
}

type QueryIsApprovedWrapper struct {
	IsApprovedParam IsApprovedParam `json:"is_approved_param"`
}

type IsApprovedParam struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	Approver   sdk.AccAddress `json:"approver"`
}

type QueryApproversWrapper struct {
	ApproversParam ApproversParam `json:"approvers_param"`
}

type ApproversParam struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
}
