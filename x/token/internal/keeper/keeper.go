package keeper

import (
	"fmt"

	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/contract"
	"github.com/line/lbm-sdk/x/token/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	accountKeeper  types.AccountKeeper
	storeKey       sdk.StoreKey
	contractKeeper contract.Keeper
	cdc            *codec.Codec
}

func NewKeeper(cdc *codec.Codec, accountKeeper types.AccountKeeper, contractKeeper contract.Keeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		accountKeeper:  accountKeeper,
		storeKey:       storeKey,
		contractKeeper: contractKeeper,
		cdc:            cdc,
	}
}

func (k Keeper) NewContractID(ctx sdk.Context) string {
	return k.contractKeeper.NewContractID(ctx)
}
func (k Keeper) HasContractID(ctx sdk.Context) bool {
	return k.contractKeeper.HasContractID(ctx, k.getContractID(ctx))
}
func (k Keeper) getContractID(ctx sdk.Context) string {
	contractI := ctx.Context().Value(contract.CtxKey{})
	if contractI == nil {
		panic("contract id does not set on the context")
	}
	return contractI.(string)
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
