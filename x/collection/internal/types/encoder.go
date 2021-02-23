package types

import (
	"encoding/json"

	sdk "github.com/line/lbm-sdk/types"
)

const (
	EncodeRouterKey = "collectionencode"
)

type MsgRoute string

const (
	RCreateCollection = MsgRoute("create")
	RIssueNFT         = MsgRoute("issue_nft")
	RIssueFT          = MsgRoute("issue_ft")
	RMintFT           = MsgRoute("mint_ft")
	RMintNFT          = MsgRoute("mint_nft")
	RBurnNFT          = MsgRoute("burn_nft")
	RBurnFT           = MsgRoute("burn_ft")
	RBurnFTFrom       = MsgRoute("burn_ft_from")
	RBurnNFTFrom      = MsgRoute("burn_nft_from")
	RTransferFT       = MsgRoute("transfer_ft")
	RTransferNFT      = MsgRoute("transfer_nft")
	RTransferFTFrom   = MsgRoute("transfer_ft_from")
	RTransferNFTFrom  = MsgRoute("transfer_nft_from")
	RModify           = MsgRoute("modify")
	RApprove          = MsgRoute("approve")
	RDisapprove       = MsgRoute("disapprove")
	RGrantPerm        = MsgRoute("grant_perm")
	RRevokePerm       = MsgRoute("revoke_perm")
	RAttach           = MsgRoute("attach")
	RDetach           = MsgRoute("detach")
	RAttachFrom       = MsgRoute("attach_from")
	RDetachFrom       = MsgRoute("detach_from")
)

type WasmCustomMsg struct {
	Route string          `json:"route"`
	Data  json.RawMessage `json:"data"`
}

type WasmCustomQuerier struct {
	Route string          `json:"route"`
	Data  json.RawMessage `json:"data"`
}

type QueryCollectionWrapper struct {
	CollectionParam CollectionParam `json:"collection_param"`
}

type QueryBalanceWrapper struct {
	BalanceParam BalanceParam `json:"balance_param"`
}

type QueryTokenTypesWrapper struct {
	TokenTypesParam TokenTypesParam `json:"tokentypes_param"`
}

type QueryTokensWrapper struct {
	TokensParam TokensParam `json:"tokens_param"`
}

type QueryTokenTypeWrapper struct {
	TokenTypeParam TokenTypeParam `json:"token_type_param"`
}

type QueryNFTCountWrapper struct {
	TokensParam TokensParam `json:"tokens_param"`
}

type QueryTotalWrapper struct {
	TotalParam TotalParam `json:"total_param"`
}

type QueryPermsWrapper struct {
	PermParam PermParam `json:"perm_param"`
}

type QueryApprovedWrapper struct {
	IsApprovedParam IsApprovedParam `json:"is_approved_param"`
}

type QueryApproversWrapper struct {
	ApproversParam ApproversParam `json:"approvers_param"`
}

type CollectionParam struct {
	ContractID string `json:"contract_id"`
}

type BalanceParam struct {
	ContractID string         `json:"contract_id"`
	TokenID    string         `json:"token_id"`
	Addr       sdk.AccAddress `json:"addr"`
}

type TokenTypesParam struct {
	ContractID string `json:"contract_id"`
	TokenID    string `json:"token_id"`
}

type TokensParam struct {
	ContractID string `json:"contract_id"`
	TokenID    string `json:"token_id"`
}

type TokenTypeParam struct {
	ContractID string `json:"contract_id"`
	TokenType  string `json:"token_type"`
}

type TotalParam struct {
	ContractID string `json:"contract_id"`
	TokenID    string `json:"token_id"`
}

type PermParam struct {
	ContractID string         `json:"contract_id"`
	Address    sdk.AccAddress `json:"address"`
}

type ApproversParam struct {
	ContractID string         `json:"contract_id"`
	Proxy      sdk.AccAddress `json:"proxy"`
}

type IsApprovedParam struct {
	ContractID string         `json:"contract_id"`
	Proxy      sdk.AccAddress `json:"proxy"`
	Approver   sdk.AccAddress `json:"approver"`
}
