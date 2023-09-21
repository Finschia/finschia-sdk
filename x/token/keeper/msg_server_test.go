package keeper_test

import (
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/Finschia/finschia-sdk/types"
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
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventSent", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x22}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventSent", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x22}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventRevokedOperator", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventAuthorizedOperator", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssue() {
	ownerAddr := s.vendor.String()
	toAddr := s.vendor.String()

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
				sdk.Event{Type: "lbm.token.v1.EventIssued", Attributes: []abci.EventAttribute{{Key: []uint8("contract_id"), Value: []uint8("\"fee15a74\""), Index: false}, {Key: []uint8("creator"), Value: []uint8(fmt.Sprintf("\"%s\"", ownerAddr)), Index: false}, {Key: []uint8("decimals"), Value: []uint8("0"), Index: false}, {Key: []uint8("meta"), Value: []uint8("\"\""), Index: false}, {Key: []uint8("mintable"), Value: []uint8("true"), Index: false}, {Key: []uint8("name"), Value: []uint8("\"test\""), Index: false}, {Key: []uint8("symbol"), Value: []uint8("\"TT\""), Index: false}, {Key: []uint8("uri"), Value: []uint8("\"\""), Index: false}}},
				sdk.Event{Type: "lbm.token.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8("contract_id"), Value: []uint8("\"fee15a74\""), Index: false}, {Key: []uint8("grantee"), Value: []uint8(fmt.Sprintf("\"%s\"", toAddr)), Index: false}, {Key: []uint8("granter"), Value: []uint8("\"\""), Index: false}, {Key: []uint8("permission"), Value: []uint8("\"PERMISSION_MODIFY\""), Index: false}}},
				sdk.Event{Type: "lbm.token.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8("contract_id"), Value: []uint8("\"fee15a74\""), Index: false}, {Key: []uint8("grantee"), Value: []uint8(fmt.Sprintf("\"%s\"", toAddr)), Index: false}, {Key: []uint8("granter"), Value: []uint8("\"\""), Index: false}, {Key: []uint8("permission"), Value: []uint8("\"PERMISSION_MINT\""), Index: false}}},
				sdk.Event{Type: "lbm.token.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8("contract_id"), Value: []uint8("\"fee15a74\""), Index: false}, {Key: []uint8("grantee"), Value: []uint8(fmt.Sprintf("\"%s\"", toAddr)), Index: false}, {Key: []uint8("granter"), Value: []uint8("\"\""), Index: false}, {Key: []uint8("permission"), Value: []uint8("\"PERMISSION_BURN\""), Index: false}}},
				sdk.Event{Type: "lbm.token.v1.EventMinted", Attributes: []abci.EventAttribute{{Key: []uint8("amount"), Value: []uint8("\"10\""), Index: false}, {Key: []uint8("contract_id"), Value: []uint8("\"fee15a74\""), Index: false}, {Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", ownerAddr)), Index: false}, {Key: []uint8("to"), Value: []uint8(fmt.Sprintf("\"%s\"", toAddr)), Index: false}}},
			},
		},
		"mintable false": {
			mintable: false,
			amount:   sdk.NewInt(10),
			events: sdk.Events{
				sdk.Event{Type: "lbm.token.v1.EventIssued", Attributes: []abci.EventAttribute{{Key: []uint8("contract_id"), Value: []uint8("\"fee15a74\""), Index: false}, {Key: []uint8("creator"), Value: []uint8(fmt.Sprintf("\"%s\"", ownerAddr)), Index: false}, {Key: []uint8("decimals"), Value: []uint8("0"), Index: false}, {Key: []uint8("meta"), Value: []uint8("\"\""), Index: false}, {Key: []uint8("mintable"), Value: []uint8("false"), Index: false}, {Key: []uint8("name"), Value: []uint8("\"test\""), Index: false}, {Key: []uint8("symbol"), Value: []uint8("\"TT\""), Index: false}, {Key: []uint8("uri"), Value: []uint8("\"\""), Index: false}}},
				sdk.Event{Type: "lbm.token.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8("contract_id"), Value: []uint8("\"fee15a74\""), Index: false}, {Key: []uint8("grantee"), Value: []uint8(fmt.Sprintf("\"%s\"", ownerAddr)), Index: false}, {Key: []uint8("granter"), Value: []uint8("\"\""), Index: false}, {Key: []uint8("permission"), Value: []uint8("\"PERMISSION_MODIFY\""), Index: false}}},
				sdk.Event{Type: "lbm.token.v1.EventMinted", Attributes: []abci.EventAttribute{{Key: []uint8("amount"), Value: []uint8("\"10\""), Index: false}, {Key: []uint8("contract_id"), Value: []uint8("\"fee15a74\""), Index: false}, {Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", ownerAddr)), Index: false}, {Key: []uint8("to"), Value: []uint8(fmt.Sprintf("\"%s\"", toAddr)), Index: false}}},
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
		"valid request": {
			contractID: s.contractID,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x4f, 0x44, 0x49, 0x46, 0x59, 0x22}, Index: false}}}},
		},
		"contract not found": {
			contractID: "fee1dead",
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
			err:        class.ErrContractNotExist,
		},
		"granter has no permission": {
			contractID: s.contractID,
			granter:    s.customer,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
			err:        token.ErrTokenNoPermission,
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
		"valid request": {
			contractID: s.contractID,
			from:       s.operator,
			permission: token.LegacyPermissionMint.String(),
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventRenounced", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x49, 0x4e, 0x54, 0x22}, Index: false}}}},
		},
		"contract not found": {
			contractID: "fee1dead",
			from:       s.operator,
			permission: token.LegacyPermissionMint.String(),
			err:        class.ErrContractNotExist,
		},
		"not granted yet": {
			contractID: s.contractID,
			from:       s.operator,
			permission: token.LegacyPermissionModify.String(),
			err:        token.ErrTokenNoPermission,
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgMint() {
	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.operator,
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventMinted", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x22, 0x31, 0x22}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}}}},
		},
		"contract not found": {
			contractID: "fee1dead",
			grantee:    s.operator,
			err:        class.ErrContractNotExist,
		},
		"not granted": {
			contractID: s.contractID,
			grantee:    s.customer,
			err:        token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgMint{
				ContractId: tc.contractID,
				From:       tc.grantee.String(),
				To:         s.customer.String(),
				Amount:     sdk.OneInt(),
			}
			res, err := s.msgServer.Mint(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurn() {
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventBurned", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x22}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}}}},
		},
		"contract not found": {
			contractID: "fee1dead",
			from:       s.vendor,
			err:        class.ErrContractNotExist,
		},
		"not granted": {
			contractID: s.contractID,
			from:       s.customer,
			err:        token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgBurn{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Amount:     s.balance,
			}
			res, err := s.msgServer.Burn(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorBurn() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventBurned", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x22}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}}}},
		},
		"contract not found": {
			contractID: "fee1dead",
			operator:   s.operator,
			from:       s.customer,
			err:        class.ErrContractNotExist,
		},
		"not approved": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			err:        token.ErrTokenNotApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgOperatorBurn{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				Amount:     s.balance,
			}
			res, err := s.msgServer.OperatorBurn(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgModify() {
	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.vendor,
			events:     sdk.Events{sdk.Event{Type: "lbm.token.v1.EventModified", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73}, Value: []uint8{0x5b, 0x7b, 0x22, 0x6b, 0x65, 0x79, 0x22, 0x3a, 0x22, 0x75, 0x72, 0x69, 0x22, 0x2c, 0x22, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3a, 0x22, 0x75, 0x72, 0x69, 0x22, 0x7d, 0x5d}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}}}},
		},
		"contract not found": {
			contractID: "fee1dead",
			grantee:    s.vendor,
			err:        class.ErrContractNotExist,
		},
		"not granted": {
			contractID: s.contractID,
			grantee:    s.operator,
			err:        token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgModify{
				ContractId: tc.contractID,
				Owner:      tc.grantee.String(),
				Changes:    []token.Attribute{{Key: token.AttributeKeyImageURI.String(), Value: "uri"}},
			}
			res, err := s.msgServer.Modify(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}
