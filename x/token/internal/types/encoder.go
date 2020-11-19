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
	QueryTokenParam QueryTokenParam `json:"query_token_param"`
}

type QueryTokenParam struct {
	ContractID string `json:"contract_id"`
}

type QueryBalanceWrapper struct {
	QueryBalanceParam QueryBalanceParam `json:"query_balance_param"`
}

type QueryBalanceParam struct {
	ContractID string         `json:"contract_id"`
	Address    sdk.AccAddress `json:"address"`
}

type QueryTotalWrapper struct {
	QueryTotalParam QueryTotalParam `json:"query_total_param"`
}

type QueryTotalParam struct {
	ContractID string `json:"contract_id"`
	Target     string `json:"target"`
}

type QueryPermWrapper struct {
	QueryPermParam QueryPermParam `json:"query_perm_param"`
}

type QueryPermParam struct {
	ContractID string         `json:"contract_id"`
	Address    sdk.AccAddress `json:"address"`
}

type QueryIsApprovedWrapper struct {
	QueryIsApprovedParam QueryIsApprovedParam `json:"query_is_approved_param"`
}

type QueryIsApprovedParam struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	Approver   sdk.AccAddress `json:"approver"`
}

type QueryApproversWrapper struct {
	QueryApproversParam QueryApproversParam `json:"query_approvers_param"`
}

type QueryApproversParam struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
}
