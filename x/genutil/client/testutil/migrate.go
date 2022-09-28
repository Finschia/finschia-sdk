package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/testutil"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/x/genutil/client/cli"
)

func TestGetMigrationCallback(t *testing.T) {
	for _, version := range cli.GetMigrationVersions() {
		require.NotNil(t, cli.GetMigrationCallback(version))
	}
}

func (s *IntegrationTestSuite) TestMigrateGenesis() {
	val0 := s.network.Validators[0]

	testCases := []struct {
		name      string
		genesis   string
		target    string
		expErr    bool
		expErrMsg string
		check     func(jsonOut string)
	}{
		{
			"migrate 0.42 to 0.43(result error)",
			v040Valid,
			"v0.43",
			true, "",
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			genesisFile := testutil.WriteToNewTempFile(s.T(), tc.genesis)
			jsonOutput, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cli.MigrateGenesisCmd(), []string{tc.target, genesisFile.Name()})
			if tc.expErr {
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				tc.check(jsonOutput.String())
			}
		})
	}
}
