package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

type SupplyKeeper interface {
	GetTotalInt(ctx sdk.Context, contractID, target string) (sdk.Int, sdk.Error)
	MintSupply(ctx sdk.Context, contractID string, to sdk.AccAddress, amount sdk.Int) sdk.Error
	BurnSupply(ctx sdk.Context, contractID string, from sdk.AccAddress, amount sdk.Int) sdk.Error
}

var _ SupplyKeeper = (*Keeper)(nil)

func (k Keeper) GetTotalInt(ctx sdk.Context, contractID, target string) (sdk.Int, sdk.Error) {
	supply, err := k.getSupply(ctx, contractID)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	switch target {
	case types.QuerySupply:
		return supply.GetTotalSupply(), nil
	case types.QueryBurn:
		return supply.GetTotalBurn(), nil
	case types.QueryMint:
		return supply.GetTotalMint(), nil
	default:
		return sdk.ZeroInt(), sdk.ErrInternal(fmt.Sprintf("invalid request target to query total %s", target))
	}
}

func (k Keeper) getSupply(ctx sdk.Context, contractID string) (supply types.Supply, err sdk.Error) {
	if _, err := k.GetToken(ctx, contractID); err != nil {
		return nil, err
	}
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.SupplyKey(contractID))
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

func (k Keeper) setSupply(ctx sdk.Context, supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(types.SupplyKey(supply.GetContractID()), b)
}

func (k Keeper) MintSupply(ctx sdk.Context, contractID string, to sdk.AccAddress, amount sdk.Int) sdk.Error {
	_, err := k.addBalance(ctx, contractID, to, amount)
	if err != nil {
		return err
	}

	supply, err := k.getSupply(ctx, contractID)
	if err != nil {
		return err
	}
	oldSupplyAmount := supply.GetTotalSupply()
	newSupplyAmount := oldSupplyAmount.Add(amount)
	if newSupplyAmount.IsNegative() {
		return types.ErrInsufficientSupply(types.DefaultCodespace, fmt.Sprintf("insufficient supply for token [%s]; %s < %s", contractID, oldSupplyAmount, amount))
	}
	supply = supply.Inflate(amount)
	k.setSupply(ctx, supply)

	return nil
}

func (k Keeper) BurnSupply(ctx sdk.Context, contractID string, from sdk.AccAddress, amount sdk.Int) sdk.Error {
	_, err := k.subtractBalance(ctx, contractID, from, amount)
	if err != nil {
		return err
	}

	supply, err := k.getSupply(ctx, contractID)
	if err != nil {
		return err
	}
	oldSupplyAmount := supply.GetTotalSupply()
	newSupplyAmount := oldSupplyAmount.Sub(amount)
	if newSupplyAmount.IsNegative() {
		return types.ErrInsufficientSupply(types.DefaultCodespace, fmt.Sprintf("insufficient supply for token [%s]; %s < %s", contractID, oldSupplyAmount, amount))
	}
	supply = supply.Deflate(amount)
	k.setSupply(ctx, supply)

	return nil
}
