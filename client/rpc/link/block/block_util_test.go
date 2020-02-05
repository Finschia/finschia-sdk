package block

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/golang/mock/gomock"
	. "github.com/line/link/client/rpc/mock"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	checkMark = "\u2713"
	txType    = "cosmos-sdk/StdTx"
	txMemo    = "memo"
	tx        = `{"type":"` + txType + `","value": {"memo":"` + txMemo + `"} }`
)

var (
	fromHeight = int64(1)
)

func TestValidateBlock(t *testing.T) {

	t.Log("TrustNode is true", checkMark)
	{
		_, mockTendermint, mockCliCtx, rb, bu, _, _, _ := prepare(t)

		mockCliCtx.EXPECT().TrustNode().Return(true).Times(1)
		mockTendermint.EXPECT().ValidateBlock(gomock.Any(), gomock.Any()).Times(0)
		mockTendermint.EXPECT().ValidateBlockMeta(gomock.Any(), gomock.Any()).Times(0)
		err := bu.ValidateBlock(rb)
		require.Equal(t, nil, err)
	}

	t.Log("TrustNode is false", checkMark)
	{
		check, mockTendermint, mockCliCtx, rb, bu, _, _, _ := prepare(t)

		mockCliCtx.EXPECT().TrustNode().Return(false).Times(1)
		mockCliCtx.EXPECT().Verify(rb.Block.Height).Return(check, nil).Times(1)
		mockTendermint.EXPECT().ValidateBlockMeta(rb.BlockMeta, check).Times(1)
		mockTendermint.EXPECT().ValidateBlock(rb.Block, check).Times(1)
		err := bu.ValidateBlock(rb)
		require.Equal(t, nil, err)
	}
}

func TestValidateBlockFail(t *testing.T) {

	t.Log("TrustNode is false and Verify return error")
	{
		check, mockTendermint, mockCliCtx, rb, bu, _, _, _ := prepare(t)

		mockCliCtx.EXPECT().TrustNode().Return(false).Times(1)
		verifyErr := fmt.Errorf("verify failed")
		mockCliCtx.EXPECT().Verify(rb.Block.Height).Return(check, verifyErr).Times(1)
		mockTendermint.EXPECT().ValidateBlock(gomock.Any(), gomock.Any()).Times(0)
		mockTendermint.EXPECT().ValidateBlockMeta(gomock.Any(), gomock.Any()).Times(0)
		err := bu.ValidateBlock(rb)
		require.Equal(t, verifyErr, err)
	}
	t.Log("TrustNode is false and ValidateBlockMeta return error")
	{
		check, mockTendermint, mockCliCtx, rb, bu, _, _, _ := prepare(t)

		mockCliCtx.EXPECT().TrustNode().Return(false).Times(1)
		validateMetaErr := fmt.Errorf("ValidateBlockMeta failed")
		mockCliCtx.EXPECT().Verify(rb.Block.Height).Return(check, nil).Times(1)
		mockTendermint.EXPECT().ValidateBlockMeta(rb.BlockMeta, check).Return(validateMetaErr).Times(1)
		mockTendermint.EXPECT().ValidateBlock(gomock.Any(), gomock.Any()).Times(0)
		err := bu.ValidateBlock(rb)
		require.Equal(t, validateMetaErr, err)
	}
	t.Log("TrustNode is false and ValidateBlock return error")
	{
		check, mockTendermint, mockCliCtx, rb, bu, _, _, _ := prepare(t)

		mockCliCtx.EXPECT().TrustNode().Return(false).Times(1)
		mockCliCtx.EXPECT().Verify(rb.Block.Height).Return(check, nil).Times(1)
		mockTendermint.EXPECT().ValidateBlockMeta(rb.BlockMeta, check).Return(nil).Times(1)
		validateBlockErr := fmt.Errorf("ValidateBlock failed")
		mockTendermint.EXPECT().ValidateBlock(rb.Block, check).Return(validateBlockErr).Times(1)
		err := bu.ValidateBlock(rb)
		require.Equal(t, validateBlockErr, err)
	}
}

