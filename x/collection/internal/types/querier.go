package types

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QuerierRoute             = ModuleName
	QueryBalances            = "balances"
	QueryBalance             = "balance"
	QueryTokens              = "tokens"
	QueryTokensWithTokenType = "tokensWithTokenType"
	QueryTokenTypes          = "tokentypes"
	QueryPerms               = "perms"
	QueryCollections         = "collections"
	QuerySupply              = "supply"
	QueryMint                = "mint"
	QueryBurn                = "burn"
	QueryNFTCount            = "nftcount"
	QueryNFTMint             = "nftmint"
	QueryNFTBurn             = "nftburn"
	QueryParent              = "parent"
	QueryRoot                = "root"
	QueryChildren            = "children"
	QueryIsApproved          = "approved"
	QueryApprovers           = "approver"
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
	WithHeight(height int64) context.CLIContext
}

type QueryTokenIDParams struct {
	TokenID string `json:"token_id"`
}

func NewQueryTokenIDParams(tokenID string) QueryTokenIDParams {
	return QueryTokenIDParams{TokenID: tokenID}
}

type QueryTokenTypeParams struct {
	TokenType string `json:"token_type"`
}

func NewQueryTokenTypeParams(tokenType string) QueryTokenTypeParams {
	return QueryTokenTypeParams{TokenType: tokenType}
}

type QueryTokenIDAccAddressParams struct {
	TokenID string         `json:"token_id"`
	Addr    sdk.AccAddress `json:"addr"`
}

func NewQueryTokenIDAccAddressParams(tokenID string, addr sdk.AccAddress) QueryTokenIDAccAddressParams {
	return QueryTokenIDAccAddressParams{TokenID: tokenID, Addr: addr}
}

type QueryAccAddressParams struct {
	Addr sdk.AccAddress `json:"addr"`
}

func NewQueryAccAddressParams(addr sdk.AccAddress) QueryAccAddressParams {
	return QueryAccAddressParams{Addr: addr}
}

type QueryIsApprovedParams struct {
	Proxy    sdk.AccAddress `json:"proxy"`
	Approver sdk.AccAddress `json:"approver"`
}

func NewQueryIsApprovedParams(proxy sdk.AccAddress, approver sdk.AccAddress) QueryIsApprovedParams {
	return QueryIsApprovedParams{
		Proxy:    proxy,
		Approver: approver,
	}
}

type QueryProxyParams struct {
	Proxy sdk.AccAddress `json:"proxy"`
}

func NewQueryApproverParams(proxy sdk.AccAddress) QueryProxyParams {
	return QueryProxyParams{Proxy: proxy}
}
