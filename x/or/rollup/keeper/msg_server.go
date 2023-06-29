package keeper

import (
	"context"

	sdk "github.com/Finschia/finschia-sdk/types"

	"github.com/Finschia/finschia-sdk/x/or/rollup/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateRollup(goCtx context.Context, msg *types.MsgCreateRollup) (*types.MsgCreateRollupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, found := k.GetRollup(ctx, msg.RollupName); found {
		return nil, types.ErrExistRollupName
	}

	rollup := types.Rollup{
		RollupName:            msg.RollupName,
		Creator:               msg.Creator,
		MaxSequencers:         msg.MaxSequencers,
		PermissionedAddresses: msg.PermissionedAddresses,
	}

	k.SetRollup(ctx, rollup)

	return &types.MsgCreateRollupResponse{}, nil
}

func (k msgServer) RegisterSequencer(goCtx context.Context, msg *types.MsgRegisterSequencer) (*types.MsgRegisterSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, found := k.GetSequencer(ctx, msg.Creator); found {
		return nil, types.ErrSequencerExists
	}

	rollup, found := k.GetRollup(ctx, msg.RollupName)
	if found {
		permissionedAddresses := rollup.PermissionedAddresses.Addresses
		if len(permissionedAddresses) > 0 {
			permissioned := false
			for i := range permissionedAddresses {
				if permissionedAddresses[i] == msg.Creator {
					permissioned = true
					break
				}
			}
			if !permissioned {
				return nil, types.ErrSequencerNotPermissioned
			}
		}
	} else {
		return nil, types.ErrNotExistRollupName
	}

	sequencersByRollup, found := k.GetSequencersByRollupName(ctx, msg.RollupName)
	sequencer := types.Sequencer{
		SequencerAddress: msg.Creator,
		Pubkey:           msg.Pubkey,
		RollupName:       msg.RollupName,
	}

	if found {
		// TODO: Modify when rollup:sequencer=1:N
		return nil, types.SequencersByRollupExists
	} else {
		sequencersByRollup.RollupName = msg.RollupName
		sequencersByRollup.Sequencers = append(sequencersByRollup.Sequencers, sequencer)

		k.SetSequencersByRollup(ctx, sequencersByRollup)
	}

	k.SetSequencer(ctx, sequencer)

	return &types.MsgRegisterSequencerResponse{}, nil
}

func (k msgServer) RemoveSequencer(goCtx context.Context, msg *types.MsgRemoveSequencer) (*types.MsgRemoveSequencerResponse, error) {
	panic("implement me")
}
