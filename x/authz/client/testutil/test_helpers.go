package testutil

import (
	"github.com/line/lbm-sdk/testutil"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	"github.com/line/lbm-sdk/x/authz/client/cli"
)

func ExecGrant(val *network.Validator, args []string) (testutil.BufferWriter, error) {
	cmd := cli.NewCmdGrantAuthorization()
	clientCtx := val.ClientCtx
	return clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
}
