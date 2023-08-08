package keeper

import (
	"fmt"

	"github.com/Finschia/finschia-rdk/codec"
	sdk "github.com/Finschia/finschia-rdk/types"
	authtypes "github.com/Finschia/finschia-rdk/x/auth/types"
	"github.com/Finschia/finschia-rdk/x/or/rollup/types"
	paramtypes "github.com/Finschia/finschia-rdk/x/params/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	bankKeeper types.BankKeeper
	authKeeper types.AccountKeeper
	storeKey   sdk.StoreKey
	memKey     sdk.StoreKey
	paramstore paramtypes.Subspace
}

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

func (k Keeper) Slash(ctx sdk.Context, rollupName string, sequencerAddress string, value sdk.Coin) error {
	_, found := k.GetRollup(ctx, rollupName)
	if !found {
		return types.ErrNotExistRollupName
	}
	_, found = k.GetSequencer(ctx, sequencerAddress)
	if !found {
		return types.ErrNotExistSequencer
	}

	deposit, found := k.GetDeposit(ctx, rollupName, sequencerAddress)
	if !found {
		return types.ErrNotFoundDeposit
	}

	slashAmount := value
	if deposit.Value.Amount.LT(value.Amount) {
		slashAmount = deposit.Value
		deposit.Value = sdk.NewCoin(value.Denom, sdk.NewInt(0))
	} else {
		deposit.Value.Sub(value)
	}

	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(slashAmount))
	if err != nil {
		return err
	}

	k.SetDeposit(ctx, deposit)

	return nil
}
