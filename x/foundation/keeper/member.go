package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) updateDecisionPolicy(ctx sdk.Context, policy foundation.DecisionPolicy) error {
	info := k.GetFoundationInfo(ctx)
	info.SetDecisionPolicy(policy)
	info.Version++
	k.setFoundationInfo(ctx, info)

	return nil
}

func (k Keeper) GetFoundationInfo(ctx sdk.Context) foundation.FoundationInfo {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(foundationInfoKey)
	if len(bz) == 0 {
		panic("the foundation info must have been registered")
	}

	var info foundation.FoundationInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info
}

func (k Keeper) setFoundationInfo(ctx sdk.Context, info foundation.FoundationInfo) error {
	bz, err := k.cdc.Marshal(&info)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(foundationInfoKey, bz)

	return nil
}

func (k Keeper) updateMembers(ctx sdk.Context, members []foundation.Member) error {
	weightUpdate := sdk.ZeroDec()
	for _, new := range members {
		if new.Weight.IsZero() { // Delete
			old, err := k.GetMember(ctx, sdk.AccAddress(new.Address))
			if err != nil {
				return err
			}

			k.deleteMember(ctx, sdk.AccAddress(old.Address))
			weightUpdate = weightUpdate.Sub(old.Weight)
		}

		old, err := k.GetMember(ctx, sdk.AccAddress(new.Address))
		if err == nil {
			weightUpdate = weightUpdate.Sub(old.Weight)
		}

		k.setMember(ctx, new)
		weightUpdate = weightUpdate.Add(new.Weight)
	}

	info := k.GetFoundationInfo(ctx)
	info.TotalWeight.Add(weightUpdate)
	info.Version++
	k.setFoundationInfo(ctx, info)

	return nil
}

func (k Keeper) GetMember(ctx sdk.Context, address sdk.AccAddress) (*foundation.Member, error) {
	store := ctx.KVStore(k.storeKey)
	key := memberKey(address)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "No such member: %s", address.String())
	}

	var member foundation.Member
	if err := k.cdc.Unmarshal(bz, &member); err != nil {
		return nil, err
	}
	return &member, nil
}

func (k Keeper) setMember(ctx sdk.Context, member foundation.Member) error {
	bz, err := k.cdc.Marshal(&member)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	key := memberKey(sdk.AccAddress(member.Address))
	store.Set(key, bz)

	return nil
}

func (k Keeper) deleteMember(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := memberKey(address)
	store.Delete(key)
}

func (k Keeper) GetOperator(ctx sdk.Context) sdk.AccAddress {
	info := k.GetFoundationInfo(ctx)
	return sdk.AccAddress(info.Operator)
}

func (k Keeper) GetAdmin(ctx sdk.Context) sdk.AccAddress {
	return k.authKeeper.GetModuleAccount(ctx, foundation.AdministratorName).GetAddress()
}
