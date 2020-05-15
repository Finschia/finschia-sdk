package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/cli"
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
	)
	rootCmd.PersistentFlags().Bool(cli.FlagTestnet, true, "Set whether the target chain is a testnet or not.")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig(cmd *cobra.Command) error {
	testnet, err := cmd.PersistentFlags().GetBool(cli.FlagTestnet)
	if err != nil {
		return err
	}

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(testnet), linktypes.Bech32PrefixAccPub(testnet))
	config.SetBech32PrefixForValidator(linktypes.Bech32PrefixValAddr(testnet), linktypes.Bech32PrefixValPub(testnet))
	config.SetBech32PrefixForConsensusNode(linktypes.Bech32PrefixConsAddr(testnet), linktypes.Bech32PrefixConsPub(testnet))
	config.SetCoinType(linktypes.CoinType)
	config.SetFullFundraiserPath(linktypes.FullFundraiserPath)
	config.Seal()

	return nil
}
