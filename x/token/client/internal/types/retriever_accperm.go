package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	context "github.com/line/link/client"
	"github.com/line/link/x/token/internal/types"
)

type AccountPermissionRetriever struct {
	querier types.NodeQuerier
}

func NewAccountPermissionRetriever(querier types.NodeQuerier) AccountPermissionRetriever {
	return AccountPermissionRetriever{querier: querier}
}

func (ar AccountPermissionRetriever) GetAccountPermission(ctx context.CLIContext, addr sdk.AccAddress) (types.Permissions, error) {
	pms, _, err := ar.GetAccountPermissionWithHeight(ctx, addr)
	return pms, err
}

func (ar AccountPermissionRetriever) GetAccountPermissionWithHeight(ctx context.CLIContext, addr sdk.AccAddress) (types.Permissions, int64, error) {
	var pms types.Permissions
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryAccountPermissionParams(addr))
	if err != nil {
		return pms, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPerms), bs)
	if err != nil {
		return pms, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &pms); err != nil {
		return pms, height, err
	}

	return pms, height, nil
}
