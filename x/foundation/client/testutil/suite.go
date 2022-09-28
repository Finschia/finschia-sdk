package testutil

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/crypto/hd"
	"github.com/line/lbm-sdk/crypto/keyring"
	"github.com/line/lbm-sdk/testutil/network"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	sdk "github.com/line/lbm-sdk/types"
	bankcli "github.com/line/lbm-sdk/x/bank/client/cli"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/client/cli"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	setupHeight int64

	operator      sdk.AccAddress
	comingMember  sdk.AccAddress
	leavingMember sdk.AccAddress
	stranger      sdk.AccAddress
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100)))),
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	genesisState := s.cfg.GenesisState

	var foundationData foundation.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[foundation.ModuleName], &foundationData))

	// enable foundation
	params := foundation.Params{
		FoundationTax: sdk.MustNewDecFromStr("0.2"),
		CensoredMsgTypeUrls: []string{
			sdk.MsgTypeURL((*stakingtypes.MsgCreateValidator)(nil)),
			sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
		},
	}
	foundationData.Params = params

	var operatorMnemonic string
	operatorMnemonic, s.operator = s.createMnemonic("operator")
	info := &foundation.FoundationInfo{
		Operator: s.operator.String(),
		Version:  1,
	}
	err := info.SetDecisionPolicy(&foundation.ThresholdDecisionPolicy{
		Threshold: sdk.OneDec(),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: time.Hour,
		},
	})
	s.Require().NoError(err)
	foundationData.Foundation = info

	var strangerMnemonic string
	strangerMnemonic, s.stranger = s.createMnemonic("stranger")
	var leavingMemberMnemonic string
	leavingMemberMnemonic, s.leavingMember = s.createMnemonic("leavingmember")

	grantees := []sdk.AccAddress{s.stranger, s.leavingMember}
	foundationData.Authorizations = make([]foundation.GrantAuthorization, len(grantees))
	for i, grantee := range grantees {
		ga := foundation.GrantAuthorization{
			Grantee: grantee.String(),
		}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{})
		s.Require().NotNil(ga)
		foundationData.Authorizations[i] = *ga
	}

	foundationDataBz, err := s.cfg.Codec.MarshalJSON(&foundationData)
	s.Require().NoError(err)
	genesisState[foundation.ModuleName] = foundationDataBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	var comingMemberMnemonic string
	comingMemberMnemonic, s.comingMember = s.createMnemonic("comingmember")

	s.createAccount("operator", operatorMnemonic)
	s.createAccount("stranger", strangerMnemonic)
	s.createAccount("comingmember", comingMemberMnemonic)
	s.createAccount("leavingmember", leavingMemberMnemonic)

	s.addMembers([]sdk.AccAddress{s.leavingMember})
	id := s.submitProposal(&foundation.MsgWithdrawFromTreasury{
		Operator: s.operator.String(),
		To:       s.network.Validators[0].Address.String(),
		Amount:   sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.OneInt())),
	}, false)
	s.vote(id, []sdk.AccAddress{s.network.Validators[0].Address, s.leavingMember})
	s.Require().NoError(s.network.WaitForNextBlock())

	s.setupHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

// submit a proposal
func (s *IntegrationTestSuite) addMembers(members []sdk.AccAddress) {
	val := s.network.Validators[0]

	updates := make([]json.RawMessage, len(members))
	for i, member := range members {
		update := foundation.MemberRequest{
			Address: member.String(),
		}
		bz, err := s.cfg.Codec.MarshalJSON(&update)
		s.Require().NoError(err)

		updates[i] = bz
	}
	updateArg, err := json.Marshal(updates)
	s.Require().NoError(err)

	args := append([]string{
		s.operator.String(),
		string(updateArg),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.operator),
	}, commonArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdUpdateMembers(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().Zero(res.Code, out.String())
}

// submit a proposal
func (s *IntegrationTestSuite) submitProposal(msg sdk.Msg, try bool) uint64 {
	val := s.network.Validators[0]

	proposers := []string{val.Address.String()}
	proposersBz, err := json.Marshal(&proposers)
	s.Require().NoError(err)

	args := append([]string{
		"test proposal",
		string(proposersBz),
		s.msgToString(msg),
	}, commonArgs...)
	if try {
		args = append(args, fmt.Sprintf("--%s=%s", cli.FlagExec, cli.ExecTry))
	}
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdSubmitProposal(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().Zero(res.Code, out.String())

	events := res.Logs[0].Events
	proposalEvent, _ := sdk.TypedEventToEvent(&foundation.EventSubmitProposal{})
	for _, e := range events {
		if e.Type == proposalEvent.Type {
			var proposal foundation.Proposal
			err := val.ClientCtx.Codec.UnmarshalJSON([]byte(e.Attributes[0].Value), &proposal)
			s.Require().NoError(err)

			return proposal.Id
		}
	}
	panic("You cannot reach here")
}

func (s *IntegrationTestSuite) vote(proposalID uint64, voters []sdk.AccAddress) {
	val := s.network.Validators[0]

	for _, voter := range voters {
		args := append([]string{
			fmt.Sprint(proposalID),
			voter.String(),
			foundation.VOTE_OPTION_YES.String(),
			"test vote",
		}, commonArgs...)
		out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdVote(), args)
		s.Require().NoError(err)

		var res sdk.TxResponse
		s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
		s.Require().Zero(res.Code, out.String())
	}
}

func (s *IntegrationTestSuite) msgToString(msg sdk.Msg) string {
	anyJSON, err := s.cfg.Codec.MarshalInterfaceJSON(msg)
	s.Require().NoError(err)

	cliMsgs := []json.RawMessage{anyJSON}
	msgsBz, err := json.Marshal(cliMsgs)
	s.Require().NoError(err)

	return string(msgsBz)
}

// creates an account
func (s *IntegrationTestSuite) createMnemonic(uid string) (string, sdk.AccAddress) {
	cstore := keyring.NewInMemory()
	info, mnemonic, err := cstore.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)

	return mnemonic, info.GetAddress()
}

// creates an account and send some coins to it for the future transactions.
func (s *IntegrationTestSuite) createAccount(uid, mnemonic string) {
	val := s.network.Validators[0]
	info, err := val.ClientCtx.Keyring.NewAccount(uid, mnemonic, keyring.DefaultBIP39Passphrase, sdk.FullFundraiserPath, hd.Secp256k1)
	s.Require().NoError(err)

	addr := info.GetAddress()
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
	s.Require().Zero(res.Code, out.String())
}
