package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/slashing/types"
)

// NewQuerier creates a new querier for slashing clients.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, k, legacyQuerierCdc)

		case types.QuerySigningInfo:
			return querySigningInfo(ctx, req, k, legacyQuerierCdc)

		case types.QuerySigningInfos:
			return querySigningInfos(ctx, req, k, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySigningInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySigningInfoRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// https://github.com/cosmos/cosmos-sdk/issues/12573
	// Will be removed, but fix this
	addr, err := sdk.ConsAddressFromBech32(params.ConsAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	signingInfo, found := k.GetValidatorSigningInfo(ctx, addr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoSigningInfoFound, params.ConsAddress)
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, signingInfo)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySigningInfos(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySigningInfosParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var signingInfos []types.ValidatorSigningInfo

	k.IterateValidatorSigningInfos(ctx, func(consAddr sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool) {
		signingInfos = append(signingInfos, info)
		return false
	})

	start, end := client.Paginate(len(signingInfos), params.Page, params.Limit, int(k.sk.MaxValidators(ctx)))
	if start < 0 || end < 0 {
		signingInfos = []types.ValidatorSigningInfo{}
	} else {
		signingInfos = signingInfos[start:end]
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, signingInfos)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
