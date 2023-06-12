package keeper

import (
	"fmt"

	"github.com/Finschia/finschia-sdk/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"github.com/Finschia/finschia-sdk/x/or/rollup/types"
	paramtypes "github.com/Finschia/finschia-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		bankKeeper types.BankKeeper
		authKeeper types.AccountKeeper
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(cdc codec.BinaryCodec, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, storeKey, memKey sdk.StoreKey, ps paramtypes.Subspace) Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr.Empty() {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:        cdc,
		bankKeeper: bankKeeper,
		authKeeper: accountKeeper,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) GetRollupAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
}
