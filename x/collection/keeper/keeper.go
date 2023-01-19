package keeper

import (
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

// Keeper defines the collection module Keeper
type Keeper struct {
	classKeeper collection.ClassKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc codec.Codec
}

// NewKeeper returns a collection keeper
func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	ck collection.ClassKeeper,
) Keeper {
	return Keeper{
		classKeeper: ck,
		storeKey:    key,
		cdc:         cdc,
	}
}
