package keeper

import (
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"

	"context"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/stakingplus/types"
)

type msgServer struct {
	stakingtypes.MsgServer

	ck types.ConsortiumKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper stakingkeeper.Keeper, ck types.ConsortiumKeeper) stakingtypes.MsgServer {
	return &msgServer{
		MsgServer: stakingkeeper.NewMsgServerImpl(keeper),
		ck       : ck,
	}
}

var _ stakingtypes.MsgServer = msgServer{}

func (k msgServer) CreateValidator(goCtx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.ck.GetEnabled(ctx) {
		valAddr := sdk.ValAddress(msg.ValidatorAddress)
		if auth, err := k.ck.GetValidatorAuth(ctx, valAddr); err != nil || !auth.CreationAllowed {
			return nil, sdkerrors.ErrInvalidRequest
		}
	}

	return k.MsgServer.CreateValidator(goCtx, msg)
}
