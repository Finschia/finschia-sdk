package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal"
)

type Keeper struct {
	impl internal.Keeper
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
	subspace paramstypes.Subspace,
) Keeper {
	return Keeper{
		impl: internal.NewKeeper(
			cdc,
			storeService,
			router,
			authKeeper,
			bankKeeper,
			feeCollectorName,
			config,
			authority,
			subspace,
		),
	}
}

// GetAuthority returns the x/foundation module's authority.
func (k Keeper) GetAuthority() string {
	return k.impl.GetAuthority()
}

func (k Keeper) Accept(ctx sdk.Context, grantee sdk.AccAddress, msg sdk.Msg) error {
	return k.impl.Accept(ctx, grantee, msg)
}

func (k Keeper) InitGenesis(ctx sdk.Context, gs *foundation.GenesisState) error {
	return k.impl.InitGenesis(ctx, gs)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *foundation.GenesisState {
	return k.impl.ExportGenesis(ctx)
}

func (k Keeper) Hooks() internal.Hooks {
	return k.impl.Hooks()
}

func NewMsgServer(k Keeper) foundation.MsgServer {
	return internal.NewMsgServer(k.impl)
}

func NewQueryServer(k Keeper) foundation.QueryServer {
	return internal.NewQueryServer(k.impl)
}

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	internal.RegisterInvariants(ir, k.impl)
}

func BeginBlocker(ctx sdk.Context, k Keeper) error {
	return internal.BeginBlocker(ctx, k.impl)
}

func EndBlocker(ctx sdk.Context, k Keeper) error {
	return internal.EndBlocker(ctx, k.impl)
}

func NewFoundationProposalsHandler(k Keeper) govtypes.Handler {
	return internal.NewFoundationProposalsHandler(k.impl)
}

type Migrator = internal.Migrator

func NewMigrator(k Keeper) Migrator {
	return internal.NewMigrator(k.impl)
}
