package cmd

import (
	"github.com/Finschia/finschia-rdk/baseapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/ulbqb/statelesskit/querier"
)

func SetRoute(db dbm.DB) func(*baseapp.BaseApp) {
	return func(bapp *baseapp.BaseApp) {
		preimageQuerier := func(_ sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
			return querier.MakePreimageQuerier(db)(path, req)
		}
		bapp.QueryRouter().AddRoute("preimage", preimageQuerier)
	}
}
