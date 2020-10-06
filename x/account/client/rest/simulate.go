package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// SimulateReq defines a tx simulating request.
type SimulateReq struct {
	Tx            types.StdTx `json:"tx" yaml:"tx"`
	GasAdjustment string      `json:"gas_adjustment"`
}

type ABCIErrorResponse struct {
	Codespace string `json:"codespace"`
	Code      uint32 `json:"code"`
	Error     string `json:"error"`
}

// SimulateTxRequest implements a tx simulating handler that is responsible
// for simulating a valid and signed tx.
func SimulateTxRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SimulateReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, req.GasAdjustment, flags.DefaultGasAdjustment)
		if !ok {
			// ParseFloat64OrReturnBadRequest has already written error response.
			return
		}

		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid gas adjustment: %g", gasAdj))
			return
		}

		_, adjusted, err := utils.CalculateGas(cliCtx.QueryWithData, cliCtx.Codec, txBytes, gasAdj)
		if err != nil {
			if ctxtErr, ok := err.(*context.Error); ok {
				WriteABCIErrorResponse(w, http.StatusInternalServerError, cliCtx.Codec, ctxtErr)
			} else {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		rest.WriteSimulationResponse(w, cliCtx.Codec, adjusted)
	}
}

// nolint: errcheck
func WriteABCIErrorResponse(w http.ResponseWriter, status int, cdc *codec.Codec, err *context.Error) {
	errBody := cdc.MustMarshalJSON(
		ABCIErrorResponse{
			Codespace: err.Codespace,
			Code:      err.Code,
			Error:     err.Message,
		},
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(errBody)
}
