package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	QueryCollectionParam QueryCollectionParam `json:"query_collection_param"`
}

type QueryBalanceWrapper struct {
	QueryBalanceParam QueryBalanceParam `json:"query_collection_balance_param"`
}

type QueryTokenTypesWrapper struct {
	QueryTokenTypesParam QueryTokenTypesParam `json:"query_tokentypes_param"`
}

type QueryTokensWrapper struct {
	QueryTokensParam QueryTokensParam `json:"query_tokens_with_collection_param"`
}

type QueryTokenTypeWrapper struct {
	QueryTokenTypeParam QueryTokenTypeParam `json:"query_token_type_param"`
}

type QueryNFTCountWrapper struct {
	QueryNFTCountParam QueryNFTCountParam `json:"query_nft_count_param"`
}

type QueryTotalWrapper struct {
	QueryTotalParam QueryTotalParam `json:"query_total_with_collection_param"`
}

type QueryPermsWrapper struct {
	QueryPermParam QueryPermParam `json:"query_perm_param"`
}

type QueryApprovedWrapper struct {
	QueryApprovedParam QueryApprovedParam `json:"query_approved_param"`
}

type QueryApproversWrapper struct {
	QueryProxyParam QueryProxyParam `json:"query_proxy_param"`
}

type QueryCollectionParam struct {
	ContractID string `json:"contract_id"`
}

type QueryBalanceParam struct {
	ContractID string         `json:"contract_id"`
	TokenID    string         `json:"token_id"`
	Addr       sdk.AccAddress `json:"addr"`
}

type QueryTokenTypesParam struct {
	ContractID string `json:"contract_id"`
	TokenID    string `json:"token_id"`
}

type QueryTokensParam struct {
	ContractID string `json:"contract_id"`
	TokenID    string `json:"token_id"`
}

type QueryTokenTypeParam struct {
	ContractID string `json:"contract_id"`
	TokenType  string `json:"token_type"`
}

type QueryNFTCountParam struct {
	ContractID string `json:"contract_id"`
	TokenID    string `json:"token_id"`
	Target     string `json:"target"`
}

type QueryTotalParam struct {
	ContractID string `json:"contract_id"`
	TokenID    string `json:"token_id"`
	Target     string `json:"target"`
}

type QueryPermParam struct {
	ContractID string         `json:"contract_id"`
	Address    sdk.AccAddress `json:"address"`
}

type QueryProxyParam struct {
	ContractID string         `json:"contract_id"`
	Proxy      sdk.AccAddress `json:"proxy"`
}

type QueryApprovedParam struct {
	ContractID string         `json:"contract_id"`
	Proxy      sdk.AccAddress `json:"proxy"`
	Approver   sdk.AccAddress `json:"approver"`
}
