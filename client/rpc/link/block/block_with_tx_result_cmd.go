package block

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"strconv"
)

var (
	DefaultBlockFetchSize int8 = 20
)

func WithTxResultCommand(cdc *amino.Codec) *cobra.Command {

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
		RunE: printByParams(cdc),
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

func printByParams(cdc *amino.Codec) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cliCtx := context.NewCLIContext().WithCodec(cdc)

		data, err := process(args, NewBlockUtil(cliCtx))
		if err != nil {
			return err
		}

		fmt.Println(string(data))
		return nil
	}
}

func process(args []string, util *Util) ([]byte, error) {
	fromBlockHeight, fetchSize := parseCmdParams(args)
	client, err := util.lcliCtx.GetNode()
	if err != nil {
		return nil, err
	}
	latestBlock, err := client.Block(nil)
	if err != nil {
		return nil, err
	}

	blockWithTxResults, err := util.fetchByBlockHeights(&latestBlock.Block.Height, &fromBlockHeight, &fetchSize)
	if err != nil {
		return nil, err
	}

	return util.IndentJSON(blockWithTxResults)
}

func parseCmdParams(args []string) (fromBlockHeight int64, fetchSizeInt8 int8) {
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
	fetchSizeInt8 = int8(fetchSize)
	if fetchSizeInt8 > DefaultBlockFetchSize {
		fetchSizeInt8 = DefaultBlockFetchSize
	}
	return
}
