package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-sdk/testutil"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/token"
	"github.com/Finschia/finschia-sdk/x/token/class"
)

func (s *KeeperTestSuite) TestMsgSend() {
	testCases := map[string]struct {
		contractID string
		amount     sdk.Int
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     s.balance,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.W(s.balance), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.vendor), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "fee1dead",
			amount:     sdk.OneInt(),
			err:        class.ErrContractNotExist,
		},
		"insufficient funds": {
			contractID: s.contractID,
			amount:     s.balance.Add(sdk.OneInt()),
			err:        token.ErrInsufficientBalance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgSend{
				ContractId: tc.contractID,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     tc.amount,
			}
			res, err := s.msgServer.Send(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorSend() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		amount     sdk.Int
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount:     s.balance,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.W(s.balance), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.operator), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "fee1dead",
			operator:   s.operator,
			from:       s.customer,
			amount:     s.balance,
			err:        class.ErrContractNotExist,
		},
		"not approved": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			amount:     s.balance,
			err:        token.ErrTokenNotApproved,
		},
		"insufficient funds": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount:     s.balance.Add(sdk.OneInt()),
			err:        token.ErrInsufficientBalance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgOperatorSend{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				Amount:     tc.amount,
			}
			res, err := s.msgServer.OperatorSend(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokeOperator() {
	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.operator,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventRevokedOperator",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("holder"), Value: testutil.W(s.customer), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.operator), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "fee1dead",
			holder:     s.customer,
			operator:   s.operator,
			err:        class.ErrContractNotExist,
		},
		"no authorization": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.vendor,
			err:        token.ErrTokenNotApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgRevokeOperator{
				ContractId: tc.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}
			res, err := s.msgServer.RevokeOperator(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgAuthorizeOperator() {
	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventAuthorizedOperator",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("holder"), Value: testutil.W(s.customer), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "fee1dead",
			holder:     s.customer,
			operator:   s.vendor,
			err:        class.ErrContractNotExist,
		},
		"already approved": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.operator,
			err:        token.ErrTokenAlreadyApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgAuthorizeOperator{
				ContractId: tc.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}
			res, err := s.msgServer.AuthorizeOperator(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssue() {
	testCases := map[string]struct {
		mintable bool
		amount   sdk.Int
		err      error
		events   sdk.Events
	}{
		"mintable true": {
			mintable: true,
			amount:   sdk.NewInt(10),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventIssued",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("ca8bfd79"), Index: false},
						{Key: []uint8("creator"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("decimals"), Value: []byte("0"), Index: false},
						{Key: []uint8("meta"), Value: testutil.W(""), Index: false},
						{Key: []uint8("mintable"), Value: []byte("true"), Index: false},
						{Key: []uint8("name"), Value: testutil.W("test"), Index: false},
						{Key: []uint8("symbol"), Value: testutil.W("TT"), Index: false},
						{Key: []uint8("uri"), Value: testutil.W(""), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.token.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("ca8bfd79"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("granter"), Value: testutil.W(""), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_MODIFY"), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.token.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("ca8bfd79"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("granter"), Value: testutil.W(""), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_MINT"), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.token.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("ca8bfd79"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("granter"), Value: testutil.W(""), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_BURN"), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.token.v1.EventMinted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: testutil.W("10"), Index: false},
						{Key: []uint8("contract_id"), Value: testutil.W("ca8bfd79"), Index: false},
						{Key: []uint8("operator"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("to"), Value: testutil.W(s.vendor), Index: false},
					},
				},
			},
		},
		"mintable false": {
			mintable: false,
			amount:   sdk.NewInt(10),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventIssued",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("ca8bfd79"), Index: false},
						{Key: []uint8("creator"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("decimals"), Value: []byte("0"), Index: false},
						{Key: []uint8("meta"), Value: testutil.W(""), Index: false},
						{Key: []uint8("mintable"), Value: []byte("false"), Index: false},
						{Key: []uint8("name"), Value: testutil.W("test"), Index: false},
						{Key: []uint8("symbol"), Value: testutil.W("TT"), Index: false},
						{Key: []uint8("uri"), Value: testutil.W(""), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.token.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("ca8bfd79"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("granter"), Value: testutil.W(""), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_MODIFY"), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.token.v1.EventMinted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: testutil.W(sdk.NewInt(10)), Index: false},
						{Key: []uint8("contract_id"), Value: testutil.W("ca8bfd79"), Index: false},
						{Key: []uint8("operator"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("to"), Value: testutil.W(s.vendor), Index: false},
					},
				},
			},
		},
	}

	// define a function to check MsgIssue result
	checkerIssueResult := func(ctx sdk.Context, contractId string, expectedMintable bool, expectedAmount sdk.Int) {
		// check contract
		contract, err := s.queryServer.Contract(sdk.WrapSDKContext(ctx), &token.QueryContractRequest{ContractId: contractId})
		s.Require().NoError(err)
		s.Require().Equal(expectedMintable, contract.Contract.Mintable)

		// check supply
		supply, err := s.queryServer.Supply(sdk.WrapSDKContext(ctx), &token.QuerySupplyRequest{ContractId: contractId})
		s.Require().NoError(err)
		s.Require().Equal(expectedAmount, supply.Amount)

		// check mint
		mint, err := s.queryServer.Minted(sdk.WrapSDKContext(ctx), &token.QueryMintedRequest{ContractId: contractId})
		s.Require().NoError(err)
		s.Require().Equal(expectedAmount, mint.Amount)

		// check burnt
		burn, err := s.queryServer.Burnt(sdk.WrapSDKContext(ctx), &token.QueryBurntRequest{ContractId: contractId})
		s.Require().NoError(err)
		s.Require().Equal(sdk.ZeroInt(), burn.Amount)

		// check owner balance
		balance, err := s.queryServer.Balance(sdk.WrapSDKContext(ctx), &token.QueryBalanceRequest{
			ContractId: contractId,
			Address:    s.vendor.String(),
		})
		s.Require().NoError(err)
		s.Require().Equal(expectedAmount, balance.Amount)
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgIssue{
				Owner:    s.vendor.String(),
				To:       s.vendor.String(),
				Mintable: tc.mintable,
				Name:     "test",
				Symbol:   "TT",
				Amount:   tc.amount,
			}
			res, err := s.msgServer.Issue(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())

			// check result status
			checkerIssueResult(ctx, res.ContractId, tc.mintable, tc.amount)

			// Second request for the same request
			res2, err := s.msgServer.Issue(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			// check result status
			checkerIssueResult(ctx, res2.ContractId, tc.mintable, tc.amount)
			s.Require().NotEqual(res.ContractId, res2.ContractId)
		})
	}
}

func (s *KeeperTestSuite) TestMsgGrantPermission() {
	testCases := map[string]struct {
		contractID string
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission string
		err        error
		events     sdk.Events
	}{
		"contract not found": {
			contractID: "fee1dead",
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
			err:        class.ErrContractNotExist,
		},
		"contract has no permission - MINT": {
			contractID: s.unmintableContractId,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionMint.String(),
			err:        token.ErrTokenNoPermission,
		},
		"contract has no permission - BURN": {
			contractID: s.unmintableContractId,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionBurn.String(),
			err:        token.ErrTokenNoPermission,
		},
		"granter has no permission - MINT": {
			contractID: s.contractID,
			granter:    s.customer,
			grantee:    s.stranger,
			permission: token.LegacyPermissionMint.String(),
			err:        token.ErrTokenNoPermission,
		},
		"granter has no permission - BURN": {
			contractID: s.contractID,
			granter:    s.customer,
			grantee:    s.stranger,
			permission: token.LegacyPermissionBurn.String(),
			err:        token.ErrTokenNoPermission,
		},
		"granter has no permission - MODIFY": {
			contractID: s.contractID,
			granter:    s.customer,
			grantee:    s.stranger,
			permission: token.LegacyPermissionModify.String(),
			err:        token.ErrTokenNoPermission,
		},
		"valid request - MINT": {
			contractID: s.contractID,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionMint.String(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("9be17165"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.operator), Index: false},
						{Key: []uint8("granter"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_MINT"), Index: false},
					},
				},
			},
		},
		"valid request - BURN": {
			contractID: s.contractID,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionBurn.String(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("9be17165"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.operator), Index: false},
						{Key: []uint8("granter"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_BURN"), Index: false},
					},
				},
			},
		},
		"valid request - MODIFY": {
			contractID: s.contractID,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("9be17165"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.operator), Index: false},
						{Key: []uint8("granter"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_MODIFY"), Index: false},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgGrantPermission{
				ContractId: tc.contractID,
				From:       tc.granter.String(),
				To:         tc.grantee.String(),
				Permission: tc.permission,
			}
			res, err := s.msgServer.GrantPermission(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())

			// check to grant permission
			per, err := s.queryServer.GranteeGrants(sdk.WrapSDKContext(ctx), &token.QueryGranteeGrantsRequest{
				ContractId: tc.contractID,
				Grantee:    tc.grantee.String(),
				Pagination: nil,
			})
			s.Require().NoError(err)
			s.Require().NotNil(per)
			expectPermission := token.Grant{
				Grantee:    tc.grantee.String(),
				Permission: token.Permission(token.LegacyPermissionFromString(tc.permission)),
			}
			s.Require().Contains(per.Grants, expectPermission)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokePermission() {
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		permission string
		err        error
		events     sdk.Events
	}{
		"contract not found": {
			contractID: "fee1dead",
			from:       s.operator,
			permission: token.LegacyPermissionMint.String(),
			err:        class.ErrContractNotExist,
		},
		"contract has no permission - MINT": {
			contractID: s.unmintableContractId,
			from:       s.operator,
			permission: token.LegacyPermissionMint.String(),
			err:        token.ErrTokenNoPermission,
		},
		"contract has no permission - BURN": {
			contractID: s.unmintableContractId,
			from:       s.operator,
			permission: token.LegacyPermissionBurn.String(),
			err:        token.ErrTokenNoPermission,
		},
		"grantee has no permission - MINT": {
			contractID: s.contractID,
			from:       s.customer,
			permission: token.LegacyPermissionMint.String(),
			err:        token.ErrTokenNoPermission,
		},
		"grantee has no permission - BURN": {
			contractID: s.contractID,
			from:       s.customer,
			permission: token.LegacyPermissionBurn.String(),
			err:        token.ErrTokenNoPermission,
		},
		"grantee has no permission - MODIFY": {
			contractID: s.contractID,
			from:       s.customer,
			permission: token.LegacyPermissionModify.String(),
			err:        token.ErrTokenNoPermission,
		},
		"valid request - revoke MINT": {
			contractID: s.contractID,
			from:       s.operator,
			permission: token.LegacyPermissionMint.String(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventRenounced",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("9be17165"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.operator), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_MINT"), Index: false},
					},
				},
			},
		},
		"valid request - revoke BURN": {
			contractID: s.contractID,
			from:       s.operator,
			permission: token.LegacyPermissionBurn.String(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventRenounced",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("9be17165"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.operator), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_BURN"), Index: false},
					},
				},
			},
		},
		"valid request - revoke MODIFY": {
			contractID: s.contractID,
			from:       s.vendor,
			permission: token.LegacyPermissionModify.String(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventRenounced",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("contract_id"), Value: testutil.W("9be17165"), Index: false},
						{Key: []uint8("grantee"), Value: testutil.W(s.vendor), Index: false},
						{Key: []uint8("permission"), Value: testutil.W("PERMISSION_MODIFY"), Index: false},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgRevokePermission{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Permission: tc.permission,
			}
			res, err := s.msgServer.RevokePermission(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())

			// check to remove permission
			per, err := s.queryServer.GranteeGrants(sdk.WrapSDKContext(ctx), &token.QueryGranteeGrantsRequest{
				ContractId: tc.contractID,
				Grantee:    tc.from.String(),
				Pagination: nil,
			})
			s.Require().NoError(err)
			s.Require().NotNil(per)
			expectPermission := token.Grant{
				Grantee:    tc.from.String(),
				Permission: token.Permission(token.LegacyPermissionFromString(tc.permission)),
			}
			s.Require().NotContains(per.Grants, expectPermission)
		})
	}
}

