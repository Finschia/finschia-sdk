package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) Send(ctx sdk.Context, classID string, from, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.subtractToken(ctx, classID, from, amount); err != nil {
		return err
	}

	if err := k.addToken(ctx, classID, to, amount); err != nil {
		return err
	}

	if !k.accountKeeper.HasAccount(ctx, to) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, to))
	}

	return nil
}

func (k Keeper) AuthorizeOperator(ctx sdk.Context, classID string, approver, proxy sdk.AccAddress) error {
	if k.GetAuthorization(ctx, classID, approver, proxy) != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("Already authorized")
	}
	if _, err := k.GetClass(ctx, classID); err != nil {
		return sdkerrors.ErrNotFound.Wrapf("ID not exists: %s", classID)
	}

	k.setAuthorization(ctx, classID, approver, proxy)
	return nil
}

func (k Keeper) RevokeOperator(ctx sdk.Context, classID string, approver, proxy sdk.AccAddress) error {
	if k.GetAuthorization(ctx, classID, approver, proxy) == nil {
		return sdkerrors.ErrNotFound.Wrap("No authorization")
	}
	if _, err := k.GetClass(ctx, classID); err != nil {
		return sdkerrors.ErrNotFound.Wrapf("ID not exists: %s", classID)
	}

	k.setAuthorization(ctx, classID, approver, proxy)
	return nil
}

func (k Keeper) GetAuthorization(ctx sdk.Context, classID string, approver, proxy sdk.AccAddress) *token.Authorization {
	store := ctx.KVStore(k.storeKey)
	if store.Has(authorizationKey(classID, proxy, approver)) {
		return &token.Authorization{
			Approver: approver.String(),
			Proxy:    proxy.String(),
		}
	}
	return nil
}

func (k Keeper) setAuthorization(ctx sdk.Context, classID string, approver, proxy sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := authorizationKey(classID, proxy, approver)
	store.Set(key, []byte{})
}

func (k Keeper) deleteAuthorization(ctx sdk.Context, classID string, approver, proxy sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := authorizationKey(classID, proxy, approver)
	store.Delete(key)
}

func (k Keeper) subtractToken(ctx sdk.Context, classID string, addr sdk.AccAddress, amount sdk.Int) error {
	if amount.IsNegative() {
		return sdkerrors.ErrInvalidCoins.Wrap(amount.String())
	}

	balance := k.GetBalance(ctx, classID, addr)
	newBalance := balance.Sub(amount)
	if newBalance.IsNegative() {
		return sdkerrors.ErrInsufficientFunds.Wrapf("%s is smaller than %s", balance, amount)
	}

	if err := k.setBalance(ctx, classID, addr, newBalance); err != nil {
		return err
	}

	// Emit an event on token spend
	// Since: finschia
	return ctx.EventManager().EmitTypedEvent(&token.EventSpent{
		ContractId: classID,
		Spender:    addr.String(),
		Amount:     amount,
	})
}

func (k Keeper) addToken(ctx sdk.Context, classID string, addr sdk.AccAddress, amount sdk.Int) error {
	if amount.IsNegative() {
		return sdkerrors.ErrInvalidCoins.Wrap(amount.String())
	}

	balance := k.GetBalance(ctx, classID, addr)
	newBalance := balance.Add(amount)

	if err := k.setBalance(ctx, classID, addr, newBalance); err != nil {
		return err
	}

	// Emit an event on token receive
	// Since: finschia
	return ctx.EventManager().EmitTypedEvent(&token.EventReceived{
		ContractId: classID,
		Receiver:   addr.String(),
		Amount:     amount,
	})
}

func (k Keeper) GetBalance(ctx sdk.Context, classID string, addr sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	amount := sdk.ZeroInt()
	bz := store.Get(balanceKey(classID, addr))
	if bz != nil {
		if err := amount.Unmarshal(bz); err != nil {
			panic(err)
		}
	}
	return amount
}

// setBalance sets balance.
// The caller must validate `balance`.
func (k Keeper) setBalance(ctx sdk.Context, classID string, addr sdk.AccAddress, balance sdk.Int) error {
	store := ctx.KVStore(k.storeKey)
	key := balanceKey(classID, addr)
	if balance.IsZero() {
		store.Delete(key)
	} else {
		bz, err := balance.Marshal()
		if err != nil {
			return err
		}
		store.Set(key, bz)
	}

	return nil
}
