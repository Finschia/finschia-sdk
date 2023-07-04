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
	}
	sequencersByRollup.RollupName = msg.RollupName
	sequencersByRollup.Sequencers = append(sequencersByRollup.Sequencers, sequencer)

	k.SetSequencersByRollup(ctx, sequencersByRollup)

	k.SetSequencer(ctx, sequencer)

	depositMsg := types.NewMsgDeposit(msg.RollupName, msg.Creator, msg.Value)
	_, err := k.Deposit(goCtx, depositMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterSequencerResponse{}, nil
}

func (k msgServer) RemoveSequencer(goCtx context.Context, msg *types.MsgRemoveSequencer) (*types.MsgRemoveSequencerResponse, error) {
	panic("implement me")
}

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := k.GetRollup(ctx, msg.RollupName)
	if !found {
		return nil, types.ErrNotExistRollupName
	}

	sequencersInfo, found := k.GetSequencersByRollupName(ctx, msg.RollupName)
	if !found {
		return nil, types.ErrNotExistSequencer
	}
	sequencerContains := false
	for _, v := range sequencersInfo.Sequencers {
		if v.SequencerAddress == msg.SequencerAddress {
			sequencerContains = true
			break
		}
	}
	if !sequencerContains {
		return nil, types.ErrNotExistSequencer
	}

	addr, err := sdk.AccAddressFromBech32(msg.SequencerAddress)
	if err != nil {
		return nil, err
	}
	amount := sdk.NewCoins(msg.Value)
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, amount)
	if err != nil {
		return nil, err
	}
	deposit := types.Deposit{
		RollupName:       msg.RollupName,
		SequencerAddress: msg.SequencerAddress,
		Value:            msg.Value,
	}

	k.SetDeposit(ctx, deposit)

	return &types.MsgDepositResponse{}, nil
}

func (k msgServer) Slash(ctx sdk.Context, rollupName string, sequencerAddress string, value sdk.Coin) error {
	_, found := k.GetRollup(ctx, rollupName)
	if !found {
		return types.ErrNotExistRollupName
	}
	_, found = k.GetSequencer(ctx, sequencerAddress)
	if !found {
		return types.ErrNotExistSequencer
	}

	deposit, found := k.GetDeposit(ctx, rollupName, sequencerAddress)
	if !found {
		return types.ErrNotFoundDeposit
	}

	slashAmount := value
	if deposit.Value.Amount.LT(value.Amount) {
		slashAmount = deposit.Value
		deposit.Value = sdk.NewCoin(value.Denom, sdk.NewInt(0))
	} else {
		deposit.Value.Sub(value)
	}

	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(slashAmount))
	if err != nil {
		return err
	}

	k.SetDeposit(ctx, deposit)

	return nil
}

func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := k.GetRollup(ctx, msg.RollupName)
	if !found {
		return nil, types.ErrNotExistRollupName
	}
	_, found = k.GetSequencer(ctx, msg.SequencerAddress)
	if !found {
		return nil, types.ErrNotExistSequencer
	}

	deposit, found := k.GetDeposit(ctx, msg.RollupName, msg.SequencerAddress)
	if !found {
		return nil, types.ErrNotFoundDeposit
	}

	afterWithdrawAmount := deposit.Value.Sub(msg.Value)

	if afterWithdrawAmount.Amount.LT(sdk.NewInt(types.MinDepositAmount)) {
		return nil, types.ErrIsNotEnoughDepositAmount
	}

	addr, err := sdk.AccAddressFromBech32(msg.SequencerAddress)
	if err != nil {
		return nil, err
	}
	amount := sdk.NewCoins(msg.Value)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, amount)
	if err != nil {
		return nil, err
	}
	deposit = types.Deposit{
		RollupName:       msg.RollupName,
		SequencerAddress: msg.SequencerAddress,
		Value:            afterWithdrawAmount,
	}

	k.SetDeposit(ctx, deposit)

	return &types.MsgWithdrawResponse{}, nil

}
