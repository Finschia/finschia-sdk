package baseapp

import (
	"fmt"
	"io"

	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/store"
	storetypes "github.com/line/lbm-sdk/store/types"
	sdk "github.com/line/lbm-sdk/types"
)

// File for storing in-package BaseApp optional functions,
// for options that need access to non-exported fields of the BaseApp

// SetPruning sets a pruning option on the multistore associated with the app
func SetPruning(opts sdk.PruningOptions) func(*BaseApp) {
	return func(bap *BaseApp) { bap.cms.SetPruning(opts) }
}

// SetLazyLoading sets a lazyLoading option to load a version lazily
func SetLazyLoading(lazyLoading bool) func(*BaseApp) {
	return func(bap *BaseApp) { bap.cms.SetLazyLoading(lazyLoading) }
}

// SetMinGasPrices returns an option that sets the minimum gas prices on the app.
func SetMinGasPrices(gasPricesStr string) func(*BaseApp) {
	gasPrices, err := sdk.ParseDecCoins(gasPricesStr)
	if err != nil {
		panic(fmt.Sprintf("invalid minimum gas prices: %v", err))
	}

	return func(bap *BaseApp) { bap.setMinGasPrices(gasPrices) }
}

// SetHaltHeight returns a BaseApp option function that sets the halt block height.
func SetHaltHeight(blockHeight uint64) func(*BaseApp) {
	return func(bap *BaseApp) { bap.setHaltHeight(blockHeight) }
}

// SetHaltTime returns a BaseApp option function that sets the halt block time.
func SetHaltTime(haltTime uint64) func(*BaseApp) {
	return func(bap *BaseApp) { bap.setHaltTime(haltTime) }
}

// SetInterBlockCache provides a BaseApp option function that sets the
// inter-block cache.
func SetInterBlockCache(cache sdk.MultiStorePersistentCache) func(*BaseApp) {
	return func(app *BaseApp) { app.setInterBlockCache(cache) }
}

// SetTrace will turn on or off trace flag
func SetTrace(trace bool) func(*BaseApp) {
	return func(app *BaseApp) { app.setTrace(trace) }
}

// SetIAVLCacheSize sets the maximum number of entries in the all iavl node caches
func SetIAVLCacheSize(cacheSize int) func(*BaseApp) {
	return func(app *BaseApp) { app.cms.SetIAVLCacheSize(cacheSize) }
}

// SetMetrics sets prometheus metrics
func SetMetrics(bappMetrics *Metrics, storeMetrics *storetypes.Metrics, iavlMetricsProvider iavl.MetricsProvider) func(*BaseApp) {
	return func(app *BaseApp) {
		app.setMetrics(bappMetrics)
		app.cms.SetMetrics(storeMetrics, iavlMetricsProvider)
	}
}

func (app *BaseApp) SetName(name string) {
	if app.sealed {
		panic("SetName() on sealed BaseApp")
	}
	app.name = name
}

// SetAppVersion sets the application's version string.
func (app *BaseApp) SetAppVersion(v string) {
	if app.sealed {
		panic("SetAppVersion() on sealed BaseApp")
	}
	app.appVersion = v
}

func (app *BaseApp) SetDB(db dbm.DB) {
	if app.sealed {
		panic("SetDB() on sealed BaseApp")
	}
	app.db = db
}

func (app *BaseApp) SetCMS(cms store.CommitMultiStore) {
	if app.sealed {
		panic("SetEndBlocker() on sealed BaseApp")
	}
	app.cms = cms
}

func (app *BaseApp) SetInitChainer(initChainer sdk.InitChainer) {
	if app.sealed {
		panic("SetInitChainer() on sealed BaseApp")
	}
	app.initChainer = initChainer
}

func (app *BaseApp) SetBeginBlocker(beginBlocker sdk.BeginBlocker) {
	if app.sealed {
		panic("SetBeginBlocker() on sealed BaseApp")
	}
	app.beginBlocker = beginBlocker
}

func (app *BaseApp) SetEndBlocker(endBlocker sdk.EndBlocker) {
	if app.sealed {
		panic("SetEndBlocker() on sealed BaseApp")
	}
	app.endBlocker = endBlocker
}

func (app *BaseApp) SetAnteHandler(ah sdk.AnteHandler) {
	if app.sealed {
		panic("SetAnteHandler() on sealed BaseApp")
	}
	app.anteHandler = ah
}

func (app *BaseApp) SetAddrPeerFilter(pf sdk.PeerFilter) {
	if app.sealed {
		panic("SetAddrPeerFilter() on sealed BaseApp")
	}
	app.addrPeerFilter = pf
}

func (app *BaseApp) SetIDPeerFilter(pf sdk.PeerFilter) {
	if app.sealed {
		panic("SetIDPeerFilter() on sealed BaseApp")
	}
	app.idPeerFilter = pf
}

func (app *BaseApp) SetFauxMerkleMode() {
	if app.sealed {
		panic("SetFauxMerkleMode() on sealed BaseApp")
	}
	app.fauxMerkleMode = true
}

// SetCommitMultiStoreTracer sets the store tracer on the BaseApp's underlying
// CommitMultiStore.
func (app *BaseApp) SetCommitMultiStoreTracer(w io.Writer) {
	app.cms.SetTracer(w)
}

// SetStoreLoader allows us to customize the rootMultiStore initialization.
func (app *BaseApp) SetStoreLoader(loader StoreLoader) {
	if app.sealed {
		panic("SetStoreLoader() on sealed BaseApp")
	}
	app.storeLoader = loader
}

// SetRouter allows us to customize the router.
func (app *BaseApp) SetRouter(router sdk.Router) {
	if app.sealed {
		panic("SetRouter() on sealed BaseApp")
	}
	app.router = router
}
