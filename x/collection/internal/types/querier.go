package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/client"
)

const (
	QuerierRoute     = ModuleName
	QueryBalance     = "balance"
	QueryTokens      = "tokens"
	QueryTokenTypes  = "tokentypes"
	QueryPerms       = "perms"
	QueryCollections = "collections"
	QuerySupply      = "supply"
	QueryMint        = "mint"
	QueryBurn        = "burn"
	QueryNFTCount    = "nftcount"
	QueryNFTMint     = "nftmint"
	QueryNFTBurn     = "nftburn"
	QueryParent      = "parent"
	QueryRoot        = "root"
	QueryChildren    = "children"
	QueryIsApproved  = "approved"
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
	WithHeight(height int64) client.CLIContext
}

type QueryContractIDParams struct {
	ContractID string `json:"contract_id"`
}

func NewQueryContractIDParams(contractID string) QueryContractIDParams {
	return QueryContractIDParams{ContractID: contractID}
}

type QueryContractIDTokenIDParams struct {
	ContractID string `json:"contract_id"`
	TokenID    string `json:"token_id"`
}

func NewQueryContractIDTokenIDParams(contractID, tokenID string) QueryContractIDTokenIDParams {
	return QueryContractIDTokenIDParams{ContractID: contractID, TokenID: tokenID}
}

type QueryContractIDTokenIDAccAddressParams struct {
	ContractID string         `json:"contract_id"`
	TokenID    string         `json:"token_id"`
	Addr       sdk.AccAddress `json:"addr"`
}

func NewQueryContractIDTokenIDAccAddressParams(contractID, tokenID string, addr sdk.AccAddress) QueryContractIDTokenIDAccAddressParams {
	return QueryContractIDTokenIDAccAddressParams{ContractID: contractID, TokenID: tokenID, Addr: addr}
}

type QueryContractIDAccAddressParams struct {
	ContractID string         `json:"contract_id"`
	Addr       sdk.AccAddress `json:"addr"`
}

func NewQueryContractIDAccAddressParams(contractID string, addr sdk.AccAddress) QueryContractIDAccAddressParams {
	return QueryContractIDAccAddressParams{ContractID: contractID, Addr: addr}
}

type QueryIsApprovedParams struct {
	ContractID string         `json:"contract_id"`
	Proxy      sdk.AccAddress `json:"proxy"`
	Approver   sdk.AccAddress `json:"approver"`
}

func NewQueryIsApprovedParams(contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) QueryIsApprovedParams {
	return QueryIsApprovedParams{
		ContractID: contractID,
		Proxy:      proxy,
		Approver:   approver,
	}
}
