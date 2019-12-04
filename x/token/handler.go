package token

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktypes "github.com/link-chain/link/types"
	"github.com/link-chain/link/x/token/internal/keeper"
	"github.com/link-chain/link/x/token/internal/types"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgIssue:
			return handleMsgIssue(ctx, keeper, msg)
		case MsgIssueCollection:
			return handleMsgIssueCollection(ctx, keeper, msg)
		case MsgIssueNFT:
			return handleMsgIssueNFT(ctx, keeper, msg)
		case MsgIssueNFTCollection:
			return handleMsgIssueNFTCollection(ctx, keeper, msg)
		case MsgMint:
			return handleMsgMint(ctx, keeper, msg)
		case MsgBurn:
			return handleMsgBurn(ctx, keeper, msg)
		case MsgGrantPermission:
			return handleMsgGrant(ctx, keeper, msg)
		case MsgRevokePermission:
			return handleMsgRevoke(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized  Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func checkPermissionAndOccupyIfEmpty(ctx sdk.Context, keeper keeper.Keeper, symbol string, owner sdk.AccAddress) sdk.Error {
	if !keeper.HasPermission(ctx, owner, NewIssuePermission(symbol)) {

		ownerStr := owner.String()
		if symbol[len(symbol)-linktypes.AccAddrSuffixLen:] != ownerStr[len(ownerStr)-linktypes.AccAddrSuffixLen:] {
			return ErrInvalidTokenSymbol(DefaultCodespace)
		}

		err := keeper.OccupySymbol(ctx, symbol, owner)
		if err != nil {
			return ErrTokenPermission(DefaultCodespace)
		}
	}
	return nil
}

func handleMsgIssue(ctx sdk.Context, keeper keeper.Keeper, msg MsgIssue) sdk.Result {
	if err := checkPermissionAndOccupyIfEmpty(ctx, keeper, msg.Symbol, msg.Owner); err != nil {
		return err.Result()
	}

	token := NewFT(msg.Name, msg.Symbol, msg.Decimals, msg.Mintable)
	err := keeper.IssueFT(ctx, token, msg.Amount, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyTokenType, types.AttributeValueTokenTypeFT),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgIssueCollection(ctx sdk.Context, keeper keeper.Keeper, msg MsgIssueCollection) sdk.Result {
	if err := checkPermissionAndOccupyIfEmpty(ctx, keeper, msg.Symbol, msg.Owner); err != nil {
		return err.Result()
	}

	symbol := linktypes.SymbolCollectionToken(msg.Symbol, msg.TokenID)
	token := NewFT(msg.Name, symbol, msg.Decimals, msg.Mintable)
	err := keeper.IssueFT(ctx, token, msg.Amount, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyTokenType, types.AttributeValueTokenTypeCFT),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgIssueNFT(ctx sdk.Context, keeper keeper.Keeper, msg MsgIssueNFT) sdk.Result {
	if err := checkPermissionAndOccupyIfEmpty(ctx, keeper, msg.Symbol, msg.Owner); err != nil {
		return err.Result()
	}

	token := NewNFT(msg.Name, msg.Symbol, msg.TokenURI)
	err := keeper.IssueNFT(ctx, token, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyTokenType, types.AttributeValueTokenTypeNFT),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgIssueNFTCollection(ctx sdk.Context, keeper keeper.Keeper, msg MsgIssueNFTCollection) sdk.Result {
	if err := checkPermissionAndOccupyIfEmpty(ctx, keeper, msg.Symbol, msg.Owner); err != nil {
		return err.Result()
	}

	symbol := linktypes.SymbolCollectionToken(msg.Symbol, msg.TokenID)
	token := NewNFT(msg.Name, symbol, msg.TokenURI)
	err := keeper.IssueNFT(ctx, token, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyTokenType, types.AttributeValueTokenTypeCNFT),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgMint(ctx sdk.Context, keeper keeper.Keeper, msg MsgMint) sdk.Result {
	err := keeper.MintTokens(ctx, msg.Amount, msg.To)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.To.String()),
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
