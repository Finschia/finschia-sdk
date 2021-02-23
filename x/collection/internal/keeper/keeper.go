package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/params/subspace"

	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/line/lbm-sdk/x/contract"
)

type Keeper struct {
	accountKeeper  types.AccountKeeper
	contractKeeper contract.Keeper
	storeKey       sdk.StoreKey
	paramsSpace    subspace.Subspace
	cdc            *codec.Codec
}

func NewKeeper(
	cdc *codec.Codec,
	accountKeeper types.AccountKeeper,
	contractKeeper contract.Keeper,
	paramsSpace subspace.Subspace,
	storeKey sdk.StoreKey,
) Keeper {
	return Keeper{
		accountKeeper:  accountKeeper,
		contractKeeper: contractKeeper,
		storeKey:       storeKey,
		paramsSpace:    paramsSpace.WithKeyTable(types.ParamKeyTable()),
		cdc:            cdc,
	}
}
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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

func (k Keeper) UnmarshalJSON(bz []byte, ptr interface{}) error {
	return k.cdc.UnmarshalJSON(bz, ptr)
}

func (k Keeper) MarshalJSON(o interface{}) ([]byte, error) {
	return k.cdc.MarshalJSON(o)
}

func (k Keeper) MarshalJSONIndent(o interface{}) ([]byte, error) {
	return k.cdc.MarshalJSONIndent(o, "", "  ")
}

func (k Keeper) mustEncodeString(str string) []byte {
	return k.cdc.MustMarshalBinaryBare(str)
}

func (k Keeper) mustDecodeString(bz []byte) (str string) {
	k.cdc.MustUnmarshalBinaryBare(bz, &str)
	return str
}
