package testutil

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/crypto/hd"
	"github.com/Finschia/finschia-sdk/crypto/keyring"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	"github.com/Finschia/finschia-sdk/testutil/network"
	sdk "github.com/Finschia/finschia-sdk/types"
	bankcli "github.com/Finschia/finschia-sdk/x/bank/client/cli"
	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/client/cli"
	abci "github.com/tendermint/tendermint/abci/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	setupHeight int64

	vendor   sdk.AccAddress
	operator sdk.AccAddress
	customer sdk.AccAddress
	stranger sdk.AccAddress

	contractID string
	ftClassID  string
	nftClassID string

	balance sdk.Int

	lenChain int
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))).String()),
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	var gs collection.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(s.cfg.GenesisState[collection.ModuleName], &gs))

	params := collection.Params{
		DepthLimit: 4,
		WidthLimit: 4,
	}
	gs.Params = params

	gsBz, err := s.cfg.Codec.MarshalJSON(&gs)
	s.Require().NoError(err)
	s.cfg.GenesisState[collection.ModuleName] = gsBz

	s.network = network.New(s.T(), s.cfg)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.vendor = s.createAccount("vendor")
	s.operator = s.createAccount("operator")
	s.customer = s.createAccount("customer")
	s.stranger = s.createAccount("stranger")

	s.balance = sdk.NewInt(1000000)

	// vendor creates 2 token classes
	s.contractID = s.createContract(s.vendor)
	s.ftClassID = s.createFTClass(s.contractID, s.vendor)
	s.nftClassID = s.createNFTClass(s.contractID, s.vendor)

	// mint & burn fts
	for _, to := range []sdk.AccAddress{s.customer, s.operator, s.vendor, s.stranger} {
		s.mintFT(s.contractID, s.vendor, to, s.ftClassID, s.balance)

		if to.Equals(s.vendor) {
			tokenID := collection.NewFTID(s.ftClassID)
			amount := collection.NewCoins(collection.NewCoin(tokenID, s.balance))
			s.burnFT(s.contractID, s.vendor, amount)
			s.mintFT(s.contractID, s.vendor, to, s.ftClassID, s.balance)
		}
	}

	// mint nfts
	s.lenChain = 2
	for _, to := range []sdk.AccAddress{s.customer, s.operator, s.vendor, s.stranger} {
		// mint N chains per account
		numChains := 3
		for n := 0; n < numChains; n++ {
			ids := make([]string, s.lenChain)
			for i := range ids {
				ids[i] = s.mintNFT(s.contractID, s.vendor, to, s.nftClassID)
			}

			for i := range ids[1:] {
				r := len(ids) - 1 - i
				subject := ids[r]
				target := ids[r-1]
				s.attach(s.contractID, to, subject, target)
			}
		}
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
	s.authorizeOperator(s.contractID, s.customer, s.operator)
	s.authorizeOperator(s.contractID, s.vendor, s.operator)
	// for the revocation.
	s.authorizeOperator(s.contractID, s.operator, s.vendor)

	s.setupHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) pickEvent(events []abci.Event, event proto.Message, fn func(event proto.Message)) {
	getType := func(msg proto.Message) string {
		return proto.MessageName(msg)
	}

	for _, e := range events {
		if e.Type == getType(event) {
			msg, err := sdk.ParseTypedEvent(e)
			s.Require().NoError(err)

			fn(msg)
			return
		}
	}

	s.Require().Failf("event not found", "%s", events)
}

func (s *IntegrationTestSuite) createContract(creator sdk.AccAddress) string {
	val := s.network.Validators[0]
	args := append([]string{
		creator.String(),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdCreateContract(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	var event collection.EventCreatedContract
	s.pickEvent(res.Events, &event, func(e proto.Message) {
		event = *e.(*collection.EventCreatedContract)
	})
	return event.ContractId
}

func (s *IntegrationTestSuite) createFTClass(contractID string, operator sdk.AccAddress) string {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator.String(),
		fmt.Sprintf("--%s=%s", cli.FlagName, "tibetian fox"),
		fmt.Sprintf("--%s=%s", cli.FlagTo, operator),
		fmt.Sprintf("--%s=%v", cli.FlagMintable, true),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdIssueFT(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	var event collection.EventCreatedFTClass
	s.pickEvent(res.Events, &event, func(e proto.Message) {
		event = *e.(*collection.EventCreatedFTClass)
	})
	return collection.SplitTokenID(event.TokenId)
}

func (s *IntegrationTestSuite) createNFTClass(contractID string, operator sdk.AccAddress) string {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator.String(),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdIssueNFT(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	var event collection.EventCreatedNFTClass
	s.pickEvent(res.Events, &event, func(e proto.Message) {
		event = *e.(*collection.EventCreatedNFTClass)
	})
	return event.TokenType
}

func (s *IntegrationTestSuite) mintFT(contractID string, operator, to sdk.AccAddress, classID string, amount sdk.Int) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator.String(),
		to.String(),
		classID,
		amount.String(),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdMintFT(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}

func (s *IntegrationTestSuite) mintNFT(contractID string, operator, to sdk.AccAddress, classID string) string {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		operator.String(),
		to.String(),
		classID,
		fmt.Sprintf("--%s=%s", cli.FlagName, "arctic fox"),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdMintNFT(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	var event collection.EventMintedNFT
	s.pickEvent(res.Events, &event, func(e proto.Message) {
		event = *e.(*collection.EventMintedNFT)
	})

	s.Require().Equal(1, len(event.Tokens))
	return event.Tokens[0].TokenId
}

func (s *IntegrationTestSuite) burnFT(contractID string, from sdk.AccAddress, amount collection.Coins) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		from.String(),
		amount.String(),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdBurnFT(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}

func (s *IntegrationTestSuite) attach(contractID string, holder sdk.AccAddress, subject, target string) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		holder.String(),
		subject,
		target,
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdAttach(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}

func (s *IntegrationTestSuite) grant(contractID string, granter, grantee sdk.AccAddress, permission collection.Permission) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		granter.String(),
		grantee.String(),
		collection.LegacyPermission(permission).String(),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdGrantPermission(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}

// creates an account and send some coins to it for the future transactions.
func (s *IntegrationTestSuite) createAccount(uid string) sdk.AccAddress {
	val := s.network.Validators[0]
	keyInfo, _, err := val.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)
	addr := keyInfo.GetAddress()

	fee := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000000)))
	args := append([]string{
		val.Address.String(),
		addr.String(),
		fee.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
	}, commonArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, bankcli.NewSendTxCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	return addr
}

func (s *IntegrationTestSuite) authorizeOperator(contractID string, holder, operator sdk.AccAddress) {
	val := s.network.Validators[0]
	args := append([]string{
		contractID,
		holder.String(),
		operator.String(),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdAuthorizeOperator(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
