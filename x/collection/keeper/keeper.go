package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Finschia/finschia-sdk/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/token/class"
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
