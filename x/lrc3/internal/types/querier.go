package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type QueryLRC3Params struct {
	Denom string
}

// NewQueryLRC3Params creates a new instance of QuerySupplyParams
func NewQueryLRC3Params(denom string) QueryLRC3Params {
	return QueryLRC3Params{Denom: denom}
}

// Bytes exports the Denom as bytes
func (q QueryLRC3Params) Bytes() []byte {
	return []byte(q.Denom)
}

type QueryApproveParams struct {
	Denom   string
	TokenID string
}

func NewQueryApproveParams(denom string, tokenId string) QueryApproveParams {
	return QueryApproveParams{
		Denom:   denom,
		TokenID: tokenId,
	}
}

// QueryBalanceParams params for query 'custom/nfts/balance'
type QueryBalanceParams struct {
	Owner sdk.AccAddress
	Denom string
}

// NewQueryBalanceParams creates a new instance of QuerySupplyParams
func NewQueryBalanceParams(owner sdk.AccAddress, denom string) QueryBalanceParams {
	return QueryBalanceParams{
		Owner: owner,
		Denom: denom,
	}
}

type QueryOwnerOfParams struct {
	Denom   string
	TokenID string
}

func NewQueryOwnerOfParams(denom string, tokenId string) QueryOwnerOfParams {
	return QueryOwnerOfParams{
		Denom:   denom,
		TokenID: tokenId,
	}
}
