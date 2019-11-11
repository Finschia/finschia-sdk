package keeper

import (
	"fmt"
	"regexp"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/link-chain/link/x/lrc3/internal/types"
	nft "github.com/link-chain/link/x/nft"
	"github.com/link-chain/link/x/nft/exported"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	NFTKeeper nft.Keeper

	cdc *codec.Codec // The amino codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nft Keeper
func NewKeeper(cdc *codec.Codec, nftKeeper nft.Keeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey:  storeKey,
		NFTKeeper: nftKeeper,
		cdc:       cdc,
	}
}

func (k Keeper) SetLRC3(ctx sdk.Context, denom string, collection nft.Collection) sdk.Error {
	_, found := k.NFTKeeper.GetCollection(ctx, denom)
	if found {
		return types.ErrAlreadyExistLRC3(types.DefaultCodespace)
	}

	if validateDenom(denom) != nil {
		return types.ErrInvalidDenom(types.DefaultCodespace)
	}

	k.NFTKeeper.SetCollection(ctx, denom, collection)
	return nil
}

func (k Keeper) SetApproval(ctx sdk.Context, denom, tokenId string, recipient sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetApprovalKey(denom, tokenId)
	approval := types.NewApproval(tokenId, recipient)

	store.Set(key, k.cdc.MustMarshalBinaryLengthPrefixed(approval))
}

func (k Keeper) GetApproval(ctx sdk.Context, denom string, tokenId string) (approval types.Approval, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetApprovalKey(denom, tokenId))
	if b == nil {
		return approval, types.ErrNotExistApprove(types.DefaultCodespace)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &approval)
	return approval, nil
}

func (k Keeper) SetOperatorApprovals(ctx sdk.Context, denom string, ownerAddress sdk.AccAddress, operators types.Operators) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOperatorApprovalKey(denom, ownerAddress)

	store.Set(key, k.cdc.MustMarshalBinaryLengthPrefixed(operators))
}

func (k Keeper) GetOperatorApprovals(ctx sdk.Context, denom string, ownerAddress sdk.AccAddress) (operators types.Operators, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetOperatorApprovalKey(denom, ownerAddress))
	if b == nil {
		return operators, types.ErrNotExistOperatorApprovals(types.DefaultCodespace)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &operators)
	return operators, nil
}

func (k Keeper) ClearApproval(ctx sdk.Context, denom string, tokenId string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetApprovalKey(denom, tokenId)

	store.Delete(key)
}

func (k Keeper) IsApprovedOrOwner(ctx sdk.Context, spender sdk.AccAddress, symbol string, nft exported.NFT) bool {
	owner := nft.GetOwner()
	if owner.String() == "" {
		return false
	}
	tokenApprove, _ := k.GetApproval(ctx, symbol, nft.GetID())
	operators, _ := k.GetOperatorApprovals(ctx, symbol, owner)

	return owner.Equals(spender) || tokenApprove.ApprovedAddress.Equals(spender) || operators.Find(spender)
}

var (
	// Denominations can be 3 ~ 16 characters long.
	reDnmString = `[a-z][a-z0-9][a-z0-9\-\_]{2,15}`
	reDnm       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reDnmString))
)

func validateDenom(denom string) error {
	if !reDnm.MatchString(denom) {
		return fmt.Errorf("invalid denom: %s", denom)
	}
	return nil
}
