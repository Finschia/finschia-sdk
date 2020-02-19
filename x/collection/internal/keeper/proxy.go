package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

var (
	ApprovedValue = []byte{0x01}
)

type ProxyKeeper interface {
	IsApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress, symbol string) bool
	SetApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress, symbol string) sdk.Error
	DeleteApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress, symbol string) sdk.Error
}

func (k Keeper) IsApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	approvedKey := types.CollectionApprovedKey(proxy, approver, symbol)
	return store.Has(approvedKey)
}

func (k Keeper) SetApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress, symbol string) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(symbol)) {
		return types.ErrCollectionNotExist(types.DefaultCodespace, symbol)
	}
	approvedKey := types.CollectionApprovedKey(proxy, approver, symbol)
	if store.Has(approvedKey) {
		return types.ErrCollectionAlreadyApproved(types.DefaultCodespace, proxy.String(), approver.String(), symbol)
	}
	store.Set(approvedKey, ApprovedValue)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeApproveCollection,
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyApprover, approver.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
		),
	})

	return nil
}

func (k Keeper) DeleteApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress, symbol string) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(symbol)) {
		return types.ErrCollectionNotExist(types.DefaultCodespace, symbol)
	}
	approvedKey := types.CollectionApprovedKey(proxy, approver, symbol)
	if !store.Has(approvedKey) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), approver.String(), symbol)
	}
	store.Delete(approvedKey)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDisapproveCollection,
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyApprover, approver.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
		),
	})

	return nil
}
