// +build !integration

package loadgenerator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/stretchr/testify/require"
)

func TestLoadHandler(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()

	var testCases = []struct {
		testName    string
		tartgetType string
	}{
		{
			"Query Account",
			types.QueryAccount,
		},
		{
			"Tx Send",
			types.TxSend,
		},
	}

	for i, tc := range testCases {
		t.Logf("Test %s", tc.testName)
		{
			// Given LoadGenerator
			sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
			lg := NewLoadGenerator()
			// Given Router
			r := mux.NewRouter()
			RegisterHandlers(lg, r)
			// Given Request
			config := types.Config{
				MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
				TPS:               tests.TestTPS,
				Duration:          tests.TestDuration,
				MaxWorkers:        tests.TestMaxWorkers,
				PacerType:         types.ConstantPacer,
				TargetURL:         server.URL,
				ChainID:           tests.TestChainID,
				CoinName:          tests.TestCoinName,
				Mnemonic:          tests.TestMnemonic,
			}
			res := httptest.NewRecorder()
			request := types.NewLoadRequest(tc.tartgetType, config)
			body, err := json.Marshal(request)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", "/target/load", bytes.NewBuffer(body))
			require.NoError(t, err)

			// When
			r.ServeHTTP(res, req)

			// Then
			require.Equal(t, http.StatusOK, res.Code)
			require.Equal(t, tests.TestTPS, lg.config.TPS)
			require.Equal(t, tests.TestDuration, lg.config.Duration)
			require.Equal(t, server.URL, lg.targetBuilder.LCDURL)
			require.Equal(t, tests.ExpectedNumTargetsConstant*i, mock.GetCallCounter(server.URL).QueryAccountCallCount)
		}
	}
}

func TestLoadHandlerWithInvalidParameters(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()

	var testCases = []struct {
		testName   string
		targetType string
		config     types.Config
		error      string
	}{
		{
			"with empty chain id",
			types.QueryAccount,
			types.Config{
				MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
				TPS:               tests.TestTPS,
				Duration:          tests.TestDuration,
				MaxWorkers:        tests.TestMaxWorkers,
				PacerType:         types.ConstantPacer,
				TargetURL:         server.URL,
				ChainID:           "",
				CoinName:          tests.TestCoinName,
				Mnemonic:          tests.TestMnemonic,
			},
			"Invalid Load Parameter Error: invalid parameter of load handler\n",
		},
		{
			"with invalid tps",
			types.QueryAccount,
			types.Config{
				MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
				TPS:               -1,
				Duration:          tests.TestDuration,
				MaxWorkers:        tests.TestMaxWorkers,
				PacerType:         types.ConstantPacer,
				TargetURL:         server.URL,
				ChainID:           tests.TestChainID,
				CoinName:          tests.TestCoinName,
				Mnemonic:          tests.TestMnemonic,
			},
			"Invalid Load Parameter Error: invalid parameter of load handler\n",
		},
		{
			"with invalid duration",
			types.QueryAccount,
			types.Config{
				MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
				TPS:               tests.TestTPS,
				Duration:          -1,
				MaxWorkers:        tests.TestMaxWorkers,
				PacerType:         types.ConstantPacer,
				TargetURL:         server.URL,
				ChainID:           tests.TestChainID,
				CoinName:          tests.TestCoinName,
				Mnemonic:          tests.TestMnemonic,
			},
			"Invalid Load Parameter Error: invalid parameter of load handler\n",
		},
		{
			"with invalid mnemonic",
			types.QueryAccount,
			types.Config{
				MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
				TPS:               tests.TestTPS,
				Duration:          tests.TestDuration,
				MaxWorkers:        tests.TestMaxWorkers,
				PacerType:         types.ConstantPacer,
				TargetURL:         server.URL,
				ChainID:           tests.TestChainID,
				CoinName:          tests.TestCoinName,
				Mnemonic:          "invalid mnemonic",
			},
			"Invalid mnemonic\n",
		},
		{
			"with invalid pacer",
			types.QueryAccount,
			types.Config{
				MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
				TPS:               tests.TestTPS,
				Duration:          tests.TestDuration,
				MaxWorkers:        tests.TestMaxWorkers,
				PacerType:         "invalid",
				TargetURL:         server.URL,
				ChainID:           tests.TestChainID,
				CoinName:          tests.TestCoinName,
				Mnemonic:          tests.TestMnemonic,
			},
			"Invalid pacer type: invalid\n",
		},
		{
			"with invalid target type",
			"invalid type",
			types.Config{
				MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
				TPS:               tests.TestTPS,
				Duration:          tests.TestDuration,
				MaxWorkers:        tests.TestMaxWorkers,
				PacerType:         types.ConstantPacer,
				TargetURL:         server.URL,
				ChainID:           tests.TestChainID,
				CoinName:          tests.TestCoinName,
				Mnemonic:          tests.TestMnemonic,
			},
			"Invalid target Type Error: invalid target type\n",
		},
	}

	for _, tc := range testCases {
		t.Logf("Test %s", tc.testName)
		{
			// Given LoadGenerator
			sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
			lg := NewLoadGenerator()
			// Given Router
			r := mux.NewRouter()
			RegisterHandlers(lg, r)
			// Given Request
			res := httptest.NewRecorder()
			request := types.NewLoadRequest(tc.targetType, tc.config)
			body, err := json.Marshal(request)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", "/target/load", bytes.NewBuffer(body))
			require.NoError(t, err)

			// When
			r.ServeHTTP(res, req)

			body, err = ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			// Then
			require.Equal(t, http.StatusBadRequest, res.Code)
			require.Equal(t, tc.error, string(body))
		}
	}
}

func TestFireHandler(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given LoadGenerator
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	lg := NewLoadGenerator()
	// Given Router
	r := mux.NewRouter()
	RegisterHandlers(lg, r)
	// Given Config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		MaxWorkers:        tests.TestMaxWorkers,
		PacerType:         types.ConstantPacer,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		Mnemonic:          tests.TestMnemonic,
	}

	t.Log("Test with constant pacer")
	{
		// Load Targets
		res := httptest.NewRecorder()
		request := types.NewLoadRequest(types.TxSend, config)
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/target/load", bytes.NewBuffer(body))
		require.NoError(t, err)
		r.ServeHTTP(res, req)
		// Given Fire Request
		res = httptest.NewRecorder()
		req, err = http.NewRequest("POST", "/target/fire", nil)
		require.NoError(t, err)

		// When
		r.ServeHTTP(res, req)

		// Then
		require.Equal(t, http.StatusOK, res.Code)
		require.Equal(t, tests.ExpectedNumTargetsConstant, mock.GetCallCounter(server.URL).BroadcastTxCallCount)
	}
	t.Log("Test with linear pacer")
	{
		// Given
		config.PacerType = types.LinearPacer
		// Load Targets
		res := httptest.NewRecorder()
		request := types.NewLoadRequest(types.TxSend, config)
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/target/load", bytes.NewBuffer(body))
		require.NoError(t, err)
		r.ServeHTTP(res, req)
		// Given Fire Request
		res = httptest.NewRecorder()
		req, err = http.NewRequest("POST", "/target/fire", nil)
		require.NoError(t, err)
		// Clear Call Counter
		mock.ClearCallCounter(server.URL)

		// When
		r.ServeHTTP(res, req)

		// Then
		require.Equal(t, http.StatusOK, res.Code)
		require.Equal(t, tests.ExpectedNumTargetsLinear, mock.GetCallCounter(server.URL).BroadcastTxCallCount)
	}
}
