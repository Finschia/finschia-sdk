package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/Finschia/finschia-sdk/x/fbridge/testutil"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestSetParams(t *testing.T) {
	key, memKey, ctx, encCfg, authKeeper, bankKeeper, _ := testutil.PrepareFbridgeTest(t, 0)
	keeper := NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, types.DefaultAuthority().String())

	tcs := map[string]struct {
		malleate func() types.Params
		isErr    bool
	}{
		"invalid guardian trust level": {
			malleate: func() types.Params {
				params := types.Params{}
				params.GuardianTrustLevel = types.Fraction{Numerator: 4, Denominator: 3}
				return params
			},
			isErr: true,
		},
		"invalid operator trust level": {
			malleate: func() types.Params {
				params := types.Params{}
				params.GuardianTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.OperatorTrustLevel = types.Fraction{Numerator: 4, Denominator: 3}
				return params
			},
			isErr: true,
		},
		"invalid judge trust level": {
			malleate: func() types.Params {
				params := types.Params{}
				params.GuardianTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.OperatorTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.JudgeTrustLevel = types.Fraction{Numerator: 4, Denominator: 3}
				return params
			},
			isErr: true,
		},
		"invalid proposal period": {
			malleate: func() types.Params {
				params := types.Params{}
				params.GuardianTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.OperatorTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.JudgeTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.ProposalPeriod = 0
				return params
			},
			isErr: true,
		},
		"invalid timelock period": {
			malleate: func() types.Params {
				params := types.Params{}
				params.GuardianTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.OperatorTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.JudgeTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.ProposalPeriod = 10
				params.TimelockPeriod = 0
				return params
			},
			isErr: true,
		},
		"invalid target denom": {
			malleate: func() types.Params {
				params := types.Params{}
				params.GuardianTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.OperatorTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.JudgeTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.ProposalPeriod = 10
				params.TimelockPeriod = 20
				params.TargetDenom = ""
				return params
			},
			isErr: true,
		},
		"missing some fields": {
			malleate: func() types.Params {
				params := types.Params{}
				params.GuardianTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.OperatorTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.JudgeTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.TimelockPeriod = 20
				return params
			},
			isErr: true,
		},
		"valid": {
			malleate: func() types.Params {
				params := types.Params{}
				params.GuardianTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.OperatorTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.JudgeTrustLevel = types.Fraction{Numerator: 2, Denominator: 3}
				params.ProposalPeriod = 10
				params.TimelockPeriod = 20
				params.TargetDenom = "stake"
				return params
			},
			isErr: false,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			params := tc.malleate()
			if tc.isErr {
				require.Error(t, keeper.SetParams(ctx, params))
			} else {
				require.NoError(t, keeper.SetParams(ctx, params))
			}
		})
	}
}
