package stakingplus

import (
	"fmt"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
)

func (s *E2ETestSuite) TestNewTxCmdCreateValidator() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"grantee msg": {
			[]string{
				"./grantee.json",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(s.grantee)),
			},
			true,
		},
		"stranger msg": {
			[]string{
				"./stranger.json",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(s.stranger)),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewCreateValidatorCmd(s.valAddressCodec)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)

			err = clitestutil.CheckTxCode(s.network, val.ClientCtx, res.TxHash, 0)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
		})
	}
}
