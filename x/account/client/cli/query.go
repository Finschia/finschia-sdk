package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/line/link-modules/x/account/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"

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
		Use:   "txs",
		Short: "Query for paginated transactions that match a set of tags",
		Long: strings.TrimSpace(`
Search for transactions that match the exact given tags where results are paginated.

Example:
$ <appcli> query txs --tags 'message.action:send&message.sender:yoshi' --page 1 --limit 30
$ <appcli> query txs --tags 'message.action:send&message.sender:yoshi' --height-from 77 --height-to 79

You can also search by height range without tags:
$ <appcli> query txs --height-from 77 --height-to 79
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			tagsStr := viper.GetString(flagTags)
			tagsStr = strings.Trim(tagsStr, "'")

			var tags []string
			if len(tagsStr) > 0 {
				tags = strings.Split(tagsStr, "&")
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

var (
	DefaultBlockFetchSize int64 = 20
)

func QueryBlockWithTxResponsesCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-with-tx-result [from_block_height] [fetchsize]",
		Short: "Get verified data for the block and tx and tx_result from given `from_block_height` to `fetchsize`.",
		Long: "Up to 20 Items can be returned, and more are ignored. \n" +
			"The Default fetchsize is 20 and if there are not enough blocks in the fetchsize requested from from_block_height, \n" +
			"It will respond to the latest block height from from_block_height param. \n" +
			"You can know latest block height by latest_block_height property of result. \n" +
			"The direction of hasMore is from low to high blockHeight. \n" +
			"Usage:\n" +
			"  linkcli query block-with-tx-result 1 10\n" +
			"  linkcli query block-with-tx-result 1 30 (it will return 20 Items)\n" +
			"  linkcli query block-with-tx-result 1",

		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			fromBlockHeight, fetchSize := parseCmdParams(args)

			latestBlockHeight, err := utils.LatestBlockHeight(cliCtx)
			if err != nil {
				return err
			}

			if fromBlockHeight >= latestBlockHeight {
				return fmt.Errorf("the block height does not exist. Requested: %d, Latest: %d", fromBlockHeight, latestBlockHeight)
			}

			blockWithTxReponses, err := utils.BlockWithTxResponses(cliCtx, latestBlockHeight, fromBlockHeight, fetchSize)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(blockWithTxReponses)
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

func parseCmdParams(args []string) (int64, int64) {
	fromBlockHeight, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		panic(err)
	}
	if len(args) == 1 {
		return fromBlockHeight, DefaultBlockFetchSize
	}
	fetchSize, err := strconv.ParseInt(args[1], 10, 8)
	if err != nil {
		panic(err)
	}
	if fetchSize > DefaultBlockFetchSize {
		fetchSize = DefaultBlockFetchSize
	}
	return fromBlockHeight, fetchSize
}
