package testutil

import (
	"encoding/json"
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/Finschia/finschia-sdk/client/flags"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	sdk "github.com/Finschia/finschia-sdk/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/fswap/client/cli"
)

func (s *IntegrationTestSuite) TestCmdTxMsgSwap() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// avoid printing as yaml from CLI command
	clientCtx.OutputFormat = jsonOutputFormat
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
	}{
		{
			"valid transaction",
			[]string{
				val.Address.String(),
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)).String(),
				s.toDenom.Base,
			},
			false,
			&sdk.TxResponse{},
		},
		{
			"invalid request (wrong number of args)",
			[]string{
				val.Address.String(),
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)).String(),
				s.toDenom.Base,
				"extra",
			},
			true,
			nil,
		},
		{
			"invalid request (invalid from address)",
			[]string{
				"invalidAddress",
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)).String(),
				s.toDenom.Base,
			},
			true,
			nil,
		},
		{
			"invalid request (invalid from coin amount)",
			[]string{
				val.Address.String(),
				"",
				s.toDenom.Base,
			},
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.CmdTxMsgSwap()
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, commonArgs...))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), tc.respType), bz.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdTxMsgSwapAll() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// avoid printing as yaml from CLI command
	clientCtx.OutputFormat = jsonOutputFormat
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
	}{
		{
			"valid transaction",
			[]string{
				val.Address.String(),
				s.cfg.BondDenom,
				s.toDenom.Base,
			},
			false,
			&sdk.TxResponse{},
		},
		{
			"invalid request (wrong number of args)",
			[]string{
				val.Address.String(),
				s.cfg.BondDenom,
				s.toDenom.Base,
				"extra",
			},
			true,
			nil,
		},
		{
			"invalid request (invalid from address)",
			[]string{
				"invalidAddress",
				s.cfg.BondDenom,
				s.toDenom.Base,
			},
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.CmdTxMsgSwapAll()
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, commonArgs...))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), tc.respType), bz.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestMsgSetSwap() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// avoid printing as yaml from CLI command
	clientCtx.OutputFormat = jsonOutputFormat

	denomMeta := struct {
		Metadata banktypes.Metadata `json:"metadata"`
	}{
		Metadata: s.toDenom,
	}
	jsonBytes, err := json.Marshal(denomMeta)
	s.Require().NoError(err)
	denomMetaString := string(jsonBytes)

	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := []struct {
		name string

		args      []string
		expectErr bool
	}{
		{
			"valid transaction (generateOnly)",
			[]string{
				s.authority.String(),
				denomMetaString,
				fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
				fmt.Sprintf("--%s=%s", cli.FlagFromDenom, s.dummySwap.FromDenom),
				fmt.Sprintf("--%s=%s", cli.FlagToDenom, s.dummySwap.ToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagAmountCapForToDenom, s.dummySwap.AmountCapForToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagSwapRate, s.dummySwap.SwapRate),
			},
			false,
		},
		{
			"invalid transaction (without generateOnly)",
			[]string{
				s.authority.String(),
				denomMetaString,
				fmt.Sprintf("--%s=false", flags.FlagGenerateOnly),
				fmt.Sprintf("--%s=%s", cli.FlagFromDenom, s.dummySwap.FromDenom),
				fmt.Sprintf("--%s=%s", cli.FlagToDenom, s.dummySwap.ToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagAmountCapForToDenom, s.dummySwap.AmountCapForToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagSwapRate, s.dummySwap.SwapRate),
			},
			true,
		},
		{
			"extra args",
			[]string{
				s.authority.String(),
				denomMetaString,
				"extra",
				fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
				fmt.Sprintf("--%s=%s", cli.FlagFromDenom, s.dummySwap.FromDenom),
				fmt.Sprintf("--%s=%s", cli.FlagToDenom, s.dummySwap.ToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagAmountCapForToDenom, s.dummySwap.AmountCapForToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagSwapRate, s.dummySwap.SwapRate),
			},
			true,
		},
		{
			"invalid authority",
			[]string{
				"invalid authority",
				denomMetaString,
				fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
				fmt.Sprintf("--%s=%s", cli.FlagFromDenom, s.dummySwap.FromDenom),
				fmt.Sprintf("--%s=%s", cli.FlagToDenom, s.dummySwap.ToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagAmountCapForToDenom, s.dummySwap.AmountCapForToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagSwapRate, s.dummySwap.SwapRate),
			},
			true,
		},
		{
			"invalid json",
			[]string{
				s.authority.String(),
				"invalid json",
				fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
				fmt.Sprintf("--%s=%s", cli.FlagFromDenom, s.dummySwap.FromDenom),
				fmt.Sprintf("--%s=%s", cli.FlagToDenom, s.dummySwap.ToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagAmountCapForToDenom, s.dummySwap.AmountCapForToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagSwapRate, s.dummySwap.SwapRate),
			},
			true,
		},
		{
			"invalid amountCapForToDenom",
			[]string{
				s.authority.String(),
				denomMetaString,
				fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
				fmt.Sprintf("--%s=%s", cli.FlagFromDenom, s.dummySwap.FromDenom),
				fmt.Sprintf("--%s=%s", cli.FlagToDenom, s.dummySwap.ToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagAmountCapForToDenom, "123.456"),
				fmt.Sprintf("--%s=%s", cli.FlagSwapRate, s.dummySwap.SwapRate),
			},
			true,
		},
		{
			"invalid swapRate",
			[]string{
				s.authority.String(),
				denomMetaString,
				fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
				fmt.Sprintf("--%s=%s", cli.FlagFromDenom, s.dummySwap.FromDenom),
				fmt.Sprintf("--%s=%s", cli.FlagToDenom, s.dummySwap.ToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagAmountCapForToDenom, s.dummySwap.AmountCapForToDenom),
				fmt.Sprintf("--%s=%s", cli.FlagSwapRate, "abc.123"),
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.CmdMsgSetSwap()
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, commonArgs...))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				_, err := s.cfg.TxConfig.TxJSONDecoder()(bz.Bytes())
				s.Require().NoError(err)
			}
		})
	}
}
