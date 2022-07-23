package testutil

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	ostcli "github.com/line/ostracon/libs/cli"

	"github.com/line/lbm-sdk/client/flags"
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

func (s *IntegrationTestSuite) TestNewQueryCmdSupply() {
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
				s.ftClassID,
			},
			true,
			&collection.QuerySupplyResponse{
				Supply: s.balance.Mul(sdk.NewInt(4)),
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.ftClassID,
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
				s.ftClassID,
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
			cmd := cli.NewQueryCmdSupply()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QuerySupplyResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdMinted() {
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
				s.ftClassID,
			},
			true,
			&collection.QueryMintedResponse{
				Minted: s.balance.Mul(sdk.NewInt(5)),
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.ftClassID,
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
				s.ftClassID,
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
			cmd := cli.NewQueryCmdMinted()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryMintedResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdBurnt() {
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
				s.ftClassID,
			},
			true,
			&collection.QueryBurntResponse{
				Burnt: s.balance,
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.ftClassID,
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
				s.ftClassID,
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
			cmd := cli.NewQueryCmdBurnt()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryBurntResponse
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
				Contract: collection.Contract{ContractId: s.contractID},
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

func (s *IntegrationTestSuite) TestNewQueryCmdContracts() {
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
		"query all": {
			[]string{},
			true,
			&collection.QueryContractsResponse{
				Contracts:  []collection.Contract{{ContractId: s.contractID}},
				Pagination: &query.PageResponse{},
			},
		},
		"extra args": {
			[]string{
				"extra",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdContracts()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryContractsResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdFTClass() {
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
				s.ftClassID,
			},
			true,
			&collection.QueryFTClassResponse{
				Class: collection.FTClass{
					Id:       s.ftClassID,
					Name:     "tibetian fox",
					Decimals: 8,
					Mintable: true,
				},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.ftClassID,
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
		"class not found": {
			[]string{
				s.contractID,
				"00bab10c",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdFTClass()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryFTClassResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdNFTClass() {
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
			&collection.QueryNFTClassResponse{
				Class: collection.NFTClass{Id: s.nftClassID},
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
		"class not found": {
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
			cmd := cli.NewQueryCmdNFTClass()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryNFTClassResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdNFT() {
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
			&collection.QueryNFTResponse{
				Token: collection.NFT{
					Id: tokenID,
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
			cmd := cli.NewQueryCmdNFT()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryNFTResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
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
					Id: collection.NewNFTID(s.nftClassID, 1),
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
					Id: collection.NewNFTID(s.nftClassID, 1),
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
					Id: collection.NewNFTID(s.nftClassID, 2),
				}},
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

func (s *IntegrationTestSuite) TestNewQueryCmdGrant() {
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
				collection.PermissionIssue.String(),
			},
			true,
			&collection.QueryGrantResponse{
				Grant: collection.Grant{
					Grantee:    s.operator.String(),
					Permission: collection.PermissionIssue,
				},
			},
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				collection.PermissionIssue.String(),
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
		"no permission found": {
			[]string{
				s.contractID,
				s.customer.String(),
				collection.PermissionIssue.String(),
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdGrant()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryGrantResponse
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

func (s *IntegrationTestSuite) TestNewQueryCmdAuthorization() {
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
			&collection.QueryAuthorizationResponse{
				Authorization: collection.Authorization{
					Holder:   s.customer.String(),
					Operator: s.operator.String(),
				},
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
		"no authorization found": {
			[]string{
				s.contractID,
				s.customer.String(),
				s.vendor.String(),
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdAuthorization()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryAuthorizationResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdOperatorAuthorizations() {
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
			&collection.QueryOperatorAuthorizationsResponse{
				Authorizations: []collection.Authorization{
					{
						Holder:   s.operator.String(),
						Operator: s.vendor.String(),
					},
				},
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
			cmd := cli.NewQueryCmdOperatorAuthorizations()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual collection.QueryOperatorAuthorizationsResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}
