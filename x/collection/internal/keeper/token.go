package keeper

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type TokenKeeper interface {
	GetToken(ctx sdk.Context, contractID, tokenID string) (types.Token, sdk.Error)
	HasToken(ctx sdk.Context, contractID, tokenID string) bool
	SetToken(ctx sdk.Context, token types.Token) sdk.Error
	DeleteToken(ctx sdk.Context, contractID, tokenID string) sdk.Error
	UpdateToken(ctx sdk.Context, token types.Token) sdk.Error
	GetTokens(ctx sdk.Context, contractID string) (tokens types.Tokens, err sdk.Error)
	GetFT(ctx sdk.Context, contractID, tokenID string) (types.FT, sdk.Error)
	GetFTs(ctx sdk.Context, contractID string) (tokens types.Tokens, err sdk.Error)
	GetNFT(ctx sdk.Context, contractID, tokenID string) (types.NFT, sdk.Error)
	GetNFTCount(ctx sdk.Context, contractID, tokenType string) (sdk.Int, sdk.Error)
	GetNFTCountInt(ctx sdk.Context, contractID, tokenType, target string) (sdk.Int, sdk.Error)
	GetNFTs(ctx sdk.Context, contractID, tokenType string) (tokens types.Tokens, err sdk.Error)
	GetNextTokenIDFT(ctx sdk.Context, contractID string) (string, sdk.Error)
	GetNextTokenIDNFT(ctx sdk.Context, contractID, tokenType string) (string, sdk.Error)
}

var _ TokenKeeper = (*Keeper)(nil)

func (k Keeper) GetToken(ctx sdk.Context, contractID, tokenID string) (types.Token, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(contractID, tokenID)
	bz := store.Get(tokenKey)
	if bz == nil {
		return nil, types.ErrTokenNotExist(types.DefaultCodespace, contractID, tokenID)
	}
	token := k.mustDecodeToken(bz)
	return token, nil
}
func (k Keeper) HasToken(ctx sdk.Context, contractID, tokenID string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(contractID, tokenID)
	return store.Has(tokenKey)
}

func (k Keeper) SetToken(ctx sdk.Context, token types.Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(token.GetContractID(), token.GetTokenID())
	if store.Has(tokenKey) {
		return types.ErrTokenExist(types.DefaultCodespace, token.GetContractID(), token.GetTokenID())
	}
	store.Set(tokenKey, k.mustEncodeToken(token))
	tokenType := token.GetTokenType()
	if tokenType[0] == types.FungibleFlag[0] {
		k.setNextTokenTypeFT(ctx, token.GetContractID(), tokenType)
	} else {
		k.setNextTokenIndexNFT(ctx, token.GetContractID(), tokenType, token.GetTokenIndex())
	}
	return nil
}

func (k Keeper) UpdateToken(ctx sdk.Context, token types.Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(token.GetContractID(), token.GetTokenID())
	if !store.Has(tokenKey) {
		return types.ErrTokenNotExist(types.DefaultCodespace, token.GetContractID(), token.GetTokenID())
	}
	store.Set(tokenKey, k.mustEncodeToken(token))
	return nil
}

func (k Keeper) DeleteToken(ctx sdk.Context, contractID, tokenID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(contractID, tokenID)
	if !store.Has(tokenKey) {
		return types.ErrTokenNotExist(types.DefaultCodespace, contractID, tokenID)
	}
	store.Delete(tokenKey)
	return nil
}

func (k Keeper) GetTokens(ctx sdk.Context, contractID string) (tokens types.Tokens, err sdk.Error) {
	_, err = k.GetCollection(ctx, contractID)
	if err != nil {
		return nil, err
	}
	k.iterateToken(ctx, contractID, "", false, func(t types.Token) bool {
		tokens = append(tokens, t)
		return false
	})
	return tokens, nil
}

func (k Keeper) GetFTs(ctx sdk.Context, contractID string) (tokens types.Tokens, err sdk.Error) {
	_, err = k.GetCollection(ctx, contractID)
	if err != nil {
		return nil, err
	}
	k.iterateToken(ctx, contractID, types.FungibleFlag, false, func(t types.Token) bool {
		tokens = append(tokens, t)
		return false
	})
	return tokens, nil
}

func (k Keeper) GetFT(ctx sdk.Context, contractID, tokenID string) (types.FT, sdk.Error) {
	token, err := k.GetToken(ctx, contractID, tokenID)
	if err != nil {
		return nil, err
	}
	ft, ok := token.(types.FT)
	if !ok {
		return nil, types.ErrTokenNotNFT(types.DefaultCodespace, token.GetTokenID())
	}
	return ft, nil
}
func (k Keeper) GetNFT(ctx sdk.Context, contractID, tokenID string) (types.NFT, sdk.Error) {
	token, err := k.GetToken(ctx, contractID, tokenID)
	if err != nil {
		return nil, err
	}
	nft, ok := token.(types.NFT)
	if !ok {
		return nil, types.ErrTokenNotNFT(types.DefaultCodespace, token.GetTokenID())
	}
	return nft, nil
}

func (k Keeper) GetNFTs(ctx sdk.Context, contractID, tokenType string) (tokens types.Tokens, err sdk.Error) {
	_, err = k.GetCollection(ctx, contractID)
	if err != nil {
		return nil, err
	}
	k.iterateToken(ctx, contractID, tokenType, false, func(t types.Token) bool {
		tokens = append(tokens, t)
		return false
	})
	return tokens, nil
}
func (k Keeper) GetNFTCount(ctx sdk.Context, contractID, tokenType string) (sdk.Int, sdk.Error) {
	_, err := k.GetCollection(ctx, contractID)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	tokens, err := k.GetNFTs(ctx, contractID, tokenType)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return sdk.NewInt(int64(len(tokens))), nil
}

