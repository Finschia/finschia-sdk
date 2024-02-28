package collection

import (
	"fmt"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/client/cli"
)

func (s *E2ETestSuite) TestNewTxCmdSendNFT() {
	val := s.network.Validators[0]
	tokenID := s.tokenIDs[s.stranger]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.stranger,
				s.customer,
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.stranger,
				s.customer,
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.stranger,
				s.customer,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.stranger,
				s.customer,
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdSendNFT()
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
	tokenID := s.tokenIDs[s.customer]

	testCases := map[string]struct {
		args    []string
		valid   bool
		success bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator,
				s.customer,
				s.vendor,
				tokenID,
			},
			true,
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator,
				s.customer,
				s.vendor,
				tokenID,
				"extra",
			},
			false,
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator,
				s.customer,
				s.vendor,
			},
			false,
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator,
				s.customer,
				s.vendor,
				tokenID,
			},
			false,
			false,
		},
		"invalid operator": {
			[]string{
				s.contractID,
				s.stranger,
				s.customer,
				s.vendor,
				tokenID,
			},
			true,
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorSendNFT()
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
				s.vendor,
			},
			true,
		},
		"extra args": {
			[]string{
				s.vendor,
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
				s.operator,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator,
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
				s.operator,
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
				s.operator,
				s.customer,
				s.nftClassID,
				fmt.Sprintf("--%s=%s", cli.FlagName, "arctic fox"),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator,
				s.customer,
				s.nftClassID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator,
				s.customer,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator,
				s.customer,
				s.nftClassID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdMintNFT()
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
	tokenID := s.tokenIDs[s.operator]

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator,
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator,
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator,
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
	tokenID := s.tokenIDs[s.vendor]

	testCases := map[string]struct {
		args    []string
		valid   bool
		success bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator,
				s.vendor,
				tokenID,
			},
			true,
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator,
				s.vendor,
				tokenID,
				"extra",
			},
			false,
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator,
				s.vendor,
			},
			false,
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator,
				s.vendor,
				tokenID,
			},
			false,
			false,
		},
		"invalid operator": {
			[]string{
				s.contractID,
				s.stranger,
				s.vendor,
				tokenID,
			},
			true,
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorBurnNFT()
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
				s.operator,
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
				s.operator,
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
				s.operator,
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator,
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
				s.operator,
				s.customer,
				collection.LegacyPermissionMint.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator,
				s.customer,
				collection.LegacyPermissionMint.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator,
				s.customer,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdGrantPermission()
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
				s.vendor,
				collection.LegacyPermissionModify.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor,
				collection.LegacyPermissionModify.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor,
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
				s.vendor,
				s.customer,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor,
				s.customer,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdAuthorizeOperator()
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
				s.operator,
				s.vendor,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator,
				s.vendor,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdRevokeOperator()
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
