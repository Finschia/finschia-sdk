package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgMintNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgMintNFT) sdk.Result {
	_, err := keeper.GetCollection(ctx, msg.ContractID)
	if err != nil {
		return err.Result()
	}

	tokenID, err := keeper.GetNextTokenIDNFT(ctx, msg.ContractID, msg.TokenType)
	if err != nil {
		return err.Result()
	}

	token := types.NewNFT(msg.ContractID, tokenID, msg.Name, msg.To)
	err = keeper.MintNFT(ctx, msg.From, token)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgMintFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgMintFT) sdk.Result {
	err := keeper.MintFT(ctx, msg.ContractID, msg.From, msg.To, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
