package main

import (
	"fmt"
	"os"
	"path"

	"github.com/line/link/types"
	"github.com/line/link/x/account"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/line/link/app"
	"github.com/line/link/client"
	"github.com/line/link/version"
	authclient "github.com/line/link/x/auth/client"
	"github.com/line/link/x/bank"
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
		Use:   "linkcli",
		Short: "Command line interface for interacting with linkd",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().Bool(flagTestnet, false, "Run with testnet mode. The address prefix becomes tlink if this flag is set.")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		client.StatusCommand(),
		client.MempoolCmd(cdc),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		client.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		client.Commands(),
		client.LineBreak,
		client.NewCompletionCmd(rootCmd, true),
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
		client.QueryGenesisAccountCmd(cdc),
		client.LineBreak,
		client.ValidatorCommand(cdc),
		client.BlockCommand(cdc),
		client.BlockWithResultCommand(cdc),
		client.QueryGenesisTxCmd(cdc),
		authclient.QueryTxsByEventsCmd(cdc),
		authclient.QueryTxCmd(cdc),
		client.LineBreak,
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
		bank.SendTxCmd(cdc),
		client.LineBreak,
		account.CreateAccountTxCmd(cdc),
		account.EmptyTxCmd(cdc),
		client.LineBreak,
		authclient.GetSignCommand(cdc),
		authclient.GetMultiSignCommand(cdc),
		client.LineBreak,
		authclient.GetBroadcastCommand(cdc),
		authclient.GetEncodeCommand(cdc),
		client.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank and account commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName || cmd.Use == account.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *client.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
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

	if viper.GetBool(flagTestnet) {
		types.SetTestnetMode()
	}

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr(), types.Bech32PrefixAccPub())
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr(), types.Bech32PrefixValPub())
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr(), types.Bech32PrefixConsPub())
	config.SetCoinType(types.CoinType)
	config.SetFullFundraiserPath(types.FullFundraiserPath)
	config.Seal()

	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
