package safetybox

import (
	"fmt"

	"github.com/line/link/x/safetybox/internal/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QuerySafetyBox:
			return querySafetyBox(ctx, req, keeper)
		case types.QueryAccountRole:
			return queryAccountRole(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown safetybox query endpoint")
		}
	}
}

func querySafetyBox(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySafetyBoxParams
	var err error
	if err = types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	sb, err := keeper.GetSafetyBox(ctx, params.SafetyBoxID)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Can not get the safety box", err.Error()))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, sb)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAccountRole(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryAccountRoleParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	pms, err := keeper.GetPermissions(ctx, params.SafetyBoxID, params.Role, params.Address)
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(types.ModuleCdc, pms)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}
