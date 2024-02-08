package keeper

import (
	addresscodec "cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection"
)

// Keeper defines the collection module Keeper
type Keeper struct {
	cdc          codec.Codec
	addressCodec addresscodec.Codec
	storeService store.KVStoreService
}

// NewKeeper returns a collection keeper
func NewKeeper(
	cdc codec.Codec,
	kvStoreService store.KVStoreService,
) Keeper {
	k := Keeper{
		cdc:          cdc,
		storeService: kvStoreService,
	}
	k.addressCodec = cdc.InterfaceRegistry().SigningContext().AddressCodec()
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+collection.ModuleName)
}

func ValidateLegacyContract(k Keeper, ctx sdk.Context, contractID string) error {
	if !k.HasID(ctx, contractID) {
		return collection.ErrContractNotExist.Wrap(contractID)
	}

	if _, err := k.GetContract(ctx, contractID); err != nil {
		return err
	}

	return nil
}
