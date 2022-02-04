package keeper

import (
	"context"

	// sdk "github.com/line/lbm-sdk/types"
	// sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

type queryServer struct {
	keeper Keeper
}

// NewQueryServer returns an implementation of the token QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) token.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

var _ token.QueryServer = queryServer{}

// Balance queries the number of tokens of a given class owned by the owner.
func (s queryServer) Balance(context.Context, *token.QueryBalanceRequest) (*token.QueryBalanceResponse, error) {
	panic("Not implemented")
}

// Supply queries the number of tokens from the given class id.
func (s queryServer) Supply(context.Context, *token.QuerySupplyRequest) (*token.QuerySupplyResponse, error) {
	panic("Not implemented")
}

// Token queries an token metadata based on its class id.
func (s queryServer) Token(context.Context, *token.QueryTokenRequest) (*token.QueryTokenResponse, error) {
	panic("Not implemented")
}

// Tokens queries all token metadata.
func (s queryServer) Tokens(context.Context, *token.QueryTokensRequest) (*token.QueryTokensResponse, error) {
	panic("Not implemented")
}
