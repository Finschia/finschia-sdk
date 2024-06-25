package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/fbridge/keeper"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

type IntegrationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
	msgServer   types.MsgServer
	guardians   []sdk.AccAddress
	operator    sdk.AccAddress
	ethAddr     string
}

func (s *IntegrationTestSuite) SetupTest() {
	s.app = simapp.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
	s.queryClient = types.NewQueryClient(queryHelper)
	types.RegisterQueryServer(queryHelper, s.app.FbridgeKeeper)
	s.msgServer = keeper.NewMsgServer(s.app.FbridgeKeeper)

	s.guardians = simapp.AddTestAddrs(s.app, s.ctx, 3, sdk.NewInt(1000000000))
	for _, guardian := range s.guardians {
		_, err := s.app.FbridgeKeeper.RegisterRoleProposal(s.ctx, types.DefaultAuthority(), guardian, types.RoleGuardian)
		s.Require().NoError(err)
	}
	s.operator = simapp.AddTestAddrs(s.app, s.ctx, 1, sdk.NewInt(1000000000))[0]
	_, err := s.app.FbridgeKeeper.RegisterRoleProposal(s.ctx, types.DefaultAuthority(), s.operator, types.RoleOperator)
	s.Require().NoError(err)
	s.app.FbridgeKeeper.EndBlocker(s.ctx)

	s.ethAddr = "0x1A7C26B0437Aa2d3c8454383650a5D3c35087f91"
}

func (s *IntegrationTestSuite) TestInactiveQuries() {
	goctx := sdk.WrapSDKContext(s.ctx)

	s.Require().Panics(func() {
		_, _ = s.queryClient.GreatestSeqByOperator(goctx, &types.QueryGreatestSeqByOperatorRequest{})
	})

	s.Require().Panics(func() {
		_, _ = s.queryClient.GreatestConsecutiveConfirmedSeq(goctx, &types.QueryGreatestConsecutiveConfirmedSeqRequest{})
	})

	s.Require().Panics(func() {
		_, _ = s.queryClient.SubmittedProvision(goctx, &types.QuerySubmittedProvisionRequest{})
	})

	s.Require().Panics(func() {
		_, _ = s.queryClient.ConfirmedProvision(goctx, &types.QueryConfirmedProvisionRequest{})
	})

	s.Require().Panics(func() {
		_, _ = s.queryClient.NeededSubmissionSeqs(goctx, &types.QueryNeededSubmissionSeqsRequest{})
	})

	s.Require().Panics(func() {
		_, _ = s.queryClient.Commitments(goctx, &types.QueryCommitmentsRequest{})
	})
}

