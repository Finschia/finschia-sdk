package keeper

import (
	addresscodec "cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection-token/class"
)

// Keeper defines the collection module Keeper
type Keeper struct {
	cdc          codec.Codec
	addressCodec addresscodec.Codec
	storeService store.KVStoreService
	classKeeper  collection.ClassKeeper
}

// NewKeeper returns a collection keeper
func NewKeeper(
	cdc codec.Codec,
	kvStoreService store.KVStoreService,
	ck collection.ClassKeeper,
) Keeper {
	k := Keeper{
		cdc:          cdc,
		classKeeper:  ck,
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
	if !k.classKeeper.HasID(ctx, contractID) {
		return class.ErrContractNotExist.Wrap(contractID)
	}

	if _, err := k.GetContract(ctx, contractID); err != nil {
		return err
	}

	return nil
}
