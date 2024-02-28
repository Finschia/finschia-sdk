package collection

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/address"
	cmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/client/cli"
)

type E2ETestSuite struct {
	suite.Suite

	cfg         network.Config
	ac          address.Codec
	network     *network.Network
	setupHeight int64

	commonArgs []string

	vendor     string
	operator   string
	customer   string
	stranger   string
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
	s.ac = s.network.Config.InterfaceRegistry.SigningContext().AddressCodec()

	s.commonArgs = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, cmath.NewInt(100))).String()),
	}

	s.vendor, err = s.ac.BytesToString(s.network.Validators[0].Address)
	s.Require().NoError(err)
	s.operator = s.createAccount("operator")
	s.customer = s.createAccount("customer")
	s.stranger = s.createAccount("stranger")

	// vendor creates nft token class
	s.contractID = s.createContract(s.vendor)
	s.nftClassID = s.createNFTClass(s.contractID, s.vendor)

	// mint nfts
	s.tokenIDs = make(map[string]string, 4)
	for _, to := range []string{s.customer, s.operator, s.vendor, s.stranger} {
		s.tokenIDs[to] = s.mintNFT(s.contractID, s.vendor, to, s.nftClassID)
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

	s.Require().NoError(s.network.WaitForNextBlock())
	s.setupHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)
}

func (s *E2ETestSuite) TearDownSuite() {
	s.T().Log("tearing down collection e2e test suite")
	s.network.Cleanup()
}

func (s *E2ETestSuite) createContract(creator string) string {
	val := s.network.Validators[0]
	args := append([]string{
		creator,
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

func (s *E2ETestSuite) createNFTClass(contractID, operator string) string {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator,
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

func (s *E2ETestSuite) mintNFT(contractID, operator, to, classID string) string {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator,
		to,
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

func (s *E2ETestSuite) burnNFT(contractID, operator, tokenID string) collection.Coins {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator,
		tokenID,
	}, s.commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdBurnNFT(), args)
	s.Require().NoError(err)
	txResp := s.getTxResp(out, 0)
	var event collection.EventBurned
	s.pickEvent(txResp.Events, &event, func(e proto.Message) {
		event = *e.(*collection.EventBurned)
	})

	s.Require().GreaterOrEqual(len(event.Amount), 1)
	return event.Amount
}

func (s *E2ETestSuite) grant(contractID, granter, grantee string, permission collection.Permission) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		granter,
		grantee,
		collection.LegacyPermission(permission).String(),
	}, s.commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdGrantPermission(), args)
	s.Require().NoError(err)
	_ = s.getTxResp(out, 0)
}

func (s *E2ETestSuite) authorizeOperator(contractID, holder, operator string) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		holder,
		operator,
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
func (s *E2ETestSuite) createAccount(uid string) string {
	val := s.network.Validators[0]
	keyInfo, _, err := val.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)
	addr, err := keyInfo.GetAddress()
	s.Require().NoError(err)

	out, err := clitestutil.MsgSendExec(val.ClientCtx, val.Address, addr, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, cmath.NewInt(1000000))), s.ac, s.commonArgs...)
	s.Require().NoError(err)
	s.getTxResp(out, 0)
	a, err := s.ac.BytesToString(addr)
	s.Require().NoError(err)
	return a
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
