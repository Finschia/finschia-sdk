package types

import (
	"fmt"
	context "github.com/line/link/client"
	"github.com/line/link/x/token/internal/types"
)

type TokenRetriever struct {
	querier types.NodeQuerier
}

func NewTokenRetriever(querier types.NodeQuerier) TokenRetriever {
	return TokenRetriever{querier: querier}
}

func (ar TokenRetriever) GetToken(ctx context.CLIContext, symbol, tokenID string) (types.Token, error) {
	token, _, err := ar.GetTokenWithHeight(ctx, symbol, tokenID)
	return token, err
}

func (ar TokenRetriever) GetAllTokens(ctx context.CLIContext) (types.Tokens, error) {
	tokens, _, err := ar.GetAllTokensWithHeight(ctx)
	return tokens, err
}

func (ar TokenRetriever) GetTokenWithHeight(ctx context.CLIContext, symbol, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTokens), bs)
	if err != nil {
		return nil, height, err
	}

	var token types.Token
	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return nil, height, err
	}

	return token, height, nil
}

func (ar TokenRetriever) GetAllTokensWithHeight(ctx context.CLIContext) (types.Tokens, int64, error) {
	var tokens types.Tokens

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTokens), nil)
	if err != nil {
		return tokens, 0, err
	}

	err = ctx.Codec.UnmarshalJSON(res, &tokens)
	if err != nil {
		return tokens, 0, err
	}

	return tokens, height, nil
}

func (ar TokenRetriever) EnsureExists(ctx context.CLIContext, symbol, tokenID string) error {
	if _, err := ar.GetToken(ctx, symbol, tokenID); err != nil {
		return err
	}
	return nil
}

func (ar TokenRetriever) GetParent(ctx context.CLIContext, symbol, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParent), bs)
	if res == nil {
		return nil, 0, err
	}

	var token types.Token
	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return nil, 0, err
	}

	return token, height, nil
}

func (ar TokenRetriever) GetRoot(ctx context.CLIContext, symbol, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRoot), bs)
	if res == nil {
		return nil, 0, err
	}

	var token types.Token
	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return nil, 0, err
	}

	return token, height, nil
}

func (ar TokenRetriever) GetChildren(ctx context.CLIContext, symbol, tokenID string) (types.Tokens, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySymbolTokenIDParams(symbol, tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryChildren), bs)
	if res == nil {
		return nil, 0, err
	}

	var tokens types.Tokens
	if err := ctx.Codec.UnmarshalJSON(res, &tokens); err != nil {
		return nil, 0, err
	}

	return tokens, height, nil
}
