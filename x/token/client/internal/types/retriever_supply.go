package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	context "github.com/line/link/client"
	"github.com/line/link/x/token/internal/types"
)

type SupplyRetriever struct {
	querier types.NodeQuerier
}

func NewSupplyRetriever(querier types.NodeQuerier) SupplyRetriever {
	return SupplyRetriever{querier: querier}
}

func (ar SupplyRetriever) GetSupply(ctx context.CLIContext, symbol string) (sdk.Int, error) {
	supply, _, err := ar.GetSupplyWithHeight(ctx, symbol)
	return supply, err
}

func (ar SupplyRetriever) GetSupplyWithHeight(ctx context.CLIContext, symbol string) (sdk.Int, int64, error) {
	bs, err := ctx.Codec.MarshalJSON(types.NewQuerySupplyParams(symbol))
	if err != nil {
		return sdk.NewInt(0), 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySupply), bs)
	if err != nil {
		return sdk.NewInt(0), height, err
	}

	var supply sdk.Int
	if err := ctx.Codec.UnmarshalJSON(res, &supply); err != nil {
		return sdk.NewInt(0), height, err
	}

	return supply, height, nil
}

func (ar SupplyRetriever) EnsureExists(ctx context.CLIContext, symbol string) error {
	if _, err := ar.GetSupply(ctx, symbol); err != nil {
		return err
	}
	return nil
}
