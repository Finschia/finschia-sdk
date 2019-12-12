package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

type SupplyRetriever struct {
	querier types.NodeQuerier
}

func NewSupplyRetriever(querier types.NodeQuerier) SupplyRetriever {
	return SupplyRetriever{querier: querier}
}

func (ar SupplyRetriever) GetSupply(symbol string) (sdk.Int, error) {
	supply, _, err := ar.GetSupplyWithHeight(symbol)
	return supply, err
}

func (ar SupplyRetriever) GetSupplyWithHeight(symbol string) (sdk.Int, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQuerySupplyParams(symbol))
	if err != nil {
		return sdk.NewInt(0), 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySupply), bs)
	if err != nil {
		return sdk.NewInt(0), height, err
	}

	var supply sdk.Int
	if err := types.ModuleCdc.UnmarshalJSON(res, &supply); err != nil {
		return sdk.NewInt(0), height, err
	}

	return supply, height, nil
}

func (ar SupplyRetriever) EnsureExists(symbol string) error {
	if _, err := ar.GetSupply(symbol); err != nil {
		return err
	}
	return nil
}
