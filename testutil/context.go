package testutil

import (
	"github.com/line/ostracon/libs/log"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/line/tm-db/v2/memdb"

	"github.com/line/lbm-sdk/store"
	sdk "github.com/line/lbm-sdk/types"
)

// DefaultContext creates a sdk.Context with a fresh MemDB that can be used in tests.
func DefaultContext(key sdk.StoreKey) sdk.Context {
	db := memdb.NewDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	// cms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdk.NewContext(cms, ocproto.Header{}, false, log.NewNopLogger())

	return ctx
}
