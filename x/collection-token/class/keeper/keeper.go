package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Keeper defines the class module Keeper
type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
}

// NewKeeper returns a class keeper
func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
	}
}
