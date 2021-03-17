package testutil

import (
	"context"
	"fmt"

	ostcfg "github.com/line/ostracon/config"
	"github.com/line/ostracon/libs/cli"
	"github.com/line/ostracon/libs/log"
	"github.com/spf13/viper"

	"github.com/line/lbm-sdk/v2/client"
	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/server"
	"github.com/line/lbm-sdk/v2/testutil"
	"github.com/line/lbm-sdk/v2/types/module"
	genutilcli "github.com/line/lbm-sdk/v2/x/genutil/client/cli"
)

func ExecInitCmd(testMbm module.BasicManager, home string, cdc codec.JSONMarshaler) error {
	logger := log.NewNopLogger()
	cfg, err := CreateDefaultTendermintConfig(home)
	if err != nil {
		return err
	}

	cmd := genutilcli.InitCmd(testMbm, home)
	serverCtx := server.NewContext(viper.New(), cfg, logger)
	clientCtx := client.Context{}.WithJSONMarshaler(cdc).WithHomeDir(home)

	_, out := testutil.ApplyMockIO(cmd)
	clientCtx = clientCtx.WithOutput(out)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
	ctx = context.WithValue(ctx, server.ServerContextKey, serverCtx)

	cmd.SetArgs([]string{"appnode-test", fmt.Sprintf("--%s=%s", cli.HomeFlag, home)})

	return cmd.ExecuteContext(ctx)
}

func CreateDefaultTendermintConfig(rootDir string) (*ostcfg.Config, error) {
	conf := ostcfg.DefaultConfig()
	conf.SetRoot(rootDir)
	ostcfg.EnsureRoot(rootDir)

	if err := conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	return conf, nil
}
