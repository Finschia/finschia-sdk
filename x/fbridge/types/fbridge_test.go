package types_test

import (
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsValidRole(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		role    types.Role
		expPass bool
	}{
		"valid role - guardian": {
			role:    types.RoleGuardian,
			expPass: true,
		},
		"valid role - operator": {
			role:    types.RoleOperator,
			expPass: true,
		},
		"valid role - judge": {
			role:    types.RoleJudge,
			expPass: true,
		},
		"invalid role": {
			role:    types.Role(10),
			expPass: false,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := types.IsValidRole(tc.role)
			if tc.expPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestIsValidVoteOption(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		option  types.VoteOption
		expPass bool
	}{
		"valid option - yes": {
			option:  types.OptionYes,
			expPass: true,
		},
		"valid option - no": {
			option:  types.OptionNo,
			expPass: true,
		},
		"invalid option": {
			option:  types.VoteOption(10),
			expPass: false,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := types.IsValidVoteOption(tc.option)
			if tc.expPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestIsValidBridgeStatus(t *testing.T) {
	tcs := map[string]struct {
		status  types.BridgeStatus
		expPass bool
	}{
		"valid status - active": {
			status:  types.StatusActive,
			expPass: true,
		},
		"valid status - inactive": {
			status:  types.StatusInactive,
			expPass: true,
		},
		"invalid status": {
			status:  types.BridgeStatus(10),
			expPass: false,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := types.IsValidBridgeStatus(tc.status)
			if tc.expPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
