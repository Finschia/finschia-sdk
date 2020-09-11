package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/genesis/internal/types"
)

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

func (k Keeper) SetGenesisMessage(ctx sdk.Context, genesisMessage string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GenesisKeyPrefix, []byte(genesisMessage))
}

func (k Keeper) GetGenesisMessage(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GenesisKeyPrefix)
	if bz == nil {
		return ""
	}
	return string(bz)
}
