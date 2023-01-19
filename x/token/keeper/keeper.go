package keeper

import (
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

// Keeper defines the token module Keeper
type Keeper struct {
	classKeeper token.ClassKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc codec.Codec
}

// NewKeeper returns a token keeper
func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	ck token.ClassKeeper,
) Keeper {
	return Keeper{
		classKeeper: ck,
		storeKey:    key,
		cdc:         cdc,
	}
}
