package testutil

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	ostcli "github.com/line/ostracon/libs/cli"

	"github.com/line/lbm-sdk/client/flags"
	codectypes "github.com/line/lbm-sdk/codec/types"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/collection"
	"github.com/line/lbm-sdk/x/collection/client/cli"
)

func (s *IntegrationTestSuite) TestNewQueryCmdBalance() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				s.customer.String(),
				fmt.Sprintf("--%s=%s", cli.FlagTokenID, collection.NewFTID(s.ftClassID)),
			},
			true,
			&collection.QueryBalanceResponse{
				Balance: collection.NewFTCoin(s.ftClassID, s.balance),
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.customer.String(),
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"invalid address": {
			[]string{
				s.contractID,
				"",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdBalances()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryBalanceResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdFTSupply() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	tokenID := collection.NewFTID(s.ftClassID)
	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				tokenID,
			},
			true,
			&collection.QueryFTSupplyResponse{
				Supply: s.balance.Mul(sdk.NewInt(4)),
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				tokenID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"invalid contract id": {
			[]string{
				"",
				tokenID,
			},
			false,
			nil,
		},
		"invalid class id": {
			[]string{
				s.contractID,
				"",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdFTSupply()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryFTSupplyResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdFTMinted() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	tokenID := collection.NewFTID(s.ftClassID)
	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				tokenID,
			},
			true,
			&collection.QueryFTMintedResponse{
				Minted: s.balance.Mul(sdk.NewInt(5)),
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				tokenID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"invalid contract id": {
			[]string{
				"",
				tokenID,
			},
			false,
			nil,
		},
		"invalid class id": {
			[]string{
				s.contractID,
				"",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdFTMinted()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryFTMintedResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdFTBurnt() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	tokenID := collection.NewFTID(s.ftClassID)
	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				tokenID,
			},
			true,
			&collection.QueryFTBurntResponse{
				Burnt: s.balance,
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				tokenID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"invalid contract id": {
			[]string{
				"",
				tokenID,
			},
			false,
			nil,
		},
		"invalid class id": {
			[]string{
				s.contractID,
				"",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdFTBurnt()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryFTBurntResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdNFTSupply() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				s.nftClassID,
			},
			true,
			&collection.QueryNFTSupplyResponse{
				Supply: sdk.NewInt(24),
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.nftClassID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"invalid contract id": {
			[]string{
				"",
				s.nftClassID,
			},
			false,
			nil,
		},
		"invalid class id": {
			[]string{
				s.contractID,
				"",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdNFTSupply()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryNFTSupplyResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdNFTMinted() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				s.nftClassID,
			},
			true,
			&collection.QueryNFTMintedResponse{
				Minted: sdk.NewInt(24),
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.nftClassID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"invalid contract id": {
			[]string{
				"",
				s.nftClassID,
			},
			false,
			nil,
		},
		"invalid class id": {
			[]string{
				s.contractID,
				"",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdNFTMinted()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryNFTMintedResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdNFTBurnt() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				s.nftClassID,
			},
			true,
			&collection.QueryNFTBurntResponse{
				Burnt: sdk.ZeroInt(),
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.nftClassID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"invalid contract id": {
			[]string{
				"",
				s.nftClassID,
			},
			false,
			nil,
		},
		"invalid class id": {
			[]string{
				s.contractID,
				"",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdNFTBurnt()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryNFTBurntResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdContract() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
			},
			true,
			&collection.QueryContractResponse{
				Contract: collection.Contract{Id: s.contractID},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdContract()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryContractResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdTokenType() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				s.nftClassID,
			},
			true,
			&collection.QueryTokenTypeResponse{
				TokenType: collection.TokenType{
					ContractId: s.contractID,
					TokenType:  s.nftClassID,
				},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.nftClassID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"token not found": {
			[]string{
				s.contractID,
				"deadbeef",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdTokenType()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryTokenTypeResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdToken() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	tokenID := collection.NewNFTID(s.nftClassID, 1)
	token, err := codectypes.NewAnyWithValue(&collection.OwnerNFT{
		ContractId: s.contractID,
		TokenId:    tokenID,
		Name:       "arctic fox",
		Owner:      s.customer.String(),
	})
	s.Require().NoError(err)

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				tokenID,
			},
			true,
			&collection.QueryTokenResponse{
				Token: *token,
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				tokenID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"token not found": {
			[]string{
				s.contractID,
				collection.NewNFTID("deadbeef", 1),
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdToken()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryTokenResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			err = collection.TokenUnpackInterfaces(&actual.Token, val.ClientCtx.InterfaceRegistry)
			s.Require().NoError(err)
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdRoot() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	tokenID := collection.NewNFTID(s.nftClassID, 2)

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				tokenID,
			},
			true,
			&collection.QueryRootResponse{
				Root: collection.NFT{
					TokenId: collection.NewNFTID(s.nftClassID, 1),
					Name:    "arctic fox",
				},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				tokenID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"token not found": {
			[]string{
				s.contractID,
				collection.NewNFTID("deadbeef", 1),
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdRoot()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryRootResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdParent() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	tokenID := collection.NewNFTID(s.nftClassID, 2)

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				tokenID,
			},
			true,
			&collection.QueryParentResponse{
				Parent: collection.NFT{
					TokenId: collection.NewNFTID(s.nftClassID, 1),
					Name:    "arctic fox",
				},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				tokenID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
		"token not found": {
			[]string{
				s.contractID,
				collection.NewNFTID("deadbeef", 1),
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdParent()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryParentResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdChildren() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	tokenID := collection.NewNFTID(s.nftClassID, 1)

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				tokenID,
			},
			true,
			&collection.QueryChildrenResponse{
				Children: []collection.NFT{{
					TokenId: collection.NewNFTID(s.nftClassID, 2),
					Name:    "arctic fox",
				}},
				Pagination: &query.PageResponse{},
			},
		},
		"token not found": {
			[]string{
				s.contractID,
				collection.NewNFTID("deadbeef", 1),
			},
			true,
			&collection.QueryChildrenResponse{
				Children: []collection.NFT{},
				Pagination: &query.PageResponse{},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				tokenID,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdChildren()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryChildrenResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdGranteeGrants() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				s.operator.String(),
			},
			true,
			&collection.QueryGranteeGrantsResponse{
				Grants: []collection.Grant{
					{
						Grantee:    s.operator.String(),
						Permission: collection.PermissionIssue,
					},
					{
						Grantee:    s.operator.String(),
						Permission: collection.PermissionModify,
					},
					{
						Grantee:    s.operator.String(),
						Permission: collection.PermissionMint,
					},
					{
						Grantee:    s.operator.String(),
						Permission: collection.PermissionBurn,
					},
				},
				Pagination: &query.PageResponse{
					Total: 4,
				},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdGranteeGrants()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryGranteeGrantsResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdIsOperatorFor() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
			},
			true,
			&collection.QueryIsOperatorForResponse{
				Authorized: true,
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdIsOperatorFor()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryIsOperatorForResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdHoldersByOperator() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.contractID,
				s.vendor.String(),
			},
			true,
			&collection.QueryHoldersByOperatorResponse{
				Holders:    []string{s.operator.String()},
				Pagination: &query.PageResponse{},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdHoldersByOperator()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryHoldersByOperatorResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}
