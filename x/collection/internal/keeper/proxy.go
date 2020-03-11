package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/collection/internal/types"
)

var (
	ApprovedValue = []byte{0x01}
)

type ProxyKeeper interface {
	IsApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) bool
	SetApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) error
	DeleteApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) error
}

func (k Keeper) IsApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	approvedKey := types.CollectionApprovedKey(contractID, proxy, approver)
	return store.Has(approvedKey)
}

func (k Keeper) SetApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(contractID)) {
		return sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", contractID)
	}
	approvedKey := types.CollectionApprovedKey(contractID, proxy, approver)
	if store.Has(approvedKey) {
		return sdkerrors.Wrapf(types.ErrCollectionAlreadyApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), approver.String(), contractID)
	}
	store.Set(approvedKey, ApprovedValue)

	// Set Account if not exists yet
	account := k.accountKeeper.GetAccount(ctx, proxy)
	if account == nil {
		account = k.accountKeeper.NewAccountWithAddress(ctx, proxy)
		k.accountKeeper.SetAccount(ctx, account)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeApproveCollection,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyApprover, approver.String()),
		),
	})

	return nil
}

func (k Keeper) DeleteApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(contractID)) {
		return sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", contractID)
	}
	approvedKey := types.CollectionApprovedKey(contractID, proxy, approver)
	if !store.Has(approvedKey) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), approver.String(), contractID)
	}
	store.Delete(approvedKey)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDisapproveCollection,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyApprover, approver.String()),
		),
	})

	return nil
}
