package keeper

import (
	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper/internal"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

type Keeper struct {
	impl internal.Keeper
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
) Keeper {
	return Keeper{
		impl: internal.NewKeeper(
			cdc,
			key,
			router,
			authKeeper,
			bankKeeper,
			feeCollectorName,
			config,
			authority,
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

func NewMsgServer(k Keeper) foundation.MsgServer {
	return internal.NewMsgServer(k.impl)
}

func NewQueryServer(k Keeper) foundation.QueryServer {
	return internal.NewQueryServer(k.impl)
}

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	internal.RegisterInvariants(ir, k.impl)
}

func BeginBlocker(ctx sdk.Context, k Keeper) {
	internal.BeginBlocker(ctx, k.impl)
}

func EndBlocker(ctx sdk.Context, k Keeper) {
	internal.EndBlocker(ctx, k.impl)
}

func NewFoundationProposalsHandler(k Keeper) govtypes.Handler {
	return internal.NewFoundationProposalsHandler(k.impl)
}
