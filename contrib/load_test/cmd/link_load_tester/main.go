package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/cli"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b) + "/../.."
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(basepath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	rootCmd := &cobra.Command{
		Use:   "link-load-tester",
		Short: "Command line interface for load tester of link",
	}

	rootCmd.AddCommand(
		cli.RunSlaveCmd(),
		cli.PrepareCmd(),
		cli.StartCmd(),
		cli.ReportCmd(),
	)
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig()
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() error {
	config := sdk.GetConfig()
	types.SetBech32Prefix(config, true)
	config.SetCoinType(linktypes.CoinType)
	config.SetFullFundraiserPath(linktypes.FullFundraiserPath)
	return nil
}
