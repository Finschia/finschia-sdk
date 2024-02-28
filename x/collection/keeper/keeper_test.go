package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/keeper"
	"github.com/Finschia/finschia-sdk/x/collection/module"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx          sdk.Context
	keeper       keeper.Keeper
	queryServer  collection.QueryServer
	msgServer    collection.MsgServer
	addressCodec address.Codec

	vendor   sdk.AccAddress
	operator sdk.AccAddress
	customer sdk.AccAddress
	stranger sdk.AccAddress

	contractID string
	nftClassID string

	balance math.Int

	numNFTs    int
	issuedNFTs map[string][]collection.NFT
}

func (s *KeeperTestSuite) createRandomAccounts(accNum int) []sdk.AccAddress {
	seenAddresses := make(map[string]bool, accNum)
	addresses := make([]sdk.AccAddress, accNum)
	for i := range addresses {
		var addr sdk.AccAddress
		for {
			pk := secp256k1.GenPrivKey().PubKey()
			addr = sdk.AccAddress(pk.Address())
			if !seenAddresses[s.bytesToString(addr)] {
				seenAddresses[s.bytesToString(addr)] = true
				break
			}
		}
		addresses[i] = addr
	}
	return addresses
}

func (s *KeeperTestSuite) SetupTest() {
	s.prepareInitialSetup()

	s.balance = math.NewInt(1000000)
	s.contractID = s.keeper.CreateContract(s.ctx, s.vendor, collection.Contract{
		Name: "fox",
	})

	for _, permission := range []collection.Permission{
		collection.PermissionMint,
		collection.PermissionBurn,
	} {
		s.keeper.Grant(s.ctx, s.contractID, s.vendor, s.operator, permission)
	}

	// create a non-fungible token class
	nftClassID, err := s.keeper.CreateTokenClass(s.ctx, s.contractID, &collection.NFTClass{
		Name: "fennec fox",
	})
	s.Require().NoError(err)
	s.nftClassID = *nftClassID

	// mint nfts
	newParams := func(classID string, size int) []collection.MintNFTParam {
		res := make([]collection.MintNFTParam, size)
		for i := range res {
			res[i] = collection.MintNFTParam{
				TokenType: s.nftClassID,
			}
		}
		return res
	}
	s.numNFTs = 4
	s.issuedNFTs = make(map[string][]collection.NFT)
	for _, to := range []sdk.AccAddress{s.customer, s.operator, s.vendor} {
		nfts, err := s.keeper.MintNFT(s.ctx, s.contractID, to, newParams(s.nftClassID, s.numNFTs))
		s.issuedNFTs[s.bytesToString(to)] = nfts
		s.Require().NoError(err)
	}

	// authorize
	err = s.keeper.AuthorizeOperator(s.ctx, s.contractID, s.customer, s.operator)
	s.Require().NoError(err)
	err = s.keeper.AuthorizeOperator(s.ctx, s.contractID, s.customer, s.stranger)
	s.Require().NoError(err)

	// not token contract
	notTokenContractID := s.keeper.NewID(s.ctx)
	err = keeper.ValidateLegacyContract(s.keeper, s.ctx, notTokenContractID)
	s.Require().ErrorIs(err, collection.ErrCollectionNotExist)
}

func (s *KeeperTestSuite) prepareInitialSetup() {
	// Create Store for test
	key := storetypes.NewKVStoreKey(collection.StoreKey)
	tkey := storetypes.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContextWithDB(s.T(), key, tkey)
	kvStoreService := runtime.NewKVStoreService(key)
	s.ctx = testCtx.Ctx

	// Create EncodingConfig
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	encCfg.InterfaceRegistry = codectestutil.CodecOptions{
		AccAddressPrefix: "link",
		ValAddressPrefix: "linkvaloper",
	}.NewInterfaceRegistry()
	encCfg.Codec = codec.NewProtoCodec(encCfg.InterfaceRegistry)

	collection.RegisterInterfaces(encCfg.InterfaceRegistry)
	testdata.RegisterInterfaces(encCfg.InterfaceRegistry)

	// Create BaseApp
	bapp := baseapp.NewBaseApp(
		"collection",
		log.NewNopLogger(),
		testCtx.DB,
		encCfg.TxConfig.TxDecoder(),
	)
	bapp.SetInterfaceRegistry(encCfg.InterfaceRegistry)

	s.keeper = keeper.NewKeeper(encCfg.Codec, kvStoreService)
	s.addressCodec = encCfg.Codec.InterfaceRegistry().SigningContext().AddressCodec()
	s.keeper.InitGenesis(s.ctx, collection.DefaultGenesisState())
	s.keeper.SetParams(s.ctx, collection.Params{})

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	addresses := []*sdk.AccAddress{
		&s.vendor,
		&s.operator,
		&s.customer,
		&s.stranger,
	}
	for i, addr := range s.createRandomAccounts(len(addresses)) {
		*addresses[i] = addr
	}
}

func (s *KeeperTestSuite) bytesToString(address []byte) string {
	addr, err := s.addressCodec.BytesToString(address)
	s.Require().NoError(err)
	return addr
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
