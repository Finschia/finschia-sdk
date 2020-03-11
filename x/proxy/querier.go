package proxy

import (
	"github.com/line/link/x/proxy/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// creates a querier for REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryAllowance:
			return queryProxyAllowance(ctx, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown proxy query endpoint")
		}
	}
}

func queryProxyAllowance(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryProxyAllowance
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	allowance, err := keeper.GetProxyAllowance(ctx, types.NewProxyDenom(params.Proxy, params.OnBehalfOf, params.Denom))
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(types.ModuleCdc, allowance)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}
