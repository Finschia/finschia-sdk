package cmd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	svrcmd "github.com/Finschia/finschia-rdk/server/cmd"
	"github.com/Finschia/finschia-rdk/simapp"
	"github.com/Finschia/finschia-rdk/simapp/simd/cmd"
	"github.com/Finschia/finschia-sdk/x/genutil/client/cli"
)

func TestInitCmd(t *testing.T) {
	t.Skipf("🔬 The rollkit/cosmos-sdk also remains faulty.")
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",        // Test the init cmd
		"simapp-test", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
	})

	require.NoError(t, svrcmd.Execute(rootCmd, simapp.DefaultNodeHome))
}
