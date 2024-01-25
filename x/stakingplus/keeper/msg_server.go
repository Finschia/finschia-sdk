package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/stakingplus"
)

type msgServer struct {
	stakingtypes.MsgServer

	fk stakingplus.FoundationKeeper
}

// NewMsgServerImpl returns an implementation of the staking MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *stakingkeeper.Keeper, fk stakingplus.FoundationKeeper) stakingtypes.MsgServer {
	return &msgServer{
		MsgServer: stakingkeeper.NewMsgServerImpl(keeper),
		fk:        fk,
	}
}

var _ stakingtypes.MsgServer = msgServer{}

func (k msgServer) CreateValidator(goCtx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, errors.ErrInvalidAddress.Wrapf("invalid validator address: %s", msg.ValidatorAddress)
	}

	grantee := sdk.AccAddress(valAddr)

	if err := k.fk.Accept(ctx, grantee, msg); err != nil {
		return nil, err
	}

	return k.MsgServer.CreateValidator(goCtx, msg)
}
