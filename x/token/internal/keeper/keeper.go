package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	supplyKeeper  types.SupplyKeeper
	iamKeeper     types.IamKeeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
}

func NewKeeper(cdc *codec.Codec, supplyKeeper types.SupplyKeeper, iamKeeper types.IamKeeper, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		supplyKeeper:  supplyKeeper,
		iamKeeper:     iamKeeper.WithPrefix(types.ModuleName),
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) IssueFT(ctx sdk.Context, token types.FT, amount sdk.Int, owner sdk.AccAddress) sdk.Error {
	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetDenom(), amount)), owner)
	if err != nil {
		return err
	}

	mintPerm := types.NewMintPermission(token.GetDenom())
	if token.GetMintable() {
		k.AddPermission(ctx, owner, mintPerm)
	}

	tokenUriModifyPerm := types.NewModifyTokenURIPermission(token.GetDenom())
	k.AddPermission(ctx, owner, tokenUriModifyPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyDenom, token.GetDenom()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.GetMintable())),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.GetDecimals().String()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, mintPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, mintPerm.GetAction()),
		),
		sdk.NewEvent(
			types.EventTypeModifyTokenURIPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, tokenUriModifyPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, tokenUriModifyPerm.GetAction()),
		),
	})

	return nil
}

func (k Keeper) IssueNFT(ctx sdk.Context, token types.NFT, owner sdk.AccAddress) sdk.Error {

	//TODO: move it to the invariant check https://github.com/line/link/issues/322
	//if !k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(token.Symbol).IsZero() {
	//	return types.ErrTokenNFTExist(types.DefaultCodespace)
	//}
	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetDenom(), sdk.NewInt(1))), owner)
	if err != nil {
		return err
	}

	tokenUriModifyPerm := types.NewModifyTokenURIPermission(token.GetDenom())
	k.AddPermission(ctx, owner, tokenUriModifyPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyDenom, token.GetDenom()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
		sdk.NewEvent(
			types.EventTypeModifyTokenURIPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, tokenUriModifyPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, tokenUriModifyPerm.GetAction()),
		),
	})

	return nil
}

func (k Keeper) CreateCollection(ctx sdk.Context, collection types.Collection, owner sdk.AccAddress) sdk.Error {
	err := k.SetCollection(ctx, collection)
	if err != nil {
		return err
	}

	perm := types.NewIssuePermission(collection.GetSymbol())
	k.AddPermission(ctx, owner, perm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCollection,
			sdk.NewAttribute(types.AttributeKeyName, collection.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, collection.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
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

func (k Keeper) ModifyTokenURI(ctx sdk.Context, owner sdk.AccAddress, symbol, tokenID, tokenURI string) sdk.Error {
	token, err := k.GetToken(ctx, symbol, tokenID)
	if err != nil {
		return err
	}
	tokenURIModifyPerm := types.NewModifyTokenURIPermission(token.GetDenom())
	if !k.HasPermission(ctx, owner, tokenURIModifyPerm) {
		return types.ErrTokenPermission(types.DefaultCodespace, owner, tokenURIModifyPerm)
	}
	token.SetTokenURI(tokenURI)

	err = k.ModifyToken(ctx, token)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyTokenURI,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyDenom, token.GetDenom()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
	})
	return nil
}
