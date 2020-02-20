package block

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/line/link/client/rpc/mock"
	"github.com/stretchr/testify/require"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestRequestHandlerFn(t *testing.T) {
	height := "1"
	uri := "/blocks/" + height + "/?"
	pathVariable := map[string]string{
		"height": height,
	}
	t.Logf("normal case, return statusCode %d", http.StatusOK)
	{
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, pathVariable)
		mockClient, mockCdc, mockTendermint, mockCliCtx, rs, fbh, rb := prepareForBlockRest(t)

		mockCliCtx.EXPECT().GetNode().Return(mockClient, nil).Times(1)
		mockClient.EXPECT().Status().Return(rs, nil).Times(1)
		mockCliCtx.EXPECT().GetNode().Return(mockClient, nil).Times(1)
		mockClient.EXPECT().Block(&fbh).Return(rb, nil).Times(1)
		mockCliCtx.EXPECT().TrustNode().Return(true).Times(1)
		mockCliCtx.EXPECT().Verify(fbh).Times(1)
		mockTendermint.EXPECT().ValidateBlockMeta(gomock.Any(), gomock.Any()).Times(1)
		mockTendermint.EXPECT().ValidateBlock(gomock.Any(), gomock.Any()).Times(1)
		mockCliCtx.EXPECT().Indent().Return(true).Times(1)
		expectedOutput := []byte(`string`)
		mockCdc.EXPECT().MarshalJSONIndent(rb).Return(expectedOutput, nil).Times(1)
		statusCode, output, hasErr := processBlockFetchReq(req, &Util{
			lcdc:    mockCdc,
			ltmtl:   mockTendermint,
			lcliCtx: mockCliCtx,
		})
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, expectedOutput, output)
		require.Equal(t, false, hasErr)
	}

	t.Logf("could not fetch latestBlock, return statusCode %d", http.StatusInternalServerError)
	{
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, pathVariable)
		mockClient, mockCodec, mockTendermint, mockCliContext, _, _, _ := prepareForBlockRest(t)

		mockCliContext.EXPECT().GetNode().Return(mockClient, nil).Times(1)
		errMsg := "fetching latestblock of status"
		mockClient.EXPECT().Status().Return(nil, fmt.Errorf(errMsg)).Times(1)
		statusCode, output, hasErr := processBlockFetchReq(req, &Util{
			lcdc:    mockCodec,
			ltmtl:   mockTendermint,
			lcliCtx: mockCliContext,
		})
		require.Equal(t, http.StatusInternalServerError, statusCode)
		require.Equal(t, true, hasErr)
		require.Equal(t, fmt.Sprintf("failed to parse chain height. because of %s", errMsg), output.(string))
	}
}

func prepareForBlockRest(t *testing.T) (*mock.MockClient, *mock.MockCodec, *mock.MockTendermint, *mock.MockCLIContext,
	*ctypes.ResultStatus, int64, *ctypes.ResultBlock) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockClient(ctrl)
	mockCodec := mock.NewMockCodec(ctrl)
	mockTendermint := mock.NewMockTendermint(ctrl)
	mockCliContext := mock.NewMockCLIContext(ctrl)
	fromBlockHeightInt64 := int64(1)
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
	return mockClient, mockCodec, mockTendermint, mockCliContext, rs, fromBlockHeightInt64, rb
}
