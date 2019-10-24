package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
)

//BlockCommand returns the verified block data for a given heights
func BlockCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [height] [extended]",
		Short: "Get verified data for a the block at given height. Extended is an optional to get translated txs.",
		Long: "Usage:\n" +
			"  linkcli query block (the latest block)\n" +
			"  linkcli query block 23 (the block #23)\n" +
			"  linkcli query block 7 extended (the block #7 with translated transaction information)\n" +
			"  linkcli query block extended (the latest block with translated transaction information)",
		Args: cobra.MaximumNArgs(2),
		RunE: getPrintBlock(cdc),
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

func getBlock(cliCtx context.CLIContext, height *int64, extended bool) ([]byte, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(height)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		check, err := cliCtx.Verify(res.Block.Height)
		if err != nil {
			return nil, err
		}

		err = tmliteProxy.ValidateBlockMeta(res.BlockMeta, check)
		if err != nil {
			return nil, err
		}

		err = tmliteProxy.ValidateBlock(res.Block, check)
		if err != nil {
			return nil, err
		}
	}

	var blockResponse []byte
	if cliCtx.Indent {
		blockResponse, err = codec.Cdc.MarshalJSONIndent(res, "", "  ")
	} else {
		blockResponse, err = codec.Cdc.MarshalJSON(res)
	}

	// deserialize transactions from Amino and serialize to JSON
	if extended {
		var translatedTxs [][]byte
		translatedTxs, err = aminoToJsonTxs(&res.Block.Data.Txs, cliCtx)
		if err != nil {
			return nil, err
		}

		// deserialize transactions from Amino and serialize to JSON
		return injectTranslatedTxs(blockResponse, translatedTxs)
	}

	return blockResponse, err
}

// inject translated transactions to block data
func injectTranslatedTxs(blockResponse []byte, translatedTxs [][]byte) ([]byte, error) {
	var block map[string]interface{}
	var toInject []map[string]interface{}

	// load block response as a map
	e := json.Unmarshal(blockResponse, &block)
	if e != nil {
		return nil, e
	}

	// load translated txs as a map
	for _, translatedTx := range translatedTxs {
		var translated map[string]interface{}

		e = json.Unmarshal(translatedTx, &translated)
		if e != nil {
			return nil, e
		}
		// generate a slice to inject
		toInject = append(toInject, translated)
	}

	// inject the translated transactions
	block["block"].(map[string]interface{})["data"].(map[string]interface{})["txs"] = toInject

	return json.Marshal(block)
}

// deserialize transactions from Amino and serialize to JSON
func aminoToJsonTxs(txs *types.Txs, cliCtx context.CLIContext) ([][]byte, error) {
	var result [][]byte
	cdc := cliCtx.Codec

	// for each transaction in the transaction slice
	for _, aminoSerializedTx := range *txs {
		// tx is amino serialized but not base64 encoded in memory
		// tx is encoded to base64 while serializing to JSON
		aminoDeserializedTx, err := parseTx(cdc, aminoSerializedTx)
		if err != nil {
			return nil, err
		}

		// Serialize to JSON
		bz, err := cdc.MarshalJSON(aminoDeserializedTx)
		if err != nil {
			return nil, err
		}

		// collect as a slice
		result = append(result, bz)
	}

	return result, nil
}

func parseTx(cdc *codec.Codec, txBytes []byte) (auth.StdTx, error) {
	var tx auth.StdTx

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return tx, err
	}

	return tx, nil
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

// CMD

func getPrintBlock(cdc *amino.Codec) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// optional params
		height, extended := parseCmdParamsOptional(args)

		output, err := getBlock(context.NewCLIContext().WithCodec(cdc), height, extended)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	}
}

// REST

// REST handler to get a block
func BlockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				"couldn't parse block height. Assumed format is '/block/{height}'.")
			return
		}

		chainHeight, err := GetChainHeight(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "failed to parse chain height")
			return
		}

		if height > chainHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound, "requested block height is bigger then the chain length")
			return
		}

		output, err := getBlock(cliCtx, &height, parseExtended(r))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}

// REST handler to get the latest block
func LatestBlockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getBlock(cliCtx, nil, parseExtended(r))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}

func parseExtended(r *http.Request) bool {
	return r.URL.Query().Get("extended") == "true"
}
