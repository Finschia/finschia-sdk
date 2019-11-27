package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	flagHash = "hash"

	defaultLimit = 30
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
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
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

// REST handler for num unconfirmed txs
func NumUnconfirmedTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getNumUnconfirmedTxsCmd(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
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
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
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

	var txs []json.RawMessage
	for _, tx := range res.Txs {
		var txJSON json.RawMessage
		if hash {
			txHash := common.HexBytes(tx.Hash())
			txJSON = cdc.MustMarshalJSON(txHash)
		} else {
			var stdTx auth.StdTx
			cdc.MustUnmarshalBinaryLengthPrefixed(tx, &stdTx)
			txJSON = cdc.MustMarshalJSON(stdTx)
		}
		txs = append(txs, txJSON)
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

// REST handler for unconfirmed txs
func UnconfirmedTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, hash, err := parseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		output, err := getUnconfirmedTxsCmd(cliCtx, limit, hash)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}

func parseHTTPArgs(r *http.Request) (limit int, hash bool, err error) {
	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return limit, hash, err
		}
	}

	hashStr := r.FormValue("hash")
	if hashStr == "true" {
		hash = true
	}

	return limit, hash, nil
}
