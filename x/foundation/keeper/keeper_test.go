package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
)

var (
	delPk   = ed25519.GenPrivKey().PubKey()
	delAddr = sdk.BytesToAccAddress(delPk.Address())
	valAddr = delAddr.ToValAddress()
)

func TestCleanup(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	k := app.FoundationKeeper

	// add auths
	auth := foundation.ValidatorAuth{
		OperatorAddress: valAddr.String(),
		CreationAllowed: true,
	}
	require.NoError(t, k.SetValidatorAuth(ctx, auth))

	// cleanup
	k.Cleanup(ctx)
	require.Empty(t, k.GetValidatorAuths(ctx))
}

type KeeperTestSuite struct {
	suite.Suite
	ctx  sdk.Context

	app *simapp.SimApp
	keeper keeper.Keeper
	queryServer foundation.QueryServer
	msgServer foundation.MsgServer

	operator sdk.AccAddress
	members []sdk.AccAddress
	stranger sdk.AccAddress

	activeProposal uint64
	votedProposal uint64
	abortedProposal uint64
	invalidProposal uint64

	balance sdk.Int
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	s.app = simapp.Setup(checkTx)
	s.ctx = s.app.BaseApp.NewContext(checkTx, ocproto.Header{})
	s.keeper = s.app.FoundationKeeper

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	createAddress := func() sdk.AccAddress {
		return sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	s.operator = s.keeper.GetOperator(s.ctx)
	s.members = make([]sdk.AccAddress, foundation.DefaultConfig().MinThreshold.TruncateInt64())
	for i := range s.members {
		s.members[i] = createAddress()
	}
	s.stranger = createAddress()

	s.balance = sdk.NewInt(1000000)
	holders := []sdk.AccAddress{
		s.stranger,
		s.app.AccountKeeper.GetModuleAccount(s.ctx, foundation.TreasuryName).GetAddress(),
	}
	for _, holder := range holders {
		err := s.app.BankKeeper.SetBalance(s.ctx, holder, sdk.NewCoin(sdk.DefaultBondDenom, s.balance))
		s.Require().NoError(err)
	}

	updates := make([]foundation.Member, len(s.members))
	for i, member := range s.members {
		updates[i] = foundation.Member{
			Address: member.String(),
			Weight: sdk.OneDec(),
		}
	}
	err := s.keeper.UpdateMembers(s.ctx, updates)
	s.Require().NoError(err)

	// create a proposal
	s.activeProposal, err = s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To: s.stranger.String(),
			Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
		},
	})
	s.Require().NoError(err)
	for _, member := range s.members[1:] {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.activeProposal,
			Voter: member.String(),
			Option: foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create a proposal voted by all members
	s.votedProposal, err = s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To: s.stranger.String(),
			Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
		},
	})
	s.Require().NoError(err)
	for _, member := range s.members {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.votedProposal,
			Voter: member.String(),
			Option: foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create an aborted proposal
	s.abortedProposal, err = s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To: s.stranger.String(),
			Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
		},
	})
	s.Require().NoError(err)
	err = s.keeper.WithdrawProposal(s.ctx, s.abortedProposal)
	s.Require().NoError(err)

	// create an invalid proposal which contains invalid message
	s.invalidProposal, err = s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To: s.stranger.String(),
			Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance.Add(sdk.OneInt()))),
		},
	})
	s.Require().NoError(err)
	for _, member := range s.members {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.invalidProposal,
			Voter: member.String(),
			Option: foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// grant stranger to receive foundation treasury
	err = s.keeper.Grant(s.ctx, foundation.ModuleName, s.stranger, &foundation.ReceiveFromTreasuryAuthorization{})
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
