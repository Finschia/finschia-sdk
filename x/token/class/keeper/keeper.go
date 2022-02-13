package keeper

import (
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
)

// Keeper defines the class module Keeper
type Keeper struct {
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc codec.BinaryMarshaler
}

// NewKeeper returns a class keeper
func NewKeeper(
	cdc codec.BinaryMarshaler,
	key sdk.StoreKey,
) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}
