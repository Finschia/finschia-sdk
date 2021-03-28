package keeper

import (
	"fmt"

	"github.com/line/lbm-sdk/v2/x/collection/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) AddNFTOwner(ctx sdk.Context, addr sdk.AccAddress, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	tokenOwnerKey := types.AccountOwnNFTKey(k.getContractID(ctx), addr, tokenID)
	if store.Has(tokenOwnerKey) {
		panic(fmt.Sprintf("account: %s already has the token: %s", addr.String(), tokenID))
	}
	store.Set(tokenOwnerKey, []byte(tokenID))
}

func (k Keeper) DeleteNFTOwner(ctx sdk.Context, addr sdk.AccAddress, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	tokenOwnerKey := types.AccountOwnNFTKey(k.getContractID(ctx), addr, tokenID)
	if !store.Has(tokenOwnerKey) {
		panic(fmt.Sprintf("account: %s has not the token: %s", addr.String(), tokenID))
	}
	store.Delete(tokenOwnerKey)
}

func (k Keeper) HasNFTOwner(ctx sdk.Context, addr sdk.AccAddress, tokenID string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenOwnerKey := types.AccountOwnNFTKey(k.getContractID(ctx), addr, tokenID)
	return store.Has(tokenOwnerKey)
}

func (k Keeper) ChangeNFTOwner(ctx sdk.Context, from, to sdk.AccAddress, tokenID string) error {
	if !k.HasNFTOwner(ctx, from, tokenID) {
		return sdkerrors.Wrapf(types.ErrInsufficientToken, "insufficient account funds[%s]; account has no coin", k.getContractID(ctx))
	}

	k.DeleteNFTOwner(ctx, from, tokenID)
	k.AddNFTOwner(ctx, to, tokenID)
	return nil
}

func (k Keeper) GetNFTsOwner(ctx sdk.Context, addr sdk.AccAddress) (tokenIDs []string) {
	store := ctx.KVStore(k.storeKey)
	var iter = sdk.KVStorePrefixIterator(store, types.AccountOwnNFTKey(k.getContractID(ctx), addr, ""))
	defer iter.Close()
	for {
		if !iter.Valid() {
			break
		}

		val := iter.Value()
		tokenIDs = append(tokenIDs, string(val))
		iter.Next()
	}
	return tokenIDs
}
