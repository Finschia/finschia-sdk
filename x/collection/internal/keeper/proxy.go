package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

var (
	ApprovedValue = []byte{0x01}
)

type ProxyKeeper interface {
	IsApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) bool
	SetApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) sdk.Error
	DeleteApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) sdk.Error
}

func (k Keeper) IsApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	approvedKey := types.CollectionApprovedKey(contractID, proxy, approver)
	return store.Has(approvedKey)
}

func (k Keeper) SetApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(contractID)) {
		return types.ErrCollectionNotExist(types.DefaultCodespace, contractID)
	}
	approvedKey := types.CollectionApprovedKey(contractID, proxy, approver)
	if store.Has(approvedKey) {
		return types.ErrCollectionAlreadyApproved(types.DefaultCodespace, proxy.String(), approver.String(), contractID)
	}
	store.Set(approvedKey, ApprovedValue)

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

func (k Keeper) DeleteApproved(ctx sdk.Context, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(contractID)) {
		return types.ErrCollectionNotExist(types.DefaultCodespace, contractID)
	}
	approvedKey := types.CollectionApprovedKey(contractID, proxy, approver)
	if !store.Has(approvedKey) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), approver.String(), contractID)
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
