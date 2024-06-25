package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestValidateParams(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		malleate func(p *types.Params)
		expErr   bool
	}{
		"valid params": {
			expErr: false,
		},
		"invalid guardian trust level": {
			malleate: func(p *types.Params) {
				p.GuardianTrustLevel = types.Fraction{Numerator: 3, Denominator: 2}
			},
			expErr: true,
		},
		"invalid operator trust level": {
			malleate: func(p *types.Params) {
				p.OperatorTrustLevel = types.Fraction{Numerator: 0, Denominator: 2}
			},
			expErr: true,
		},
		"invalid judge trust level": {
			malleate: func(p *types.Params) {
				p.JudgeTrustLevel = types.Fraction{Numerator: 0, Denominator: 0}
			},
			expErr: true,
		},
		"invalid proposal period": {
			malleate: func(p *types.Params) {
				p.ProposalPeriod = 0
			},
			expErr: true,
		},
		"invalid timelock period": {
			malleate: func(p *types.Params) {
				p.TimelockPeriod = 0
			},
			expErr: true,
		},
		"invalid target denom": {
			malleate: func(p *types.Params) {
				p.TargetDenom = "invalid_denom"
			},
			expErr: true,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			p := types.DefaultParams()
			if tc.malleate != nil {
				tc.malleate(&p)
			}
			if tc.expErr {
				err := p.ValidateParams()
				require.Error(t, err)
			} else {
				require.NoError(t, p.ValidateParams())
			}
		})

	}
}

func TestCheckTrustLevelThreshold(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		total      uint64
		current    uint64
		trustLevel types.Fraction
		isPanic    bool
		isValid    bool
	}{
		"meet the trust level": {
			current:    3,
			total:      4,
			trustLevel: types.Fraction{Numerator: 2, Denominator: 3},
			isValid:    true,
		},
		"not meet the trust level": {
			current:    1,
			total:      2,
			trustLevel: types.Fraction{Numerator: 2, Denominator: 3},
			isValid:    false,
		},
		"total is 0": {
			total:      0,
			current:    3,
			trustLevel: types.Fraction{Numerator: 2, Denominator: 3},
			isValid:    false,
		},
		"invalid trust level - 1": {
			total:      10,
			current:    8,
			trustLevel: types.Fraction{Numerator: 3, Denominator: 2},
			isPanic:    true,
		},
		"invalid trust level - 2": {
			total:      10,
			current:    8,
			trustLevel: types.Fraction{Numerator: 3, Denominator: 0},
			isPanic:    true,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			if tc.isPanic {
				require.Panics(t, func() { types.CheckTrustLevelThreshold(tc.total, tc.current, tc.trustLevel) })
			} else if tc.isValid {
				require.True(t, types.CheckTrustLevelThreshold(tc.total, tc.current, tc.trustLevel))
			} else {
				require.False(t, types.CheckTrustLevelThreshold(tc.total, tc.current, tc.trustLevel))
			}
		})
	}
}
