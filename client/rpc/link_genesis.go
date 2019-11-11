package rpc

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	genutilrest "github.com/cosmos/cosmos-sdk/x/genutil/client/rest"

	"github.com/link-chain/link/x/auth/client/utils"

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
	_ = viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	cmd.Flags().Bool(flags.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	_ = viper.BindPFlag(flags.FlagTrustNode, cmd.Flags().Lookup(flags.FlagTrustNode))

	return cmd
}

// QueryGenesisAccountCmd returns a command to get genesis accounts.
func QueryGenesisAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genesis-accounts",
		Short: "Query for genesis accounts.",
		Long: strings.TrimSpace(`
		Query genesis accounts that occurred when the chain first started.
		Example:
		$ <appcli> query genesis-account
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
	_ = viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	cmd.Flags().Bool(flags.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	_ = viper.BindPFlag(flags.FlagTrustNode, cmd.Flags().Lookup(flags.FlagTrustNode))

	cmd.Flags().Int32(flagPage, rest.DefaultPage, "Query a specific page of paginated results")
	cmd.Flags().Int32(flagLimit, rest.DefaultLimit, "Query number of transactions results per page returned")

	return cmd
}

// QueryGenesisTxRequestHandlerFn implements a REST handler to get the genesis
func QueryGenesisTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genutilrest.QueryGenesisTxs(cliCtx, w)
	}
}

// QueryGenesisAccountRequestHandlerFn implements a REST handler to get the genesis accounts
func QueryGenesisAccountRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, page, limit, err := rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		output, err := utils.QueryGenesisAccount(cliCtx, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
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
