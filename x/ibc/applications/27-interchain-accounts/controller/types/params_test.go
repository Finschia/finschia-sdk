package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/controller/types"
)

func TestValidateParams(t *testing.T) {
	require.NoError(t, types.DefaultParams().Validate())
	require.NoError(t, types.NewParams(false).Validate())
}
