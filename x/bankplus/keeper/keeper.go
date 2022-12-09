package keeper

import (
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	bankkeeper "github.com/line/lbm-sdk/x/bank/keeper"
	"github.com/line/lbm-sdk/x/bank/types"
	paramtypes "github.com/line/lbm-sdk/x/params/types"
)

var _ Keeper = (*BaseKeeper)(nil)

type Keeper interface {
	bankkeeper.Keeper

	AddToInactiveAddr(ctx sdk.Context, address sdk.AccAddress)
	DeleteFromInactiveAddr(ctx sdk.Context, address sdk.AccAddress)
	IsInactiveAddr(address sdk.AccAddress) bool

	InitializeBankPlus(ctx sdk.Context)
}

type BaseKeeper struct {
	bankkeeper.BaseKeeper

	ak            types.AccountKeeper
	cdc           codec.Codec
	storeKey      sdk.StoreKey
	inactiveAddrs map[string]bool
}

func NewBaseKeeper(
	cdc codec.Codec, storeKey sdk.StoreKey, ak types.AccountKeeper, paramSpace paramtypes.Subspace,
	blockedAddr map[string]bool,
) BaseKeeper {
	return BaseKeeper{
		BaseKeeper:    bankkeeper.NewBaseKeeper(cdc, storeKey, ak, paramSpace, blockedAddr),
		ak:            ak,
		cdc:           cdc,
		storeKey:      storeKey,
		inactiveAddrs: map[string]bool{},
	}
}

func (keeper BaseKeeper) InitializeBankPlus(ctx sdk.Context) {
	keeper.loadAllInactiveAddrs(ctx)
}

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist.
func (keeper BaseKeeper) SendCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
) error {
	senderAddr := keeper.ak.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	return keeper.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another.
// It will panic if either module account does not exist.
func (keeper BaseKeeper) SendCoinsFromModuleToModule(
	ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
) error {
	senderAddr := keeper.ak.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := keeper.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return keeper.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (keeper BaseKeeper) SendCoinsFromAccountToModule(
	ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins,
) error {
	recipientAcc := keeper.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return keeper.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

func (keeper BaseKeeper) isInactiveAddr(addr sdk.AccAddress) bool {
	return keeper.inactiveAddrs[addr.String()]
}

// SendCoins transfers amt coins from a sending account to a receiving account.
// This is wrapped bank the `SendKeeper` interface of `bank` module,
// and checks if `toAddr` is a inactiveAddr managed by the module.
func (keeper BaseKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	// if toAddr is smart contract, check the status of contract.
	if keeper.isInactiveAddr(toAddr) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", toAddr)
	}

	return keeper.BaseSendKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// AddToInactiveAddr adds the address to `inactiveAddr`.
func (keeper BaseKeeper) AddToInactiveAddr(ctx sdk.Context, address sdk.AccAddress) {
	if !keeper.inactiveAddrs[address.String()] {
		keeper.inactiveAddrs[address.String()] = true

		keeper.addToInactiveAddr(ctx, address)
	}
}

// DeleteFromInactiveAddr removes the address from `inactiveAddr`.
func (keeper BaseKeeper) DeleteFromInactiveAddr(ctx sdk.Context, address sdk.AccAddress) {
	if keeper.inactiveAddrs[address.String()] {
		delete(keeper.inactiveAddrs, address.String())

		keeper.deleteFromInactiveAddr(ctx, address)
	}
}

// IsInactiveAddr returns if the address is added in inactiveAddr.
func (keeper BaseKeeper) IsInactiveAddr(address sdk.AccAddress) bool {
	return keeper.inactiveAddrs[address.String()]
}
