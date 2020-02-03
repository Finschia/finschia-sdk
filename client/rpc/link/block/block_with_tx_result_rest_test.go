package block

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	cdc "github.com/line/link/client/rpc/link/block/codec"
	"github.com/line/link/client/rpc/mock"
	"github.com/stretchr/testify/require"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestWithTxResultRequestHandlerFn(t *testing.T) {
	fromBlockHeight := "1"
	uri := "/blocks_with_tx_results/" + fromBlockHeight + "/?"
	pathVariable := map[string]string{
		"from_height": fromBlockHeight,
	}
	queryParam := make(url.Values)
	queryParam["fetchsize"] = []string{strconv.Itoa(int(DefaultBlockFetchSize))}

	t.Logf("Fetch LatestBlock with statusCode %d", http.StatusOK)
	{
		req, err := http.NewRequest("GET", uri+queryParam.Encode(), nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, pathVariable)
		mockClient, mockCodec, mockTendermint, mockCliContext, rs, rb, fromBlockHeightInt64 := prepareForRest(t, fromBlockHeight)
		mockCliContext.EXPECT().GetNode().Return(mockClient, nil).Times(1)

		mockCliContext.EXPECT().GetNode().Return(mockClient, nil).Times(1)
		mockClient.EXPECT().Status().Return(rs, nil).Times(1)
		mockClient.EXPECT().Block(nil).Return(rb, nil).Times(1)
		mockClient.EXPECT().Block(&rb.Block.Height).Return(rb, nil).Times(1)
		mockCliContext.EXPECT().TrustNode().Return(true).Times(1)
		statusCode, output, hasErr := processReq(req, &Util{
			lcdc:    mockCodec,
			ltmtl:   mockTendermint,
			lcliCtx: mockCliContext,
		})
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, false, hasErr)
		require.Equal(t, false, output.(*cdc.HasMoreResponseWrapper).HasMore)
		require.Equal(t, fromBlockHeightInt64, output.(*cdc.HasMoreResponseWrapper).Items[0].ResultBlock.Block.Height)
	}

	t.Logf("could not fetch latestBlock, return statusCode %d", http.StatusInternalServerError)
	{
		req, err := http.NewRequest("GET", uri+queryParam.Encode(), nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, pathVariable)
		mockClient, mockCodec, mockTendermint, mockCliContext, _, _, _ := prepareForRest(t, fromBlockHeight)

		mockCliContext.EXPECT().GetNode().Return(mockClient, nil).Times(1)
		errMsg := "errorMsg"
		mockClient.EXPECT().Status().Return(nil, fmt.Errorf(errMsg)).Times(1)
		statusCode, output, hasErr := processReq(req, &Util{
			lcdc:    mockCodec,
			ltmtl:   mockTendermint,
			lcliCtx: mockCliContext,
		})
		require.Equal(t, http.StatusInternalServerError, statusCode)
		require.Equal(t, true, hasErr)
		require.Equal(t, fmt.Sprintf("couldn't get latestBlockHeight. because of %s", errMsg), output.(string))
	}
}

func prepareForRest(t *testing.T, fromBlockHeight string) (*mock.MockClient, *mock.MockCodec, *mock.MockTendermint,
	*mock.MockCLIContext, *ctypes.ResultStatus, *ctypes.ResultBlock, int64) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockClient(ctrl)
	mockCodec := mock.NewMockCodec(ctrl)
	mockTendermint := mock.NewMockTendermint(ctrl)
	mockCliContext := mock.NewMockCLIContext(ctrl)

	fromBlockHeightInt64, _ := strconv.ParseInt(fromBlockHeight, 10, 64)
	rs := &ctypes.ResultStatus{
		SyncInfo: ctypes.SyncInfo{
			LatestBlockHeight: fromBlockHeightInt64,
		},
	}
	rb := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: fromBlockHeightInt64},
		},
		BlockMeta: &tmtypes.BlockMeta{
			BlockID: tmtypes.BlockID{},
			Header:  tmtypes.Header{},
		},
	}
	return mockClient, mockCodec, mockTendermint, mockCliContext, rs, rb, fromBlockHeightInt64
}
