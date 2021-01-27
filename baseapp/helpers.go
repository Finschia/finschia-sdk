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

func (app *BaseApp) Simulate(txBytes []byte, tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	err := validateBasicTxMsgs(tx.GetMsgs())
	if err != nil {
		return sdk.GasInfo{}, nil, err
	}

	return app.runTx(txBytes, tx, true)
}

func (app *BaseApp) Deliver(tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	err := validateBasicTxMsgs(tx.GetMsgs())
	if err != nil {
		return sdk.GasInfo{}, nil, err
	}

	return app.runTx(nil, tx, false)
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
