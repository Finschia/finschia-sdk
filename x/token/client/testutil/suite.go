package testutil

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	abci "github.com/line/ostracon/abci/types"

	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/crypto/hd"
	"github.com/line/lbm-sdk/crypto/keyring"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	sdk "github.com/line/lbm-sdk/types"
	bankcli "github.com/line/lbm-sdk/x/bank/client/cli"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	setupHeight int64

	vendor   sdk.AccAddress
	customer sdk.AccAddress

	classes []token.TokenClass

	balance sdk.Int
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

	s.network = network.New(s.T(), s.cfg)
	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.vendor = s.createAccount("vendor")
	s.customer = s.createAccount("customer")

	s.classes = []token.TokenClass{
		{
			Name:     "test",
			Symbol:   "ZERO",
			Decimals: 8,
			Mintable: true,
		},
		{
			Name:     "test",
			Symbol:   "ONE",
			Decimals: 8,
			Mintable: true,
		},
	}

	s.balance = sdk.NewInt(1000)

	// vendor creates 2 tokens
	s.classes[1].ContractId = s.createClass(s.vendor, s.vendor, s.classes[1].Name, s.classes[1].Symbol, s.balance, s.classes[1].Mintable)
	s.classes[0].ContractId = s.createClass(s.vendor, s.customer, s.classes[0].Name, s.classes[0].Symbol, s.balance, s.classes[0].Mintable)

	// customer approves vendor to manipulate its tokens of the both classes, so vendor can do OperatorXXX (Send or Burn) later.
	for _, class := range s.classes {
		s.authorizeOperator(class.ContractId, s.customer, s.vendor)
	}

	s.setupHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) createClass(owner, to sdk.AccAddress, name, symbol string, supply sdk.Int, mintable bool) string {
	val := s.network.Validators[0]
	args := append([]string{
		owner.String(),
		to.String(),
		name,
		symbol,
		fmt.Sprintf("--%s=%v", cli.FlagMintable, mintable),
		fmt.Sprintf("--%s=%s", cli.FlagSupply, supply),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, owner),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdIssue(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	var event token.EventIssued
	s.pickEvent(res.Events, &event, func(e proto.Message) {
		event = *e.(*token.EventIssued)
	})
	return event.ContractId
}

// creates an account and send some coins to it for the future transactions.
func (s *IntegrationTestSuite) createAccount(uid string) sdk.AccAddress {
	val := s.network.Validators[0]
	keyInfo, _, err := val.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)
	addr := keyInfo.GetAddress()

	fee := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000)))
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
		fmt.Sprintf("--%s=%s", flags.FlagFrom, holder),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdApprove(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
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

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
