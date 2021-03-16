package baseapp

import (
	ostproto "github.com/line/ostracon/proto/ostracon/types"

	sdk "github.com/line/lbm-sdk/v2/types"
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
)

func (app *BaseApp) Check(txEncoder sdk.TxEncoder, tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	// runTx expects tx bytes as argument, so we encode the tx argument into
	// bytes. Note that runTx will actually decode those bytes again. But since
	// this helper is only used in tests/simulation, it's fine.
	bz, err := txEncoder(tx)
	if err != nil {
		return sdk.GasInfo{}, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "%s", err)
	}
	return app.runTx(runTxModeCheck, bz)
}

func (app *BaseApp) Simulate(txBytes []byte) (sdk.GasInfo, *sdk.Result, error) {
	return app.runTx(runTxModeSimulate, txBytes)
}

func (app *BaseApp) Deliver(txEncoder sdk.TxEncoder, tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	// See comment for Check().
	bz, err := txEncoder(tx)
	if err != nil {
		return sdk.GasInfo{}, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "%s", err)
	}
	return app.runTx(runTxModeDeliver, bz)
}

// Context with current {check, deliver}State of the app used by tests.
func (app *BaseApp) NewContext(isCheckTx bool, header ostproto.Header) sdk.Context {
	if isCheckTx {
		return sdk.NewContext(app.checkState.ms, header, true, app.logger).
			WithMinGasPrices(app.minGasPrices)
	}

	return sdk.NewContext(app.deliverState.ms, header, false, app.logger)
}

func (app *BaseApp) NewUncachedContext(isCheckTx bool, header ostproto.Header) sdk.Context {
	return sdk.NewContext(app.cms, header, isCheckTx, app.logger)
}
