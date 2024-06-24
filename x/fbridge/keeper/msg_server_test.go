package keeper_test

import (
	"fmt"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (s *IntegrationTestSuite) TestInactiveTxs() {
	goctx := sdk.WrapSDKContext(s.ctx)

	s.Require().Panics(func() {
		_, _ = s.msgServer.Provision(goctx, &types.MsgProvision{})
	})

	s.Require().Panics(func() {
		_, _ = s.msgServer.HoldTransfer(goctx, &types.MsgHoldTransfer{})
	})

	s.Require().Panics(func() {
		_, _ = s.msgServer.ReleaseTransfer(goctx, &types.MsgReleaseTransfer{})
	})

	s.Require().Panics(func() {
		_, _ = s.msgServer.RemoveProvision(goctx, &types.MsgRemoveProvision{})
	})

	s.Require().Panics(func() {
		_, _ = s.msgServer.ClaimBatch(goctx, &types.MsgClaimBatch{})
	})

	s.Require().Panics(func() {
		_, _ = s.msgServer.Claim(goctx, &types.MsgClaim{})
	})
}

func (s *IntegrationTestSuite) TestUpdateParams() {
	tcs := map[string]struct {
		msg    types.MsgUpdateParams
		expErr bool
	}{
		"valid request": {
			msg: types.MsgUpdateParams{
				Authority: types.DefaultAuthority().String(),
				Params:    types.DefaultParams(),
			},
			expErr: false,
		},
		"invalid authority": {
			msg: types.MsgUpdateParams{
				Authority: "invalid",
				Params:    types.DefaultParams(),
			},
			expErr: true,
		},
		"invalid params": {
			msg: types.MsgUpdateParams{
				Authority: types.DefaultAuthority().String(),
				Params:    types.Params{},
			},
			expErr: true,
		},
	}

	goctx := sdk.WrapSDKContext(s.ctx)
	for name, tc := range tcs {
		s.Run(name, func() {
			_, err := s.msgServer.UpdateParams(goctx, &tc.msg)
			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestTransfer() {
	var msg types.MsgTransfer
	tcs := map[string]struct {
		malleate func()
		postExec func()
		expErr   bool
	}{
		"valid request": {
			malleate: func() {
				msg = types.MsgTransfer{
					Sender:   s.guardians[0].String(),
					Receiver: "0x1A7C26B0437Aa2d3c8454383650a5D3c35087f91",
					Amount:   sdk.NewInt(100),
				}
			},
			expErr: false,
		},
		"invalid sender": {
			malleate: func() {
				msg = types.MsgTransfer{
					Sender:   "invalid",
					Receiver: "0x1A7C26B0437Aa2d3c8454383650a5D3c35087f91",
					Amount:   sdk.NewInt(100),
				}
			},
			expErr: true,
		},
		"invalid receiver": {
			malleate: func() {
				msg = types.MsgTransfer{
					Sender:   s.guardians[0].String(),
					Receiver: "invalid",
					Amount:   sdk.NewInt(100),
				}
			},
			expErr: true,
		},
		"insufficient balance": {
			malleate: func() {
				msg = types.MsgTransfer{
					Sender:   s.guardians[0].String(),
					Receiver: "0x1A7C26B0437Aa2d3c8454383650a5D3c35087f91",
					Amount:   sdk.NewInt(int64(0x7FFFFFFFFFFFFFFF)),
				}
			},
			expErr: true,
		},
		"bridge halted": {
			malleate: func() {
				msg = types.MsgTransfer{
					Sender:   s.guardians[0].String(),
					Receiver: "0x1A7C26B0437Aa2d3c8454383650a5D3c35087f91",
					Amount:   sdk.NewInt(100),
				}

				_, err := s.msgServer.SetBridgeStatus(sdk.WrapSDKContext(s.ctx), &types.MsgSetBridgeStatus{
					Guardian: s.guardians[0].String(),
					Status:   types.StatusInactive,
				})
				s.Require().NoError(err)
				_, err = s.msgServer.SetBridgeStatus(sdk.WrapSDKContext(s.ctx), &types.MsgSetBridgeStatus{
					Guardian: s.guardians[1].String(),
					Status:   types.StatusInactive,
				})
				s.Require().NoError(err)
			},
			postExec: func() {
				_, err := s.msgServer.SetBridgeStatus(sdk.WrapSDKContext(s.ctx), &types.MsgSetBridgeStatus{
					Guardian: s.guardians[0].String(),
					Status:   types.StatusActive,
				})
				s.Require().NoError(err)
				_, err = s.msgServer.SetBridgeStatus(sdk.WrapSDKContext(s.ctx), &types.MsgSetBridgeStatus{
					Guardian: s.guardians[1].String(),
					Status:   types.StatusActive,
				})
				s.Require().NoError(err)
			},
			expErr: true,
		},
	}

	goctx := sdk.WrapSDKContext(s.ctx)
	for name, tc := range tcs {
		s.Run(name, func() {
			tc.malleate()
			_, err := s.msgServer.Transfer(goctx, &msg)
			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

			if tc.postExec != nil {
				tc.postExec()
			}
		})
	}
}

func (s *IntegrationTestSuite) TestSuggestRole() {
	var msg types.MsgSuggestRole
	tcs := map[string]struct {
		malleate func()
		expErr   bool
	}{
		"valid request": {
			malleate: func() {
				msg = types.MsgSuggestRole{
					From:   s.guardians[0].String(),
					Target: s.guardians[2].String(),
					Role:   types.RoleOperator,
				}
			},
			expErr: false,
		},
		"invalid proposer": {
			malleate: func() {
				msg = types.MsgSuggestRole{
					From:   "invalid",
					Target: s.guardians[2].String(),
					Role:   types.RoleOperator,
				}
			},
			expErr: true,
		},
		"invalid target address": {
			malleate: func() {
				msg = types.MsgSuggestRole{
					From:   s.guardians[0].String(),
					Target: "invalid",
					Role:   types.RoleOperator,
				}
			},
			expErr: true,
		},
		"unsupported role": {
			malleate: func() {
				msg = types.MsgSuggestRole{
					From:   s.guardians[0].String(),
					Target: s.guardians[1].String(),
					Role:   types.Role(10),
				}
			},
			expErr: true,
		},
		"target already has same role": {
			malleate: func() {
				msg = types.MsgSuggestRole{
					From:   s.guardians[0].String(),
					Target: s.guardians[1].String(),
					Role:   types.RoleGuardian,
				}
			},
			expErr: true,
		},
	}

	goctx := sdk.WrapSDKContext(s.ctx)
	for name, tc := range tcs {
		s.Run(name, func() {
			tc.malleate()
			_, err := s.msgServer.SuggestRole(goctx, &msg)
			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestAddVoteForRole() {
	goctx := sdk.WrapSDKContext(s.ctx)
	_, err := s.msgServer.SuggestRole(goctx, &types.MsgSuggestRole{
		From:   s.guardians[0].String(),
		Target: s.guardians[2].String(),
		Role:   types.RoleOperator,
	})
	s.Require().NoError(err)
	const proposalID = 5

	var msg types.MsgAddVoteForRole
	tcs := map[string]struct {
		malleate func()
		expErr   bool
	}{
		"valid request": {
			malleate: func() {
				msg = types.MsgAddVoteForRole{
					From:       s.guardians[0].String(),
					ProposalId: proposalID,
					Option:     types.OptionYes,
				}
			},
			expErr: false,
		},
		"invalid voter": {
			malleate: func() {
				msg = types.MsgAddVoteForRole{
					From:       "invalid",
					ProposalId: 0,
					Option:     types.OptionYes,
				}
			},
			expErr: true,
		},
		"unauthorized voter": {
			malleate: func() {
				msg = types.MsgAddVoteForRole{
					From:       s.operator.String(),
					ProposalId: proposalID,
					Option:     types.OptionYes,
				}
			},
			expErr: true,
		},
		"invalid proposal id": {
			malleate: func() {
				msg = types.MsgAddVoteForRole{
					From:       s.guardians[0].String(),
					ProposalId: 100,
					Option:     types.OptionYes,
				}
			},
			expErr: true,
		},
		"invalid option": {
			malleate: func() {
				msg = types.MsgAddVoteForRole{
					From:       s.guardians[0].String(),
					ProposalId: 0,
					Option:     types.VoteOption(10),
				}
			},
			expErr: true,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			tc.malleate()
			_, err := s.msgServer.AddVoteForRole(goctx, &msg)
			if tc.expErr {
				s.Require().Error(err)
				fmt.Println(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestSetBridgeStatus() {
	var msg types.MsgSetBridgeStatus
	tcs := map[string]struct {
		malleate func()
		expErr   bool
	}{
		"1. valid request - inactive": {
			malleate: func() {
				msg = types.MsgSetBridgeStatus{
					Guardian: s.guardians[0].String(),
					Status:   types.StatusInactive,
				}
			},
			expErr: false,
		},
		"2. valid request - active": {
			malleate: func() {
				msg = types.MsgSetBridgeStatus{
					Guardian: s.guardians[0].String(),
					Status:   types.StatusActive,
				}
			},
			expErr: false,
		},
		"invalid guardian address": {
			malleate: func() {
				msg = types.MsgSetBridgeStatus{
					Guardian: "invalid",
					Status:   types.StatusInactive,
				}
			},
			expErr: true,
		},
		"equal to current status": {
			malleate: func() {
				msg = types.MsgSetBridgeStatus{
					Guardian: s.guardians[1].String(),
					Status:   types.StatusActive,
				}
			},
			expErr: true,
		},
		"invalid bridge status": {
			malleate: func() {
				msg = types.MsgSetBridgeStatus{
					Guardian: s.guardians[1].String(),
					Status:   types.BridgeStatus(10),
				}
			},
			expErr: true,
		},
	}

	goctx := sdk.WrapSDKContext(s.ctx)
	for name, tc := range tcs {
		s.Run(name, func() {
			tc.malleate()
			_, err := s.msgServer.SetBridgeStatus(goctx, &msg)
			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