func (s *KeeperTestSuite) TestMsgMint() {
	testCases := map[string]struct {
		isNegativeCase bool
		req            *token.MsgMint
		expectedEvents sdk.Events
		expectedError  *sdkerrors.Error
	}{
		"mint(contractID, from, to, 10)": {
			req: &token.MsgMint{
				ContractId: s.contractID,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     sdk.NewInt(10),
			},
			expectedEvents: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventMinted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.W(sdk.NewInt(10)), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer), Index: false},
					},
				},
			},
		},
		"mint(contractID, from, from, 10)": {
			req: &token.MsgMint{
				ContractId: s.contractID,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     sdk.NewInt(10),
			},
			expectedEvents: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventMinted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.W(sdk.NewInt(10)), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer), Index: false},
					},
				},
			},
		},
		"mint(contractID, vendor, customer, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgMint{
				ContractId: s.unmintableContractId,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     sdk.OneInt(),
			},
			expectedError: token.ErrTokenNoPermission,
		},
		"mint(nonExistingContractId, from, to, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgMint{
				ContractId: "fee1dead",
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     sdk.OneInt(),
			},
			expectedError: class.ErrContractNotExist,
		},
		"mint(contractID, from, unauthorized account, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgMint{
				ContractId: s.contractID,
				From:       s.stranger.String(),
				To:         s.vendor.String(),
				Amount:     sdk.OneInt(),
			},
			expectedError: token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			s.Require().NoError(tc.req.ValidateBasic())
			from, err := sdk.AccAddressFromBech32(tc.req.From)
			s.Require().NoError(err)
			to, err := sdk.AccAddressFromBech32(tc.req.To)
			s.Require().NoError(err)
			ctx, _ := s.ctx.CacheContext()
			prevFrom := s.keeper.GetBalance(ctx, tc.req.ContractId, from)
			prevTo := s.keeper.GetBalance(ctx, tc.req.ContractId, to)
			prevMint := s.keeper.GetMinted(ctx, tc.req.ContractId)
			prevSupplyAmount := s.keeper.GetSupply(ctx, tc.req.ContractId)

			// Act
			res, err := s.msgServer.Mint(sdk.WrapSDKContext(ctx), tc.req)
			if tc.isNegativeCase {
				s.Require().Nil(res)
				s.Require().ErrorIs(err, tc.expectedError)
				s.Require().Equal(0, len(ctx.EventManager().Events()))
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			// Assert
			events := ctx.EventManager().Events()
			s.Require().Equal(tc.expectedEvents, events)
			mintAmount := tc.req.Amount
			curMinted := s.keeper.GetMinted(ctx, tc.req.ContractId)
			curSupply := s.keeper.GetSupply(ctx, tc.req.ContractId)
			curToAmount := s.keeper.GetBalance(ctx, s.contractID, to)
			s.Require().Equal(prevMint.Add(mintAmount), curMinted)
			s.Require().Equal(prevSupplyAmount.Add(mintAmount), curSupply)
			s.Require().Equal(prevTo.Add(mintAmount), curToAmount)
			if !from.Equals(to) {
				curFrom := s.keeper.GetBalance(ctx, s.contractID, from)
				s.Require().Equal(prevFrom, curFrom)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurn() {
	testCases := map[string]struct {
		isNegativeCase bool
		req            *token.MsgBurn
		expectedEvents sdk.Events
		expectedError  *sdkerrors.Error
	}{
		"burn(contractID, from, amount)": {
			req: &token.MsgBurn{
				ContractId: s.contractID,
				From:       s.vendor.String(),
				Amount:     sdk.OneInt(),
			},
			expectedEvents: sdk.Events{
				sdk.Event{
					Type: "lbm.token.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.W(sdk.OneInt()), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.vendor), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor), Index: false},
					},
				},
			},
		},
		"burn(nonExistingContractId, from, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgBurn{
				ContractId: "fee1dead",
				From:       s.vendor.String(),
				Amount:     sdk.OneInt(),
			},
			expectedError: class.ErrContractNotExist,
		},
		"burn(contractID, from, unauthorized account, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgBurn{
				ContractId: s.contractID,
				From:       s.stranger.String(),
				Amount:     sdk.OneInt(),
			},
			expectedError: token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			from, err := sdk.AccAddressFromBech32(tc.req.From)
			s.Require().NoError(err)
			ctx, _ := s.ctx.CacheContext()
			prevFrom := s.keeper.GetBalance(ctx, tc.req.ContractId, from)
			prevBurnt := s.keeper.GetBurnt(ctx, tc.req.ContractId)
			prevSupplyAmount := s.keeper.GetSupply(ctx, tc.req.ContractId)
			s.Require().NoError(tc.req.ValidateBasic())

			// Act
			res, err := s.msgServer.Burn(sdk.WrapSDKContext(ctx), tc.req)
			if tc.isNegativeCase {
				s.Require().Nil(res)
				s.Require().ErrorIs(err, tc.expectedError)
				s.Require().Equal(0, len(ctx.EventManager().Events()))
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			// Assert
			events := ctx.EventManager().Events()
			s.Require().Equal(tc.expectedEvents, events)

			curBurnt := s.keeper.GetBurnt(ctx, tc.req.ContractId)
			curSupply := s.keeper.GetSupply(ctx, tc.req.ContractId)
			curFromAmount := s.keeper.GetBalance(ctx, s.contractID, from)
			burnAmount := tc.req.Amount
			s.Require().Equal(prevBurnt.Add(burnAmount), curBurnt)
			s.Require().Equal(prevSupplyAmount.Sub(burnAmount), curSupply)
			s.Require().Equal(prevFrom.Sub(burnAmount), curFromAmount)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorBurn() {
	testCases := map[string]struct {
		isNegativeCase bool
		req            *token.MsgOperatorBurn
		expectedEvent  sdk.Event
		expectedError  *sdkerrors.Error
	}{
		"operatorBurn(contractID, operator, from, 1)": {
			req: &token.MsgOperatorBurn{
				ContractId: s.contractID,
				Operator:   s.operator.String(),
				From:       s.customer.String(),
				Amount:     sdk.OneInt(),
			},
			expectedEvent: sdk.Event{
				Type: "lbm.token.v1.EventBurned",
				Attributes: []abci.EventAttribute{
					{Key: []byte("amount"), Value: testutil.W(sdk.OneInt()), Index: false},
					{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
					{Key: []byte("from"), Value: testutil.W(s.customer), Index: false},
					{Key: []byte("operator"), Value: testutil.W(s.operator), Index: false},
				},
			},
		},
		"operatorBurn(nonExistingContractId, operator, from, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgOperatorBurn{
				ContractId: "fee1dead",
				Operator:   s.operator.String(),
				From:       s.customer.String(),
				Amount:     sdk.OneInt(),
			},
			expectedError: class.ErrContractNotExist,
		},
		"operatorBurn(contractID, operator, unauthorized account, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgOperatorBurn{
				ContractId: s.contractID,
				Operator:   s.operator.String(),
				From:       s.stranger.String(),
				Amount:     sdk.OneInt(),
			},
			expectedError: token.ErrTokenNotApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			operator, err := sdk.AccAddressFromBech32(tc.req.Operator)
			s.Require().NoError(err)
			from, err := sdk.AccAddressFromBech32(tc.req.From)
			s.Require().NoError(err)
			prevOperator := s.keeper.GetBalance(s.ctx, tc.req.ContractId, operator)
			prevFrom := s.keeper.GetBalance(s.ctx, tc.req.ContractId, from)
			prevBurnt := s.keeper.GetBurnt(s.ctx, tc.req.ContractId)
			prevSupplyAmount := s.keeper.GetSupply(s.ctx, tc.req.ContractId)
			s.Require().NoError(tc.req.ValidateBasic())
			prevEvtCnt := len(s.ctx.EventManager().Events())

			// Act
			res, err := s.msgServer.OperatorBurn(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.isNegativeCase {
				s.Require().Nil(res)
				s.Require().ErrorIs(err, tc.expectedError)
				s.Require().Equal(prevEvtCnt, len(s.ctx.EventManager().Events()))
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			// Assert
			events := s.ctx.EventManager().Events()
			s.Require().Equal(events[len(events)-1], tc.expectedEvent)
			s.Require().Greater(len(s.ctx.EventManager().Events()), prevEvtCnt)

			curBurnt := s.keeper.GetBurnt(s.ctx, tc.req.ContractId)
			curSupply := s.keeper.GetSupply(s.ctx, tc.req.ContractId)
			curFromAmount := s.keeper.GetBalance(s.ctx, s.contractID, from)
			burnAmount := tc.req.Amount
			s.Require().Equal(prevBurnt.Add(burnAmount), curBurnt)
			s.Require().Equal(prevSupplyAmount.Sub(burnAmount), curSupply)
			s.Require().Equal(prevFrom.Sub(burnAmount), curFromAmount)
			if !from.Equals(operator) {
				curOperator := s.keeper.GetBalance(s.ctx, s.contractID, operator)
				s.Require().Equal(prevOperator, curOperator)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgModify() {
	testCases := map[string]struct {
		isNegativeCase bool
		req            *token.MsgModify
		expectedEvents sdk.Events
		expectedError  *sdkerrors.Error
	}{
		"modify(contractID, owner, changes:uri,name)": {
			req: &token.MsgModify{
				ContractId: s.contractID,
				Owner:      s.vendor.String(),
				Changes: []token.Attribute{
					{Key: token.AttributeKeyURI.String(), Value: "uri"},
					{Key: token.AttributeKeyName.String(), Value: "NA<ENDSLSDN"},
				},
			},
			expectedEvents: []sdk.Event{
				{
					Type: "lbm.token.v1.EventModified",
					Attributes: []abci.EventAttribute{
						{Key: []byte("changes"), Value: testutil.MustJSONMarshal([]token.Attribute{{Key: token.AttributeKeyURI.String(), Value: "uri"}, {Key: token.AttributeKeyName.String(), Value: "NA<ENDSLSDN"}}), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor), Index: false},
					},
				},
			},
		},
		"modify(contractID, owner, changes:uri)": {
			req: &token.MsgModify{
				ContractId: s.contractID,
				Owner:      s.vendor.String(),
				Changes:    []token.Attribute{{Key: token.AttributeKeyURI.String(), Value: "uri222"}},
			},
			expectedEvents: []sdk.Event{
				{
					Type: "lbm.token.v1.EventModified",
					Attributes: []abci.EventAttribute{
						{Key: []byte("changes"), Value: testutil.MustJSONMarshal([]token.Attribute{{Key: token.AttributeKeyURI.String(), Value: "uri222"}}), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor), Index: false},
					},
				},
			},
		},
		"modify(nonExistingContractId, from, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgModify{
				ContractId: "fee1dead",
				Owner:      s.vendor.String(),
				Changes:    []token.Attribute{{Key: token.AttributeKeyURI.String(), Value: "uri"}},
			},
			expectedError: class.ErrContractNotExist,
		},
		"modify(contractID, from, unauthorized account, 1) -> error": {
			isNegativeCase: true,
			req: &token.MsgModify{
				ContractId: s.contractID,
				Owner:      s.stranger.String(),
				Changes:    []token.Attribute{{Key: token.AttributeKeyURI.String(), Value: "uri"}},
			},
			expectedError: token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			s.Require().NoError(tc.req.ValidateBasic())
			ctx, _ := s.ctx.CacheContext()

			// Act
			res, err := s.msgServer.Modify(sdk.WrapSDKContext(ctx), tc.req)
			if tc.isNegativeCase {
				s.Require().Nil(res)
				s.Require().ErrorIs(err, tc.expectedError)
				s.Require().Equal(0, len(ctx.EventManager().Events()))
				return
			}
			s.Require().NotNil(res)
			s.Require().NoError(err)

			// Assert
			events := ctx.EventManager().Events()
			s.Require().Equal(tc.expectedEvents, events)
		})
	}
}
