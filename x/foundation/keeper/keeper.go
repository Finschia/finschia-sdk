package keeper

import (
	"github.com/line/ostracon/libs/log"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/stakingplus"
)

// Keeper defines the foundation module Keeper
type Keeper struct {
	// The codec for binary encoding/decoding.
	cdc codec.Codec

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	router *baseapp.MsgServiceRouter

	// keepers
	authKeeper    foundation.AuthKeeper
	bankKeeper    foundation.BankKeeper
	stakingKeeper foundation.StakingKeeper

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
	stakingKeeper foundation.StakingKeeper,
	feeCollectorName string,
	config foundation.Config,
) Keeper {
	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		router:           router,
		authKeeper:       authKeeper,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		feeCollectorName: feeCollectorName,
		config:           config,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+foundation.ModuleName)
}

// Cleaning up the states
func (k Keeper) Cleanup(ctx sdk.Context) {
	msgTypeURL := stakingplus.CreateValidatorAuthorization{}.MsgTypeURL()
	var createValidatorGrantees []sdk.AccAddress
	k.iterateAuthorizations(ctx, func(grantee sdk.AccAddress, authorization foundation.Authorization) (stop bool) {
		if authorization.MsgTypeURL() == msgTypeURL {
			createValidatorGrantees = append(createValidatorGrantees, grantee)
		}
		return false
	})

	for _, grantee := range createValidatorGrantees {
		if err := k.Revoke(ctx, grantee, msgTypeURL); err != nil {
			panic(err)
		}
	}
}
