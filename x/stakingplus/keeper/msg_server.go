package keeper

import (
	"github.com/Finschia/finschia-sdk/types/errors"
	stakingkeeper "github.com/Finschia/finschia-sdk/x/staking/keeper"
	stakingtypes "github.com/Finschia/finschia-sdk/x/staking/types"

	"context"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/stakingplus"
)

type msgServer struct {
	stakingtypes.MsgServer

	fk stakingplus.FoundationKeeper
}

// NewMsgServerImpl returns an implementation of the staking MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper stakingkeeper.Keeper, fk stakingplus.FoundationKeeper) stakingtypes.MsgServer {
	return &msgServer{
		MsgServer: stakingkeeper.NewMsgServerImpl(keeper),
		fk:        fk,
	}
}

var _ stakingtypes.MsgServer = msgServer{}

func (k msgServer) CreateValidator(goCtx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	grantee, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, errors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", msg.DelegatorAddress)
	}

	if err := k.fk.Accept(ctx, grantee, msg); err != nil {
		return nil, err
	}

	return k.MsgServer.CreateValidator(goCtx, msg)
}
