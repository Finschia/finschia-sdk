package keeper

import (
	"context"

	"cosmossdk.io/core/address"
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

	ak      types.AccountKeeper
	cdc     codec.Codec
	addrCdc address.Codec

	storeService   store.KVStoreService
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
		addrCdc:        cdc.InterfaceRegistry().SigningContext().AddressCodec(),
	}
}

func (k BaseKeeper) InitializeBankPlus(ctx context.Context) {
	k.loadAllInactiveAddrs(ctx)
}

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist.
func (k BaseKeeper) SendCoinsFromModuleToAccount(
	ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if k.BlockedAddr(recipientAddr) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr)
	}

	return k.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another.
// It will panic if either module account does not exist.
func (k BaseKeeper) SendCoinsFromModuleToModule(
	ctx context.Context, senderModule, recipientModule string, amt sdk.Coins,
) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (k BaseKeeper) SendCoinsFromAccountToModule(
	ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins,
) error {
	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

func (k BaseKeeper) isInactiveAddr(addr sdk.AccAddress) bool {
	addrString, err := k.addrCdc.BytesToString(addr)
	if err != nil {
		panic(err)
	}
	return k.inactiveAddrs[addrString]
}

// SendCoins transfers amt coins from a sending account to a receiving account.
// This is wrapped bank the `SendKeeper` interface of `bank` module,
// and checks if `toAddr` is a inactiveAddr managed by the module.
func (k BaseKeeper) SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	// if toAddr is smart contract, check the status of contract.
	if k.isInactiveAddr(toAddr) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", toAddr)
	}

	return k.BaseSendKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// AddToInactiveAddr adds the address to `inactiveAddr`.
func (k BaseKeeper) AddToInactiveAddr(ctx context.Context, addr sdk.AccAddress) {
	addrString, err := k.addrCdc.BytesToString(addr)
	if err != nil {
		panic(err)
	}
	if !k.inactiveAddrs[addrString] {
		k.inactiveAddrs[addrString] = true

		k.addToInactiveAddr(ctx, addr)
	}
}

// DeleteFromInactiveAddr removes the address from `inactiveAddr`.
func (k BaseKeeper) DeleteFromInactiveAddr(ctx context.Context, addr sdk.AccAddress) {
	addrString, err := k.addrCdc.BytesToString(addr)
	if err != nil {
		panic(err)
	}
	if k.inactiveAddrs[addrString] {
		delete(k.inactiveAddrs, addrString)

		k.deleteFromInactiveAddr(ctx, addr)
	}
}

// IsInactiveAddr returns if the address is added in inactiveAddr.
func (k BaseKeeper) IsInactiveAddr(addr sdk.AccAddress) bool {
	addrString, err := k.addrCdc.BytesToString(addr)
	if err != nil {
		panic(err)
	}
	return k.inactiveAddrs[addrString]
}

func (k BaseKeeper) InputOutputCoins(ctx context.Context, input types.Input, outputs []types.Output) error {
	if k.deactMultiSend {
		return sdkerrors.ErrNotSupported.Wrap("MultiSend was deactivated")
	}

	for _, out := range outputs {
		if k.inactiveAddrs[out.Address] {
			return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", out.Address)
		}
	}

	return k.BaseSendKeeper.InputOutputCoins(ctx, input, outputs)
}
