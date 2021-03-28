package utils

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/line/lbm-sdk/v2/x/account/client/utils/mock"
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

	atypes "github.com/line/lbm-sdk/v2/x/account/client/types"
)

type mockNodeResponses struct {
	resTx       *ctypes.ResultTx
	resBlock    *ctypes.ResultBlock
	resTxSearch *ctypes.ResultTxSearch
}

// nolint:unparam
func setupMockNodeResponses(
	t *testing.T,
	cdc *codec.Codec,
	hashString string,
	height int64,
	index uint32,
	code uint32,
	codespace string,
) mockNodeResponses {
	hash, err := hex.DecodeString(hashString)
	assert.NoError(t, err)

	stdTx := &auth.StdTx{
		Memo: "empty tx",
	}

	bz, err := cdc.MarshalBinaryLengthPrefixed(stdTx)
	assert.NoError(t, err)
	resTx := &ctypes.ResultTx{
		Hash:   hash,
		Height: height,
		Index:  index,
		TxResult: abci.ResponseDeliverTx{
			Code:      code,
			Codespace: codespace,
			Log:       "[]",
		},
		Tx: bz,
	}
	resBlock := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{
				Height: height,
				Time:   time.Time{},
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

func TestQueryTxsByEventsResponseContainsIndexCodeCodespace(t *testing.T) {
	// nolint:goconst
	hashString := "15E23C9F72602046D86BC9F0ECAE53E43A8206C113A29D94454476B9887AAB7F"
	height := int64(100)
	index := uint32(10)
	code := uint32(0)
	// nolint:goconst
	codespace := "codespace"

	cdc := setupCodec()
	nodeResponses := setupMockNodeResponses(t, cdc, hashString, height, index, code, codespace)

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}
	mockClient.EXPECT().TxSearch("tx.height=0", !cliCtx.TrustNode, 1, 30, "").Return(nodeResponses.resTxSearch, nil)

	mockClient.EXPECT().Block(&nodeResponses.resTx.Height).Return(nodeResponses.resBlock, nil)

	res, err := QueryTxsByEvents(cliCtx, []string{"tx.height=0"}, 1, 30)
	assert.NoError(t, err)
	assert.Equal(t, 1, res.Count)
	assertTxResponse(t, hashString, height, index, code, codespace, res.Txs[0])
}

func TestQueryTxResponseContainsIndexCodeCodespace(t *testing.T) {
	hashString := "15E23C9F72602046D86BC9F0ECAE53E43A8206C113A29D94454476B9887AAB7F"
	height := int64(100)
	index := uint32(10)
	code := uint32(0)
	codespace := "codespace"

	cdc := setupCodec()
	nodeResponses := setupMockNodeResponses(t, cdc, hashString, height, index, code, codespace)

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	hash, err := hex.DecodeString(hashString)
	assert.NoError(t, err)
	mockClient.EXPECT().Tx(hash, !cliCtx.TrustNode).Return(nodeResponses.resTx, nil)
	mockClient.EXPECT().Block(&nodeResponses.resTx.Height).Return(nodeResponses.resBlock, nil)

	res, err := QueryTx(cliCtx, hashString)
	assert.NoError(t, err)
	assertTxResponse(t, hashString, height, index, code, codespace, res)
}

func assertTxResponse(
	t *testing.T,
	hashString string,
	height int64,
	index, code uint32,
	codespace string,
	res atypes.TxResponse,
) {
	assert.Equal(t, height, res.Height)
	assert.Equal(t, index, res.Index)
	assert.Equal(t, hashString, res.TxHash)
	assert.Equal(t, code, res.Code)
	assert.Equal(t, codespace, res.Codespace)
}

// nolint:dupl
func TestQueryTxMarshalledResponseContainsIndexCodeCodespace(t *testing.T) {
	hashString := "15E23C9F72602046D86BC9F0ECAE53E43A8206C113A29D94454476B9887AAB7F"
	height := int64(100)
	index := uint32(10)
	code := uint32(1)
	codespace := "codespace"

	cdc := setupCodec()
	nodeResponses := setupMockNodeResponses(t, cdc, hashString, height, index, code, codespace)

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	hash, err := hex.DecodeString(hashString)
	assert.NoError(t, err)
	mockClient.EXPECT().Tx(hash, !cliCtx.TrustNode).Return(nodeResponses.resTx, nil)
	mockClient.EXPECT().Block(&nodeResponses.resTx.Height).Return(nodeResponses.resBlock, nil)

	res, err := QueryTx(cliCtx, hashString)
	assert.NoError(t, err)

	out, err := cdc.MarshalJSONIndent(res, "", "  ")
	assert.NoError(t, err)

	var m map[string]interface{}
	err = json.Unmarshal(out, &m)
	assert.NoError(t, err)

	assert.Contains(t, m, "index")
	assert.Contains(t, m, "code")
	assert.Contains(t, m, "codespace")
}

// nolint:dupl
func TestQueryTxMarshalledResponseEmptyIndexCodeCodespace(t *testing.T) {
	hashString := "15E23C9F72602046D86BC9F0ECAE53E43A8206C113A29D94454476B9887AAB7F"
	height := int64(100)
	index := uint32(0)
	code := uint32(0)
	codespace := ""

	cdc := setupCodec()
	nodeResponses := setupMockNodeResponses(t, cdc, hashString, height, index, code, codespace)

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	hash, err := hex.DecodeString(hashString)
	assert.NoError(t, err)
	mockClient.EXPECT().Tx(hash, !cliCtx.TrustNode).Return(nodeResponses.resTx, nil)
	mockClient.EXPECT().Block(&nodeResponses.resTx.Height).Return(nodeResponses.resBlock, nil)

	res, err := QueryTx(cliCtx, hashString)
	assert.NoError(t, err)

	out, err := cdc.MarshalJSONIndent(res, "", "  ")
	assert.NoError(t, err)

	var m map[string]interface{}
	err = json.Unmarshal(out, &m)
	assert.NoError(t, err)

	assert.Contains(t, m, "index")
	assert.NotContains(t, m, "code")
	assert.NotContains(t, m, "codespace")
}

func TestQueryGenesisTxs(t *testing.T) {
	cdc := setupCodec()
	mockClient := mock.NewMockClient(gomock.NewController(t))
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
	mockClient.EXPECT().Genesis().Return(&genesisResult, nil)

	genesisTxs, err := QueryGenesisTx(cliCtx)
	assert.NoError(t, err)
	assert.NotEmpty(t, genesisTxs)

	// not exist genesis tx
	genesisDoc = tmtypes.GenesisDoc{
		AppState: json.RawMessage(`{"genutil":{"gentxs":[]}`),
	}
	genesisResult = ctypes.ResultGenesis{
		Genesis: &genesisDoc,
	}
	mockClient.EXPECT().Genesis().Return(&genesisResult, nil)
	genesisTxs, err = QueryGenesisTx(cliCtx)
	assert.Empty(t, genesisTxs)
	assert.Error(t, err)
}

func TestQueryGenesisAccount(t *testing.T) {
	testQueryGenesisAccount(t, "link")
	testQueryGenesisAccount(t, "tlink")
}

func testQueryGenesisAccount(t *testing.T, prefix string) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(prefix, prefix+sdk.PrefixPublic)

	cdc := setupCodec()
	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	// exist genesis account
	var addr string
	if config.GetBech32AccountAddrPrefix() == "tlink" {
		addr = "tlink19rqsvml8ldr0yrhaewgv9smcdvrew5panjryj3"
	} else {
		addr = "link19rqsvml8ldr0yrhaewgv9smcdvrew5pah9j5t5"
	}
	genesisDoc := tmtypes.GenesisDoc{
		AppState: json.RawMessage(`{"auth":{"accounts":[{"type":"cosmos-sdk/Account","value":{"address":"` + addr + `","coins":[]}}]}}`),
	}
	genesisResult := ctypes.ResultGenesis{
		Genesis: &genesisDoc,
	}
	mockClient.EXPECT().Genesis().Return(&genesisResult, nil).AnyTimes()
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
	reqE0 := mustNewGetRequest(t, "/")

	req0 := mustNewGetRequest(t, "/?foo=faa")

	req1 := mustNewGetRequest(t, "/?foo=faa&limit=5")
	req2 := mustNewGetRequest(t, "/?foo=faa&page=5")
	req3 := mustNewGetRequest(t, "/?foo=faa&page=5&limit=5")

	reqE1 := mustNewGetRequest(t, "/?foo=faa&page=-1")
	reqE2 := mustNewGetRequest(t, "/?foo=faa&limit=-1")
	reqE3 := mustNewGetRequest(t, "/?page=5&limit=5")

	req4 := mustNewGetRequest(t, "/?foo=faa&height.from=1")
	req5 := mustNewGetRequest(t, "/?foo=faa&height.to=1")

	reqE4 := mustNewGetRequest(t, "/?foo=faa&height.from=-1")
	reqE5 := mustNewGetRequest(t, "/?foo=faa&height.to=-1")

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

func mustNewGetRequest(t *testing.T, url string) *http.Request {
	return mustNewRequest(t, "GET", url, nil)
}
