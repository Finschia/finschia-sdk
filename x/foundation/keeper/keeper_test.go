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
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	minttypes "github.com/line/lbm-sdk/x/mint/types"
	"github.com/line/lbm-sdk/x/stakingplus"
)

var (
	delPk   = ed25519.GenPrivKey().PubKey()
	delAddr = sdk.AccAddress(delPk.Address())
	valAddr = sdk.ValAddress(delAddr)
)

func TestCleanup(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	k := app.FoundationKeeper

	// add grant
	require.NoError(t, k.Grant(ctx, govtypes.ModuleName, delAddr, &stakingplus.CreateValidatorAuthorization{}))

	// cleanup
	k.Cleanup(ctx)
	require.Empty(t, k.GetGrants(ctx))
}

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

	activeProposal  uint64
	votedProposal   uint64
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
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
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

	updates := make([]foundation.Member, len(s.members))
	for i, member := range s.members {
		updates[i] = foundation.Member{
			Address:       member.String(),
			Participating: true,
		}
	}
	err := s.keeper.UpdateMembers(s.ctx, updates)
	s.Require().NoError(err)

	// create a proposal
	s.activeProposal, err = s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To:       s.stranger.String(),
			Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
		},
	})
	s.Require().NoError(err)
	for _, member := range s.members[1:] {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.activeProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create a proposal voted by all members
	s.votedProposal, err = s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To:       s.stranger.String(),
			Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
		},
	})
	s.Require().NoError(err)
	for _, member := range s.members {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.votedProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create an aborted proposal
	s.abortedProposal, err = s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To:       s.stranger.String(),
			Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
		},
	})
	s.Require().NoError(err)
	err = s.keeper.WithdrawProposal(s.ctx, s.abortedProposal)
	s.Require().NoError(err)

	// create an invalid proposal which contains invalid message
	s.invalidProposal, err = s.keeper.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Operator: s.operator.String(),
			To:       s.stranger.String(),
			Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance.Add(sdk.OneInt()))),
		},
	})
	s.Require().NoError(err)
	for _, member := range s.members {
		err := s.keeper.Vote(s.ctx, foundation.Vote{
			ProposalId: s.invalidProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
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
