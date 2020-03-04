package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err sdk.Error)
	GetOrNewAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err sdk.Error)
	GetAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err sdk.Error)
	SetAccount(ctx sdk.Context, acc types.Account) sdk.Error
	UpdateAccount(ctx sdk.Context, acc types.Account) sdk.Error
	GetBalance(ctx sdk.Context, contractID, tokenID string, addr sdk.AccAddress) (sdk.Int, sdk.Error)
}

var _ AccountKeeper = (*Keeper)(nil)

func (k Keeper) GetBalance(ctx sdk.Context, contractID, tokenID string, addr sdk.AccAddress) (sdk.Int, sdk.Error) {
	acc, err := k.GetAccount(ctx, contractID, addr)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return acc.GetCoins().AmountOf(tokenID), nil
}

func (k Keeper) NewAccountWithAddress(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err sdk.Error) {
	acc = types.NewBaseAccountWithAddress(contractID, addr)
	if err = k.SetAccount(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (k Keeper) SetAccount(ctx sdk.Context, acc types.Account) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	accKey := types.AccountKey(acc.GetContractID(), acc.GetAddress())
	if store.Has(accKey) {
		return types.ErrAccountExist(types.DefaultCodespace, acc.GetAddress())
	}
	store.Set(accKey, k.cdc.MustMarshalBinaryBare(acc))

	// Set Account if not exists yet
	account := k.accountKeeper.GetAccount(ctx, acc.GetAddress())
	if account == nil {
		account = k.accountKeeper.NewAccountWithAddress(ctx, acc.GetAddress())
		k.accountKeeper.SetAccount(ctx, account)
	}

	return nil
}

func (k Keeper) UpdateAccount(ctx sdk.Context, acc types.Account) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	accKey := types.AccountKey(acc.GetContractID(), acc.GetAddress())
	if !store.Has(accKey) {
		return types.ErrAccountNotExist(types.DefaultCodespace, acc.GetAddress())
	}
	store.Set(accKey, k.cdc.MustMarshalBinaryBare(acc))
	return nil
}

func (k Keeper) GetOrNewAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err sdk.Error) {
	acc, err = k.GetAccount(ctx, contractID, addr)
	if err != nil {
		acc, err = k.NewAccountWithAddress(ctx, contractID, addr)
		if err != nil {
			return nil, err
		}
	}
	return acc, nil
}

func (k Keeper) GetAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	accKey := types.AccountKey(contractID, addr)
	if !store.Has(accKey) {
		return nil, types.ErrAccountNotExist(types.DefaultCodespace, addr)
	}
	bz := store.Get(accKey)
	return k.mustDecodeAccount(bz), nil
}

func (k Keeper) mustDecodeAccount(bz []byte) (acc types.Account) {
	err := k.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}
