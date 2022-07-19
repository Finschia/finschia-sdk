package keeper

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/telemetry"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

// Keeper defines the collection module Keeper
type Keeper struct {
	accountKeeper collection.AccountKeeper
	classKeeper   collection.ClassKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc codec.Codec
}

// NewKeeper returns a collection keeper
func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	ak collection.AccountKeeper,
	ck collection.ClassKeeper,
) Keeper {
	return Keeper{
		accountKeeper: ak,
		classKeeper:   ck,
		storeKey:      key,
		cdc:           cdc,
	}
}

func (k Keeper) createAccountOnAbsence(ctx sdk.Context, address sdk.AccAddress) {
	if !k.accountKeeper.HasAccount(ctx, address) {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, address))
	}
}
