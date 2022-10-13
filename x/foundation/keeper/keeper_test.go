package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
	minttypes "github.com/line/lbm-sdk/x/mint/types"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context

	app         *simapp.SimApp
	keeper      keeper.Keeper
	queryServer foundation.QueryServer
	msgServer   foundation.MsgServer

	operator sdk.AccAddress
	members  []sdk.AccAddress
	stranger sdk.AccAddress

	activeProposal    uint64
	votedProposal     uint64
	withdrawnProposal uint64
	invalidProposal   uint64
	noHandlerProposal uint64
	nextProposal      uint64

	balance sdk.Int
}

func newMsgCreateDog(name string) sdk.Msg {
	return &testdata.MsgCreateDog{
		Dog: &testdata.Dog{
			Name: name,
		},
	}
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	s.app = simapp.Setup(checkTx)
	testdata.RegisterInterfaces(s.app.InterfaceRegistry())
	testdata.RegisterMsgServer(s.app.MsgServiceRouter(), testdata.MsgServerImpl{})

	s.ctx = s.app.BaseApp.NewContext(checkTx, ocproto.Header{})
	s.keeper = s.app.FoundationKeeper

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	s.keeper.SetParams(s.ctx, foundation.Params{
		FoundationTax: sdk.OneDec(),
		CensoredMsgTypeUrls: []string{
			sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
		},
	})

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	s.operator = s.keeper.GetOperator(s.ctx)
	s.members = make([]sdk.AccAddress, 10)
	for i := range s.members {
		s.members[i] = createAddress()
		member := foundation.Member{
			Address: s.members[i].String(),
		}
		s.keeper.SetMember(s.ctx, member)
	}
	s.stranger = createAddress()

	info := foundation.DefaultFoundation()
	info.TotalWeight = sdk.NewDec(int64(len(s.members)))
	err := info.SetDecisionPolicy(workingPolicy())
	s.Require().NoError(err)
	s.keeper.SetFoundationInfo(s.ctx, info)

	s.balance = sdk.NewInt(1000000)
	s.keeper.SetPool(s.ctx, foundation.Pool{
		Treasury: sdk.NewDecCoinsFromCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
	})
	holders := []sdk.AccAddress{
		s.stranger,
		s.app.AccountKeeper.GetModuleAccount(s.ctx, foundation.TreasuryName).GetAddress(),
		s.app.AccountKeeper.GetModuleAccount(s.ctx, authtypes.FeeCollectorName).GetAddress(),
	}
	for _, holder := range holders {
		amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance))

		// using minttypes here introduces dependency on x/mint
		// the work around would be registering a new module account on this suite
		// because x/bank already has dependency on x/mint, and we must have dependency
		// on x/bank, it's OK to use x/mint here.
		minterName := minttypes.ModuleName
		err := s.app.BankKeeper.MintCoins(s.ctx, minterName, amount)
		s.Require().NoError(err)

		minter := s.app.AccountKeeper.GetModuleAccount(s.ctx, minterName).GetAddress()
		err = s.app.BankKeeper.SendCoins(s.ctx, minter, holder, amount)
		s.Require().NoError(err)
	}

	// create a proposal
	activeProposal, err := s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{newMsgCreateDog("shiba1")})
	s.Require().NoError(err)
	s.activeProposal = *activeProposal

	for _, member := range s.members[1:] {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.activeProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create a proposal voted by all members
	votedProposal, err := s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{newMsgCreateDog("shiba2")})
	s.Require().NoError(err)
	s.votedProposal = *votedProposal

	for _, member := range s.members {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.votedProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_NO,
		})
		s.Require().NoError(err)
	}

	// create an withdrawn proposal
	withdrawnProposal, err := s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{newMsgCreateDog("shiba3")})
	s.Require().NoError(err)
	s.withdrawnProposal = *withdrawnProposal

	err = s.keeper.WithdrawProposal(s.ctx, s.withdrawnProposal)
	s.Require().NoError(err)

	// create an invalid proposal which contains invalid message
	invalidProposal, err := s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To:       s.stranger.String(),
			Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance.Add(sdk.OneInt()))),
		},
	})
	s.Require().NoError(err)
	s.invalidProposal = *invalidProposal

	for _, member := range s.members {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.invalidProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create an invalid proposal which contains invalid message
	noHandlerProposal, err := s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{testdata.NewTestMsg(s.operator)})
	s.Require().NoError(err)
	s.noHandlerProposal = *noHandlerProposal

	for _, member := range s.members {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.noHandlerProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// next proposal is the proposal id for the upcoming proposal
	s.nextProposal = s.noHandlerProposal + 1

	// grant stranger to receive foundation treasury
	err = s.keeper.Grant(s.ctx, s.stranger, &foundation.ReceiveFromTreasuryAuthorization{})
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
