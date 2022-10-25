package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

const (
	moduleAccountInvariant = "module-accounts"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(foundation.ModuleName, moduleAccountInvariant, ModuleAccountInvariant(k))
}

func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		treasuryAcc := k.authKeeper.GetModuleAccount(ctx, foundation.TreasuryName)
		balance := k.bankKeeper.GetAllBalances(ctx, treasuryAcc.GetAddress())

		treasury := k.GetTreasury(ctx)
		msg := fmt.Sprintf("coins in the treasury; expected %s, got %s\n", treasury, balance)
		broken := !treasury.IsEqual(sdk.NewDecCoinsFromCoins(balance...))

		return sdk.FormatInvariant(foundation.ModuleName, moduleAccountInvariant, msg), broken
	}
}
