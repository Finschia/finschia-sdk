package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) issue(ctx sdk.Context, class token.Token, owner, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.setClass(ctx, class); err != nil {
		return err
	}

	// zero check?
	amounts := []token.FT{
		token.FT{
			ClassId: class.Id,
			Amount: amount,
		},
	}
	if err := k.mintTokens(ctx, to, amounts); err != nil {
		return err
	}

	// TODO: emit events
	return nil
}

func (k Keeper) GetClass(ctx sdk.Context, classId string) (*token.Token, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(classKey(classId))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "No information for %s", classId)
	}

	var class token.Token
	if err := k.cdc.UnmarshalBinaryBare(bz, &class); err != nil {
		return nil, err
	}
	return &class, nil
}

func (k Keeper) setClass(ctx sdk.Context, class token.Token) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryBare(&class)
	if err != nil {
		return err
	}

	store.Set(classKey(class.Id), bz)
	return nil
}

func (k Keeper) mint(ctx sdk.Context, grantee, to sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		if ok := k.GetGrant(ctx, grantee, amount.ClassId, token.ActionMint); !ok {
			return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for mint %s tokens", grantee, amount.ClassId)
		}
	}

	if err := k.mintTokens(ctx, to, amounts); err != nil {
		return err
	}

	// TODO: emit events

	return nil
}

func (k Keeper) mintTokens(ctx sdk.Context, addr sdk.AccAddress, amounts []token.FT) error {
	if err := k.addTokens(ctx, addr, amounts); err != nil {
		return err
	}

	for _, amount := range amounts {
		supplyMint := k.GetSupply(ctx, token.SupplyMint, amount.ClassId)
		supplyMint.Amount = supplyMint.Amount.Add(amount.Amount)
		if err := k.setSupply(ctx, token.SupplyMint, supplyMint); err != nil {
			return err
		}

		supplyTotal := k.GetSupply(ctx, token.SupplyTotal, amount.ClassId)
		supplyTotal.Amount = supplyTotal.Amount.Add(amount.Amount)
		if err := k.setSupply(ctx, token.SupplyTotal, supplyTotal); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) burn(ctx sdk.Context, from sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		if ok := k.GetGrant(ctx, from, amount.ClassId, token.ActionBurn); !ok {
			return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for burn %s tokens", from, amount.ClassId)
		}
	}

	if err := k.burnTokens(ctx, from, amounts); err != nil {
		return err
	}

	// TODO: emit events
	return nil
}

func (k Keeper) burnFrom(ctx sdk.Context, proxy, from sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		if ok := k.GetGrant(ctx, proxy, amount.ClassId, token.ActionBurn); ok {
			continue
		} else if ok := k.GetProxy(ctx, from, proxy, amount.ClassId); ok {
			continue
		}
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for burn %s tokens of %s", proxy, amount.ClassId, from)
	}

	if err := k.burnTokens(ctx, from, amounts); err != nil {
		return err
	}

	// TODO: emit events

	return nil
}

func (k Keeper) burnTokens(ctx sdk.Context, addr sdk.AccAddress, amounts []token.FT) error {
	if err := k.subtractTokens(ctx, addr, amounts); err != nil {
		return err
	}

	for _, amount := range amounts {
		supplyBurn := k.GetSupply(ctx, token.SupplyBurn, amount.ClassId)
		supplyBurn.Amount = supplyBurn.Amount.Add(amount.Amount)
		if err := k.setSupply(ctx, token.SupplyBurn, supplyBurn); err != nil {
			return err
		}

		supplyTotal := k.GetSupply(ctx, token.SupplyTotal, amount.ClassId)
		supplyTotal.Amount = supplyTotal.Amount.Sub(amount.Amount)
		if err := k.setSupply(ctx, token.SupplyTotal, supplyTotal); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) GetSupply(ctx sdk.Context, supplyType, classId string) token.FT {
	store := ctx.KVStore(k.storeKey)
	amount := sdk.ZeroInt()
	bz := store.Get(supplyKey(supplyType, classId))
	if bz != nil {
		if err := amount.Unmarshal(bz); err != nil {
			panic(err)
		}
	}

	return token.FT{
		ClassId: classId,
		Amount: amount,
	}
}

// The caller must validate `supply`.
func (k Keeper) setSupply(ctx sdk.Context, supplyType string, supply token.FT) error {
	store := ctx.KVStore(k.storeKey)
	key := supplyKey(supplyType, supply.ClassId)
	if supply.Amount.IsZero() {
		store.Delete(key)
	} else {
		bz, err := supply.Amount.Marshal()
		if err != nil {
			return err
		}
		store.Set(key, bz)
	}

	return nil
}

func (k Keeper) modify() {
}

func (k Keeper) GetGrant(ctx sdk.Context, grantee sdk.AccAddress, classId, action string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(grantKey(grantee, classId, action))
}

func (k Keeper) setGrant(ctx sdk.Context, grantee sdk.AccAddress, classId, action string, set bool) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(grantee, classId, action)
	if set {
		store.Set(key, []byte{})
	} else {
		store.Delete(key)
	}
}
