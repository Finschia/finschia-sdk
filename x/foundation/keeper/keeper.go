package keeper

import (
	"github.com/line/ostracon/libs/log"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

// Keeper defines the foundation module Keeper
type Keeper struct {
	// The codec for binary encoding/decoding.
	cdc codec.Codec

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	router *baseapp.MsgServiceRouter

	// keepers
	authKeeper foundation.AuthKeeper
	bankKeeper foundation.BankKeeper

	feeCollectorName string

	config foundation.Config
}

// NewKeeper returns a foundation keeper. It handles:
// - updating validator auths.
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	router *baseapp.MsgServiceRouter,
	authKeeper foundation.AuthKeeper,
	bankKeeper foundation.BankKeeper,
	feeCollectorName string,
	config foundation.Config,
) Keeper {
	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		router:           router,
		authKeeper:       authKeeper,
		bankKeeper:       bankKeeper,
		feeCollectorName: feeCollectorName,
		config:           config,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+foundation.ModuleName)
}
