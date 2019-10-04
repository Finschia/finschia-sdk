package types

import (
	"fmt"
)

const (
	QuerierRoute = "token"
	QueryToken   = "tokens"
)

type QueryTokenParams struct {
	Symbol string `json:"symbol"`
}

func (r QueryTokenParams) String() string {
	return r.Symbol
}

func NewQueryTokenParams(symbol string) QueryTokenParams {
	return QueryTokenParams{Symbol: symbol}
}

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

type TokenRetriever struct {
	querier NodeQuerier
}

func NewTokenRetriever(querier NodeQuerier) TokenRetriever {
	return TokenRetriever{querier: querier}
}

func (ar TokenRetriever) GetToken(symbol string) (Token, error) {
	token, _, err := ar.GetTokenWithHeight(symbol)
	return token, err
}

func (ar TokenRetriever) GetTokenWithHeight(symbol string) (Token, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryTokenParams(symbol))
	if err != nil {
		return Token{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryToken), bs)
	if err != nil {
		return Token{}, height, err
	}

	var token Token
	if err := ModuleCdc.UnmarshalJSON(res, &token); err != nil {
		return Token{}, height, err
	}

	return token, height, nil
}

func (ar TokenRetriever) EnsureExists(symbol string) error {
	if _, err := ar.GetToken(symbol); err != nil {
		return err
	}
	return nil
}
