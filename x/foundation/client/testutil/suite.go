package testutil

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/crypto/hd"
	"github.com/Finschia/finschia-sdk/crypto/keyring"
	"github.com/Finschia/finschia-sdk/testutil/network"
	"github.com/Finschia/finschia-sdk/testutil/testdata"

	"github.com/Finschia/finschia-sdk/client/flags"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	sdk "github.com/Finschia/finschia-sdk/types"
	bankcli "github.com/Finschia/finschia-sdk/x/bank/client/cli"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	setupHeight int64

	authority       sdk.AccAddress
	comingMember    sdk.AccAddress
	leavingMember   sdk.AccAddress
	permanentMember sdk.AccAddress
	stranger        sdk.AccAddress

	proposalID uint64
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
	testdata.RegisterInterfaces(s.cfg.InterfaceRegistry)

	genesisState := s.cfg.GenesisState

	var foundationData foundation.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[foundation.ModuleName], &foundationData))

	// enable foundation tax
	params := foundation.Params{
		FoundationTax: sdk.MustNewDecFromStr("0.2"),
	}
	foundationData.Params = params

	var strangerMnemonic string
	strangerMnemonic, s.stranger = s.createMnemonic("stranger")
	var leavingMemberMnemonic string
	leavingMemberMnemonic, s.leavingMember = s.createMnemonic("leavingmember")
	var permanentMemberMnemonic string
	permanentMemberMnemonic, s.permanentMember = s.createMnemonic("permanentmember")

	foundationData.Members = []foundation.Member{
		{
			Address:  s.leavingMember.String(),
			Metadata: "leaving member",
		},
		{
			Address:  s.permanentMember.String(),
			Metadata: "permanent member",
		},
	}

	info := foundation.DefaultFoundation()
	info.TotalWeight = sdk.NewDecFromInt(sdk.NewInt(int64(len(foundationData.Members))))
	info.SetDecisionPolicy(&foundation.ThresholdDecisionPolicy{
		Threshold: sdk.OneDec(),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: 7 * 24 * time.Hour,
		},
	})
	foundationData.Foundation = info

	// enable censorship
	censorships := []foundation.Censorship{
		{
			MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
			Authority:  foundation.CensorshipAuthorityFoundation,
		},
	}
	foundationData.Censorships = censorships

	treasuryReceivers := []sdk.AccAddress{s.stranger, s.leavingMember}
	for _, receiver := range treasuryReceivers {
		ga := foundation.GrantAuthorization{
			Grantee: receiver.String(),
		}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{})
		s.Require().NotNil(ga)
		foundationData.Authorizations = append(foundationData.Authorizations, *ga)
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

	s.authority = foundation.DefaultAuthority()
	s.createAccount("stranger", strangerMnemonic)
	s.createAccount("comingmember", comingMemberMnemonic)
	s.createAccount("leavingmember", leavingMemberMnemonic)
	s.createAccount("permanentmember", permanentMemberMnemonic)

	s.proposalID = s.submitProposal(testdata.NewTestMsg(s.authority), false)
	s.vote(s.proposalID, []sdk.AccAddress{s.leavingMember, s.permanentMember})
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
func (s *IntegrationTestSuite) submitProposal(msg sdk.Msg, try bool) uint64 {
	val := s.network.Validators[0]

	proposers := []string{s.permanentMember.String()}
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
