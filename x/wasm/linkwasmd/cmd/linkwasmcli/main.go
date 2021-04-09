package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/line/lbm-sdk/v2/x/account"
	authclient "github.com/line/lbm-sdk/v2/x/account/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	"github.com/line/ostracon/libs/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/line/lbm-sdk/v2/x/coin"
	"github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/app"
	"github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/types"
)

const (
	flagTestnet = "testnet"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	rootCmd := &cobra.Command{
		Use:   "linkwasmcli",
		Short: "Command line interface for interacting with linkwasmd",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().Bool(flagTestnet, false, "Run with testnet mode. The address prefix becomes tlink if this flag is set.")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		flags.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		flags.NewCompletionCmd(rootCmd, true),
		version.Cmd,
	)

	// Add flags and prefix all env exposed with GA
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authclient.GetAccountCmd(cdc),
		flags.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authclient.QueryTxsByEventsCmd(cdc),
		authclient.QueryTxCmd(cdc),
		authclient.QueryBlockWithTxResponsesCommand(cdc),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		coin.SendTxCmd(cdc),
		flags.LineBreak,
		account.CreateAccountTxCmd(cdc),
		account.EmptyTxCmd(cdc),
		flags.LineBreak,
		authclient.GetSignCommand(cdc),
		authclient.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authclient.GetBroadcastCommand(cdc),
		authclient.GetEncodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank and account commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == coin.ModuleName || cmd.Use == account.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	authclient.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	testnet := viper.GetBool(flagTestnet)

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAcc(testnet), types.Bech32PrefixAccPub(testnet))
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr(testnet), types.Bech32PrefixValPub(testnet))
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr(testnet), types.Bech32PrefixConsPub(testnet))
	config.SetCoinType(types.CoinType)
	config.SetFullFundraiserPath(types.FullFundraiserPath)
	config.Seal()

	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
