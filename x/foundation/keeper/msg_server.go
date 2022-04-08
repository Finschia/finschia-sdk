package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServer returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) foundation.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ foundation.MsgServer = msgServer{}

// FundTreasury defines a method to fund the treasury.
func (s msgServer) FundTreasury(c context.Context, req *foundation.MsgFundTreasury) (*foundation.MsgFundTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.fundTreasury(ctx, sdk.AccAddress(req.From), req.Amount); err != nil {
		return nil, err
	}

	return &foundation.MsgFundTreasuryResponse{}, nil
}

// WithdrawFromTreasury defines a method to withdraw coins from the treasury.
func (s msgServer) WithdrawFromTreasury(c context.Context, req *foundation.MsgWithdrawFromTreasury) (*foundation.MsgWithdrawFromTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.withdrawFromTreasury(ctx, sdk.AccAddress(req.Operator), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	return &foundation.MsgWithdrawFromTreasuryResponse{}, nil
}
