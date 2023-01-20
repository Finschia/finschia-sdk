package keeper

import (
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/class"
)

// Keeper defines the token module Keeper
type Keeper struct {
	accountKeeper token.AccountKeeper
	classKeeper   token.ClassKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc codec.Codec
}

// NewKeeper returns a token keeper
func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	ak token.AccountKeeper,
	ck token.ClassKeeper,
) Keeper {
	return Keeper{
		accountKeeper: ak,
		classKeeper:   ck,
		storeKey:      key,
		cdc:           cdc,
	}
}

func ValidateLegacyContract(k Keeper, ctx sdk.Context, contractID string) error {
	if !k.classKeeper.HasID(ctx, contractID) {
		return class.ErrContractNotExist.Wrap(contractID)
	}

	if _, err := k.GetClass(ctx, contractID); err != nil {
		return err
	}

	return nil
}
