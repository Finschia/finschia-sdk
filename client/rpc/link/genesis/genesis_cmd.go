package genesis

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/line/link/x/account/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagPage  = "page"
	flagLimit = "limit"
)

// QueryGenesisTxCmd returns a command to get genesis transaction.
func QueryGenesisTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genesis-txs",
		Short: "Query for genesis transactions.",
		Long: strings.TrimSpace(`
		Query genesis transactions that occurred when the chain first started.
		Example:
		$ <appcli> query genesis-txs
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txs, err := utils.QueryGenesisTx(cliCtx)
			if err != nil {
				return err
			}

			return print(cdc, txs, cliCtx.Indent)
		},
	}

	cmd.Flags().StringP(flags.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	err := viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	if err != nil {
		panic(err)
	}

	cmd.Flags().Bool(flags.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	err = viper.BindPFlag(flags.FlagTrustNode, cmd.Flags().Lookup(flags.FlagTrustNode))
	if err != nil {
		panic(err)
	}

	return cmd
}

// QueryGenesisAccountCmd returns a command to get genesis accounts.
func QueryGenesisAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genesis-accounts",
		Short: "Query for paginated genesis accounts.",
		Long: strings.TrimSpace(`
		Query genesis accounts that occurred when the chain first started.
		Example:
		$ <appcli> query genesis-account --page 1 --limit 30
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			page := viper.GetInt(flagPage)
			limit := viper.GetInt(flagLimit)

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accounts, err := utils.QueryGenesisAccount(cliCtx, page, limit)
			if err != nil {
				return err
			}

			return print(cdc, accounts, cliCtx.Indent)
		},
	}

	cmd.Flags().StringP(flags.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	err := viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	if err != nil {
		panic(err)
	}
	cmd.Flags().Bool(flags.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	err = viper.BindPFlag(flags.FlagTrustNode, cmd.Flags().Lookup(flags.FlagTrustNode))
	if err != nil {
		panic(err)
	}

	cmd.Flags().Int32(flagPage, rest.DefaultPage, "Query a specific page of paginated results")
	cmd.Flags().Int32(flagLimit, rest.DefaultLimit, "Query number of transactions results per page returned")

	return cmd
}

func print(cdc *codec.Codec, i interface{}, indent bool) error {
	var (
		output []byte
		err    error
	)
	if indent {
		output, err = cdc.MarshalJSONIndent(i, "", "  ")
	} else {
		output, err = cdc.MarshalJSON(i)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
