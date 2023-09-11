package cmd

import (
	"github.com/Finschia/finschia-rdk/baseapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/ulbqb/statelesskit/querier"
)

func SetCustomStoreQueryRoute(db dbm.DB) func(*baseapp.BaseApp) {
	return func(bapp *baseapp.BaseApp) {
		customStoreQuerier := func(_ sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
			return querier.MakeCustomStoreQuerier(db)(path, req)
		}
		bapp.QueryRouter().AddRoute("store", customStoreQuerier)
	}
}