func TestIndentJSONRB(t *testing.T) {
	_, _, mockCliCtx, _, bu, _, mockCodecUtil, _ := prepare(t)

	expectedJSON := []byte("good")
	var expectedErr error = nil
	resultBlock := &ctypes.ResultBlock{}

	t.Log("Indent is false", checkMark)
	{

		mockCliCtx.EXPECT().Indent().Return(false).Times(1)
		mockCodecUtil.EXPECT().MarshalJSONIndent(gomock.Any()).Times(0)
		mockCodecUtil.EXPECT().MarshalJSON(resultBlock).Return(expectedJSON, expectedErr).Times(1)
		notIndentedJSON, err := bu.Indent(resultBlock)
		require.Equal(t, expectedJSON, notIndentedJSON)
		require.Equal(t, expectedErr, err)
	}

	t.Log("Indent is true", checkMark)
	{
		mockCliCtx.EXPECT().Indent().Return(true).Times(1)
		mockCodecUtil.EXPECT().MarshalJSONIndent(gomock.Any()).Return(expectedJSON, expectedErr).Times(1)
		mockCodecUtil.EXPECT().MarshalJSONIndent(resultBlock).Times(0)
		indentedJSON, err := bu.Indent(resultBlock)
		require.Equal(t, expectedJSON, indentedJSON)
		require.Equal(t, expectedErr, err)
	}
}

func TestInjectByteToJsonTxs(t *testing.T) {
	_, _, _, _, bu, _, _, _ := prepare(t)
	bs := []byte(`{
	"block": {
		"data": {
			"txs": [{"type":"cosmos-sdk/StdTx","value":{"fee":{"amount":[]}}}]
		}
	}}`)
	var byteTxa [][]byte = nil
	byteTxa = append(byteTxa, []byte(tx))

	block, err := bu.InjectByteToJsonTxs(bs, byteTxa)
	require.NoError(t, err)
	actual, err := json.Marshal(block["block"].(map[string]interface{})["data"].(map[string]interface{})["txs"])
	require.NoError(t, err)
	var result []map[string]interface{}
	err = json.Unmarshal(actual, &result)
	require.NoError(t, err)
	require.Equal(t, txType, result[0]["type"])
	require.Equal(t, txMemo, result[0]["value"].(map[string]interface{})["memo"])
	require.Equal(t, nil, err)
}

func TestCalcFetchBlockHeight(t *testing.T) {

	t.Log("ChainBlockHeight greaterThanEqual request ", checkMark)
	{
		latestBlockHeight := int64(21)
		actual := NewFetchInfo(&latestBlockHeight, &fromHeight, &DefaultBlockFetchSize)
		require.Equal(t, fromHeight, actual.inclusiveFromHeight)
		require.Equal(t, int64(21), actual.exclusiveToHeight)
		require.Equal(t, true, actual.hasMore)
		require.Equal(t, DefaultBlockFetchSize, actual.fetchItemCnt)
		require.Equal(t, int(DefaultBlockFetchSize), len(actual.fetchItemRange))
	}

	t.Log("ChainBlockHeight Equal request ", checkMark)
	{
		latestBlockHeight := int64(20)
		actual := NewFetchInfo(&latestBlockHeight, &fromHeight, &DefaultBlockFetchSize)
		require.Equal(t, fromHeight, actual.inclusiveFromHeight)
		require.Equal(t, int64(21), actual.exclusiveToHeight)
		require.Equal(t, false, actual.hasMore)
		require.Equal(t, DefaultBlockFetchSize, actual.fetchItemCnt)
		require.Equal(t, int(DefaultBlockFetchSize), len(actual.fetchItemRange))
	}

	t.Log("ChainBlockHeight LessThan request ", checkMark)
	{
		latestBlockHeight := int64(19)
		actual := NewFetchInfo(&latestBlockHeight, &fromHeight, &DefaultBlockFetchSize)
		require.Equal(t, fromHeight, actual.inclusiveFromHeight)
		require.Equal(t, int64(20), actual.exclusiveToHeight)
		require.Equal(t, false, actual.hasMore)
		require.Equal(t, DefaultBlockFetchSize-1, actual.fetchItemCnt)
		require.Equal(t, int(DefaultBlockFetchSize-1), len(actual.fetchItemRange))
		require.Equal(t, int64(1), actual.fetchItemRange[0])
		require.Equal(t, int64(19), actual.fetchItemRange[18])
	}
}

