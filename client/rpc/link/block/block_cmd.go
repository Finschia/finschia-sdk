package block

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

//Command returns the verified block data for a given heights
func Command(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [height] [extended]",
		Short: "Get verified data for a block at a given height. Extended is optional to get translated txs.",
		Long: "Usage:\n" +
			"  linkcli query block (the latest block)\n" +
			"  linkcli query block 23 (the block #23)\n" +
			"  linkcli query block 7 extended (the block #7 with translated transaction information)\n" +
			"  linkcli query block extended (the latest block with translated transaction information)",
		Args: cobra.MaximumNArgs(2),
		RunE: printBlock(cdc),
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

// CMD
func printBlock(cdc *amino.Codec) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// optional params
		height, isConvertBlockTxToJSON := parseCmdParamsOptional(args)

		cliCtx := context.NewCLIContext().WithCodec(cdc)
		output, err := GetBlock(NewBlockUtil(cliCtx), height, isConvertBlockTxToJSON)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	}
}

// parse optional arguments from CMD
// returns nil and an empty string if not set
func parseCmdParamsOptional(args []string) (*int64, bool) {
	var height *int64
	var extended bool

	if len(args) == 1 {
		// [extended] only
		if args[0] == "extended" {
			return nil, true
		} else {
			extended = false
		}

		// [height] only
		h, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, extended
		} else if h > 0 {
			tmp := int64(h)
			height = &tmp
			return height, extended
		}
	} else if len(args) == 2 {
		// [height] [extended]
		h, err := strconv.Atoi(args[0])
		if err != nil {
			height = nil
		} else if h > 0 {
			tmp := int64(h)
			height = &tmp
		}

		if args[1] == "extended" {
			extended = true
		}

		return height, extended
	}

	return nil, false
}

func GetBlock(util *Util, cursor *int64, isConvertBlockTxToJSON bool) ([]byte, error) {
	// get the node
	node, err := util.lcliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	block, err := node.Block(cursor)
	if err != nil {
		return nil, err
	}

	err = util.ValidateBlock(block)
	if err != nil {
		return nil, err
	}

	blockResponse, err := util.Indent(block)
	if err != nil {
		return blockResponse, err
	}
	// deserialize transactions from Amino and serialize to JSON
	if isConvertBlockTxToJSON {

		byteTxs, err := aminoToJsonTxs(util.lcliCtx.Codec(), &block.Block.Data.Txs)
		if err != nil {
			return nil, err
		}

		// deserialize transactions from Amino and serialize to JSON
		txInjectedBlockResponse, err := util.InjectByteToJsonTxs(blockResponse, byteTxs)
		if err != nil {
			return nil, err
		}
		return json.Marshal(txInjectedBlockResponse)
	}
	return blockResponse, err
}

// deserialize transactions from Amino and serialize to JSON
func aminoToJsonTxs(cdc *codec.Codec, txs *types.Txs) (result [][]byte, err error) {
	// for each transaction in the transaction slice
	for _, aminoSerializedTx := range *txs {
		// tx is amino serialized but not base64 encoded in memory
		// tx is encoded to base64 while serializing to JSON
		aminoDeserializedTx, err := parseTx(cdc, aminoSerializedTx)
		if err != nil {
			return nil, err
		}

		// Serialize to JSON
		if bz, err := cdc.MarshalJSON(aminoDeserializedTx); err != nil {
			return nil, err
		} else {
			// collect as a slice
			result = append(result, bz)
		}
	}
	return
}

func parseTx(cdc *codec.Codec, txBytes []byte) (tx auth.StdTx, err error) {
	if err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx); err != nil {
		return tx, err
	}
	return
}
