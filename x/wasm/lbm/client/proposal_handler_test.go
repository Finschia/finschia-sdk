package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
)

func TestGovRestHandlers(t *testing.T) {
	type dict map[string]interface{}
	var (
		anyAddress = "link18vd8fpwxzck93qlwghaj6arh4p7c5n89fvcmzu"
		aBaseReq   = dict{
			"from":           anyAddress,
			"memo":           "rest test",
			"chain_id":       "testing",
			"account_number": "1",
			"sequence":       "1",
			"fees":           []dict{{"denom": "cony", "amount": "1000000"}},
		}
	)

	encodingConfig := wasmkeeper.MakeEncodingConfig(t)
	clientCtx := client.Context{}.
		WithJSONCodec(encodingConfig.Marshaler).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithChainID("testing")

	// router setup as in gov/client/rest/tx.go
	propSubRtr := mux.NewRouter().PathPrefix("/gov/proposals").Subrouter()
	for _, ph := range ProposalHandlers {
		r := ph.RESTHandler(clientCtx)
		propSubRtr.HandleFunc(fmt.Sprintf("/%s", r.SubRoute), r.Handler).Methods("POST")
	}

	specs := map[string]struct {
		srcBody dict
		srcPath string
		expCode int
	}{
		"deactivate contract": {
			srcPath: "/gov/proposals/deactivate_contract",
			srcBody: dict{
				"title":       "Test Proposal",
				"description": "My proposal",
				"deposit":     []dict{{"denom": "cony", "amount": "10"}},
				"proposer":    "link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5",
				"contract":    "link1ghekyjucln7y67ntx7cf27m9dpuxxemnqk82wt",
				"base_req":    aBaseReq,
			},
			expCode: http.StatusOK,
		},
		"deactivate contract with no contract": {
			srcPath: "/gov/proposals/deactivate_contract",
			srcBody: dict{
				"title":       "Test Proposal",
				"description": "My proposal",
				"deposit":     []dict{{"denom": "cony", "amount": "10"}},
				"proposer":    "link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5",
				"base_req":    aBaseReq,
			},
			expCode: http.StatusBadRequest,
		},
		"activate contract": {
			srcPath: "/gov/proposals/activate_contract",
			srcBody: dict{
				"title":       "Test Proposal",
				"description": "My proposal",
				"deposit":     []dict{{"denom": "cony", "amount": "10"}},
				"proposer":    "link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5",
				"contract":    "link1ghekyjucln7y67ntx7cf27m9dpuxxemnqk82wt",
				"base_req":    aBaseReq,
			},
			expCode: http.StatusOK,
		},
		"activate contract with no contract": {
			srcPath: "/gov/proposals/activate_contract",
			srcBody: dict{
				"title":       "Test Proposal",
				"description": "My proposal",
				"deposit":     []dict{{"denom": "cony", "amount": "10"}},
				"proposer":    "link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5",
				"base_req":    aBaseReq,
			},
			expCode: http.StatusBadRequest,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			src, err := json.Marshal(spec.srcBody)
			require.NoError(t, err)

			// when
			r := httptest.NewRequest("POST", spec.srcPath, bytes.NewReader(src))
			w := httptest.NewRecorder()
			propSubRtr.ServeHTTP(w, r)

			// then
			require.Equal(t, spec.expCode, w.Code, w.Body.String())
		})
	}
}
