package keeper_test

import (
	gocontext "context"
	"fmt"
	"strconv"

	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/gov/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryProposal() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	var (
		req         *types.QueryProposalRequest
		expProposal types.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryProposalRequest{}
			},
			false,
		},
		{
			"non existing proposal request",
			func() {
				req = &types.QueryProposalRequest{ProposalId: 3}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &types.QueryProposalRequest{ProposalId: 0}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryProposalRequest{ProposalId: 1}
				testProposal := types.NewTextProposal("Proposal", "testing proposal")
				submittedProposal, err := app.GovKeeper.SubmitProposal(ctx, testProposal)
				suite.Require().NoError(err)
				suite.Require().NotEmpty(submittedProposal)

				expProposal = submittedProposal
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			proposalRes, err := queryClient.Proposal(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expProposal.String(), proposalRes.Proposal.String())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(proposalRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryProposals() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	testProposals := []types.Proposal{}

	var (
		req    *types.QueryProposalsRequest
		expRes *types.QueryProposalsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty state request",
			func() {
				req = &types.QueryProposalsRequest{}
			},
			true,
		},
		{
			"request proposals with limit 3",
			func() {
				// create 5 test proposals
				for i := 0; i < 5; i++ {
					num := strconv.Itoa(i + 1)
					testProposal := types.NewTextProposal("Proposal"+num, "testing proposal "+num)
					proposal, err := app.GovKeeper.SubmitProposal(ctx, testProposal)
					suite.Require().NotEmpty(proposal)
					suite.Require().NoError(err)
					testProposals = append(testProposals, proposal)
				}

				req = &types.QueryProposalsRequest{
					Pagination: &query.PageRequest{Limit: 3},
				}

				expRes = &types.QueryProposalsResponse{
					Proposals: testProposals[:3],
				}
			},
			true,
		},
		{
			"request 2nd page with limit 4",
			func() {
				req = &types.QueryProposalsRequest{
					Pagination: &query.PageRequest{Offset: 3, Limit: 3},
				}

				expRes = &types.QueryProposalsResponse{
					Proposals: testProposals[3:],
				}
			},
			true,
		},
		{
			"request with limit 2 and count true",
			func() {
				req = &types.QueryProposalsRequest{
					Pagination: &query.PageRequest{Limit: 2, CountTotal: true},
				}

				expRes = &types.QueryProposalsResponse{
					Proposals: testProposals[:2],
				}
			},
			true,
		},
		{
			"request with filter of status deposit period",
			func() {
				req = &types.QueryProposalsRequest{
					ProposalStatus: types.StatusDepositPeriod,
				}

				expRes = &types.QueryProposalsResponse{
					Proposals: testProposals,
				}
			},
			true,
		},
		{
			"request with filter of deposit address",
			func() {
				depositCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, app.StakingKeeper.TokensFromConsensusPower(ctx, 20)))
				deposit := types.NewDeposit(testProposals[0].ProposalId, addrs[0], depositCoins)
				app.GovKeeper.SetDeposit(ctx, deposit)

				req = &types.QueryProposalsRequest{
					Depositor: addrs[0].String(),
				}

				expRes = &types.QueryProposalsResponse{
					Proposals: testProposals[:1],
				}
			},
			true,
		},
		{
			"request with filter of deposit address",
			func() {
				testProposals[1].Status = types.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, testProposals[1])
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, testProposals[1].ProposalId, addrs[0], types.NewNonSplitVoteOption(types.OptionAbstain)))

				req = &types.QueryProposalsRequest{
					Voter: addrs[0].String(),
				}

				expRes = &types.QueryProposalsResponse{
					Proposals: testProposals[1:2],
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			proposals, err := queryClient.Proposals(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)

				suite.Require().Len(proposals.GetProposals(), len(expRes.GetProposals()))
				for i := 0; i < len(proposals.GetProposals()); i++ {
					suite.Require().Equal(proposals.GetProposals()[i].String(), expRes.GetProposals()[i].String())
				}

			} else {
				suite.Require().Error(err)
				suite.Require().Nil(proposals)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryVote() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req      *types.QueryVoteRequest
		expRes   *types.QueryVoteResponse
		proposal types.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryVoteRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &types.QueryVoteRequest{
					ProposalId: 0,
					Voter:      addrs[0].String(),
				}
			},
			false,
		},
		{
			"empty voter request",
			func() {
				req = &types.QueryVoteRequest{
					ProposalId: 1,
					Voter:      "",
				}
			},
			false,
		},
		{
			"non existed proposal",
			func() {
				req = &types.QueryVoteRequest{
					ProposalId: 3,
					Voter:      addrs[0].String(),
				}
			},
			false,
		},
		{
			"no votes present",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)

				req = &types.QueryVoteRequest{
					ProposalId: proposal.ProposalId,
					Voter:      addrs[0].String(),
				}

				expRes = &types.QueryVoteResponse{}
			},
			false,
		},
		{
			"valid request",
			func() {
				proposal.Status = types.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, proposal)
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, addrs[0], types.NewNonSplitVoteOption(types.OptionAbstain)))

				req = &types.QueryVoteRequest{
					ProposalId: proposal.ProposalId,
					Voter:      addrs[0].String(),
				}

				expRes = &types.QueryVoteResponse{Vote: types.Vote{ProposalId: proposal.ProposalId, Voter: addrs[0].String(), Option: types.OptionAbstain, Options: []types.WeightedVoteOption{{Option: types.OptionAbstain, Weight: sdk.MustNewDecFromStr("1.0")}}}}
			},
			true,
		},
		{
			"wrong voter id request",
			func() {
				req = &types.QueryVoteRequest{
					ProposalId: proposal.ProposalId,
					Voter:      addrs[1].String(),
				}

				expRes = &types.QueryVoteResponse{}
			},
			false,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			vote, err := queryClient.Vote(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, vote)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(vote)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryVotes() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(30000000))

	var (
		req      *types.QueryVotesRequest
		expRes   *types.QueryVotesResponse
		proposal types.Proposal
		votes    types.Votes
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryVotesRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &types.QueryVotesRequest{
					ProposalId: 0,
				}
			},
			false,
		},
		{
			"non existed proposals",
			func() {
				req = &types.QueryVotesRequest{
					ProposalId: 2,
				}
			},
			true,
		},
		{
			"create a proposal and get votes",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)

				req = &types.QueryVotesRequest{
					ProposalId: proposal.ProposalId,
				}
			},
			true,
		},
		{
			"request after adding 2 votes",
			func() {
				proposal.Status = types.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, proposal)

				votes = []types.Vote{
					{ProposalId: proposal.ProposalId, Voter: addrs[0].String(), Option: types.OptionAbstain, Options: types.NewNonSplitVoteOption(types.OptionAbstain)},
					{ProposalId: proposal.ProposalId, Voter: addrs[1].String(), Option: types.OptionYes, Options: types.NewNonSplitVoteOption(types.OptionYes)},
				}
				accAddr1, err1 := sdk.AccAddressFromBech32(votes[0].Voter)
				accAddr2, err2 := sdk.AccAddressFromBech32(votes[1].Voter)
				suite.Require().NoError(err1)
				suite.Require().NoError(err2)
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, accAddr1, votes[0].Options))
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, accAddr2, votes[1].Options))

				req = &types.QueryVotesRequest{
					ProposalId: proposal.ProposalId,
				}

				expRes = &types.QueryVotesResponse{
					Votes: votes,
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			votes, err := queryClient.Votes(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetVotes(), votes.GetVotes())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(votes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	queryClient := suite.queryClient

	var (
		req    *types.QueryParamsRequest
		expRes *types.QueryParamsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryParamsRequest{}
			},
			false,
		},
		{
			"deposit params request",
			func() {
				req = &types.QueryParamsRequest{ParamsType: types.ParamDeposit}
				expRes = &types.QueryParamsResponse{
					DepositParams: types.DefaultDepositParams(),
					TallyParams:   types.NewTallyParams(sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0)),
				}
			},
			true,
		},
		{
			"voting params request",
			func() {
				req = &types.QueryParamsRequest{ParamsType: types.ParamVoting}
				expRes = &types.QueryParamsResponse{
					VotingParams: types.DefaultVotingParams(),
					TallyParams:  types.NewTallyParams(sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0)),
				}
			},
			true,
		},
		{
			"tally params request",
			func() {
				req = &types.QueryParamsRequest{ParamsType: types.ParamTallying}
				expRes = &types.QueryParamsResponse{
					TallyParams: types.DefaultTallyParams(),
				}
			},
			true,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryParamsRequest{ParamsType: "wrongPath"}
				expRes = &types.QueryParamsResponse{}
			},
			false,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			params, err := queryClient.Params(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetDepositParams(), params.GetDepositParams())
				suite.Require().Equal(expRes.GetVotingParams(), params.GetVotingParams())
				suite.Require().Equal(expRes.GetTallyParams(), params.GetTallyParams())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(params)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDeposit() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req      *types.QueryDepositRequest
		expRes   *types.QueryDepositResponse
		proposal types.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDepositRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &types.QueryDepositRequest{
					ProposalId: 0,
					Depositor:  addrs[0].String(),
				}
			},
			false,
		},
		{
			"empty deposit address request",
			func() {
				req = &types.QueryDepositRequest{
					ProposalId: 1,
					Depositor:  "",
				}
			},
			false,
		},
		{
			"non existed proposal",
			func() {
				req = &types.QueryDepositRequest{
					ProposalId: 2,
					Depositor:  addrs[0].String(),
				}
			},
			false,
		},
		{
			"no deposits proposal",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)
				suite.Require().NotNil(proposal)

				req = &types.QueryDepositRequest{
					ProposalId: proposal.ProposalId,
					Depositor:  addrs[0].String(),
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				depositCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, app.StakingKeeper.TokensFromConsensusPower(ctx, 20)))
				deposit := types.NewDeposit(proposal.ProposalId, addrs[0], depositCoins)
				app.GovKeeper.SetDeposit(ctx, deposit)

				req = &types.QueryDepositRequest{
					ProposalId: proposal.ProposalId,
					Depositor:  addrs[0].String(),
				}

				expRes = &types.QueryDepositResponse{Deposit: deposit}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			deposit, err := queryClient.Deposit(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(deposit.GetDeposit(), expRes.GetDeposit())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(expRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDeposits() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req      *types.QueryDepositsRequest
		expRes   *types.QueryDepositsResponse
		proposal types.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDepositsRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &types.QueryDepositsRequest{
					ProposalId: 0,
				}
			},
			false,
		},
		{
			"non existed proposal",
			func() {
				req = &types.QueryDepositsRequest{
					ProposalId: 2,
				}
			},
			true,
		},
		{
			"create a proposal and get deposits",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)

				req = &types.QueryDepositsRequest{
					ProposalId: proposal.ProposalId,
				}
			},
			true,
		},
		{
			"get deposits with default limit",
			func() {
				depositAmount1 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, app.StakingKeeper.TokensFromConsensusPower(ctx, 20)))
				deposit1 := types.NewDeposit(proposal.ProposalId, addrs[0], depositAmount1)
				app.GovKeeper.SetDeposit(ctx, deposit1)

				depositAmount2 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, app.StakingKeeper.TokensFromConsensusPower(ctx, 30)))
				deposit2 := types.NewDeposit(proposal.ProposalId, addrs[1], depositAmount2)
				app.GovKeeper.SetDeposit(ctx, deposit2)

				deposits := types.Deposits{deposit1, deposit2}

				req = &types.QueryDepositsRequest{
					ProposalId: proposal.ProposalId,
				}

				expRes = &types.QueryDepositsResponse{
					Deposits: deposits,
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			deposits, err := queryClient.Deposits(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetDeposits(), deposits.GetDeposits())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(deposits)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryTally() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	addrs, _ := createValidators(suite.T(), ctx, app, []int64{5, 5, 5})

	var (
		req      *types.QueryTallyResultRequest
		expRes   *types.QueryTallyResultResponse
		proposal types.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryTallyResultRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &types.QueryTallyResultRequest{ProposalId: 0}
			},
			false,
		},
		{
			"query non existed proposal",
			func() {
				req = &types.QueryTallyResultRequest{ProposalId: 1}
			},
			false,
		},
		{
			"create a proposal and get tally",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)
				suite.Require().NotNil(proposal)

				req = &types.QueryTallyResultRequest{ProposalId: proposal.ProposalId}

				expRes = &types.QueryTallyResultResponse{
					Tally: types.EmptyTallyResult(),
				}
			},
			true,
		},
		{
			"request tally after few votes",
			func() {
				proposal.Status = types.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, proposal)

				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, addrs[1], types.NewNonSplitVoteOption(types.OptionYes)))
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, addrs[2], types.NewNonSplitVoteOption(types.OptionYes)))

				req = &types.QueryTallyResultRequest{ProposalId: proposal.ProposalId}

				expRes = &types.QueryTallyResultResponse{
					Tally: types.TallyResult{
						Yes: sdk.NewInt(3 * 5 * 1000000),
					},
				}
			},
			true,
		},
		{
			"request final tally after status changed",
			func() {
				proposal.Status = types.StatusPassed
				app.GovKeeper.SetProposal(ctx, proposal)
				proposal, _ = app.GovKeeper.GetProposal(ctx, proposal.ProposalId)

				req = &types.QueryTallyResultRequest{ProposalId: proposal.ProposalId}

				expRes = &types.QueryTallyResultResponse{
					Tally: proposal.FinalTallyResult,
				}
			},
			true,
		},
		{
			"proposal status failed",
			func() {
				propTime := time.Now()
				proposal := v1.Proposal{
					Id:     1,
					Status: v1.StatusFailed,
					FinalTallyResult: &v1.TallyResult{
						YesCount:         "4",
						AbstainCount:     "1",
						NoCount:          "0",
						NoWithVetoCount:  "0",
						OptionOneCount:   "4",
						OptionTwoCount:   "1",
						OptionThreeCount: "0",
						OptionFourCount:  "0",
						SpamCount:        "0",
					},
					SubmitTime:      &propTime,
					VotingStartTime: &propTime,
					VotingEndTime:   &propTime,
					Metadata:        "proposal metadata",
				}
				err := suite.govKeeper.Proposals.Set(suite.ctx, proposal.Id, proposal)
				suite.Require().NoError(err)

				req = &v1.QueryTallyResultRequest{ProposalId: proposal.Id}

				expTally = &v1.TallyResult{
					YesCount:         "4",
					AbstainCount:     "1",
					NoCount:          "0",
					NoWithVetoCount:  "0",
					OptionOneCount:   "4",
					OptionTwoCount:   "1",
					OptionThreeCount: "0",
					OptionFourCount:  "0",
					SpamCount:        "0",
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			tally, err := queryClient.TallyResult(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.String(), tally.String())
			} else {
				suite.Require().Error(err)
<<<<<<< HEAD
				suite.Require().Nil(tally)
=======
				suite.Require().Nil(tallyRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestLegacyGRPCQueryTallyResult() {
	suite.reset()
	queryClient := suite.legacyQueryClient

	var (
		req      *v1beta1.QueryTallyResultRequest
		expTally *v1beta1.TallyResult
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &v1beta1.QueryTallyResultRequest{}
			},
			false,
		},
		{
			"non existing proposal request",
			func() {
				req = &v1beta1.QueryTallyResultRequest{ProposalId: 2}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &v1beta1.QueryTallyResultRequest{ProposalId: 0}
			},
			false,
		},
		{
			"valid request with proposal status passed",
			func() {
				propTime := time.Now()
				proposal := v1.Proposal{
					Id:     1,
					Status: v1.StatusPassed,
					FinalTallyResult: &v1.TallyResult{
						YesCount:        "4",
						AbstainCount:    "1",
						NoCount:         "0",
						NoWithVetoCount: "0",
						SpamCount:       "0",
					},
					SubmitTime:      &propTime,
					VotingStartTime: &propTime,
					VotingEndTime:   &propTime,
					Metadata:        "proposal metadata",
				}
				err := suite.govKeeper.Proposals.Set(suite.ctx, proposal.Id, proposal)
				suite.Require().NoError(err)

				req = &v1beta1.QueryTallyResultRequest{ProposalId: proposal.Id}

				expTally = &v1beta1.TallyResult{
					Yes:        math.NewInt(4),
					Abstain:    math.NewInt(1),
					No:         math.NewInt(0),
					NoWithVeto: math.NewInt(0),
				}
			},
			true,
		},
		{
			"proposal status deposit",
			func() {
				propTime := time.Now()
				proposal := v1.Proposal{
					Id:              1,
					Status:          v1.StatusDepositPeriod,
					SubmitTime:      &propTime,
					VotingStartTime: &propTime,
					VotingEndTime:   &propTime,
					Metadata:        "proposal metadata",
				}
				err := suite.govKeeper.Proposals.Set(suite.ctx, proposal.Id, proposal)
				suite.Require().NoError(err)

				req = &v1beta1.QueryTallyResultRequest{ProposalId: proposal.Id}

				expTally = &v1beta1.TallyResult{
					Yes:        math.NewInt(0),
					Abstain:    math.NewInt(0),
					No:         math.NewInt(0),
					NoWithVeto: math.NewInt(0),
				}
			},
			true,
		},
		{
			"proposal is in voting period",
			func() {
				propTime := time.Now()
				proposal := v1.Proposal{
					Id:              1,
					Status:          v1.StatusVotingPeriod,
					SubmitTime:      &propTime,
					VotingStartTime: &propTime,
					VotingEndTime:   &propTime,
					Metadata:        "proposal metadata",
				}
				err := suite.govKeeper.Proposals.Set(suite.ctx, proposal.Id, proposal)
				suite.Require().NoError(err)
				req = &v1beta1.QueryTallyResultRequest{ProposalId: proposal.Id}

				expTally = &v1beta1.TallyResult{
					Yes:        math.NewInt(0),
					Abstain:    math.NewInt(0),
					No:         math.NewInt(0),
					NoWithVeto: math.NewInt(0),
				}
			},
			true,
		},
		{
			"proposal status failed",
			func() {
				propTime := time.Now()
				proposal := v1.Proposal{
					Id:     1,
					Status: v1.StatusFailed,
					FinalTallyResult: &v1.TallyResult{
						YesCount:        "4",
						AbstainCount:    "1",
						NoCount:         "0",
						NoWithVetoCount: "0",
						SpamCount:       "0",
					},
					SubmitTime:      &propTime,
					VotingStartTime: &propTime,
					VotingEndTime:   &propTime,
					Metadata:        "proposal metadata",
				}
				err := suite.govKeeper.Proposals.Set(suite.ctx, proposal.Id, proposal)
				suite.Require().NoError(err)

				req = &v1beta1.QueryTallyResultRequest{ProposalId: proposal.Id}

				expTally = &v1beta1.TallyResult{
					Yes:        math.NewInt(4),
					Abstain:    math.NewInt(1),
					No:         math.NewInt(0),
					NoWithVeto: math.NewInt(0),
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			if tc.malleate != nil {
				tc.malleate()
			}

			tallyRes, err := queryClient.TallyResult(gocontext.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotEmpty(tallyRes.Tally.String())
				suite.Require().Equal(expTally.String(), tallyRes.Tally.String())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(tallyRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestProposalVoteOptions() {
	suite.reset()

	testCases := []struct {
		name     string
		malleate func()
		req      *v1.QueryProposalVoteOptionsRequest
		expResp  *v1.QueryProposalVoteOptionsResponse
		errStr   string
	}{
		{
			name:   "invalid proposal id",
			req:    &v1.QueryProposalVoteOptionsRequest{},
			errStr: "proposal id can not be 0",
		},
		{
			name:   "proposal not found",
			req:    &v1.QueryProposalVoteOptionsRequest{ProposalId: 1},
			errStr: "proposal 1 doesn't exist",
		},
		{
			name: "non multiple choice proposal",
			malleate: func() {
				propTime := time.Now()
				proposal := v1.Proposal{
					Id:              1,
					Status:          v1.StatusVotingPeriod,
					SubmitTime:      &propTime,
					VotingStartTime: &propTime,
					VotingEndTime:   &propTime,
					Metadata:        "proposal metadata",
					ProposalType:    v1.ProposalType_PROPOSAL_TYPE_STANDARD,
				}
				err := suite.govKeeper.Proposals.Set(suite.ctx, proposal.Id, proposal)
				suite.Require().NoError(err)
			},
			req: &v1.QueryProposalVoteOptionsRequest{ProposalId: 1},
			expResp: &v1.QueryProposalVoteOptionsResponse{
				VoteOptions: &v1.ProposalVoteOptions{
					OptionOne:   "yes",
					OptionTwo:   "abstain",
					OptionThree: "no",
					OptionFour:  "no_with_veto",
					OptionSpam:  "spam",
				},
			},
		},
		{
			name: "invalid multiple choice proposal",
			req:  &v1.QueryProposalVoteOptionsRequest{ProposalId: 2},
			malleate: func() {
				propTime := time.Now()
				proposal := v1.Proposal{
					Id:              2,
					Status:          v1.StatusVotingPeriod,
					SubmitTime:      &propTime,
					VotingStartTime: &propTime,
					VotingEndTime:   &propTime,
					Metadata:        "proposal metadata",
					ProposalType:    v1.ProposalType_PROPOSAL_TYPE_MULTIPLE_CHOICE,
				}
				err := suite.govKeeper.Proposals.Set(suite.ctx, proposal.Id, proposal)
				suite.Require().NoError(err)

				// multiple choice proposal, but no vote options set
				// because the query does not check the proposal type,
				// it falls back to the default vote options
			},
			expResp: &v1.QueryProposalVoteOptionsResponse{
				VoteOptions: &v1.ProposalVoteOptions{
					OptionOne:   "yes",
					OptionTwo:   "abstain",
					OptionThree: "no",
					OptionFour:  "no_with_veto",
					OptionSpam:  "spam",
				},
			},
		},
		{
			name: "multiple choice proposal",
			req:  &v1.QueryProposalVoteOptionsRequest{ProposalId: 3},
			malleate: func() {
				propTime := time.Now()
				proposal := v1.Proposal{
					Id:              3,
					Status:          v1.StatusVotingPeriod,
					SubmitTime:      &propTime,
					VotingStartTime: &propTime,
					VotingEndTime:   &propTime,
					Metadata:        "proposal metadata",
					ProposalType:    v1.ProposalType_PROPOSAL_TYPE_MULTIPLE_CHOICE,
				}
				err := suite.govKeeper.Proposals.Set(suite.ctx, proposal.Id, proposal)
				suite.Require().NoError(err)
				err = suite.govKeeper.ProposalVoteOptions.Set(suite.ctx, proposal.Id, v1.ProposalVoteOptions{
					OptionOne:   "Vote for @tac0turle",
					OptionTwo:   "Vote for @facudomedica",
					OptionThree: "Vote for @alexanderbez",
				})
				suite.Require().NoError(err)
			},
			expResp: &v1.QueryProposalVoteOptionsResponse{
				VoteOptions: &v1.ProposalVoteOptions{
					OptionOne:   "Vote for @tac0turle",
					OptionTwo:   "Vote for @facudomedica",
					OptionThree: "Vote for @alexanderbez",
					OptionSpam:  "spam",
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.malleate != nil {
				tc.malleate()
			}

			resp, err := suite.queryClient.ProposalVoteOptions(suite.ctx, tc.req)
			if tc.errStr != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errStr)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResp, resp)
>>>>>>> d961aef76b (fix(x/gov): grpc query tally for failed proposal (#19725))
			}
		})
	}
}
