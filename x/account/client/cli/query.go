package cli

import (
	"fmt"
	"strings"

	"github.com/line/link/x/account/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

const (
	flagTags       = "tags"
	flagPage       = "page"
	flagLimit      = "limit"
	flagHeightFrom = "height-from"
	flagHeightTo   = "height-to"
)

// *****
// Original code: `github.com/cosmos/cosmos-sdk/x/auth/client/cli/query.go`
// Difference: referring import path of `utils`
// *****

// QueryTxsByEventsCmd returns a command to search through transactions by events.
func QueryTxsByEventsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs --tags '<key>:<value>[&<key>:<value>]'",
		Short: "Query for paginated transactions that match a set of tags",
		Long: strings.TrimSpace(`
Search for transactions that match the exact given tags where results are paginated.

Example:
$ <appcli> query txs --tags 'message.action:send&message.sender:yoshi' --page 1 --limit 30
$ <appcli> query txs --tags 'message.action:send&message.sender:yoshi' --height-from 77 --height-to 79
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			tagsStr := viper.GetString(flagTags)
			tagsStr = strings.Trim(tagsStr, "'")

			var tags []string
			if strings.Contains(tagsStr, "&") {
				tags = strings.Split(tagsStr, "&")
			} else {
				tags = append(tags, tagsStr)
			}

			var tmTags []string
			for _, tag := range tags {
				if !strings.Contains(tag, ":") {
					return fmt.Errorf("%s should be of the format <key>:<value>", tagsStr)
				} else if strings.Count(tag, ":") > 1 {
					return fmt.Errorf("%s should only contain one <key>:<value> pair", tagsStr)
				}

				keyValue := strings.Split(tag, ":")
				if keyValue[0] == tmtypes.TxHeightKey {
					tag = fmt.Sprintf("%s=%s", keyValue[0], keyValue[1])
				} else {
					tag = fmt.Sprintf("%s='%s'", keyValue[0], keyValue[1])
				}

				tmTags = append(tmTags, tag)
			}

			heightFrom := viper.GetInt64(flagHeightFrom)
			if heightFrom > 0 {
				tag := fmt.Sprintf("%s>=%d", tmtypes.TxHeightKey, heightFrom)
				tmTags = append(tmTags, tag)
			}

			heightTo := viper.GetInt64(flagHeightTo)
			if heightTo > 0 {
				tag := fmt.Sprintf("%s<=%d", tmtypes.TxHeightKey, heightTo)
				tmTags = append(tmTags, tag)
			}

			page := viper.GetInt(flagPage)
			limit := viper.GetInt(flagLimit)

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txs, err := utils.QueryTxsByEvents(cliCtx, tmTags, page, limit)
			if err != nil {
				return err
			}

			var output []byte
			if cliCtx.Indent {
				output, err = cdc.MarshalJSONIndent(txs, "", "  ")
			} else {
				output, err = cdc.MarshalJSON(txs)
			}

			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
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
	cmd.Flags().String(flagTags, "", "tag:value list of tags that must match")
	cmd.Flags().Uint32(flagPage, rest.DefaultPage, "Query a specific page of paginated results")
	cmd.Flags().Uint32(flagLimit, rest.DefaultLimit, "Query number of transactions results per page returned")
	err = cmd.MarkFlagRequired(flagTags)
	if err != nil {
		panic(err)
	}

	cmd.Flags().Int64(flagHeightFrom, 0, "Filter from a specific block height")
	cmd.Flags().Int64(flagHeightTo, 0, "Filter to a specific block height")

	return cmd
}

// QueryTxCmd implements the default command for a tx query.
func QueryTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [hash]",
		Short: "Query for a transaction by hash in a committed block",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			output, err := utils.QueryTx(cliCtx, args[0])
			if err != nil {
				return err
			}

			if output.Empty() {
				return fmt.Errorf("no transaction found with hash %s", args[0])
			}

			return cliCtx.PrintOutput(output)
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
