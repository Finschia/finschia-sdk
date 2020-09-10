package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/client"
)

const (
	QuerierRoute             = ModuleName
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
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
	WithHeight(height int64) client.CLIContext
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
