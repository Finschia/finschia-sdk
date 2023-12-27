package internal

import (
	addresscodec "cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

// Keeper defines the foundation module Keeper
type Keeper struct {
	// The codec for binary encoding/decoding.
	cdc codec.Codec
	addressCodec addresscodec.Codec

	storeService store.KVStoreService

	router baseapp.MessageRouter

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

	subspace paramstypes.Subspace
}

func NewKeeper(
	cdc codec.Codec,
	addressCodec addresscodec.Codec,
	storeService store.KVStoreService,
	router baseapp.MessageRouter,
	authKeeper foundation.AuthKeeper,
	bankKeeper foundation.BankKeeper,
	feeCollectorName string,
	config foundation.Config,
	authority string,
	subspace paramstypes.Subspace,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	// authority is x/foundation module account for now.
	if authority != foundation.DefaultAuthority().String() {
		panic("x/foundation authority must be the module account")
	}

	// TODO(@0Tech): remove x/params dependency
	// set KeyTable if it has not already been set
	if !subspace.HasKeyTable() {
		subspace = subspace.WithKeyTable(foundation.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		addressCodec:     addressCodec,
		storeService:     storeService,
		router:           router,
		authKeeper:       authKeeper,
		bankKeeper:       bankKeeper,
		feeCollectorName: feeCollectorName,
		config:           config,
		authority:        authority,
		subspace:         subspace,
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
