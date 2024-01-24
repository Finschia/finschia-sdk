package internal

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func validateMetadata(metadata string, config foundation.Config) error {
	if len(metadata) > int(config.MaxMetadataLen) {
		return sdkerrors.ErrInvalidRequest.Wrap("metadata is too large")
	}

	return nil
}

func (k Keeper) UpdateDecisionPolicy(ctx sdk.Context, policy foundation.DecisionPolicy) error {
	info := k.GetFoundationInfo(ctx)
	if err := info.SetDecisionPolicy(policy); err != nil {
		return err
	}
	info.Version++
	k.SetFoundationInfo(ctx, info)

	if err := policy.Validate(info, k.config); err != nil {
		return err
	}

	// invalidate active proposals
	k.abortOldProposals(ctx)

	return nil
}

func (k Keeper) GetFoundationInfo(ctx sdk.Context) foundation.FoundationInfo {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(foundationInfoKey)
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		panic("the foundation info must have been registered")
	}

	var info foundation.FoundationInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info
}

func (k Keeper) SetFoundationInfo(ctx sdk.Context, info foundation.FoundationInfo) {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&info)
	if err := store.Set(foundationInfoKey, bz); err != nil {
		panic(err)
	}
}

func (k Keeper) UpdateMembers(ctx sdk.Context, members []foundation.MemberRequest) error {
	weightUpdate := math.LegacyZeroDec()
	for _, request := range members {
		new := foundation.Member{
			Address:  request.Address,
			Metadata: request.Metadata,
			AddedAt:  ctx.BlockTime(),
		}
		if err := new.ValidateBasic(k.addressCodec()); err != nil {
			panic(err)
		}
		if err := validateMetadata(new.Metadata, k.config); err != nil {
			return err
		}

		addr, err := k.addressCodec().StringToBytes(new.Address)
		if err != nil {
			panic(err)
		}
		old, err := k.GetMember(ctx, addr)
		if err != nil && request.Remove { // the member must exist
			return err
		}
		if err == nil { // overwrite
			weightUpdate = weightUpdate.Sub(math.LegacyOneDec())
			new.AddedAt = old.AddedAt
		}

		if request.Remove {
			k.deleteMember(ctx, addr)
		} else {
			weightUpdate = weightUpdate.Add(math.LegacyOneDec())
			k.SetMember(ctx, new)
		}
	}

	info := k.GetFoundationInfo(ctx)
	info.TotalWeight = info.TotalWeight.Add(weightUpdate)
	info.Version++
	k.SetFoundationInfo(ctx, info)

	if err := info.GetDecisionPolicy().Validate(info, k.config); err != nil {
		return err
	}

	// invalidate active proposals
	k.abortOldProposals(ctx)

	return nil
}

func (k Keeper) GetMember(ctx sdk.Context, address sdk.AccAddress) (*foundation.Member, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := memberKey(address)
	bz, err := store.Get(key)
	if err != nil {
		return nil, err
	}
	if len(bz) == 0 {
		return nil, sdkerrors.ErrNotFound.Wrapf("No such member: %s", address)
	}

	var member foundation.Member
	k.cdc.MustUnmarshal(bz, &member)

	return &member, nil
}

func (k Keeper) SetMember(ctx sdk.Context, member foundation.Member) {
	store := k.storeService.OpenKVStore(ctx)
	addr, err := k.addressCodec().StringToBytes(member.Address)
	if err != nil {
		panic(err)
	}
	key := memberKey(addr)

	bz := k.cdc.MustMarshal(&member)
	if err := store.Set(key, bz); err != nil {
		panic(err)
	}
}

func (k Keeper) deleteMember(ctx sdk.Context, address sdk.AccAddress) {
	store := k.storeService.OpenKVStore(ctx)
	key := memberKey(address)
	if err := store.Delete(key); err != nil {
		panic(err)
	}
}

func (k Keeper) iterateMembers(ctx sdk.Context, fn func(member foundation.Member) (stop bool)) {
	store := k.storeService.OpenKVStore(ctx)
	prefix := memberKeyPrefix
	adapter := runtime.KVStoreAdapter(store)
	iterator := storetypes.KVStorePrefixIterator(adapter, prefix)
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

func (k Keeper) validateAuthority(authority string) error {
	if authority != k.authority {
		return sdkerrors.ErrUnauthorized.Wrapf("invalid authority; expected %s, got %s", k.authority, authority)
	}

	return nil
}

func (k Keeper) validateCensorshipAuthority(ctx sdk.Context, msgTypeURL, authority string) error {
	censorship, err := k.GetCensorship(ctx, msgTypeURL)
	if err != nil {
		return err
	}

	govAddr, err := k.addressCodec().BytesToString(authtypes.NewModuleAddress(govtypes.ModuleName))
	if err != nil {
		return err
	}
	authorityAddrs := map[foundation.CensorshipAuthority]string{
		foundation.CensorshipAuthorityGovernance: govAddr,
		foundation.CensorshipAuthorityFoundation: k.authority,
	}
	if expected := authorityAddrs[censorship.Authority]; authority != expected {
		return sdkerrors.ErrUnauthorized.Wrapf("invalid authority; expected %s, got %s", expected, authority)
	}

	return nil
}

func (k Keeper) validateMembers(ctx sdk.Context, members []string) error {
	for _, member := range members {
		addr, err := k.addressCodec().StringToBytes(member)
		if err != nil {
			panic(err)
		}
		if _, err := k.GetMember(ctx, addr); err != nil {
			return sdkerrors.ErrUnauthorized.Wrapf("%s is not a member", member)
		}
	}

	return nil
}
