package app

import (
	"io"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/baseapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/staking"
)

// nolint
var (
	genesisFile        string
	paramsFile         string
	exportParamsPath   string
	exportParamsHeight int
	exportStatePath    string
	exportStatsPath    string
	seed               int64
	initialBlockHeight int
	numBlocks          int
	blockSize          int
	enabled            bool
	verbose            bool
	lean               bool
	commit             bool
	period             int
	onOperation        bool // TODO Remove in favor of binary search for invariant violation
	allInvariants      bool
	genesisTime        int64
)

// DONTCOVER

// NewLinkAppUNSAFE is used for debugging purposes only.
//
// NOTE: to not use this function with non-test code
func NewLinkAppUNSAFE(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp),
) (gapp *LinkApp, keyMain, keyStaking *sdk.KVStoreKey, stakingKeeper staking.Keeper) {
	gapp = NewLinkApp(logger, db, traceStore, loadLatest, map[int64]bool{}, invCheckPeriod, baseAppOptions...)
	return gapp, gapp.keys[baseapp.MainStoreKey], gapp.keys[staking.StoreKey], gapp.stakingKeeper
}
