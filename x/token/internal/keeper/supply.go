package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
)

type SupplyKeeper interface {
	GetTotalInt(ctx sdk.Context, target string) (sdk.Int, error)
	MintSupply(ctx sdk.Context, to sdk.AccAddress, amount sdk.Int) error
	BurnSupply(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error
}

var _ SupplyKeeper = (*Keeper)(nil)

func (k Keeper) GetTotalInt(ctx sdk.Context, target string) (sdk.Int, error) {
	supply, err := k.getSupply(ctx)
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
		return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid request target to query total %s", target)
	}
}

func (k Keeper) getSupply(ctx sdk.Context) (supply types.Supply, err error) {
	if _, err := k.GetToken(ctx); err != nil {
		return nil, err
	}
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.SupplyKey(k.getContractID(ctx)))
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

func (k Keeper) setSupply(ctx sdk.Context, supply types.Supply) {
	if k.getContractID(ctx) != supply.GetContractID() {
		panic("cannot set supply with different contract id")
	}
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(types.SupplyKey(k.getContractID(ctx)), b)
}

func (k Keeper) MintSupply(ctx sdk.Context, to sdk.AccAddress, amount sdk.Int) (err error) {
	defer func() {
		// to recover from overflows
		if r := recover(); r != nil {
			err = types.WrapIfOverflowPanic(r)
		}
	}()

	_, err = k.addBalance(ctx, to, amount)
	if err != nil {
		return err
	}

	supply, err := k.getSupply(ctx)
	if err != nil {
		return err
	}
	oldSupplyAmount := supply.GetTotalSupply()
	newSupplyAmount := oldSupplyAmount.Add(amount)
	if newSupplyAmount.IsNegative() {
		return sdkerrors.Wrapf(types.ErrInsufficientSupply, "insufficient supply for token [%s]; %s < %s", k.getContractID(ctx), oldSupplyAmount, amount)
	}
	supply = supply.Inflate(amount)
	k.setSupply(ctx, supply)

	return nil
}

func (k Keeper) BurnSupply(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) (err error) {
	defer func() {
		// to recover from overflows
		// however, it will return insufficient fund error instead of panicking in the case
		if r := recover(); r != nil {
			err = types.WrapIfOverflowPanic(r)
		}
	}()

	_, err = k.subtractBalance(ctx, from, amount)
	if err != nil {
		return err
	}

	supply, err := k.getSupply(ctx)
	if err != nil {
		return err
	}
	oldSupplyAmount := supply.GetTotalSupply()
	newSupplyAmount := oldSupplyAmount.Sub(amount)
	if newSupplyAmount.IsNegative() {
		return sdkerrors.Wrapf(types.ErrInsufficientSupply, "insufficient supply for token [%s]; %s < %s", k.getContractID(ctx), oldSupplyAmount, amount)
	}
	supply = supply.Deflate(amount)
	k.setSupply(ctx, supply)

	return nil
}
