package keeper

import (
	"github.com/Finschia/finschia-rdk/codec"
	sdk "github.com/Finschia/finschia-rdk/types"
)

// Keeper defines the class module Keeper
type Keeper struct {
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc codec.Codec
}

// NewKeeper returns a class keeper
func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}
