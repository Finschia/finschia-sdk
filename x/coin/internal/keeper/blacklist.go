package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/coin/internal/types"
)

func (keeper Keeper) BlacklistAccountAction(ctx sdk.Context, addr sdk.AccAddress, action string) {
	store := ctx.KVStore(keeper.storeKey)

	// value is just a key w/o the module prefix
	v := addr.String() + ":" + action
	store.Set(types.BlacklistKey(addr, action), []byte(v))
}

func (keeper Keeper) IsBlacklistedAccountAction(ctx sdk.Context, addr sdk.AccAddress, action string) bool {
	store := ctx.KVStore(keeper.storeKey)
	return store.Has(types.BlacklistKey(addr, action))
}
