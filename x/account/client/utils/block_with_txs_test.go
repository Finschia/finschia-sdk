package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/golang/mock/gomock"
	"github.com/line/lbm-sdk/v2/x/account/client/types"
	"github.com/line/lbm-sdk/v2/x/account/client/utils/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func TestLatestBlockHeight(t *testing.T) {
	cdc := setupCodec()
	height := int64(100)

	resBlock := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{
				Height: height,
				Time:   time.Time{},
			},
		},
	}

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}
	mockClient.EXPECT().Block(nil).Return(resBlock, nil)

	ret, err := LatestBlockHeight(cliCtx)
	require.NoError(t, err)
	require.Equal(t, height, ret)
}

func TestBlockWithTxResponses(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	cdc := setupCodec()

	height := int64(100)

	resBlock := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{
				Height: height,
				Time:   time.Time{},
			},
		},
	}

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}
	mockClient.EXPECT().Block(&height).Return(resBlock, nil)

	resTxSearch := &ctypes.ResultTxSearch{
		Txs:        nil,
		TotalCount: 0,
	}

	const defaultLimit = 100
	mockClient.EXPECT().
		TxSearch(fmt.Sprintf("tx.height=%d", height), !cliCtx.TrustNode, 1, defaultLimit, "").
		Return(resTxSearch, nil)

	actual, err := BlockWithTxResponses(cliCtx, height+1, height, int64(1))

	results := make([]*types.ResultBlockWithTxResponses, 1)

	txResponses := make([]types.TxResponse, 0)
	results[0] = &types.ResultBlockWithTxResponses{
		ResultBlock: &types.ResultBlock{
			BlockSize: resBlock.Block.Size(),
			BlockID:   resBlock.BlockID,
			Block:     resBlock.Block,
		},
		TxResponses: txResponses,
	}
	expected := &types.HasMoreResponseWrapper{
		Items:   results,
		HasMore: true,
	}
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestLatestBlockHeightWithError(t *testing.T) {
	cdc := setupCodec()

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}
	resError := errors.New("error")
	mockClient.EXPECT().Block(nil).Return(nil, resError)

	ret, err := LatestBlockHeight(cliCtx)
	require.Error(t, err)
	require.Equal(t, int64(-1), ret)
}

func makeTxs(t *testing.T, cdc *codec.Codec, count int, height int64) []*ctypes.ResultTx {
	txs := make([]*ctypes.ResultTx, count)

	stdTx := &auth.StdTx{
		Memo: "empty tx",
	}

	emptyTx, err := cdc.MarshalBinaryLengthPrefixed(stdTx)
	require.NoError(t, err)

	for i := 0; i < count; i++ {
		bs := make([]byte, 32)
		binary.LittleEndian.PutUint32(bs, uint32(i))
		txs[i] = &ctypes.ResultTx{
			Hash:   bs,
			Height: height,
			Index:  uint32(i),
			TxResult: abci.ResponseDeliverTx{
				Code:      0,
				Codespace: "codespace",
				Log:       "[]",
			},
			Tx: emptyTx,
		}
	}
	return txs
}

func TestGetTxs(t *testing.T) {
	var parameterizedTests = []struct {
		name     string
		txsCount int
	}{
		{"multiple pages(txsCount=123)", 123},
		{"single page(txsCount=12)", 12},
		{"multiple pages fill(txsCount=300)", 300},
		{"single page fill(txsCount=100)", 100},
		{"empty(txsCount=0)", 0},
	}

	for _, pt := range parameterizedTests {
		txsCount := pt.txsCount
		t.Run(pt.name, func(t *testing.T) {
			testGetTxsWithCount(t, txsCount)
		})
	}
}

func testGetTxsWithCount(t *testing.T, txsCount int) {
	const defaultLimit = 100

	cdc := setupCodec()
	height := int64(100)
	txs := makeTxs(t, cdc, txsCount, height)

	resBlock := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{
				Height: height,
				Time:   time.Time{},
			},
		},
	}

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	mockClient.EXPECT().Block(gomock.Any()).Return(resBlock, nil).AnyTimes()

	if txsCount == 0 {
		resTxSearch := &ctypes.ResultTxSearch{
			Txs:        nil,
			TotalCount: txsCount,
		}
		mockClient.EXPECT().
			TxSearch(fmt.Sprintf("tx.height=%d", height), !cliCtx.TrustNode, 1, defaultLimit, "").
			Return(resTxSearch, nil)
	} else {
		for i := 0; i < int(math.Ceil(float64(txsCount)/defaultLimit)); i++ {
			start := i * defaultLimit
			end := start + defaultLimit
			if end > txsCount {
				// the last segment
				end = txsCount
			}
			pageTxs := txs[start:end]
			resTxSearch := &ctypes.ResultTxSearch{
				Txs:        pageTxs,
				TotalCount: txsCount,
			}
			mockClient.EXPECT().
				TxSearch(fmt.Sprintf("tx.height=%d", height), !cliCtx.TrustNode, i+1, defaultLimit, "").
				Return(resTxSearch, nil)
		}
	}

	ret, err := getTxs(cliCtx, height)
	require.NoError(t, err)
	require.Len(t, ret, txsCount)
}

func TestEmptyTxs(t *testing.T) {
	const defaultLimit = 100

	cdc := setupCodec()
	height := int64(100)

	resBlock := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{
				Height: height,
				Time:   time.Time{},
			},
		},
	}

	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}

	mockClient.EXPECT().Block(gomock.Any()).Return(resBlock, nil).AnyTimes()

	resTxSearch := &ctypes.ResultTxSearch{
		Txs:        nil,
		TotalCount: 0,
	}
	mockClient.EXPECT().
		TxSearch(fmt.Sprintf("tx.height=%d", height), !cliCtx.TrustNode, 1, defaultLimit, "").
		Return(resTxSearch, nil)

	ret, err := getTxs(cliCtx, height)
	require.NoError(t, err)
	require.NotNil(t, ret)
}
