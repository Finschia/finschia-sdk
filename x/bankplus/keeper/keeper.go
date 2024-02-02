package keeper

import (
	"context"

	"cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

var _ Keeper = (*BaseKeeper)(nil)

type Keeper interface {
	bankkeeper.Keeper

	AddToInactiveAddr(ctx context.Context, address sdk.AccAddress)
	DeleteFromInactiveAddr(ctx context.Context, address sdk.AccAddress)
	IsInactiveAddr(address sdk.AccAddress) bool

	InitializeBankPlus(ctx context.Context)
}

type BaseKeeper struct {
	bankkeeper.BaseKeeper

	ak           types.AccountKeeper
	cdc          codec.Codec
	storeService store.KVStoreService

	inactiveAddrs  map[string]bool
	deactMultiSend bool
}

func NewBaseKeeper(
	cdc codec.Codec, storeService store.KVStoreService, ak types.AccountKeeper,
	blockedAddr map[string]bool, deactMultiSend bool, authority string, logger log.Logger,
) BaseKeeper {
	return BaseKeeper{
		BaseKeeper:     bankkeeper.NewBaseKeeper(cdc, storeService, ak, blockedAddr, authority, logger),
		ak:             ak,
		cdc:            cdc,
		storeService:   storeService,
		inactiveAddrs:  map[string]bool{},
		deactMultiSend: deactMultiSend,
	}
}

func (keeper BaseKeeper) InitializeBankPlus(ctx context.Context) {
	keeper.loadAllInactiveAddrs(ctx)
}

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist.
func (keeper BaseKeeper) SendCoinsFromModuleToAccount(
	ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
) error {
	senderAddr := keeper.ak.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if keeper.BlockedAddr(recipientAddr) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr)
	}

	return keeper.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another.
// It will panic if either module account does not exist.
func (keeper BaseKeeper) SendCoinsFromModuleToModule(
	ctx context.Context, senderModule, recipientModule string, amt sdk.Coins,
) error {
	senderAddr := keeper.ak.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := keeper.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return keeper.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (keeper BaseKeeper) SendCoinsFromAccountToModule(
	ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins,
) error {
	recipientAcc := keeper.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return keeper.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

func (keeper BaseKeeper) isInactiveAddr(addr sdk.AccAddress) bool {
	return keeper.inactiveAddrs[addr.String()]
}

// SendCoins transfers amt coins from a sending account to a receiving account.
// This is wrapped bank the `SendKeeper` interface of `bank` module,
// and checks if `toAddr` is a inactiveAddr managed by the module.
func (keeper BaseKeeper) SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	// if toAddr is smart contract, check the status of contract.
	if keeper.isInactiveAddr(toAddr) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", toAddr)
	}

	return keeper.BaseSendKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// AddToInactiveAddr adds the address to `inactiveAddr`.
func (keeper BaseKeeper) AddToInactiveAddr(ctx context.Context, address sdk.AccAddress) {
	if !keeper.inactiveAddrs[address.String()] {
		keeper.inactiveAddrs[address.String()] = true

		keeper.addToInactiveAddr(ctx, address)
	}
}

// DeleteFromInactiveAddr removes the address from `inactiveAddr`.
func (keeper BaseKeeper) DeleteFromInactiveAddr(ctx context.Context, address sdk.AccAddress) {
	if keeper.inactiveAddrs[address.String()] {
		delete(keeper.inactiveAddrs, address.String())

		keeper.deleteFromInactiveAddr(ctx, address)
	}
}

// IsInactiveAddr returns if the address is added in inactiveAddr.
func (keeper BaseKeeper) IsInactiveAddr(address sdk.AccAddress) bool {
	return keeper.inactiveAddrs[address.String()]
}

func (keeper BaseKeeper) InputOutputCoins(ctx context.Context, input types.Input, outputs []types.Output) error {
	if keeper.deactMultiSend {
		return sdkerrors.ErrNotSupported.Wrap("MultiSend was deactivated")
	}

	for _, out := range outputs {
		if keeper.inactiveAddrs[out.Address] {
			return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", out.Address)
		}
	}

	return keeper.BaseSendKeeper.InputOutputCoins(ctx, input, outputs)
}
