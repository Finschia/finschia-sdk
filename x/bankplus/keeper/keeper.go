package keeper

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

var _ bankkeeper.Keeper = (*BaseKeeper)(nil)

type BaseKeeper struct {
	bankkeeper.BaseKeeper

	ak      types.AccountKeeper
	cdc     codec.Codec
	addrCdc address.Codec

	storeService  store.KVStoreService
	inactiveAddrs map[string]bool
}

func NewBaseKeeper(
	cdc codec.Codec, storeService store.KVStoreService, ak types.AccountKeeper,
	blockedAddr map[string]bool, authority string, logger log.Logger,
) BaseKeeper {
	keeper := bankkeeper.NewBaseKeeper(cdc, storeService, ak, blockedAddr, authority, logger)
	baseKeeper := BaseKeeper{
		BaseKeeper:    keeper,
		ak:            ak,
		cdc:           cdc,
		storeService:  storeService,
		inactiveAddrs: map[string]bool{},
		addrCdc:       cdc.InterfaceRegistry().SigningContext().AddressCodec(),
	}

	return baseKeeper
}
