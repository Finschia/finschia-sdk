package keeper

import (
	"fmt"

	"github.com/Finschia/ostracon/libs/log"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdktypes "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the gov module account.
	authority string

	rollupKeeper types.RollupKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority string,
	rk types.RollupKeeper,
) Keeper {

	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		authority:    authority,
		rollupKeeper: rk,
	}
}

func (k Keeper) Logger(ctx sdktypes.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) validateGovAuthority(authority string) error {
	if _, err := sdktypes.AccAddressFromBech32(authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	if k.authority != authority {
		return sdkerrors.Wrapf(sdkerrors.ErrorInvalidSigner, "invalid authority; expected %s, got %s", k.authority, authority)
	}

	return nil
}

func (k Keeper) validateSequencerAuthority(authority string) error {
	if _, err := sdktypes.AccAddressFromBech32(authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	// TODO: check if authority is a sequencer

	return nil
}