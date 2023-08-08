package testutil

import (
	"github.com/Finschia/finschia-rdk/testutil"
	clitestutil "github.com/Finschia/finschia-rdk/testutil/cli"
	"github.com/Finschia/finschia-rdk/testutil/network"
	"github.com/Finschia/finschia-rdk/x/authz/client/cli"
)

func ExecGrant(val *network.Validator, args []string) (testutil.BufferWriter, error) {
	cmd := cli.NewCmdGrantAuthorization()
	clientCtx := val.ClientCtx
	return clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
}
