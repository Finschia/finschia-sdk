package keeper

import (
	abci "github.com/line/ostracon/abci/types"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/consortium/types"
)

// NewQuerier creates a new consortium Querier instance
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryEnabled(ctx, req, keeper, legacyQuerierCdc)

		case types.QueryValidatorAuth:
			return queryValidatorAuth(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case types.QueryValidatorAuths:
			return queryValidatorAuths(ctx, path[1:], req, keeper, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryEnabled(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	enabled := keeper.GetEnabled(ctx)

	res, err := legacyQuerierCdc.MarshalJSON(&enabled)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryValidatorAuth(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	addr := path[0]
	if err := sdk.ValidateValAddress(addr); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid validator address (%s)", err)
	}
	valAddr := sdk.ValAddress(addr)
		
	auth, err := keeper.GetValidatorAuth(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	res, err := legacyQuerierCdc.MarshalJSON(&auth)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
func queryValidatorAuths(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryValidatorAuthsRequest
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addrs := keeper.GetValidatorAuths(ctx)
	if len(addrs) != 0 {
		start, end := client.Paginate(len(addrs), int(params.Pagination.Offset), int(params.Pagination.Limit), 100)
		if start < 0 || end < 0 {
			addrs = []*types.ValidatorAuth{}
		} else {
			addrs = addrs[start:end]
		}
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, addrs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
