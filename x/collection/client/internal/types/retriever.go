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

func (r Retriever) GetAccountPermission(ctx context.CLIContext, addr sdk.AccAddress) (types.Permissions, int64, error) {
	var pms types.Permissions
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryAccAddressParams(addr))
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

func (r Retriever) GetCollection(ctx context.CLIContext, symbol string) (types.BaseCollection, int64, error) {
	var collection types.BaseCollection
	bs, err := ctx.Codec.MarshalJSON(types.NewQuerySymbolParams(symbol))
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

func (r Retriever) GetCollectionNFTCount(ctx context.CLIContext, symbol, tokenID string) (sdk.Int, int64, error) {
	var nftcount sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
	if err != nil {
		return nftcount, 0, err
	}

	res, height, err := r.query(types.QueryNFTCount, bs)
	if err != nil {
		return nftcount, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &nftcount); err != nil {
		return nftcount, height, err
	}

	return nftcount, height, nil
}

func (r Retriever) GetSupply(ctx context.CLIContext, symbol, tokenID string) (sdk.Int, int64, error) {
	var supply sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
	if err != nil {
		return supply, 0, err
	}

	res, height, err := r.query(types.QuerySupply, bs)
	if err != nil {
		return supply, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &supply); err != nil {
		return supply, height, err
	}

	return supply, height, nil
}

func (r Retriever) GetToken(ctx context.CLIContext, symbol, tokenID string) (types.Token, int64, error) {
	var token types.Token
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
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

func (r Retriever) IsApproved(ctx context.CLIContext, proxy sdk.AccAddress, approver sdk.AccAddress, symbol string) (approved bool, height int64, err error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryIsApprovedParams(proxy, approver, symbol))
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

func (r Retriever) EnsureExists(ctx context.CLIContext, symbol, tokenID string) error {
	if _, _, err := r.GetToken(ctx, symbol, tokenID); err != nil {
		return err
	}
	return nil
}

func (r Retriever) GetParent(ctx context.CLIContext, symbol, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
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

func (r Retriever) GetRoot(ctx context.CLIContext, symbol, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
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

func (r Retriever) GetChildren(ctx context.CLIContext, symbol, tokenID string) (types.Tokens, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
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
