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
		case types.QueryEnabled:
			return queryEnabled(ctx, req, keeper, legacyQuerierCdc)

		case types.QueryAllowedValidator:
			return queryAllowedValidator(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case types.QueryAllowedValidators:
			return queryAllowedValidators(ctx, path[1:], req, keeper, legacyQuerierCdc)

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

func queryAllowedValidator(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	addr := path[0]
	if err := sdk.ValidateValAddress(addr); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid validator address (%s)", err)
	}
	valAddr := sdk.ValAddress(addr)
		
	allowed := keeper.GetAllowedValidator(ctx, valAddr)
	res, err := legacyQuerierCdc.MarshalJSON(&allowed)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
func queryAllowedValidators(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAllowedValidatorsRequest
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addrs := keeper.GetAllowedValidators(ctx)
	if len(addrs) != 0 {
		start, end := client.Paginate(len(addrs), int(params.Pagination.Offset), int(params.Pagination.Limit), 100)
		if start < 0 || end < 0 {
			addrs = []sdk.ValAddress{}
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
