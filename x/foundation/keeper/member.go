package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

func validateMetadata(metadata string, config foundation.Config) error {
	if len(metadata) > int(config.MaxMetadataLen) {
		return sdkerrors.ErrInvalidRequest.Wrap("metadata is too large")
	}

	return nil
}

func (k Keeper) UpdateDecisionPolicy(ctx sdk.Context, policy foundation.DecisionPolicy) error {
	if err := policy.Validate(k.config); err != nil {
		return err
	}

	info := k.GetFoundationInfo(ctx)
	info.SetDecisionPolicy(policy)
	info.Version++
	k.setFoundationInfo(ctx, info)

	// invalidate active proposals
	k.abortOldProposals(ctx)

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

func (k Keeper) UpdateMembers(ctx sdk.Context, members []foundation.Member) error {
	weightUpdate := sdk.ZeroDec()
	for _, new := range members {
		if err := validateMetadata(new.Metadata, k.config); err != nil {
			return err
		}

		new.AddedAt = ctx.BlockTime()
		old, err := k.GetMember(ctx, sdk.AccAddress(new.Address))
		if err == nil {
			weightUpdate = weightUpdate.Sub(sdk.OneDec())
			new.AddedAt = old.AddedAt
		}

		deleting := !new.Participating
		if err != nil && deleting { // the member must exist
			return err
		}

		if deleting {
			k.deleteMember(ctx, sdk.AccAddress(old.Address))
		} else {
			weightUpdate = weightUpdate.Add(sdk.OneDec())
			k.setMember(ctx, new)
		}
	}

	info := k.GetFoundationInfo(ctx)
	info.TotalWeight.Add(weightUpdate)
	info.Version++
	k.setFoundationInfo(ctx, info)

	// invalidate active proposals
	k.abortOldProposals(ctx)

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

func (k Keeper) iterateMembers(ctx sdk.Context, fn func(member foundation.Member) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := memberKeyPrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var member foundation.Member
		k.cdc.MustUnmarshal(iterator.Value(), &member)
		if stop := fn(member); stop {
			break
		}
	}
}

func (k Keeper) GetMembers(ctx sdk.Context) []foundation.Member {
	var members []foundation.Member
	k.iterateMembers(ctx, func(member foundation.Member) (stop bool) {
		members = append(members, member)
		return false
	})

	return members
}

func (k Keeper) GetOperator(ctx sdk.Context) sdk.AccAddress {
	info := k.GetFoundationInfo(ctx)
	return sdk.AccAddress(info.Operator)
}

func (k Keeper) UpdateOperator(ctx sdk.Context, operator sdk.AccAddress) error {
	info := k.GetFoundationInfo(ctx)
	if operator.String() == info.Operator {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s is already the operator", operator)
	}

	info.Operator = operator.String()
	if err := k.setFoundationInfo(ctx, info); err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAdmin(ctx sdk.Context) sdk.AccAddress {
	return k.authKeeper.GetModuleAccount(ctx, foundation.AdministratorName).GetAddress()
}

func (k Keeper) validateOperator(ctx sdk.Context, operator string) error {
	if sdk.AccAddress(operator) != k.GetOperator(ctx) {
		return sdkerrors.ErrUnauthorized.Wrapf("%s is not the operator", operator)
	}

	return nil
}

func (k Keeper) validateMembers(ctx sdk.Context, members []string) error {
	for _, member := range members {
		if _, err := k.GetMember(ctx, sdk.AccAddress(member)); err != nil {
			return sdkerrors.ErrUnauthorized.Wrapf("%s is not a member", member)
		}
	}

	return nil
}