func (s *IntegrationTestSuite) TestParams() {
	goctx := sdk.WrapSDKContext(s.ctx)
	res, err := s.queryClient.Params(goctx, &types.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().EqualValues(types.DefaultParams(), res.Params)
}

func (s *IntegrationTestSuite) TestNextSeqSend() {
	goctx := sdk.WrapSDKContext(s.ctx)
	res, err := s.queryClient.NextSeqSend(goctx, &types.QueryNextSeqSendRequest{})
	s.Require().NoError(err)
	s.Require().EqualValues(1, res.Seq)
}

func (s *IntegrationTestSuite) TestSeqToBlocknums() {
	goctx := sdk.WrapSDKContext(s.ctx)
	req := new(types.QuerySeqToBlocknumsRequest)

	tcs := map[string]struct {
		expErr   bool
		expBlock []uint64
		malleate func()
	}{
		"empty request": {
			expErr: true,
			malleate: func() {
				req = &types.QuerySeqToBlocknumsRequest{}
			},
		},
		"exceed upper bound (1000)": {
			expErr: true,
			malleate: func() {
				seqs := [1001]uint64{}
				req = &types.QuerySeqToBlocknumsRequest{Seqs: seqs[:]}
			},
		},
		"seq not found": {
			expErr: true,
			malleate: func() {
				req = &types.QuerySeqToBlocknumsRequest{Seqs: []uint64{1001}}
			},
		},
		"success": {
			expErr:   false,
			expBlock: []uint64{0, 0},
			malleate: func() {
				_, err := s.msgServer.Transfer(goctx, &types.MsgTransfer{
					Sender:   s.guardians[0].String(),
					Receiver: s.ethAddr,
					Amount:   sdk.NewInt(100),
				})
				s.Require().NoError(err)
				req = &types.QuerySeqToBlocknumsRequest{Seqs: []uint64{1, 2}}
				_, err = s.msgServer.Transfer(goctx, &types.MsgTransfer{
					Sender:   s.guardians[1].String(),
					Receiver: s.ethAddr,
					Amount:   sdk.NewInt(100),
				})
				s.Require().NoError(err)

				req = &types.QuerySeqToBlocknumsRequest{Seqs: []uint64{1, 2}}
			},
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			tc.malleate()
			res, err := s.queryClient.SeqToBlocknums(goctx, req)
			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expBlock, res.Blocknums)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestMembers() {
	goctx := sdk.WrapSDKContext(s.ctx)
	req := new(types.QueryMembersRequest)
	tcs := map[string]struct {
		expErr   bool
		expLen   int
		malleate func()
	}{
		"query all members": {
			expErr: false,
			expLen: 4,
			malleate: func() {
				req = &types.QueryMembersRequest{}
			},
		},
		"query guardian group": {
			expErr: false,
			expLen: 3,
			malleate: func() {
				req = &types.QueryMembersRequest{Role: "guardian"}
			},
		},
		"query operator group": {
			expErr: false,
			expLen: 1,
			malleate: func() {
				req = &types.QueryMembersRequest{Role: "operator"}
			},
		},
		"query judge group": {
			expErr: false,
			expLen: 0,
			malleate: func() {
				req = &types.QueryMembersRequest{Role: "judge"}
			},
		},
		"query invalid group": {
			expErr: true,
			malleate: func() {
				req = &types.QueryMembersRequest{Role: "invalid"}
			},
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			tc.malleate()
			res, err := s.queryClient.Members(goctx, req)
			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Len(res.Members, tc.expLen)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestMember() {
	goctx := sdk.WrapSDKContext(s.ctx)
	req := new(types.QueryMemberRequest)
	tcs := map[string]struct {
		expErr   bool
		expRole  string
		malleate func()
	}{
		"query a member who has a role": {
			expErr:  false,
			expRole: "GUARDIAN",
			malleate: func() {
				req = &types.QueryMemberRequest{Address: s.guardians[0].String()}
			},
		},
		"query a member who doesn't have a role": {
			expErr: true,
			malleate: func() {
				dummy := simapp.AddTestAddrs(s.app, s.ctx, 1, sdk.NewInt(1000000000))[0]
				req = &types.QueryMemberRequest{Address: dummy.String()}
			},
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			tc.malleate()
			if tc.expErr {
				_, err := s.queryClient.Member(goctx, req)
				s.Require().Error(err)
			} else {
				res, err := s.queryClient.Member(goctx, req)
				s.Require().NoError(err)
				s.Require().Equal(tc.expRole, res.Role)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestProposals() {
	goctx := sdk.WrapSDKContext(s.ctx)
	expProposalID := []uint64{5, 6}
	_, err := s.msgServer.SuggestRole(goctx, &types.MsgSuggestRole{
		From:   s.guardians[0].String(),
		Target: s.guardians[1].String(),
		Role:   types.RoleJudge,
	})
	s.Require().NoError(err)
	_, err = s.msgServer.SuggestRole(goctx, &types.MsgSuggestRole{
		From:   s.guardians[0].String(),
		Target: s.guardians[2].String(),
		Role:   types.RoleOperator,
	})
	s.Require().NoError(err)

	req := &types.QueryProposalsRequest{
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      10,
			CountTotal: true,
			Reverse:    false,
		},
	}

	res, err := s.queryClient.Proposals(goctx, req)
	s.Require().NoError(err)
	for i, proposal := range res.Proposals {
		s.Require().Equal(expProposalID[i], proposal.Id)
	}

	req2 := &types.QueryProposalRequest{
		ProposalId: expProposalID[1],
	}

	res2, err := s.queryClient.Proposal(goctx, req2)
	s.Require().NoError(err)
	s.Require().Equal(expProposalID[1], res2.Proposal.Id)

	req2.ProposalId++
	_, err = s.queryClient.Proposal(goctx, req2)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestVotes() {
	goctx := sdk.WrapSDKContext(s.ctx)
	const expProposalID uint64 = 5
	_, err := s.msgServer.SuggestRole(goctx, &types.MsgSuggestRole{
		From:   s.guardians[0].String(),
		Target: s.guardians[1].String(),
		Role:   types.RoleJudge,
	})
	s.Require().NoError(err)
	_, err = s.msgServer.AddVoteForRole(goctx, &types.MsgAddVoteForRole{
		From:       s.guardians[0].String(),
		ProposalId: expProposalID,
		Option:     types.OptionYes,
	})
	s.Require().NoError(err)

	req := &types.QueryVotesRequest{
		ProposalId: expProposalID,
	}

	res, err := s.queryClient.Votes(goctx, req)
	s.Require().NoError(err)
	s.Require().Equal(expProposalID, res.Votes[0].ProposalId)
	s.Require().Equal(s.guardians[0].String(), res.Votes[0].Voter)
	s.Require().Equal(types.OptionYes, res.Votes[0].Option)

	req.ProposalId++
	res, err = s.queryClient.Votes(goctx, req)
	s.Require().NoError(err)
	s.Require().Empty(res.Votes)

	req2 := &types.QueryVoteRequest{
		ProposalId: expProposalID,
		Voter:      s.guardians[0].String(),
	}
	res2, err := s.queryClient.Vote(goctx, req2)
	s.Require().NoError(err)
	s.Require().Equal(expProposalID, res2.Vote.ProposalId)
	s.Require().Equal(s.guardians[0].String(), res2.Vote.Voter)
	s.Require().Equal(types.OptionYes, res2.Vote.Option)

	req2.ProposalId++
	_, err = s.queryClient.Vote(goctx, req2)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestBridgeStatus() {
	goctx := sdk.WrapSDKContext(s.ctx)
	res, err := s.queryClient.BridgeStatus(goctx, &types.QueryBridgeStatusRequest{})
	s.Require().NoError(err)
	s.Require().EqualValues(types.StatusActive, res.Status)
}
