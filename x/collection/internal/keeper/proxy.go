package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
)

var (
	ApprovedValue = []byte{0x01}
)

type ProxyKeeper interface {
	IsApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) bool
	SetApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) error
	DeleteApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) error
}

func (k Keeper) IsApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	approvedKey := types.CollectionApprovedKey(k.getContractID(ctx), proxy, approver)
	return store.Has(approvedKey)
}

func (k Keeper) SetApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(k.getContractID(ctx))) {
		return sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", k.getContractID(ctx))
	}
	approvedKey := types.CollectionApprovedKey(k.getContractID(ctx), proxy, approver)
	if store.Has(approvedKey) {
		return sdkerrors.Wrapf(types.ErrCollectionAlreadyApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), approver.String(), k.getContractID(ctx))
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
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyApprover, approver.String()),
		),
	})

	return nil
}

func (k Keeper) DeleteApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(k.getContractID(ctx))) {
		return sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", k.getContractID(ctx))
	}
	approvedKey := types.CollectionApprovedKey(k.getContractID(ctx), proxy, approver)
	if !store.Has(approvedKey) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), approver.String(), k.getContractID(ctx))
	}
	store.Delete(approvedKey)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDisapproveCollection,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyApprover, approver.String()),
		),
	})

	return nil
}
