package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) Send(ctx sdk.Context, contractID string, from, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.subtractToken(ctx, contractID, from, amount); err != nil {
		return err
	}

	if err := k.addToken(ctx, contractID, to, amount); err != nil {
		return err
	}

	return nil
}

func (k Keeper) AuthorizeOperator(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) error {
	if _, err := k.GetClass(ctx, contractID); err != nil {
		return sdkerrors.ErrNotFound.Wrapf("ID not exists: %s", contractID)
	}
	if _, err := k.GetAuthorization(ctx, contractID, holder, operator); err == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("Already authorized")
	}

	k.setAuthorization(ctx, contractID, holder, operator)

	if !k.accountKeeper.HasAccount(ctx, operator) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, operator))
	}

	return nil
}

func (k Keeper) RevokeOperator(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) error {
	if _, err := k.GetClass(ctx, contractID); err != nil {
		return sdkerrors.ErrNotFound.Wrapf("ID not exists: %s", contractID)
	}
	if _, err := k.GetAuthorization(ctx, contractID, holder, operator); err != nil {
		return err
	}

	k.deleteAuthorization(ctx, contractID, holder, operator)
	return nil
}

func (k Keeper) GetAuthorization(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) (*token.Authorization, error) {
	store := ctx.KVStore(k.storeKey)
	if store.Has(authorizationKey(contractID, operator, holder)) {
		return &token.Authorization{
			Holder:   holder.String(),
			Operator: operator.String(),
		}, nil
	}
	return nil, sdkerrors.ErrNotFound.Wrapf("no authorization by %s to %s", holder, operator)
}

func (k Keeper) setAuthorization(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := authorizationKey(contractID, operator, holder)
	store.Set(key, []byte{})
}

func (k Keeper) deleteAuthorization(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := authorizationKey(contractID, operator, holder)
	store.Delete(key)
}

func (k Keeper) subtractToken(ctx sdk.Context, contractID string, addr sdk.AccAddress, amount sdk.Int) error {
	if amount.IsNegative() {
		return sdkerrors.ErrInvalidCoins.Wrap(amount.String())
	}

	balance := k.GetBalance(ctx, contractID, addr)
	newBalance := balance.Sub(amount)
	if newBalance.IsNegative() {
		return sdkerrors.ErrInsufficientFunds.Wrapf("%s is smaller than %s", balance, amount)
	}

	if err := k.setBalance(ctx, contractID, addr, newBalance); err != nil {
		return err
	}

	return nil
}

func (k Keeper) addToken(ctx sdk.Context, contractID string, addr sdk.AccAddress, amount sdk.Int) error {
	if amount.IsNegative() {
		return sdkerrors.ErrInvalidCoins.Wrap(amount.String())
	}

	balance := k.GetBalance(ctx, contractID, addr)
	newBalance := balance.Add(amount)

	if err := k.setBalance(ctx, contractID, addr, newBalance); err != nil {
		return err
	}

	if !k.accountKeeper.HasAccount(ctx, addr) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, addr))
	}

	return nil
}

func (k Keeper) GetBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	amount := sdk.ZeroInt()
	bz := store.Get(balanceKey(contractID, addr))
	if bz != nil {
		if err := amount.Unmarshal(bz); err != nil {
			panic(err)
		}
	}
	return amount
}

// setBalance sets balance.
// The caller must validate `balance`.
func (k Keeper) setBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, balance sdk.Int) error {
	store := ctx.KVStore(k.storeKey)
	key := balanceKey(contractID, addr)
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
