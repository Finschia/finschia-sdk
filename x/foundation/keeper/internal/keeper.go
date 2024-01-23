package internal

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

// Keeper defines the foundation module Keeper
type Keeper struct {
	// The codec for binary encoding/decoding.
	cdc codec.Codec

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
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	router baseapp.MessageRouter,
	authKeeper foundation.AuthKeeper,
	bankKeeper foundation.BankKeeper,
	feeCollectorName string,
	config foundation.Config,
	authority string,
) Keeper {
	addressCodec := cdc.InterfaceRegistry().SigningContext().AddressCodec()

	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	// authority is x/foundation module account for now.
	defaultAuthority, err := addressCodec.BytesToString(foundation.DefaultAuthority())
	if err != nil {
		panic(err)
	}
	if authority != defaultAuthority {
		panic("x/foundation authority must be the module account")
	}

	return Keeper{
		cdc:              cdc,
		storeService:     storeService,
		router:           router,
		authKeeper:       authKeeper,
		bankKeeper:       bankKeeper,
		feeCollectorName: feeCollectorName,
		config:           config,
		authority:        authority,
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

func (k Keeper) addressCodec() address.Codec {
	return k.cdc.InterfaceRegistry().SigningContext().AddressCodec()
}
