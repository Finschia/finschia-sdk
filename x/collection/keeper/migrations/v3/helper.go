package v3

import (
	"errors"

	gogotypes "github.com/cosmos/gogoproto/types"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection"
)

func getRootOwner(store storetypes.KVStore, cdc codec.BinaryCodec, contractID, tokenID string) sdk.AccAddress {
	id := tokenID
	err := iterateAncestors(store, cdc, contractID, tokenID, func(tokenID string) error {
		id = tokenID
		return nil
	})
	if err != nil {
		panic(err)
	}

	return getOwner(store, contractID, id)
}

func getOwner(store storetypes.KVStore, contractID, tokenID string) sdk.AccAddress {
	key := ownerKey(contractID, tokenID)
	bz := store.Get(key)
	if bz == nil {
		panic("owner must exist")
	}

	var owner sdk.AccAddress
	if err := owner.Unmarshal(bz); err != nil {
		panic(err)
	}
	return owner
}

func iterateAncestors(store storetypes.KVStore, cdc codec.BinaryCodec, contractID, tokenID string, fn func(tokenID string) error) error {
	var err error
	for id := &tokenID; err == nil; id, err = getParent(store, cdc, contractID, *id) {
		if fnErr := fn(*id); fnErr != nil {
			return fnErr
		}
	}

	return nil
}

func getParent(store storetypes.KVStore, cdc codec.BinaryCodec, contractID, tokenID string) (*string, error) {
	key := parentKey(contractID, tokenID)
	bz := store.Get(key)
	if bz == nil {
		return nil, errors.New("token is not a child of some other")
	}

	var parent gogotypes.StringValue
	cdc.MustUnmarshal(bz, &parent)
	return &parent.Value, nil
}

func addCoin(store storetypes.KVStore, contractID string, address sdk.AccAddress, amount collection.Coin) {
	key := balanceKey(contractID, address, amount.TokenId)
	var beforeBalance math.Int
	bz := store.Get(key)
	if bz == nil {
		beforeBalance = math.ZeroInt()
	}
	if err := beforeBalance.Unmarshal(bz); err != nil {
		panic(err)
	}

	afterBalance := beforeBalance.Add(amount.Amount)
	bz, err := afterBalance.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}
