package utils

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/link-chain/link/types"
	"github.com/link-chain/link/x/auth/client/utils/mocks"
)

type mockNodeResponses struct {
	resTx       *ctypes.ResultTx
	resBlock    *ctypes.ResultBlock
	resTxSearch *ctypes.ResultTxSearch
}

func setupMockNodeResponses(hashString string, height int64, index uint32, cdc *codec.Codec) mockNodeResponses {
	hash, _ := hex.DecodeString(hashString)

	stdTx := &auth.StdTx{
		Memo: "empty tx",
	}

	bz, _ := cdc.MarshalBinaryLengthPrefixed(stdTx)
	resTx := &ctypes.ResultTx{
		Hash:     hash,
		Height:   height,
		Index:    index,
		TxResult: abci.ResponseDeliverTx{},
		Tx:       bz,
	}
	resBlock := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{
				Time: time.Time{},
			},
		},
	}
	resTxSearch := &ctypes.ResultTxSearch{
		Txs:        []*ctypes.ResultTx{resTx},
		TotalCount: 1,
	}

	return mockNodeResponses{
		resTx:       resTx,
		resBlock:    resBlock,
		resTxSearch: resTxSearch,
	}
}

func setupCodec() *codec.Codec {
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	return cdc
}

//noinspection GoBoolExpressions
func TestQueryTxsByEventsResponseContainsIndexAndCode(t *testing.T) {
	hashString := "15E23C9F72602046D86BC9F0ECAE53E43A8206C113A29D94454476B9887AAB7F"
	height := int64(100)
	index := uint32(10)

	cdc := setupCodec()
	nodeResponses := setupMockNodeResponses(hashString, height, index, cdc)

	mockClient := &mocks.Client{}
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	mockClient.On("TxSearch", "tx.height=0", !cliCtx.TrustNode, 1, 30).Return(nodeResponses.resTxSearch, nil)
	mockClient.On("Block", &nodeResponses.resTx.Height).Return(nodeResponses.resBlock, nil)

	res, _ := QueryTxsByEvents(cliCtx, []string{"tx.height=0"}, 1, 30)

	assert.Equal(t, 1, res.Count)
	assert.Equal(t, height, res.Txs[0].Height)
	assert.Equal(t, index, res.Txs[0].Index)
	assert.Equal(t, hashString, res.Txs[0].TxHash)
	assert.Equal(t, uint32(0), res.Txs[0].Code)
}

//noinspection GoBoolExpressions
func TestQueryTxResponseContainsIndexAndCode(t *testing.T) {
	hashString := "15E23C9F72602046D86BC9F0ECAE53E43A8206C113A29D94454476B9887AAB7F"
	height := int64(100)
	index := uint32(10)

	cdc := setupCodec()
	nodeResponses := setupMockNodeResponses(hashString, height, index, cdc)

	mockClient := &mocks.Client{}
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	hash, _ := hex.DecodeString(hashString)
	mockClient.On("Tx", hash, !cliCtx.TrustNode).Return(nodeResponses.resTx, nil)
	mockClient.On("Block", &nodeResponses.resTx.Height).Return(nodeResponses.resBlock, nil)

	res, _ := QueryTx(cliCtx, hashString)

	assert.Equal(t, height, res.Height)
	assert.Equal(t, index, res.Index)
	assert.Equal(t, hashString, res.TxHash)
	assert.Equal(t, uint32(0), res.Code)
}

func TestQueryGenesisTxs(t *testing.T) {
	cdc := setupCodec()
	mockClient := &mocks.Client{}
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	// exist genesis tx
	genesisDoc := tmtypes.GenesisDoc{
		AppState: json.RawMessage(`{"genutil":{"gentxs":[{"type":"cosmos-sdk/StdTx","value":{"memo":"test_genesis"}}]}}`),
	}
	genesisResult := ctypes.ResultGenesis{
		Genesis: &genesisDoc,
	}
	mockClient.On("Genesis").Return(&genesisResult, nil)

	genesisTxs, _ := QueryGenesisTx(cliCtx)
	assert.NotEmpty(t, genesisTxs)

	// not exist genesis tx
	genesisDoc = tmtypes.GenesisDoc{
		AppState: json.RawMessage(`{"genutil":{"gentxs":[]}`),
	}
	genesisResult = ctypes.ResultGenesis{
		Genesis: &genesisDoc,
	}
	mockClient.On("Genesis").Return(&genesisResult, nil)
	genesisTxs, _ = QueryGenesisTx(cliCtx)
	assert.Empty(t, genesisTxs)

}

