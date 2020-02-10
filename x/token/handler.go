package token

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		keeper.Logger(ctx).Debug("message", "decoded message", msg)
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgIssue:
			return handleMsgIssue(ctx, keeper, msg)
		case MsgMint:
			return handleMsgMint(ctx, keeper, msg)
		case MsgBurn:
			return handleMsgBurn(ctx, keeper, msg)
		case MsgCreateCollection:
			return handleMsgCreateCollection(ctx, keeper, msg)
		case MsgIssueCFT:
			return handleMsgIssueCFT(ctx, keeper, msg)
		case MsgMintCNFT:
			return handleMsgMintCNFT(ctx, keeper, msg)
		case MsgIssueCNFT:
			return handleMsgIssueCNFT(ctx, keeper, msg)
		case MsgMintCFT:
			return handleMsgMintCFT(ctx, keeper, msg)
		case MsgBurnCFT:
			return handleMsgBurnCFT(ctx, keeper, msg)
		case MsgGrantPermission:
			return handleMsgGrant(ctx, keeper, msg)
		case MsgRevokePermission:
			return handleMsgRevoke(ctx, keeper, msg)
		case MsgModifyTokenURI:
			return handleMsgModifyTokenURI(ctx, keeper, msg)
		case MsgTransferFT:
			return handleMsgTransferFT(ctx, keeper, msg)
		case MsgTransferCFT:
			return handleMsgTransferCFT(ctx, keeper, msg)
		case MsgTransferCNFT:
			return handleMsgTransferCNFT(ctx, keeper, msg)
		case MsgAttach:
			return handleMsgAttach(ctx, keeper, msg)
		case MsgDetach:
			return handleMsgDetach(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized  Msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgIssue(ctx sdk.Context, keeper keeper.Keeper, msg MsgIssue) sdk.Result {
	token := types.NewFT(msg.Name, msg.Symbol, msg.TokenURI, msg.Decimals, msg.Mintable)
	err := keeper.IssueFT(ctx, token, msg.Amount, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgIssueCFT(ctx sdk.Context, keeper keeper.Keeper, msg MsgIssueCFT) sdk.Result {
	collection, err := keeper.GetCollection(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}
	perm := types.NewIssuePermission(collection.GetSymbol())
	if !keeper.HasPermission(ctx, msg.Owner, perm) {
		return types.ErrTokenPermission(DefaultCodespace, msg.Owner, perm).Result()
	}

	token := types.NewCollectiveFT(collection, msg.Name, msg.TokenURI, msg.Decimals, msg.Mintable)
	err = keeper.IssueFTCollection(ctx, token, msg.Amount, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgIssueCNFT(ctx sdk.Context, keeper keeper.Keeper, msg MsgIssueCNFT) sdk.Result {
	_, err := keeper.GetCollection(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}

	perm := types.NewIssuePermission(msg.Symbol)
	if !keeper.HasPermission(ctx, msg.Owner, perm) {
		return types.ErrTokenPermission(DefaultCodespace, msg.Owner, perm).Result()
	}

	tokenType, err := keeper.GetNextTokenTypeForCNFT(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}

	err = keeper.IssueCNFT(ctx, msg.Symbol, tokenType, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgMintCNFT(ctx sdk.Context, keeper keeper.Keeper, msg MsgMintCNFT) sdk.Result {
	collection, err := keeper.GetCollection(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}

	perm := types.NewMintPermission(msg.Symbol + msg.TokenType)
	if !keeper.HasPermission(ctx, msg.From, perm) {
		return types.ErrTokenPermission(DefaultCodespace, msg.From, perm).Result()
	}

	if !keeper.HasTokenType(ctx, msg.Symbol, msg.TokenType) {
		return types.ErrCollectionTokenTypeNotExist(DefaultCodespace, msg.Symbol, msg.TokenType).Result()
	}

	token := types.NewCollectiveNFT(collection, msg.Name, msg.TokenType, msg.TokenURI, msg.To)
	err = keeper.MintCollectionNFT(ctx, token, msg.From)
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

func handleMsgCreateCollection(ctx sdk.Context, keeper keeper.Keeper, msg MsgCreateCollection) sdk.Result {
	collection := types.NewCollection(msg.Symbol, msg.Name)
	err := keeper.CreateCollection(ctx, collection, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgMint(ctx sdk.Context, keeper keeper.Keeper, msg MsgMint) sdk.Result {
	err := keeper.MintTokens(ctx, msg.Amount, msg.From, msg.To)
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

func handleMsgBurn(ctx sdk.Context, keeper keeper.Keeper, msg MsgBurn) sdk.Result {
	err := keeper.BurnTokens(ctx, msg.Amount, msg.From)
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

func handleMsgMintCFT(ctx sdk.Context, keeper keeper.Keeper, msg MsgMintCFT) sdk.Result {
	err := keeper.MintCollectionTokens(ctx, msg.Amount, msg.From, msg.To)
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

func handleMsgBurnCFT(ctx sdk.Context, keeper keeper.Keeper, msg MsgBurnCFT) sdk.Result {
	err := keeper.BurnCollectionTokens(ctx, msg.Amount, msg.From)
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

func handleMsgGrant(ctx sdk.Context, keeper keeper.Keeper, msg MsgGrantPermission) sdk.Result {
	err := keeper.GrantPermission(ctx, msg.From, msg.To, msg.Permission)
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
func handleMsgRevoke(ctx sdk.Context, keeper keeper.Keeper, msg MsgRevokePermission) sdk.Result {
	err := keeper.RevokePermission(ctx, msg.From, msg.Permission)
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

func handleMsgModifyTokenURI(ctx sdk.Context, keeper keeper.Keeper, msg MsgModifyTokenURI) sdk.Result {
	err := keeper.ModifyTokenURI(ctx, msg.Owner, msg.Symbol, msg.TokenID, msg.TokenURI)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTransferFT(ctx sdk.Context, k keeper.Keeper, msg types.MsgTransferFT) sdk.Result {
	err := k.TransferFT(ctx, msg.FromAddress, msg.ToAddress, msg.Symbol, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTransferCFT(ctx sdk.Context, keeper keeper.Keeper, msg MsgTransferCFT) sdk.Result {
	err := keeper.TransferCFT(ctx, msg.FromAddress, msg.ToAddress, msg.Symbol, msg.TokenID, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTransferCNFT(ctx sdk.Context, keeper keeper.Keeper, msg MsgTransferCNFT) sdk.Result {
	err := keeper.TransferCNFT(ctx, msg.FromAddress, msg.ToAddress, msg.Symbol, msg.TokenID)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAttach(ctx sdk.Context, keeper keeper.Keeper, msg MsgAttach) sdk.Result {
	err := keeper.Attach(ctx, msg.FromAddress, msg.Symbol, msg.ToTokenID, msg.TokenID)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDetach(ctx sdk.Context, keeper keeper.Keeper, msg MsgDetach) sdk.Result {
	err := keeper.Detach(ctx, msg.FromAddress, msg.ToAddress, msg.Symbol, msg.TokenID)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
