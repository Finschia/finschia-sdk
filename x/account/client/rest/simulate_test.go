package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/line/lbm-sdk/client/context"
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/rest"
	"github.com/line/lbm-sdk/x/auth"
	"github.com/line/lbm-sdk/x/auth/types"

	"github.com/line/lbm-sdk/x/account/client/utils/mock"
)

type mockHTTPWriter struct {
	statusCode int
	bodyBuf    bytes.Buffer
}

func (m *mockHTTPWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockHTTPWriter) Write(body []byte) (int, error) {
	return m.bodyBuf.Write(body)
}

func (m *mockHTTPWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

var _ http.ResponseWriter = &mockHTTPWriter{}

func setupCodec() *codec.Codec {
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	return cdc
}

func TestSimulateTxRequest(t *testing.T) {
	cdc := setupCodec()

	// assumes node response is
	var gasUsed uint64 = 10000
	adjustment := 1.2

	// set up mock node response
	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}
	abciRes := &ctypes.ResultABCIQuery{
		Response: abci.ResponseQuery{
			Value: codec.Cdc.MustMarshalBinaryBare(sdk.SimulationResponse{
				GasInfo: sdk.GasInfo{
					GasUsed: gasUsed,
				},
			}),
		},
	}
	mockClient.EXPECT().ABCIQueryWithOptions("/app/simulate", gomock.Any(), gomock.Any()).Return(abciRes, nil)

	// request
	req := codec.MustMarshalJSONIndent(cdc,
		&SimulateReq{
			Tx: types.StdTx{
				Memo: "empty tx",
			},
			GasAdjustment: fmt.Sprintf("%f", adjustment),
		},
	)
	request := http.Request{Body: ioutil.NopCloser(bytes.NewReader(req))}

	writer := mockHTTPWriter{}
	SimulateTxRequest(cliCtx)(&writer, &request)

	res := rest.GasEstimateResponse{}
	cdc.MustUnmarshalJSON(writer.bodyBuf.Bytes(), &res)

	require.Equal(t, uint64(float64(gasUsed)*adjustment), res.GasEstimate)
}

func TestSimulateTxRequestWithABCIError(t *testing.T) {
	cdc := setupCodec()

	var code uint32 = 10
	codespace := "codespace"
	message := "error message"
	adjustment := 1.2

	// set up mock node response
	mockClient := mock.NewMockClient(gomock.NewController(t))
	cliCtx := context.CLIContext{
		Client:    mockClient,
		TrustNode: true,
		Codec:     cdc,
	}
	abciRes := &ctypes.ResultABCIQuery{
		Response: abci.ResponseQuery{
			Code:      code,
			Codespace: codespace,
			Log:       message,
		},
	}
	mockClient.EXPECT().ABCIQueryWithOptions("/app/simulate", gomock.Any(), gomock.Any()).Return(abciRes, nil)

	// request
	req := codec.MustMarshalJSONIndent(cdc,
		&SimulateReq{
			Tx: types.StdTx{
				Memo: "empty tx",
			},
			GasAdjustment: fmt.Sprintf("%f", adjustment),
		},
	)
	request := http.Request{Body: ioutil.NopCloser(bytes.NewReader(req))}

	writer := mockHTTPWriter{}
	SimulateTxRequest(cliCtx)(&writer, &request)

	res := ABCIErrorResponse{}
	cdc.MustUnmarshalJSON(writer.bodyBuf.Bytes(), &res)

	require.Equal(t, code, res.Code)
	require.Equal(t, codespace, res.Codespace)
	require.Equal(t, message, res.Error)
}
