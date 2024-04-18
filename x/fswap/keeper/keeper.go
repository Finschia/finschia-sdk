package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		storeKey      storetypes.StoreKey
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:           cdc,
		accountKeeper: ak,
		bankKeeper:    bk,
		storeKey:      storeKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
