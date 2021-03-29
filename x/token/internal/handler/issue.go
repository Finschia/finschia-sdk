package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/contract"
	"github.com/line/lbm-sdk/v2/x/token/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
)

func handleMsgIssue(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgIssue) (*sdk.Result, error) {
	contractI := ctx.Context().Value(contract.CtxKey{})
	if contractI == nil {
		panic("contract id does not set")
	}
	token := types.NewToken(contractI.(string), msg.Name, msg.Symbol, msg.Meta, msg.ImageURI, msg.Decimals, msg.Mintable)
	err := keeper.IssueToken(ctx, token, msg.Amount, msg.Owner, msg.To)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