func TestQueryGenesisAccount(t *testing.T) {

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)

	cdc := setupCodec()
	mockClient := &mocks.Client{}
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	// exist genesis account
	genesisDoc := tmtypes.GenesisDoc{
		AppState: json.RawMessage(`{"accounts":[{"address":"link19rqsvml8ldr0yrhaewgv9smcdvrew5pah9j5t5","coins":[]}]}`),
	}
	genesisResult := ctypes.ResultGenesis{
		Genesis: &genesisDoc,
	}
	mockClient.On("Genesis").Return(&genesisResult, nil)
	res, err := QueryGenesisAccount(cliCtx, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Accounts))

	// no exist page
	res, err = QueryGenesisAccount(cliCtx, 2, 1)
	assert.Error(t, err)
	assert.Equal(t, 0, len(res.Accounts))

	// page=0
	res, err = QueryGenesisAccount(cliCtx, 0, 1)
	assert.Error(t, err)
	assert.Equal(t, 0, len(res.Accounts))

	// page=-1
	res, err = QueryGenesisAccount(cliCtx, -1, 1)
	assert.Error(t, err)
	assert.Equal(t, 0, len(res.Accounts))

	// limit=-1
	res, err = QueryGenesisAccount(cliCtx, 1, -1)
	assert.Error(t, err)
	assert.Equal(t, 0, len(res.Accounts))
}

func TestParseHTTPArgs(t *testing.T) {
	reqE0 := mustNewRequest(t, "", "/", nil)

	req0 := mustNewRequest(t, "", "/?foo=faa", nil)

	req1 := mustNewRequest(t, "", "/?foo=faa&limit=5", nil)
	req2 := mustNewRequest(t, "", "/?foo=faa&page=5", nil)
	req3 := mustNewRequest(t, "", "/?foo=faa&page=5&limit=5", nil)

	reqE1 := mustNewRequest(t, "", "/?foo=faa&page=-1", nil)
	reqE2 := mustNewRequest(t, "", "/?foo=faa&limit=-1", nil)
	reqE3 := mustNewRequest(t, "", "/?page=5&limit=5", nil)

	req4 := mustNewRequest(t, "", "/?foo=faa&height.from=1", nil)
	req5 := mustNewRequest(t, "", "/?foo=faa&height.to=1", nil)

	reqE4 := mustNewRequest(t, "", "/?foo=faa&height.from=-1", nil)
	reqE5 := mustNewRequest(t, "", "/?foo=faa&height.to=-1", nil)

	defaultPage := rest.DefaultPage
	defaultLimit := rest.DefaultLimit

	tests := []struct {
		name  string
		req   *http.Request
		w     http.ResponseWriter
		tags  []string
		page  int
		limit int
		err   bool
	}{
		{"no params", reqE0, httptest.NewRecorder(), []string{}, defaultPage, defaultLimit, true},

		{"tags", req0, httptest.NewRecorder(), []string{"foo='faa'"}, defaultPage, defaultLimit, false},

		{"limit", req1, httptest.NewRecorder(), []string{"foo='faa'"}, defaultPage, 5, false},
		{"page", req2, httptest.NewRecorder(), []string{"foo='faa'"}, 5, defaultLimit, false},
		{"page and limit", req3, httptest.NewRecorder(), []string{"foo='faa'"}, 5, 5, false},

		{"error page 0", reqE1, httptest.NewRecorder(), []string{}, defaultPage, defaultLimit, true},
		{"error limit 0", reqE2, httptest.NewRecorder(), []string{}, defaultPage, defaultLimit, true},
		{"no tags", reqE3, httptest.NewRecorder(), []string{}, defaultPage, defaultLimit, true},

		{"height from", req4, httptest.NewRecorder(), []string{"foo='faa'", "tx.height>=1"}, defaultPage, defaultLimit, false},
		{"height to", req5, httptest.NewRecorder(), []string{"foo='faa'", "tx.height<=1"}, defaultPage, defaultLimit, false},

		{"error height from", reqE4, httptest.NewRecorder(), []string{}, defaultPage, defaultLimit, true},
		{"error height to", reqE5, httptest.NewRecorder(), []string{}, defaultPage, defaultLimit, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tags, page, limit, err := ParseHTTPArgs(tt.req)
			if tt.err {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.tags, tags)
				require.Equal(t, tt.page, page)
				require.Equal(t, tt.limit, limit)
			}
		})
	}
}

func mustNewRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	err = req.ParseForm()
	require.NoError(t, err)
	return req
}
