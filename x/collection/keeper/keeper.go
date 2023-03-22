package keeper

import (
	"github.com/line/ostracon/libs/log"

	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
	"github.com/line/lbm-sdk/x/token/class"
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

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+collection.ModuleName)
}

func ValidateLegacyContract(k Keeper, ctx sdk.Context, contractID string) error {
	if !k.classKeeper.HasID(ctx, contractID) {
		return class.ErrContractNotExist.Wrap(contractID)
	}

	if _, err := k.GetContract(ctx, contractID); err != nil {
		return err
	}

	return nil
}
