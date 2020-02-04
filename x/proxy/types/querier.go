package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QuerierRoute   = "proxy"
	QueryAllowance = "allowance"
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

type QueryProxyAllowance struct {
	ProxyDenom
}

func (qpa QueryProxyAllowance) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, qpa))
}

func NewQueryProxyAllowance(proxy, onBehalfOf sdk.AccAddress, denom string) QueryProxyAllowance {
	return QueryProxyAllowance{NewProxyDenom(proxy, onBehalfOf, denom)}
}

type ProxyAllowanceRetriever struct {
	querier NodeQuerier
}

func NewProxyAllowanceRetriever(querier NodeQuerier) ProxyAllowanceRetriever {
	return ProxyAllowanceRetriever{querier: querier}
}

func (par ProxyAllowanceRetriever) GetProxyAllowance(proxy, onBehalfOf sdk.AccAddress, denom string) (ProxyAllowance, int64, error) {
	allowance, height, err := par.GetProxyAllowanceWithHeight(proxy, onBehalfOf, denom)
	return allowance, height, err
}

func (par ProxyAllowanceRetriever) GetProxyAllowanceWithHeight(proxy, onBehalfOf sdk.AccAddress, denom string) (ProxyAllowance, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryProxyAllowance(proxy, onBehalfOf, denom))
	if err != nil {
		return ProxyAllowance{}, 0, err
	}

	res, height, err := par.querier.QueryWithData(
		fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryAllowance),
		bs,
	)
	if err != nil {
		return ProxyAllowance{}, height, err
	}

	var pxar ProxyAllowance
	if err := ModuleCdc.UnmarshalJSON(res, &pxar); err != nil {
		return ProxyAllowance{}, height, err
	}

	return pxar, height, nil
}