func (k Keeper) GetNFTCountInt(ctx sdk.Context, contractID, tokenType, target string) (sdk.Int, sdk.Error) {
	_, err := k.GetCollection(ctx, contractID)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	switch target {
	case types.QueryNFTCount:
		return k.getNFTCountTotal(ctx, contractID, tokenType), nil
	case types.QueryNFTMint:
		return k.getNFTCountMint(ctx, contractID, tokenType), nil
	case types.QueryNFTBurn:
		return k.getNFTCountBurn(ctx, contractID, tokenType), nil
	default:
		panic("invalid request target to query")
	}
}
func (k Keeper) getNFTCountTotal(ctx sdk.Context, contractID, tokenType string) sdk.Int {
	countTotal := 0
	k.iterateToken(ctx, contractID, tokenType, false, func(types.Token) bool {
		countTotal++
		return false
	})
	return sdk.NewInt(int64(countTotal))
}
func (k Keeper) getNFTCountMint(ctx sdk.Context, contractID, tokenType string) sdk.Int {
	var lastTokenIndex = "0"
	k.iterateToken(ctx, contractID, tokenType, true, func(token types.Token) bool {
		lastTokenIndex = token.GetTokenIndex()
		return true
	})
	return sdk.NewIntFromBigInt(fromHex(lastTokenIndex))
}

func (k Keeper) getNFTCountBurn(ctx sdk.Context, contractID, tokenType string) sdk.Int {
	return k.getNFTCountMint(ctx, contractID, tokenType).Sub(k.getNFTCountTotal(ctx, contractID, tokenType))
}

func (k Keeper) setNextTokenTypeFT(ctx sdk.Context, contractID, tokenType string) {
	store := ctx.KVStore(k.storeKey)
	tokenType = nextID(tokenType)
	store.Set(types.NextTokenTypeFTKey(contractID), k.mustEncodeString(tokenType))
}
func (k Keeper) setNextTokenTypeNFT(ctx sdk.Context, contractID, tokenType string) {
	store := ctx.KVStore(k.storeKey)
	tokenType = nextID(tokenType)
	store.Set(types.NextTokenTypeNFTKey(contractID), k.mustEncodeString(tokenType))
}
func (k Keeper) setNextTokenIndexNFT(ctx sdk.Context, contractID, tokenType, tokenIndex string) {
	store := ctx.KVStore(k.storeKey)
	tokenIndex = nextID(tokenIndex)
	store.Set(types.NextTokenIDNFTKey(contractID, tokenType), k.mustEncodeString(tokenIndex))
}

func (k Keeper) getNextTokenTypeFT(ctx sdk.Context, contractID string) (tokenType string, error sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextTokenTypeFTKey(contractID))
	if bz == nil {
		panic("next token type for ft should be exist")
	}
	tokenType = k.mustDecodeString(bz)
	if tokenType[0] != types.FungibleFlag[0] {
		return "", types.ErrTokenTypeFull(types.DefaultCodespace, contractID)
	}
	return tokenType, nil
}

func (k Keeper) getNextTokenTypeNFT(ctx sdk.Context, contractID string) (tokenType string, error sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextTokenTypeNFTKey(contractID))
	if bz == nil {
		panic("next token type for nft should be exist")
	}
	tokenType = k.mustDecodeString(bz)
	if tokenType == types.ReservedEmpty {
		return "", types.ErrTokenTypeFull(types.DefaultCodespace, contractID)
	}
	return tokenType, nil
}

func (k Keeper) getNextTokenIndexNFT(ctx sdk.Context, contractID, tokenType string) (tokenIndex string, error sdk.Error) {
	if !k.HasTokenType(ctx, contractID, tokenType) {
		return "", types.ErrTokenTypeNotExist(types.DefaultCodespace, contractID, tokenType)
	}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextTokenIDNFTKey(contractID, tokenType))
	if bz == nil {
		panic("next token id for nft token type should be exist")
	}
	tokenIndex = k.mustDecodeString(bz)
	if tokenIndex == types.ReservedEmpty {
		return "", types.ErrTokenIndexFull(types.DefaultCodespace, contractID, tokenType)
	}
	return tokenIndex, nil
}

func (k Keeper) GetNextTokenIDFT(ctx sdk.Context, contractID string) (string, sdk.Error) {
	if !k.ExistCollection(ctx, contractID) {
		return "", types.ErrCollectionNotExist(types.DefaultCodespace, contractID)
	}
	tokenType, err := k.getNextTokenTypeFT(ctx, contractID)
	if err != nil {
		return "", types.ErrTokenIDFull(types.DefaultCodespace, contractID)
	}
	return tokenType + types.ReservedEmpty, nil
}
func (k Keeper) GetNextTokenIDNFT(ctx sdk.Context, contractID, tokenType string) (string, sdk.Error) {
	if !k.HasTokenType(ctx, contractID, tokenType) {
		return "", types.ErrTokenTypeNotExist(types.DefaultCodespace, contractID, tokenType)
	}
	tokenIndex, err := k.getNextTokenIndexNFT(ctx, contractID, tokenType)
	if err != nil {
		return "", err
	}

	if tokenIndex == types.ReservedEmpty {
		return "", types.ErrTokenIndexFull(types.DefaultCodespace, contractID, tokenType)
	}
	return tokenType + tokenIndex, nil
}

func (k Keeper) iterateToken(ctx sdk.Context, contractID, prefix string, reverse bool, process func(types.Token) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	var iter sdk.Iterator
	if reverse {
		iter = sdk.KVStoreReversePrefixIterator(store, types.TokenKey(contractID, prefix))
	} else {
		iter = sdk.KVStorePrefixIterator(store, types.TokenKey(contractID, prefix))
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
