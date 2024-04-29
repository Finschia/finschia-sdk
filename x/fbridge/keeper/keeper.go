package keeper

import (
	"errors"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Finschia/finschia-sdk/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
	"strings"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper

	// the target denom for the bridge
	targetDenom string

	// the authority address that can execute privileged operations only if the guardian group is not set
	// - UpdateParams
	// - SuggestRole
	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	targetDenom string,
	authority string,
) Keeper {
	if addr := authKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(errors.New("fbridge module account has not been set"))
	}

	if strings.TrimSpace(authority) == "" {
		panic(errors.New("authority address cannot be empty"))
	}

	return Keeper{
		storeKey:    key,
		cdc:         cdc,
		authKeeper:  authKeeper,
		bankKeeper:  bankKeeper,
		targetDenom: targetDenom,
		authority:   authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetAuthority() string {
	return k.authority
}
