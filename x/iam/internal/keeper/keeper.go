package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/link-chain/link/x/iam/exported"
	"github.com/link-chain/link/x/iam/internal/types"
)

type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
	subname  []byte
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		subname:  []byte{},
	}
}

func (k Keeper) GetAccountPermission(ctx sdk.Context, addr sdk.AccAddress) (accPerm types.AccountPermissionI) {
	store := k.kvstore(ctx)
	bz := store.Get(types.AddressStoreKey(addr))
	if bz != nil {
		accPerm = k.mustDecodeAccountPermission(bz)
		if iam, ok := accPerm.(*types.InheritedAccountPermission); ok {
			iam.SetParent(k.GetAccountPermission(ctx, iam.ParentAddr))
		}
		return accPerm
	}
	return types.NewAccountPermission(addr)
}
func (k Keeper) GetPermissions(ctx sdk.Context, addr sdk.AccAddress) (pms []exported.PermissionI) {
	accPerm := k.GetAccountPermission(ctx, addr)
	return accPerm.GetPermissions()
}

func (k Keeper) SetAccountPermission(ctx sdk.Context, accPerm types.AccountPermissionI) {
	store := k.kvstore(ctx)
	store.Set(types.AddressStoreKey(accPerm.GetAddress()), k.cdc.MustMarshalBinaryBare(accPerm))
}

func (k Keeper) InheritPermission(ctx sdk.Context, parent, child sdk.AccAddress) {
	childPerm := k.GetAccountPermission(ctx, child)
	parentPerm := k.GetAccountPermission(ctx, parent)
	childPerm.InheritAccountPermission(parentPerm)
	k.SetAccountPermission(ctx, childPerm)
}

func (k Keeper) GrantPermission(ctx sdk.Context, addr sdk.AccAddress, p exported.PermissionI) {
	accPerm := k.GetAccountPermission(ctx, addr)
	accPerm.AddPermission(p)
	k.SetAccountPermission(ctx, accPerm)
}

func (k Keeper) RevokePermission(ctx sdk.Context, addr sdk.AccAddress, p exported.PermissionI) {
	accPerm := k.GetAccountPermission(ctx, addr)
	accPerm.RemovePermission(p)
	k.SetAccountPermission(ctx, accPerm)
}

func (k Keeper) HasPermission(ctx sdk.Context, addr sdk.AccAddress, p exported.PermissionI) bool {
	accPerm := k.GetAccountPermission(ctx, addr)
	return accPerm.HasPermission(p)
}

func (k Keeper) WithPrefix(prefix string) exported.IamKeeper {
	return Keeper{k.cdc, k.storeKey, []byte(prefix)}
}

func (k Keeper) kvstore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), append(k.subname, '/'))
}

func (k Keeper) mustDecodeAccountPermission(bz []byte) (accPerm types.AccountPermissionI) {
	err := k.cdc.UnmarshalBinaryBare(bz, &accPerm)
	if err != nil {
		panic(err)
	}
	return
}
