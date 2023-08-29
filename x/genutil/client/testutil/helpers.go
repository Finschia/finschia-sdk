package testutil

import (
	"context"
	"fmt"

	ostcfg "github.com/Finschia/ostracon/config"
	"github.com/Finschia/ostracon/libs/cli"
	"github.com/Finschia/ostracon/libs/log"
	"github.com/spf13/viper"

	"github.com/Finschia/finschia-rdk/server"
	genutilcli "github.com/Finschia/finschia-rdk/x/genutil/client/cli"
	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/testutil"
	"github.com/Finschia/finschia-sdk/types/module"
)

func ExecInitCmd(testMbm module.BasicManager, home string, cdc codec.Codec) error {
	logger := log.NewNopLogger()
	cfg, err := CreateDefaultTendermintConfig(home)
	if err != nil {
		return err
	}

	cmd := genutilcli.InitCmd(testMbm, home)
	serverCtx := server.NewContext(viper.New(), cfg, logger)
	clientCtx := client.Context{}.WithCodec(cdc).WithHomeDir(home)

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
