package token

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPublishToken:
			return handleMsgPublishToken(ctx, keeper, msg)
		case MsgMint:
			return handleMsgMint(ctx, keeper, msg)
		case MsgBurn:
			return handleMsgBurn(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized  Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgPublishToken(ctx sdk.Context, keeper Keeper, msg MsgPublishToken) sdk.Result {

	token := Token{Name: msg.Name, Symbol: msg.Symbol, Owner: msg.Owner, Mintable: msg.Mintable}
	err := keeper.SetToken(ctx, token)
	if err != nil {
		return err.Result()
	}
	newToken := sdk.NewCoin(msg.Symbol, msg.Amount)

	err = keeper.MintToken(ctx, newToken, token.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypePublishToken,
			sdk.NewAttribute(AttributeKeyName, msg.Name),
			sdk.NewAttribute(AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(AttributeKeyMintable, strconv.FormatBool(msg.Mintable)),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgMint(ctx sdk.Context, keeper Keeper, msg MsgMint) sdk.Result {
	for _, amount := range msg.Amount {
		token, err := keeper.GetToken(ctx, amount.Denom)
		if err != nil {
			return err.Result()
		}

		if !token.Owner.Equals(msg.To) {
			return ErrTokenPermissionMint(DefaultCodespace).Result()
		}

		if !token.Mintable {
			return ErrTokenNotMintable(DefaultCodespace).Result()
		}

		err = keeper.MintToken(ctx, amount, msg.To)
		if err != nil {
			return err.Result()
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeMintToken,
			sdk.NewAttribute(AttributeKeyTo, msg.To.String()),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBurn(ctx sdk.Context, keeper Keeper, msg MsgBurn) sdk.Result {
	for _, amount := range msg.Amount {
		token, err := keeper.GetToken(ctx, amount.Denom)
		if err != nil {
			return err.Result()
		}

		if !token.Owner.Equals(msg.From) {
			return ErrTokenPermissionBurn(DefaultCodespace).Result()
		}

		if !token.Mintable {
			return ErrTokenNotMintable(DefaultCodespace).Result()
		}

		err = keeper.BurnToken(ctx, amount, msg.From)
		if err != nil {
			return err.Result()
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeBurnToken,
			sdk.NewAttribute(AttributeKeyFrom, msg.From.String()),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}