func TestFetchBlock(t *testing.T) {
	t.Log("normal case", checkMark)
	{
		check, mockTendermint, mockCliCtx, rb, bu, mockClient, _, _ := prepare(t)
		latestBlockHeight := int64(19)

		mockCliCtx.EXPECT().GetNode().Return(mockClient, nil).Times(1)
		mockClient.EXPECT().Block(&latestBlockHeight).Return(rb, nil).Times(1)
		mockCliCtx.EXPECT().TrustNode().Return(false).Times(1)
		mockCliCtx.EXPECT().Verify(rb.Block.Height).Return(check, nil).Times(1)
		mockTendermint.EXPECT().ValidateBlockMeta(rb.BlockMeta, check).Return(nil).Times(1)
		mockTendermint.EXPECT().ValidateBlock(rb.Block, check).Return(nil).Times(1)
		actual, err := bu.fetchBlock(latestBlockHeight)
		require.Equal(t, rb, actual.ResultBlock)
		require.Equal(t, nil, err)
	}
}

func TestFormatTxResult(t *testing.T) {
	tn := time.Now()
	t.Log("normal case")
	{
		_, _, _, rb, bu, _, mockCodecUtil, resTx := prepare(t)
		mockCodecUtil.EXPECT().UnmarshalBinaryLengthPrefixed(gomock.Any()).Return(auth.StdTx{}, nil).Times(1)

		txRes, err := bu.formatTxResult(resTx, tn.Format(time.RFC3339))
		require.Equal(t, nil, err)
		require.Equal(t, txRes.Height, rb.Block.Height)
	}

	t.Log("UnmarshalBinaryLengthPrefixed error case")
	{
		_, _, _, _, bu, _, mockCodecUtil, resTx := prepare(t)
		unmarshalBinaryLengthPrefixedErr := fmt.Errorf("unmarshalBinaryLengthPrefixedErr")
		var stdTx auth.StdTx
		mockCodecUtil.EXPECT().UnmarshalBinaryLengthPrefixed(gomock.Any()).Return(stdTx, unmarshalBinaryLengthPrefixedErr).Times(1)

		txRes, err := bu.formatTxResult(resTx, tn.Format(time.RFC3339))
		require.Equal(t, unmarshalBinaryLengthPrefixedErr, err)
		require.Nil(t, txRes)
	}
}
func TestLatestBlockHeight(t *testing.T) {
	t.Log("normal case", checkMark)
	{
		_, _, mockCliCtx, _, bu, mockClient, _, _ := prepare(t)
		mockCliCtx.EXPECT().GetNode().Return(mockClient, nil).Times(1)
		fromBlockHeightInt64 := int64(1)
		rs := &ctypes.ResultStatus{
			SyncInfo: ctypes.SyncInfo{
				LatestBlockHeight: fromBlockHeightInt64,
			},
		}
		mockClient.EXPECT().Status().Return(rs, nil).Times(1)

		actual, err := bu.LatestBlockHeight()
		require.Equal(t, fromBlockHeightInt64, actual)
		require.Equal(t, nil, err)
	}
}
func prepare(t *testing.T) (tmtypes.SignedHeader, *MockTendermint, *MockCLIContext, *ctypes.ResultBlock, *Util, *MockClient, *MockCodec, *ctypes.ResultTx) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	check := tmtypes.SignedHeader{}
	mockTendermint := NewMockTendermint(ctrl)
	mockCliCtx := NewMockCLIContext(ctrl)
	mockClient := NewMockClient(ctrl)
	blockHeight := int64(123)
	rb := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: blockHeight},
		},
		BlockMeta: &tmtypes.BlockMeta{
			BlockID: tmtypes.BlockID{},
			Header:  tmtypes.Header{},
		},
	}
	resTx := &ctypes.ResultTx{
		Hash:     []byte(`txhash`),
		Height:   blockHeight,
		Index:    0,
		TxResult: types.ResponseDeliverTx{},
		Tx:       nil,
		Proof:    tmtypes.TxProof{},
	}

	mockCodecUtil := NewMockCodec(ctrl)
	bu := &Util{lcdc: mockCodecUtil, ltmtl: mockTendermint, lcliCtx: mockCliCtx}
	return check, mockTendermint, mockCliCtx, rb, bu, mockClient, mockCodecUtil, resTx
}
