package internal

import (
	"github.com/Finschia/ostracon/libs/log"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	paramtypes "github.com/Finschia/finschia-sdk/x/params/types"
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

	// The address capable of executing privileged messages, including:
	// - MsgUpdateParams
	// - MsgUpdateDecisionPolicy
	// - MsgUpdateMembers
	// - MsgWithdrawFromTreasury
	//
	// Typically, this should be the x/foundation module account.
	authority string

	paramSpace paramtypes.Subspace
}

func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	router *baseapp.MsgServiceRouter,
	authKeeper foundation.AuthKeeper,
	bankKeeper foundation.BankKeeper,
	feeCollectorName string,
	config foundation.Config,
	authority string,
	paramSpace paramtypes.Subspace,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	// authority is x/foundation module account for now.
	if authority != foundation.DefaultAuthority().String() {
		panic("x/foundation authority must be the module account")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(foundation.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		router:           router,
		authKeeper:       authKeeper,
		bankKeeper:       bankKeeper,
		feeCollectorName: feeCollectorName,
		config:           config,
		authority:        authority,
		paramSpace:       paramSpace,
	}
}

// GetAuthority returns the x/foundation module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+foundation.ModuleName)
}
