// nolint:unparam
package keeper

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
)

type TokenKeeper interface {
	GetToken(ctx sdk.Context, tokenID string) (types.Token, error)
	HasToken(ctx sdk.Context, tokenID string) bool
	SetToken(ctx sdk.Context, token types.Token) error
	DeleteToken(ctx sdk.Context, tokenID string) error
	UpdateToken(ctx sdk.Context, token types.Token) error
	GetTokens(ctx sdk.Context) (tokens types.Tokens, err error)
	GetFT(ctx sdk.Context, tokenID string) (types.FT, error)
	GetFTs(ctx sdk.Context) (tokens types.Tokens, err error)
	GetNFT(ctx sdk.Context, tokenID string) (types.NFT, error)
	GetNFTCount(ctx sdk.Context, tokenType string) (sdk.Int, error)
	GetNFTCountInt(ctx sdk.Context, tokenType, target string) (sdk.Int, error)
	GetNFTs(ctx sdk.Context, tokenType string) (tokens types.Tokens, err error)
	GetNextTokenIDFT(ctx sdk.Context) (string, error)
	GetNextTokenIDNFT(ctx sdk.Context, tokenType string) (string, error)
}

var _ TokenKeeper = (*Keeper)(nil)

func (k Keeper) GetToken(ctx sdk.Context, tokenID string) (types.Token, error) {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(k.getContractID(ctx), tokenID)
	bz := store.Get(tokenKey)
	if bz == nil {
		return nil, sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", k.getContractID(ctx), tokenID)
	}
	token := k.mustDecodeToken(bz)
	return token, nil
}
func (k Keeper) HasToken(ctx sdk.Context, tokenID string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(k.getContractID(ctx), tokenID)
	return store.Has(tokenKey)
}

func (k Keeper) SetToken(ctx sdk.Context, token types.Token) error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(k.getContractID(ctx), token.GetTokenID())
	if store.Has(tokenKey) {
		return sdkerrors.Wrapf(types.ErrTokenExist, "ContractID: %s, TokenID: %s", k.getContractID(ctx), token.GetTokenID())
	}
	store.Set(tokenKey, k.mustEncodeToken(token))
	tokenType := token.GetTokenType()
	if tokenType[0] == types.FungibleFlag[0] {
		k.setNextTokenTypeFT(ctx, tokenType)
	} else {
		k.setNextTokenIndexNFT(ctx, tokenType, token.GetTokenIndex())
	}
	return nil
}

func (k Keeper) UpdateToken(ctx sdk.Context, token types.Token) error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(k.getContractID(ctx), token.GetTokenID())
	if !store.Has(tokenKey) {
		return sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TotkenID: %s", k.getContractID(ctx), token.GetTokenID())
	}
	store.Set(tokenKey, k.mustEncodeToken(token))
	return nil
}

func (k Keeper) DeleteToken(ctx sdk.Context, tokenID string) error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(k.getContractID(ctx), tokenID)
	if !store.Has(tokenKey) {
		return sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TotkenID: %s", k.getContractID(ctx), tokenID)
	}
	store.Delete(tokenKey)
	return nil
}

func (k Keeper) GetTokens(ctx sdk.Context) (tokens types.Tokens, err error) {
	_, err = k.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	k.iterateToken(ctx, "", false, func(t types.Token) bool {
		tokens = append(tokens, t)
		return false
	})
	return tokens, nil
}

func (k Keeper) GetFTs(ctx sdk.Context) (tokens types.Tokens, err error) {
	_, err = k.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	k.iterateToken(ctx, types.FungibleFlag, false, func(t types.Token) bool {
		tokens = append(tokens, t)
		return false
	})
	return tokens, nil
}

func (k Keeper) GetFT(ctx sdk.Context, tokenID string) (types.FT, error) {
	token, err := k.GetToken(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	ft, ok := token.(types.FT)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrTokenNotNFT, "TokenID: %s", token.GetTokenID())
	}
	return ft, nil
}

func (k Keeper) GetNFT(ctx sdk.Context, tokenID string) (types.NFT, error) {
	token, err := k.GetToken(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	nft, ok := token.(types.NFT)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrTokenNotNFT, "TokenID: %s", token.GetTokenID())
	}
	return nft, nil
}

func (k Keeper) GetNFTs(ctx sdk.Context, tokenType string) (tokens types.Tokens, err error) {
	_, err = k.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	k.iterateToken(ctx, tokenType, false, func(t types.Token) bool {
		tokens = append(tokens, t)
		return false
	})
	return tokens, nil
}

func (k Keeper) GetNFTCount(ctx sdk.Context, tokenType string) (sdk.Int, error) {
	_, err := k.GetCollection(ctx)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	tokens, err := k.GetNFTs(ctx, tokenType)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return sdk.NewInt(int64(len(tokens))), nil
}

func (k Keeper) GetNFTCountInt(ctx sdk.Context, tokenType, target string) (sdk.Int, error) {
	_, err := k.GetCollection(ctx)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	_, err = k.GetTokenType(ctx, tokenType)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	switch target {
	case types.QueryNFTCount:
		return k.getNFTCountTotal(ctx, tokenType), nil
	case types.QueryNFTMint:
		return k.getNFTCountMint(ctx, tokenType), nil
	case types.QueryNFTBurn:
		return k.getNFTCountBurn(ctx, tokenType), nil
	default:
		panic("invalid request target to query")
	}
}
func (k Keeper) getNFTCountTotal(ctx sdk.Context, tokenType string) sdk.Int {
	return k.getNFTCountMint(ctx, tokenType).Sub(k.getNFTCountBurn(ctx, tokenType))
}
func (k Keeper) getNFTCountMint(ctx sdk.Context, tokenType string) sdk.Int {
	return k.getTokenTypeMintCount(ctx, tokenType)
}

