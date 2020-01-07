package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/tendermint/tendermint/libs/log"
	"strconv"
)

type Keeper struct {
	supplyKeeper  types.SupplyKeeper
	iamKeeper     types.IamKeeper
	accountKeeper types.AccountKeeper
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
}

func NewKeeper(cdc *codec.Codec, supplyKeeper types.SupplyKeeper, iamKeeper types.IamKeeper, accountKeeper types.AccountKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		supplyKeeper:  supplyKeeper,
		iamKeeper:     iamKeeper.WithPrefix(types.ModuleName),
		accountKeeper: accountKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) IssueFT(ctx sdk.Context, token types.Token, amount sdk.Int, owner sdk.AccAddress) sdk.Error {
	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	perm := types.NewMintPermission(token.Symbol)
	if token.Mintable {
		k.AddPermission(ctx, owner, perm)
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.Symbol, amount)), owner)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyName, token.Name),
			sdk.NewAttribute(types.AttributeKeySymbol, token.Symbol),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.Mintable)),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.Decimals.String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.TokenURI),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, perm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, perm.GetAction()),
		),
	})

	return nil
}

func (k Keeper) IssueNFT(ctx sdk.Context, token types.Token, owner sdk.AccAddress) sdk.Error {

	//TODO: move it to the invariant check https://github.com/line/link/issues/322
	//if !k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(token.Symbol).IsZero() {
	//	return types.ErrTokenNFTExist(types.DefaultCodespace)
	//}
	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.Symbol, sdk.NewInt(1))), owner)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyName, token.Name),
			sdk.NewAttribute(types.AttributeKeySymbol, token.Symbol),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, sdk.NewInt(1).String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.Mintable)),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.Decimals.String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.TokenURI),
		),
	})

	return nil
}
