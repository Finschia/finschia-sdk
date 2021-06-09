package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/line/tm-db/v2/memdb"

	"github.com/line/lfb-sdk/store/iavl"
	"github.com/line/lfb-sdk/store/rootmulti"
	storetypes "github.com/line/lfb-sdk/store/types"
)

type MerkleTestSuite struct {
	suite.Suite

	store     *rootmulti.Store
	storeKey  *storetypes.KVStoreKey
	iavlStore *iavl.Store
}

func (suite *MerkleTestSuite) SetupTest() {
	db := memdb.NewDB()
	suite.store = rootmulti.NewStore(db)

	suite.storeKey = storetypes.NewKVStoreKey("iavlStoreKey")

	suite.store.MountStoreWithDB(suite.storeKey, storetypes.StoreTypeIAVL, nil)
	suite.store.LoadVersion(0)

	suite.iavlStore = suite.store.GetCommitStore(suite.storeKey).(*iavl.Store)
}

func TestMerkleTestSuite(t *testing.T) {
	suite.Run(t, new(MerkleTestSuite))
}
