package baseapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type UpgradeHandler func(height int64, ms sdk.CommitMultiStore) (err error)

func (app *BaseApp) SetUpgradeHandler(up UpgradeHandler) {
	if app.sealed {
		panic("SetUpgrader() on sealed BaseApp")
	}
	app.upgrader = up
}

func (app *BaseApp) upgrade(req abci.RequestBeginBlock) error {
	if app.upgrader != nil {
		if err := app.upgrader(req.Header.Height, app.cms); err != nil {
			return err
		}
	}
	return nil
}
