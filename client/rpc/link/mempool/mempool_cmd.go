package mempool

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// MempoolCmd can return different results for each node.
func MempoolCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mempool",
		Short: "Query remote node for mempool",
	}

	cmd.AddCommand(
		NumUnconfirmedTxsCmd(cdc),
		UnconfirmedTxsCmd(cdc),
	)

	return cmd
}

// NumUnconfirmedTxsCmd returns the command to return the number of unconfirmed txs from the mempool.
func NumUnconfirmedTxsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "num-unconfirmed-txs",
		Short: "Get the number of unconfrimed txs from mempool",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			output, err := getNumUnconfirmedTxsCmd(cliCtx)
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
	cmd.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
	return cmd
}

func getNumUnconfirmedTxsCmd(cliCtx context.CLIContext) ([]byte, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	mempoolClient, ok := node.(client.MempoolClient)
	if !ok {
		return nil, errors.New("node does not have mempool client")
	}

	res, err := mempoolClient.NumUnconfirmedTxs()
	if err != nil {
		return nil, err
	}

	cdc := cliCtx.Codec

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(res, "", "  ")
	}

	return cdc.MarshalJSON(res)
}

// UnconfirmedTxsCmd returns the command to return the unconfirmed txs from the mempool.
func UnconfirmedTxsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unconfirmed-txs",
		Short: "Get unconfrimed txs from mempool",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			limit := viper.GetInt(flagLimit)
			hash := viper.GetBool(flagHash)

			output, err := getUnconfirmedTxsCmd(cliCtx, limit, hash)
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
	cmd.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
	cmd.Flags().Int32(flagLimit, defaultLimit, "Maximum number of unconfirmed transactions to return (max 100)")
	cmd.Flags().Bool(flagHash, false, "Return tx as hash")
	return cmd
}

func getUnconfirmedTxsCmd(cliCtx context.CLIContext, limit int, hash bool) ([]byte, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	mempoolClient, ok := node.(client.MempoolClient)
	if !ok {
		return nil, errors.New("node does not have mempool client")
	}

	res, err := mempoolClient.UnconfirmedTxs(limit)
	if err != nil {
		return nil, err
	}

	cdc := cliCtx.Codec

	txs := make([]json.RawMessage, len(res.Txs))
	for idx, tx := range res.Txs {
		var txJSON json.RawMessage
		if hash {
			txHash := common.HexBytes(tx.Hash())
			txJSON = cdc.MustMarshalJSON(txHash)
		} else {
			var stdTx auth.StdTx
			cdc.MustUnmarshalBinaryLengthPrefixed(tx, &stdTx)
			txJSON = cdc.MustMarshalJSON(stdTx)
		}
		txs[idx] = txJSON
	}

	txsJSON := cdc.MustMarshalJSON(txs)
	resJSON := cdc.MustMarshalJSON(res)

	var unconfirmedTxs map[string]json.RawMessage
	cdc.MustUnmarshalJSON(resJSON, &unconfirmedTxs)
	unconfirmedTxs["txs"] = txsJSON

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(unconfirmedTxs, "", "  ")
	}

	return cdc.MarshalJSON(unconfirmedTxs)
}
