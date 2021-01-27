package baseapp

import (
	"regexp"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func (app *BaseApp) Check(tx sdk.Tx) (sdk.GasInfo, error) {
	return app.checkTx(nil, tx, false)
}

func (app *BaseApp) Simulate(txBytes []byte, tx sdk.Tx) (gInfo sdk.GasInfo, result *sdk.Result, err error) {
	return app.validateAndRunTx(txBytes, tx, true)
}

func (app *BaseApp) Deliver(tx sdk.Tx) (gInfo sdk.GasInfo, result *sdk.Result, err error) {
	return app.validateAndRunTx(nil, tx, false)
}

func (app *BaseApp) validateAndRunTx(txBytes []byte, tx sdk.Tx, simulate bool) (gInfo sdk.GasInfo, result *sdk.Result, err error) {
	err = validateBasicTxMsgs(tx.GetMsgs())
	if err != nil {
		return gInfo, result, err
	}

	var msgsResult *sdk.MsgsResult
	gInfo, msgsResult, err = app.runTx(txBytes, tx, simulate)
	if msgsResult != nil {
		result = msgsResult.ToResult()
	}
	return gInfo, result, err
}

// Context with current {check, deliver}State of the app used by tests.
func (app *BaseApp) NewContext(isCheckTx bool, header abci.Header) sdk.Context {
	if isCheckTx {
		return sdk.NewContext(app.checkState.ms, header, true, app.logger).
			WithMinGasPrices(app.minGasPrices).
			WithConsensusParams(app.consensusParams)
	}

	return sdk.NewContext(app.deliverState.ms, header, false, app.logger)
}
