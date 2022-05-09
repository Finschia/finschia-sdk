package keeper

import (
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"

	"context"

	sdk "github.com/line/lbm-sdk/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	"github.com/line/lbm-sdk/x/stakingplus/types"
)

type msgServer struct {
	stakingtypes.MsgServer

	fk types.FoundationKeeper
}

// NewMsgServerImpl returns an implementation of the staking MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper stakingkeeper.Keeper, fk types.FoundationKeeper) stakingtypes.MsgServer {
	return &msgServer{
		MsgServer: stakingkeeper.NewMsgServerImpl(keeper),
		fk:        fk,
	}
}

var _ stakingtypes.MsgServer = msgServer{}

func (k msgServer) CreateValidator(goCtx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.fk.GetEnabled(ctx) {
		valAddr := sdk.ValAddress(msg.ValidatorAddress)
		grantee := valAddr.ToAccAddress()
		if err := k.fk.Accept(ctx, govtypes.ModuleName, grantee, msg); err != nil {
			return nil, err
		}
	}

	return k.MsgServer.CreateValidator(goCtx, msg)
}
