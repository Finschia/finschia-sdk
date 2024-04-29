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
		storeKey      storetypes.StoreKey
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		accountKeeper: ak,
		bankKeeper:    bk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) FswapInit(ctx sdk.Context, fswapInit types.FswapInit) error {
	// todo validate & check the first time (or not will reject the proposal)
	// todo add test for this keeper in keeper test
	if err := fswapInit.ValidateBasic(); err != nil {
		return err
	}

	// need confirm: Is ibcState necessary? (please ref upgrade keeper)
	if err := k.SetFswapInit(ctx, fswapInit); err != nil {
		return err
	}
	// need confirm: Is Swapped use sdk.Coin or sdk.Int
	swapped := types.Swapped{
		OldCoinAmount: sdk.Coin{
			Denom:  "",
			Amount: sdk.ZeroInt(),
		},
		NewCoinAmount: sdk.Coin{
			Denom:  "",
			Amount: sdk.ZeroInt(),
		},
	}
	if err := k.SetSwapped(ctx, swapped); err != nil {
		return err
	}
	return nil
}
