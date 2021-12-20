package stakingplus

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/stakingplus/keeper"
	"github.com/line/lbm-sdk/x/stakingplus/types"

	"github.com/line/lbm-sdk/x/staking"
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func NewHandler(k stakingkeeper.Keeper, ck types.ConsortiumKeeper) sdk.Handler {
	overriden := staking.NewHandler(k)
	msgServer := keeper.NewMsgServerImpl(k, ck)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *stakingtypes.MsgCreateValidator:
			res, err := msgServer.CreateValidator(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return overriden(ctx, msg)
		}
	}
}
