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
	info := k.GetFoundationInfo(ctx)
	info.SetDecisionPolicy(policy)
	info.Version++
	k.setFoundationInfo(ctx, info)

	if err := policy.Validate(info, k.config); err != nil {
		return err
	}

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

func (k Keeper) setFoundationInfo(ctx sdk.Context, info foundation.FoundationInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&info)
	store.Set(foundationInfoKey, bz)
}

func (k Keeper) UpdateMembers(ctx sdk.Context, members []foundation.MemberRequest) error {
	weightUpdate := sdk.ZeroDec()
	for _, request := range members {
		new := foundation.Member{
			Address:  request.Address,
			Metadata: request.Metadata,
			AddedAt:  ctx.BlockTime(),
		}
		if err := new.ValidateBasic(); err != nil {
			panic(err)
		}
		if err := validateMetadata(new.Metadata, k.config); err != nil {
			return err
		}

		addr, _ := sdk.AccAddressFromBech32(new.Address)
		old, err := k.GetMember(ctx, addr)
		if err != nil && request.Remove { // the member must exist
			return err
		}
		if err == nil { // overwrite
			weightUpdate = weightUpdate.Sub(sdk.OneDec())
			new.AddedAt = old.AddedAt
		}

		if request.Remove {
			k.deleteMember(ctx, addr)
		} else {
			weightUpdate = weightUpdate.Add(sdk.OneDec())
			k.setMember(ctx, new)
		}
	}

	info := k.GetFoundationInfo(ctx)
	info.TotalWeight = info.TotalWeight.Add(weightUpdate)
	info.Version++
	k.setFoundationInfo(ctx, info)

	if err := info.GetDecisionPolicy().Validate(info, k.config); err != nil {
		return err
	}

	// invalidate active proposals
	k.abortOldProposals(ctx)

	return nil
}

func (k Keeper) GetMember(ctx sdk.Context, address sdk.AccAddress) (*foundation.Member, error) {
	store := ctx.KVStore(k.storeKey)
	key := memberKey(address)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil, sdkerrors.ErrNotFound.Wrapf("No such member: %s", address)
	}

	var member foundation.Member
	k.cdc.MustUnmarshal(bz, &member)

	return &member, nil
}

func (k Keeper) setMember(ctx sdk.Context, member foundation.Member) {
	store := ctx.KVStore(k.storeKey)
	addr, err := sdk.AccAddressFromBech32(member.Address)
	if err != nil {
		panic(err)
	}
	key := memberKey(addr)

	bz := k.cdc.MustMarshal(&member)
	store.Set(key, bz)
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
	operator, err := sdk.AccAddressFromBech32(info.Operator)
	if err != nil {
		panic(err)
	}
	return operator
}

func (k Keeper) UpdateOperator(ctx sdk.Context, operator sdk.AccAddress) error {
	info := k.GetFoundationInfo(ctx)
	if operator.String() == info.Operator {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s is already the operator", operator)
	}

	info.Operator = operator.String()
	k.setFoundationInfo(ctx, info)

	return nil
}

func (k Keeper) GetAdmin(ctx sdk.Context) sdk.AccAddress {
	return k.authKeeper.GetModuleAccount(ctx, foundation.AdministratorName).GetAddress()
}

func (k Keeper) validateOperator(ctx sdk.Context, operator string) error {
	addr, err := sdk.AccAddressFromBech32(operator)

	if err != nil {
		return err
	}

	if !addr.Equals(k.GetOperator(ctx)) {
		return sdkerrors.ErrUnauthorized.Wrapf("%s is not the operator", operator)
	}

	return nil
}

func (k Keeper) validateMembers(ctx sdk.Context, members []string) error {
	for _, member := range members {
		addr, err := sdk.AccAddressFromBech32(member)
		if err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", member)
		}
		if _, err := k.GetMember(ctx, addr); err != nil {
			return sdkerrors.ErrUnauthorized.Wrapf("%s is not a member", member)
		}
	}

	return nil
}
