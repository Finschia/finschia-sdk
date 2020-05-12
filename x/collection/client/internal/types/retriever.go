package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	context "github.com/line/link/client"
	"github.com/line/link/x/collection/internal/types"
)

type Retriever struct {
	querier types.NodeQuerier
}

func NewRetriever(querier types.NodeQuerier) Retriever {
	return Retriever{querier: querier}
}

func (r Retriever) query(path string, data []byte) ([]byte, int64, error) {
	return r.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, path), data)
}

func (r Retriever) GetAccountBalance(ctx context.CLIContext, contractID, tokenID string, addr sdk.AccAddress) (sdk.Int, int64, error) {
	var balance sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryContractIDTokenIDAccAddressParams(contractID, tokenID, addr))
	if err != nil {
		return balance, 0, err
	}

	res, height, err := r.query(types.QueryBalance, bs)
	if err != nil {
		return balance, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &balance); err != nil {
		return balance, height, err
	}

	return balance, height, nil
}

func (r Retriever) GetAccountPermission(ctx context.CLIContext, contractID string, addr sdk.AccAddress) (types.Permissions, int64, error) {
	var pms types.Permissions
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryContractIDAccAddressParams(contractID, addr))
	if err != nil {
		return pms, 0, err
	}

	res, height, err := r.query(types.QueryPerms, bs)
	if err != nil {
		return pms, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &pms); err != nil {
		return pms, height, err
	}

	return pms, height, nil
}

func (r Retriever) GetCollection(ctx context.CLIContext, contractID string) (types.BaseCollection, int64, error) {
	var collection types.BaseCollection
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryContractIDParams(contractID))
	if err != nil {
		return collection, 0, err
	}

	res, height, err := r.query(types.QueryCollections, bs)
	if err != nil {
		return collection, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &collection); err != nil {
		return collection, height, err
	}

	return collection, height, nil
}

func (r Retriever) GetCollections(ctx context.CLIContext) (types.Collections, int64, error) {
	var collections types.Collections

	res, height, err := r.query(types.QueryCollections, nil)
	if err != nil {
		return collections, 0, err
	}

	err = ctx.Codec.UnmarshalJSON(res, &collections)
	if err != nil {
		return collections, 0, err
	}

	return collections, height, nil
}

func (r Retriever) GetCollectionNFTCount(ctx context.CLIContext, contractID, tokenID, target string) (sdk.Int, int64, error) {
	var nftcount sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryContractIDTokenIDParams(contractID, tokenID))
	if err != nil {
		return nftcount, 0, err
	}
	if target != types.QueryNFTCount && target != types.QueryNFTMint && target != types.QueryNFTBurn {
		return nftcount, 0, fmt.Errorf("invalid target : %s", target)
	}

	res, height, err := r.query(target, bs)
	if err != nil {
		return nftcount, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &nftcount); err != nil {
		return nftcount, height, err
	}

	return nftcount, height, nil
}

func (r Retriever) GetTotal(ctx context.CLIContext, contractID, tokenID, target string) (sdk.Int, int64, error) {
	var supply sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryContractIDTokenIDParams(contractID, tokenID))
	if err != nil {
		return supply, 0, err
	}

	res, height, err := r.query(target, bs)
	if err != nil {
		return supply, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &supply); err != nil {
		return supply, height, err
	}

	return supply, height, nil
}

func (r Retriever) GetToken(ctx context.CLIContext, contractID, tokenID string) (types.Token, int64, error) {
	var token types.Token
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryContractIDTokenIDParams(contractID, tokenID))
	if err != nil {
		return token, 0, err
	}

	res, height, err := r.query(types.QueryTokens, bs)
	if err != nil {
		return token, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return token, height, err
	}
	return token, height, nil
}

func (r Retriever) GetTokens(ctx context.CLIContext, contractID string) (types.Tokens, int64, error) {
	var tokens types.Tokens
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryContractIDParams(contractID))
	if err != nil {
		return tokens, 0, err
	}

	res, height, err := r.query(types.QueryTokens, bs)
	if err != nil {
		return tokens, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &tokens); err != nil {
		return tokens, height, err
	}
	return tokens, height, nil
}

func (r Retriever) GetTokenType(ctx context.CLIContext, contractID, tokenTypeID string) (types.TokenType, int64, error) {
	var tokenType types.TokenType
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryContractIDTokenIDParams(contractID, tokenTypeID))
	if err != nil {
		return tokenType, 0, err
	}

	res, height, err := r.query(types.QueryTokenTypes, bs)
	if err != nil {
		return tokenType, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &tokenType); err != nil {
		return tokenType, height, err
	}
	return tokenType, height, nil
}

func (r Retriever) GetTokenTypes(ctx context.CLIContext, contractID string) (types.TokenTypes, int64, error) {
	var tokenTypes types.TokenTypes
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryContractIDParams(contractID))
	if err != nil {
		return tokenTypes, 0, err
	}

	res, height, err := r.query(types.QueryTokenTypes, bs)
	if err != nil {
		return tokenTypes, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &tokenTypes); err != nil {
		return tokenTypes, height, err
	}
	return tokenTypes, height, nil
}

func (r Retriever) IsApproved(ctx context.CLIContext, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) (approved bool, height int64, err error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryIsApprovedParams(contractID, proxy, approver))
	if err != nil {
		return false, 0, err
	}

	res, height, err := r.query(types.QueryIsApproved, bs)
	if err != nil {
		return false, 0, err
	}

	err = ctx.Codec.UnmarshalJSON(res, &approved)
	if err != nil {
		return false, 0, err
	}

	return approved, height, nil
}

func (r Retriever) EnsureExists(ctx context.CLIContext, contractID, tokenID string) error {
	if _, _, err := r.GetToken(ctx, contractID, tokenID); err != nil {
		return err
	}
	return nil
}

func (r Retriever) GetParent(ctx context.CLIContext, contractID, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryContractIDTokenIDParams(contractID, tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := r.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParent), bs)
	if res == nil {
		return nil, 0, err
	}

	var token types.Token
	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return nil, 0, err
	}

	return token, height, nil
}

func (r Retriever) GetRoot(ctx context.CLIContext, contractID, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryContractIDTokenIDParams(contractID, tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := r.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRoot), bs)
	if res == nil {
		return nil, 0, err
	}

	var token types.Token
	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return nil, 0, err
	}

	return token, height, nil
}

func (r Retriever) GetChildren(ctx context.CLIContext, contractID, tokenID string) (types.Tokens, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryContractIDTokenIDParams(contractID, tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := r.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryChildren), bs)
	if res == nil {
		return nil, 0, err
	}

	var tokens types.Tokens
	if err := ctx.Codec.UnmarshalJSON(res, &tokens); err != nil {
		return nil, 0, err
	}

	return tokens, height, nil
}
