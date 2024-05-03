package types_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestCheckTrustLevelThreshold(t *testing.T) {
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
