package types

import (
	"fmt"
	"github.com/link-chain/link/x/token/internal/types"
)

type TokenRetriever struct {
	querier types.NodeQuerier
}

func NewTokenRetriever(querier types.NodeQuerier) TokenRetriever {
	return TokenRetriever{querier: querier}
}

func (ar TokenRetriever) GetToken(symbol string) (types.Token, error) {
	token, _, err := ar.GetTokenWithHeight(symbol)
	return token, err
}

func (ar TokenRetriever) GetAllTokens() (types.Tokens, error) {
	tokens, _, err := ar.GetAllTokensWithHeight()
	return tokens, err
}

func (ar TokenRetriever) GetTokenWithHeight(symbol string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryTokenParams(symbol))
	if err != nil {
		return types.Token{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTokens), bs)
	if err != nil {
		return types.Token{}, height, err
	}

	var token types.Token
	if err := types.ModuleCdc.UnmarshalJSON(res, &token); err != nil {
		return types.Token{}, height, err
	}

	return token, height, nil
}

func (ar TokenRetriever) GetAllTokensWithHeight() (types.Tokens, int64, error) {
	var tokens types.Tokens

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTokens), nil)
	if err != nil {
		return tokens, 0, err
	}

	err = types.ModuleCdc.UnmarshalJSON(res, &tokens)
	if err != nil {
		return tokens, 0, err
	}

	return tokens, height, nil
}

func (ar TokenRetriever) EnsureExists(symbol string) error {
	if _, err := ar.GetToken(symbol); err != nil {
		return err
	}
	return nil
}
