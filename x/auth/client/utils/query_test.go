package utils

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

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
