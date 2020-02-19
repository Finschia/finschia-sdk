package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	supplyKeeper types.SupplyKeeper
	iamKeeper    types.IamKeeper
	bankKeeper   types.BankKeeper
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
}

func NewKeeper(cdc *codec.Codec, supplyKeeper types.SupplyKeeper, iamKeeper types.IamKeeper, bankKeeper types.BankKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		supplyKeeper: supplyKeeper,
		iamKeeper:    iamKeeper.WithPrefix(types.ModuleName),
		bankKeeper:   bankKeeper,
		storeKey:     storeKey,
		cdc:          cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) UnmarshalJSON(bz []byte, ptr interface{}) error {
	return k.cdc.UnmarshalJSON(bz, ptr)
}

func (k Keeper) MarshalJSON(o interface{}) ([]byte, error) {
	return k.cdc.MarshalJSON(o)
}

func (k Keeper) MarshalJSONIndent(o interface{}) ([]byte, error) {
	return k.cdc.MarshalJSONIndent(o, "", "  ")
}
