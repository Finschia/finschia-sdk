package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) CollectFoundationTax(ctx sdk.Context) error {
	// fetch and clear the collected fees for the fund, since this is
	// called in BeginBlock, collected fees will be from the previous block
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// calculate the tax
	taxRatio := k.GetFoundationTax(ctx)
	tax, _ := feesCollected.MulDecTruncate(taxRatio).TruncateDecimal()

	// collect rewards to the foundation treasury
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, foundation.TreasuryName, tax); err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetTreasury(ctx sdk.Context) sdk.Coins {
	treasury := k.authKeeper.GetModuleAccount(ctx, foundation.TreasuryName)
	return k.bankKeeper.GetAllBalances(ctx, treasury.GetAddress())
}

func (k Keeper) FundTreasury(ctx sdk.Context, from sdk.AccAddress, amt sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, foundation.TreasuryName, amt)
}

func (k Keeper) WithdrawFromTreasury(ctx sdk.Context, to sdk.AccAddress, amt sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, foundation.TreasuryName, to, amt)
}
