package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection"
)

func legacyNFTNotFoundError(k Keeper, ctx sdk.Context, contractID, tokenID string) error {
	if err2 := collection.ValidateLegacyNFTID(tokenID); err2 == nil /* "==" is intentional */ {
		return collection.ErrTokenNotExist.Wrap(tokenID)
	}

	classID := collection.SplitTokenID(tokenID)
	if _, err2 := k.GetTokenClass(ctx, contractID, classID); err2 != nil {
		return collection.ErrTokenNotExist.Wrap(tokenID)
	}

	return collection.ErrTokenNotNFT.Wrap(tokenID)
}

func (k Keeper) hasNFT(ctx sdk.Context, contractID, tokenID string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := nftKey(contractID, tokenID)
	if ok, _ := store.Has(key); !ok {
		return legacyNFTNotFoundError(k, ctx, contractID, tokenID)
	}
	return nil
}

func (k Keeper) GetNFT(ctx sdk.Context, contractID, tokenID string) (*collection.NFT, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := nftKey(contractID, tokenID)
	bz, err := store.Get(key)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, legacyNFTNotFoundError(k, ctx, contractID, tokenID)
	}

	var token collection.NFT
	k.cdc.MustUnmarshal(bz, &token)

	return &token, nil
}

func (k Keeper) setNFT(ctx sdk.Context, contractID string, token collection.NFT) {
	store := k.storeService.OpenKVStore(ctx)
	key := nftKey(contractID, token.TokenId)

	bz, err := token.Marshal()
	if err != nil {
		panic(err)
	}
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) deleteNFT(ctx sdk.Context, contractID, tokenID string) {
	store := k.storeService.OpenKVStore(ctx)
	key := nftKey(contractID, tokenID)
	err := store.Delete(key)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) getOwner(ctx sdk.Context, contractID, tokenID string) sdk.AccAddress {
	store := k.storeService.OpenKVStore(ctx)
	key := ownerKey(contractID, tokenID)
	bz, _ := store.Get(key)
	if bz == nil {
		panic("owner must exist")
	}

	var owner sdk.AccAddress
	if err := owner.Unmarshal(bz); err != nil {
		panic(err)
	}
	return owner
}

func (k Keeper) setOwner(ctx sdk.Context, contractID, tokenID string, owner sdk.AccAddress) {
	store := k.storeService.OpenKVStore(ctx)
	key := ownerKey(contractID, tokenID)

	bz, err := owner.Marshal()
	if err != nil {
		panic(err)
	}
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) deleteOwner(ctx sdk.Context, contractID, tokenID string) {
	store := k.storeService.OpenKVStore(ctx)
	key := ownerKey(contractID, tokenID)
	err := store.Delete(key)
	if err != nil {
		panic(err)
	}
}
