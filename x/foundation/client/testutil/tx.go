package testutil

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation/client/cli"
)

func (s *IntegrationTestSuite) TestNewProposalCmdUpdateFoundationParams() {
	val := s.network.Validators[0]

	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			"with wrong # of args",
			append([]string{
				"no-args-expected",
			}, commonFlags...),
			true, 0, nil,
		},
		{
			"valid transaction",
			commonFlags,
			false, 0, &sdk.TxResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewProposalCmdUpdateFoundationParams()
			clientCtx := val.ClientCtx
			flags.AddTxFlagsToCmd(cmd)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(s.cfg.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewProposalCmdUpdateValidatorAuths() {
	val := s.network.Validators[0]

	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			"with no args",
			commonFlags,
			true, 0, nil,
		},
		{
			"with an invalid address",
			append([]string{
				fmt.Sprintf("--%s=%s",
					cli.FlagAllowedValidatorAdd,
					"this-is-an-invalid-address",
				),
			}, commonFlags...),
			true, 0, nil,
		},
		{
			"with duplicated validators in add",
			append([]string{
				fmt.Sprintf("--%s=%s,%s",
					cli.FlagAllowedValidatorAdd,
					val.ValAddress.String(),
					val.ValAddress.String(),
				),
			}, commonFlags...),
			true, 0, nil,
		},
		{
			"with same validators in both add and delete",
			append([]string{
				fmt.Sprintf("--%s=%s",
					cli.FlagAllowedValidatorAdd,
					val.ValAddress.String()),
				fmt.Sprintf("--%s=%s",
					cli.FlagAllowedValidatorDelete,
					val.ValAddress.String()),
			}, commonFlags...),
			true, 0, nil,
		},
		{
			"valid transaction",
			append([]string{
				fmt.Sprintf("--%s=%s",
					cli.FlagAllowedValidatorDelete,
					val.ValAddress.String()),
			}, commonFlags...),
			false, 0, &sdk.TxResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewProposalCmdUpdateValidatorAuths()
			clientCtx := val.ClientCtx
			flags.AddTxFlagsToCmd(cmd)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(s.cfg.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}
