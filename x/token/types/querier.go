package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QuerierRoute   = "token"
	QueryToken     = "tokens"
	QueryAllTokens = "all_tokens"
	QueryPerm      = "perms"
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

type QueryTokenParams struct {
	Symbol string `json:"symbol"`
}

func (r QueryTokenParams) String() string {
	return r.Symbol
}

func NewQueryTokenParams(symbol string) QueryTokenParams {
	return QueryTokenParams{Symbol: symbol}
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

func (ar TokenRetriever) GetAllTokens() (Tokens, error) {
	tokens, _, err := ar.GetAllTokensWithHeight()
	return tokens, err
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

func (ar TokenRetriever) GetAllTokensWithHeight() (Tokens, int64, error) {
	var tokens Tokens

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryAllTokens), nil)
	if err != nil {
		return tokens, 0, err
	}

	err = ModuleCdc.UnmarshalJSON(res, &tokens)
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

type QueryAccountPermissionParams struct {
	Addr sdk.AccAddress `json:"symbol"`
}

func NewQueryAccountPermissionParams(addr sdk.AccAddress) QueryAccountPermissionParams {
	return QueryAccountPermissionParams{Addr: addr}
}

type AccountPermissionRetriever struct {
	querier NodeQuerier
}

func NewAccountPermissionRetriever(querier NodeQuerier) AccountPermissionRetriever {
	return AccountPermissionRetriever{querier: querier}
}

func (ar AccountPermissionRetriever) GetAccountPermission(addr sdk.AccAddress) (Permissions, error) {
	pms, _, err := ar.GetAccountPermissionWithHeight(addr)
	return pms, err
}

func (ar AccountPermissionRetriever) GetAccountPermissionWithHeight(addr sdk.AccAddress) (Permissions, int64, error) {
	var pms Permissions
	bs, err := ModuleCdc.MarshalJSON(NewQueryAccountPermissionParams(addr))
	if err != nil {
		return pms, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryPerm), bs)
	if err != nil {
		return pms, height, err
	}

	if err := ModuleCdc.UnmarshalJSON(res, &pms); err != nil {
		return pms, height, err
	}

	return pms, height, nil
}
