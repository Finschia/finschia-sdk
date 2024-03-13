package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

var _ bankkeeper.Keeper = (*BaseKeeper)(nil)

type BaseKeeper struct {
	bankkeeper.BaseKeeper
	storeService store.KVStoreService
}

func NewBaseKeeper(
	cdc codec.Codec, storeService store.KVStoreService, ak types.AccountKeeper,
	blockedAddr map[string]bool, authority string, logger log.Logger,
) BaseKeeper {
	return BaseKeeper{
		BaseKeeper:   bankkeeper.NewBaseKeeper(cdc, storeService, ak, blockedAddr, authority, logger),
		storeService: storeService,
	}
}
