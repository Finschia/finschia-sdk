package baseapp

import (
	"sync"

	abci "github.com/line/ostracon/abci/types"

	sdk "github.com/line/lfb-sdk/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
)

func (app *BaseApp) startReactors() {
	go app.checkTxAsyncReactor()
}

type RequestCheckTxAsync struct {
	txBytes  []byte
	recheck  bool
	callback abci.CheckTxCallback
	prepare  *sync.WaitGroup
	tx       sdk.Tx
	err      error
}

func (app *BaseApp) checkTxAsyncReactor() {
	for req := range app.chCheckTx {
		req.prepare.Wait()
		if req.err != nil {
			req.callback(sdkerrors.ResponseCheckTx(req.err, 0, 0, app.trace))
			continue
		}

		waits, signals := app.checkAccountWGs.Register(req.tx)

		go app.checkTxAsync(req, waits, signals)
	}
}

func (app *BaseApp) prepareCheckTx(req *RequestCheckTxAsync) {
	defer req.prepare.Done()
	req.tx, req.err = app.preCheckTx(req.txBytes)
}

func (app *BaseApp) checkTxAsync(req *RequestCheckTxAsync, waits []*sync.WaitGroup, signals []*AccountWG) {
	app.checkAccountWGs.Wait(waits)
	defer app.checkAccountWGs.Done(signals)

	gInfo, err := app.checkTx(req.txBytes, req.tx, req.recheck)

	if err != nil {
		req.callback(sdkerrors.ResponseCheckTx(err, gInfo.GasWanted, gInfo.GasUsed, app.trace))
		return
	}

	req.callback(abci.ResponseCheckTx{
		GasWanted: int64(gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
	})
}
