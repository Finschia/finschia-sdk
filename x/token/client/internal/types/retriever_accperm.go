package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/link-chain/link/x/token/internal/types"
)

type AccountPermissionRetriever struct {
	querier types.NodeQuerier
}

func NewAccountPermissionRetriever(querier types.NodeQuerier) AccountPermissionRetriever {
	return AccountPermissionRetriever{querier: querier}
}

func (ar AccountPermissionRetriever) GetAccountPermission(addr sdk.AccAddress) (types.Permissions, error) {
	pms, _, err := ar.GetAccountPermissionWithHeight(addr)
	return pms, err
}

func (ar AccountPermissionRetriever) GetAccountPermissionWithHeight(addr sdk.AccAddress) (types.Permissions, int64, error) {
	var pms types.Permissions
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryAccountPermissionParams(addr))
	if err != nil {
		return pms, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPerms), bs)
	if err != nil {
		return pms, height, err
	}

	if err := types.ModuleCdc.UnmarshalJSON(res, &pms); err != nil {
		return pms, height, err
	}

	return pms, height, nil
}
