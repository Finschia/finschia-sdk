package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
)

func (k Keeper) SetBlackList(ctx sdk.Context, addr sdk.AccAddress, action string) {
	store := ctx.KVStore(k.storeKey)

	// value is just a key w/o the module prefix
	v := addr.String() + ":" + action
	store.Set(types.BlacklistKey(addr, action), []byte(v))
}

func (k Keeper) IsBlacklisted(ctx sdk.Context, addr sdk.AccAddress, action string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.BlacklistKey(addr, action))
}
