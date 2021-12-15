package keeper

import (
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	bankkeeper "github.com/line/lbm-sdk/x/bank/keeper"
	"github.com/line/lbm-sdk/x/bank/types"
	paramtypes "github.com/line/lbm-sdk/x/params/types"
)

var _ bankkeeper.Keeper = (*BaseKeeper)(nil)

type BaseKeeper struct {
	bankkeeper.BaseKeeper

	ak           types.AccountKeeper
	cdc          codec.BinaryMarshaler
	storeKey     sdk.StoreKey
	blockedAddrs map[string]bool
}

func NewBaseKeeper(
	cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, ak types.AccountKeeper, paramSpace *paramtypes.Subspace,
	blockedAddrs map[string]bool,
) BaseKeeper {
	return BaseKeeper{
		BaseKeeper:   bankkeeper.NewBaseKeeper(cdc, storeKey, ak, paramSpace, blockedAddrs),
		ak:           ak,
		cdc:          cdc,
		storeKey:     storeKey,
		blockedAddrs: map[string]bool{},
	}
}

func (keeper BaseKeeper) InitializeBankPlus(ctx sdk.Context) {
	keeper.loadAllBlockedAddrs(ctx)
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

func (keeper BaseKeeper) isBlockedAddr(addr sdk.AccAddress) bool {
	return keeper.blockedAddrs[addr.String()]
}

// SendCoins transfers amt coins from a sending account to a receiving account.
// This is wrapped bank the `SendKeeper` interface of `bank` module,
// and check if `toAddr` is a blockedAddr managed by the module.
func (keeper BaseKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	// if toAddr is smart contract, check the status of contract.
	if keeper.isBlockedAddr(toAddr) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", toAddr)
	}

	return keeper.BaseSendKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// AddBlockedAddr add the address to `blockedAddr`.
func (keeper BaseKeeper) AddBlockedAddr(ctx sdk.Context, address sdk.AccAddress) {
	if !keeper.blockedAddrs[address.String()] {
		keeper.blockedAddrs[address.String()] = true

		keeper.addBlockedAddr(ctx, address)
	}
}

// DeleteBlockedAddr remove the address fro `blockedAddr`.
func (keeper BaseKeeper) DeleteBlockedAddr(ctx sdk.Context, address sdk.AccAddress) {
	if keeper.blockedAddrs[address.String()] {
		delete(keeper.blockedAddrs, address.String())

		keeper.deleteBlockedAddr(ctx, address)
	}
}

// IsBlockedAddr return if the address is added in blockedAddr.
func (keeper BaseKeeper) IsBlockedAddr(address sdk.AccAddress) bool {
	return keeper.blockedAddrs[address.String()]
}
