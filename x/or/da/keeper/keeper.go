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

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
) *Keeper {

	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
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

func (k Keeper) validateSequencerAuthority() error {
	panic("implement me")
}
