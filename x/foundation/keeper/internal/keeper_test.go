package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/foundation"
	keeper "github.com/line/lbm-sdk/x/foundation/keeper"
	"github.com/line/lbm-sdk/x/foundation/keeper/internal"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	minttypes "github.com/line/lbm-sdk/x/mint/types"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context

	bankKeeper foundation.BankKeeper
	keeper     keeper.Keeper
	impl       internal.Keeper

	queryServer     foundation.QueryServer
	msgServer       foundation.MsgServer
	proposalHandler govtypes.Handler

	authority sdk.AccAddress
	members   []sdk.AccAddress
	stranger  sdk.AccAddress

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
	app := simapp.Setup(checkTx)
	testdata.RegisterInterfaces(app.InterfaceRegistry())
	testdata.RegisterMsgServer(app.MsgServiceRouter(), testdata.MsgServerImpl{})

	s.ctx = app.BaseApp.NewContext(checkTx, tmproto.Header{})
	s.bankKeeper = app.BankKeeper
	s.keeper = app.FoundationKeeper
	s.impl = internal.NewKeeper(
		app.AppCodec(),
		app.GetKey(foundation.ModuleName),
		app.MsgServiceRouter(),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		foundation.DefaultConfig(),
		foundation.DefaultAuthority().String(),
	)

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	s.proposalHandler = keeper.NewFoundationProposalsHandler(s.keeper)

	s.impl.SetParams(s.ctx, foundation.Params{
		FoundationTax: sdk.OneDec(),
	})

	s.impl.SetCensorship(s.ctx, foundation.Censorship{
		MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
		Authority:  foundation.CensorshipAuthorityFoundation,
	})

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	s.authority = sdk.MustAccAddressFromBech32(s.impl.GetAuthority())
	s.members = make([]sdk.AccAddress, 10)
	for i := range s.members {
		s.members[i] = createAddress()
		member := foundation.Member{
			Address: s.members[i].String(),
		}
		s.impl.SetMember(s.ctx, member)
	}
	s.stranger = createAddress()

	info := foundation.DefaultFoundation()
	info.TotalWeight = sdk.NewDec(int64(len(s.members)))
	err := info.SetDecisionPolicy(workingPolicy())
	s.Require().NoError(err)
	s.impl.SetFoundationInfo(s.ctx, info)

	s.balance = sdk.NewInt(987654321)
	s.impl.SetPool(s.ctx, foundation.Pool{
		Treasury: sdk.NewDecCoinsFromCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
	})
	holders := []sdk.AccAddress{
		s.stranger,
		app.AccountKeeper.GetModuleAccount(s.ctx, foundation.TreasuryName).GetAddress(),
		app.AccountKeeper.GetModuleAccount(s.ctx, authtypes.FeeCollectorName).GetAddress(),
	}
	for _, holder := range holders {
		amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance))

		// using minttypes here introduces dependency on x/mint
		// the work around would be registering a new module account on this suite
		// because x/bank already has dependency on x/mint, and we must have dependency
		// on x/bank, it's OK to use x/mint here.
		minterName := minttypes.ModuleName
		err := app.BankKeeper.MintCoins(s.ctx, minterName, amount)
		s.Require().NoError(err)

		minter := app.AccountKeeper.GetModuleAccount(s.ctx, minterName).GetAddress()
		err = app.BankKeeper.SendCoins(s.ctx, minter, holder, amount)
		s.Require().NoError(err)
	}

	// create an active proposal, voted yes by all members except the first member
	activeProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Authority: s.authority.String(),
			To:        s.stranger.String(),
			Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
		},
	})
	s.Require().NoError(err)
	s.activeProposal = *activeProposal

	for _, member := range s.members[1:] {
		err := s.impl.Vote(s.ctx, foundation.Vote{
			ProposalId: s.activeProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create a proposal voted no by all members
	votedProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{newMsgCreateDog("shiba1")})
	s.Require().NoError(err)
	s.votedProposal = *votedProposal

	for _, member := range s.members {
		err := s.impl.Vote(s.ctx, foundation.Vote{
			ProposalId: s.votedProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_NO,
		})
		s.Require().NoError(err)
	}

	// create an withdrawn proposal
	withdrawnProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{newMsgCreateDog("shiba2")})
	s.Require().NoError(err)
	s.withdrawnProposal = *withdrawnProposal

	err = s.impl.WithdrawProposal(s.ctx, s.withdrawnProposal)
	s.Require().NoError(err)

	// create an invalid proposal which contains invalid message
	invalidProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Authority: s.authority.String(),
			To:        s.stranger.String(),
			Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance.Add(sdk.OneInt()))),
		},
	})
	s.Require().NoError(err)
	s.invalidProposal = *invalidProposal

	for _, member := range s.members {
		err := s.impl.Vote(s.ctx, foundation.Vote{
			ProposalId: s.invalidProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create an invalid proposal which contains invalid message
	noHandlerProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.members[0].String()}, "", []sdk.Msg{testdata.NewTestMsg(s.authority)})
	s.Require().NoError(err)
	s.noHandlerProposal = *noHandlerProposal

	for _, member := range s.members {
		err := s.impl.Vote(s.ctx, foundation.Vote{
			ProposalId: s.noHandlerProposal,
			Voter:      member.String(),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// next proposal is the proposal id for the upcoming proposal
	s.nextProposal = s.noHandlerProposal + 1

	// grant stranger to receive foundation treasury
	err = s.impl.Grant(s.ctx, s.stranger, &foundation.ReceiveFromTreasuryAuthorization{})
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func TestNewKeeper(t *testing.T) {
	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}
	authority := foundation.DefaultAuthority()

	testCases := map[string]struct {
		authority sdk.AccAddress
		panics    bool
	}{
		"default authority": {
			authority: authority,
		},
		"invalid account": {
			panics: true,
		},
		"not the default authority": {
			authority: createAddress(),
			panics:    true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			newKeeper := func() keeper.Keeper {
				app := simapp.Setup(false)
				return keeper.NewKeeper(app.AppCodec(), sdk.NewKVStoreKey(foundation.StoreKey), app.MsgServiceRouter(), app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName, foundation.DefaultConfig(), tc.authority.String())
			}

			if tc.panics {
				require.Panics(t, func() { newKeeper() })
				return
			}
			require.NotPanics(t, func() { newKeeper() })

			k := newKeeper()
			require.Equal(t, authority.String(), k.GetAuthority())
		})
	}
}
