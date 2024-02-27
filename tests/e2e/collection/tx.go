package collection

import (
	"fmt"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/client/cli"
)

func (s *E2ETestSuite) TestNewTxCmdSendNFT() {
	val := s.network.Validators[0]
	tokenID := s.tokenIDs[s.stranger.String()]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.stranger.String(),
				s.customer.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdSendNFT(s.ac)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdOperatorSendNFT() {
	val := s.network.Validators[0]
	tokenID := s.tokenIDs[s.customer.String()]

	testCases := map[string]struct {
		args    []string
		valid   bool
		success bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				tokenID,
			},
			true,
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				tokenID,
				"extra",
			},
			false,
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
			},
			false,
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				tokenID,
			},
			false,
			false,
		},
		"invalid operator": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				s.vendor.String(),
				tokenID,
			},
			true,
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorSendNFT(s.ac)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			if tc.success {
				s.getTxResp(out, 0)
			} else {
				s.getTxResp(out, 29)
			}
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdCreateContract() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.vendor.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.vendor.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{},
			false,
		},
		"invalid creator": {
			[]string{
				"",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdCreateContract()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdIssueNFT() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdIssueNFT()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdMintNFT() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.nftClassID,
				fmt.Sprintf("--%s=%s", cli.FlagName, "arctic fox"),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.nftClassID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.customer.String(),
				s.nftClassID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdMintNFT(s.ac)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdBurnNFT() {
	val := s.network.Validators[0]
	tokenID := s.tokenIDs[s.operator.String()]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdBurnNFT()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdOperatorOperatorBurnNFT() {
	val := s.network.Validators[0]
	tokenID := s.tokenIDs[s.vendor.String()]

	testCases := map[string]struct {
		args    []string
		valid   bool
		success bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				tokenID,
			},
			true,
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				tokenID,
				"extra",
			},
			false,
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
			},
			false,
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.vendor.String(),
				tokenID,
			},
			false,
			false,
		},
		"invalid operator": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.vendor.String(),
				tokenID,
			},
			true,
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorBurnNFT(s.ac)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			if tc.success {
				s.getTxResp(out, 0)
			} else {
				s.getTxResp(out, 29)
			}
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdModify() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
				"tibetian fox",
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
				"tibetian fox",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
				"tibetian fox",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdModify()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdGrantPermission() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				collection.LegacyPermissionMint.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				collection.LegacyPermissionMint.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdGrantPermission(s.ac)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdRevokePermission() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.vendor.String(),
				collection.LegacyPermissionModify.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				collection.LegacyPermissionModify.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdRevokePermission()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdAuthorizeOperator() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.vendor.String(),
				s.customer.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				s.customer.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdAuthorizeOperator(s.ac)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdRevokeOperator() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdRevokeOperator(s.ac)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, s.commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.getTxResp(out, 0)
		})
	}
}