func (k Keeper) getNFTCountBurn(ctx sdk.Context, tokenType string) sdk.Int {
	return k.getTokenTypeBurnCount(ctx, tokenType)
}

func (k Keeper) setNextTokenTypeFT(ctx sdk.Context, tokenType string) {
	store := ctx.KVStore(k.storeKey)
	tokenType = nextID(tokenType)
	store.Set(types.NextTokenTypeFTKey(k.getContractID(ctx)), k.mustEncodeString(tokenType))
}
func (k Keeper) setNextTokenTypeNFT(ctx sdk.Context, tokenType string) {
	store := ctx.KVStore(k.storeKey)
	tokenType = nextID(tokenType)
	store.Set(types.NextTokenTypeNFTKey(k.getContractID(ctx)), k.mustEncodeString(tokenType))
}
func (k Keeper) setNextTokenIndexNFT(ctx sdk.Context, tokenType, tokenIndex string) {
	store := ctx.KVStore(k.storeKey)
	tokenIndex = nextID(tokenIndex)
	store.Set(types.NextTokenIDNFTKey(k.getContractID(ctx), tokenType), k.mustEncodeString(tokenIndex))
}

func (k Keeper) getNextTokenTypeFT(ctx sdk.Context) (tokenType string, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextTokenTypeFTKey(k.getContractID(ctx)))
	if bz == nil {
		panic("next token type for ft should be exist")
	}
	tokenType = k.mustDecodeString(bz)
	if tokenType[0] != types.FungibleFlag[0] {
		return "", sdkerrors.Wrapf(types.ErrTokenTypeFull, "contract id: %s", k.getContractID(ctx))
	}
	return tokenType, nil
}

func (k Keeper) getNextTokenTypeNFT(ctx sdk.Context) (tokenType string, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextTokenTypeNFTKey(k.getContractID(ctx)))
	if bz == nil {
		panic("next token type for nft should be exist")
	}
	tokenType = k.mustDecodeString(bz)
	if tokenType == types.ReservedEmpty {
		return "", sdkerrors.Wrapf(types.ErrTokenTypeFull, "contract id: %s", k.getContractID(ctx))
	}
	return tokenType, nil
}

func (k Keeper) getNextTokenIndexNFT(ctx sdk.Context, tokenType string) (tokenIndex string, error error) {
	if !k.HasTokenType(ctx, tokenType) {
		return "", sdkerrors.Wrapf(types.ErrTokenTypeNotExist, "ContractID: %s, TokenType: %s", k.getContractID(ctx), tokenType)
	}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextTokenIDNFTKey(k.getContractID(ctx), tokenType))
	if bz == nil {
		panic("next token id for nft token type should be exist")
	}
	tokenIndex = k.mustDecodeString(bz)
	if tokenIndex == types.ReservedEmpty {
		return "", sdkerrors.Wrapf(types.ErrTokenIndexFull, "ContractID: %s, TokenType: %s", k.getContractID(ctx), tokenType)
	}
	return tokenIndex, nil
}

func (k Keeper) GetNextTokenIDFT(ctx sdk.Context) (string, error) {
	if !k.ExistCollection(ctx) {
		return "", sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", k.getContractID(ctx))
	}
	tokenType, err := k.getNextTokenTypeFT(ctx)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrTokenIDFull, "ContractID: %s, TokenType: %s", k.getContractID(ctx), tokenType)
	}
	return tokenType + types.ReservedEmpty, nil
}
func (k Keeper) GetNextTokenIDNFT(ctx sdk.Context, tokenType string) (string, error) {
	if !k.ExistCollection(ctx) {
		return "", sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", k.getContractID(ctx))
	}
	tokenIndex, err := k.getNextTokenIndexNFT(ctx, tokenType)
	if err != nil {
		return "", err
	}

	if tokenIndex == types.ReservedEmpty {
		return "", sdkerrors.Wrapf(types.ErrTokenIndexFull, "ContractID: %s, TokenType: %s", k.getContractID(ctx), tokenType)
	}
	return tokenType + tokenIndex, nil
}

func (k Keeper) iterateToken(ctx sdk.Context, prefix string, reverse bool, process func(types.Token) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	var iter sdk.Iterator
	if reverse {
		iter = sdk.KVStoreReversePrefixIterator(store, types.TokenKey(k.getContractID(ctx), prefix))
	} else {
		iter = sdk.KVStorePrefixIterator(store, types.TokenKey(k.getContractID(ctx), prefix))
	}
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		token := k.mustDecodeToken(val)
		if process(token) {
			return
		}
		iter.Next()
	}
}

func (k Keeper) mustEncodeToken(token types.Token) (bz []byte) {
	return k.cdc.MustMarshalBinaryBare(token)
}
func (k Keeper) mustDecodeToken(bz []byte) (token types.Token) {
	k.cdc.MustUnmarshalBinaryBare(bz, &token)
	return token
}

func fromHex(s string) *big.Int {
	r, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("bad hex")
	}
	return r
}

func toHex(r *big.Int) string {
	return r.Text(16)
}

func nextID(id string) (nextTokenID string) {
	idInt := fromHex(id)
	idInt = idInt.Add(idInt, big.NewInt(1))
	nextTokenID = strings.Repeat("0", len(id)) + toHex(idInt)
	nextTokenID = nextTokenID[len(nextTokenID)-len(id):]
	return nextTokenID
}
