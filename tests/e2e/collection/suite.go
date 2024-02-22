package collection

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/gogoproto/proto"

	cmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/client/cli"
)

type E2ETestSuite struct {
	suite.Suite

	cfg         network.Config
	network     *network.Network
	setupHeight int64

	commonArgs []string

	vendor     sdk.AccAddress
	operator   sdk.AccAddress
	customer   sdk.AccAddress
	stranger   sdk.AccAddress
	contractID string
	nftClassID string
	tokenIDs   map[string]string
}

func NewE2ETestSuite(cfg network.Config) *E2ETestSuite {
	return &E2ETestSuite{cfg: cfg}
}

func (s *E2ETestSuite) SetupSuite() {
	s.T().Log("setting up collection e2e test suite")

	genesisState := s.cfg.GenesisState
	collectionGenesisState := new(collection.GenesisState)
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[collection.ModuleName], collectionGenesisState))

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	s.commonArgs = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, cmath.NewInt(100))).String()),
	}

	s.vendor = s.network.Validators[0].Address
	s.operator = s.createAccount("operator")
	s.customer = s.createAccount("customer")
	s.stranger = s.createAccount("stranger")

	// vendor creates nft token class
	s.contractID = s.createContract(s.vendor)
	s.nftClassID = s.createNFTClass(s.contractID, s.vendor)

	// mint nfts
	s.tokenIDs = make(map[string]string, 4)
	for _, to := range []sdk.AccAddress{s.customer, s.operator, s.vendor, s.stranger} {
		s.tokenIDs[to.String()] = s.mintNFT(s.contractID, s.vendor, to, s.nftClassID)
	}

	// grant all the permissions to operator
	for _, pv := range collection.Permission_value {
		permission := collection.Permission(pv)
		if permission == collection.PermissionUnspecified {
			continue
		}
		s.grant(s.contractID, s.vendor, s.operator, permission)
	}

	// customer and vendor approves the operator to manipulate its tokens, so vendor can do OperatorXXX (Send or Burn) later.
	s.authorizeOperator(s.contractID, s.vendor, s.operator)
	s.authorizeOperator(s.contractID, s.customer, s.operator)
	// for the revocation.
	s.authorizeOperator(s.contractID, s.operator, s.vendor)

	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
	s.setupHeight, err = s.network.LatestHeight()

}

func (s *E2ETestSuite) TearDownSuite() {
	s.T().Log("tearing down collection e2e test suite")
	s.network.Cleanup()
}

func (s *E2ETestSuite) createContract(creator sdk.AccAddress) string {
	val := s.network.Validators[0]
	args := append([]string{
		creator.String(),
	}, s.commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdCreateContract(), args)
	s.Require().NoError(err)
	txResp := s.getTxResp(out, 0)
	var event collection.EventCreatedContract
	s.pickEvent(txResp.Events, &event, func(e proto.Message) {
		event = *e.(*collection.EventCreatedContract)
	})
	return event.ContractId
}

func (s *E2ETestSuite) createNFTClass(contractID string, operator sdk.AccAddress) string {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator.String(),
	}, s.commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdIssueNFT(), args)
	s.Require().NoError(err)
	txResp := s.getTxResp(out, 0)
	var event collection.EventCreatedNFTClass
	s.pickEvent(txResp.Events, &event, func(e proto.Message) {
		event = *e.(*collection.EventCreatedNFTClass)
	})
	return event.TokenType
}

func (s *E2ETestSuite) mintNFT(contractID string, operator, to sdk.AccAddress, classID string) string {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator.String(),
		to.String(),
		classID,
		fmt.Sprintf("--%s=%s", cli.FlagName, "arctic fox"),
	}, s.commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdMintNFT(), args)
	s.Require().NoError(err)
	txResp := s.getTxResp(out, 0)
	var event collection.EventMintedNFT
	s.pickEvent(txResp.Events, &event, func(e proto.Message) {
		event = *e.(*collection.EventMintedNFT)
	})

	s.Require().Equal(1, len(event.Tokens))
	return event.Tokens[0].TokenId
}

func (s *E2ETestSuite) grant(contractID string, granter, grantee sdk.AccAddress, permission collection.Permission) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		granter.String(),
		grantee.String(),
		collection.LegacyPermission(permission).String(),
	}, s.commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdGrantPermission(), args)
	s.Require().NoError(err)
	_ = s.getTxResp(out, 0)
}

func (s *E2ETestSuite) authorizeOperator(contractID string, holder, operator sdk.AccAddress) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		holder.String(),
		operator.String(),
	}, s.commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdAuthorizeOperator(), args)
	s.Require().NoError(err)
	_ = s.getTxResp(out, 0)
}

func (s *E2ETestSuite) pickEvent(events []abci.Event, event proto.Message, fn func(event proto.Message)) {
	for _, e := range events {
		if e.Type == proto.MessageName(event) {
			msg, err := sdk.ParseTypedEvent(e)
			s.Require().NoError(err)

			fn(msg)
			return
		}
	}

	s.Require().Failf("event not found", "%s", events)
}

// creates an account and send some coins to it for the future transactions.
func (s *E2ETestSuite) createAccount(uid string) sdk.AccAddress {
	val := s.network.Validators[0]
	keyInfo, _, err := val.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)
	addr, err := keyInfo.GetAddress()
	s.Require().NoError(err)

	out, err := clitestutil.MsgSendExec(val.ClientCtx, val.Address, addr, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, cmath.NewInt(1000000))), address.NewBech32Codec("link"), s.commonArgs...)
	s.Require().NoError(err)
	s.getTxResp(out, 0)
	return addr
}

func (s *E2ETestSuite) getTxResp(out testutil.BufferWriter, expectedCode uint32) sdk.TxResponse {
	var res sdk.TxResponse
	s.Require().NoError(s.network.Validators[0].ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().NoError(s.network.WaitForNextBlock())
	txResp, err := clitestutil.GetTxResponse(s.network, s.network.Validators[0].ClientCtx, res.TxHash)
	s.Require().NoError(err)
	s.Require().EqualValues(expectedCode, txResp.Code, txResp.String())
	return txResp
}
