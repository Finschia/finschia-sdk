package types_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	ocabci "github.com/line/ostracon/abci/types"
	"github.com/line/ostracon/libs/bytes"
	ctypes "github.com/line/ostracon/rpc/core/types"

	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
)

type resultTestSuite struct {
	suite.Suite
}

func TestResultTestSuite(t *testing.T) {
	suite.Run(t, new(resultTestSuite))
}

func (s *resultTestSuite) SetupSuite() {
	s.T().Parallel()
}

func (s *resultTestSuite) TestParseABCILog() {
	logs := `[{"log":"","msg_index":1,"success":true}]`
	res, err := sdk.ParseABCILogs(logs)

	s.Require().NoError(err)
	s.Require().Len(res, 1)
	s.Require().Equal(res[0].Log, "")
	s.Require().Equal(res[0].MsgIndex, uint32(1))
}

func (s *resultTestSuite) TestABCIMessageLog() {
	cdc := codec.NewLegacyAmino()

	const maxIter = 5

	tests := []struct {
		emptyLog   bool
		emptyType  bool
		emptyKey   bool
		emptyValue bool
	}{
		{false, false, false, false},
		{true, false, false, false},
		{false, true, false, false},
		{false, false, true, false},
		{false, false, false, true},
		{false, false, true, true},
		{false, true, false, true},
		{false, true, true, false},
		{true, false, false, true},
		{true, false, true, false},
		{true, true, false, false},
		{false, true, true, true},
		{true, false, true, true},
		{true, true, false, true},
		{true, true, true, false},
		{true, true, true, true},
	}

	for _, tt := range tests {
		msgLogs := sdk.ABCIMessageLogs{}
		for numMsgs := 0; numMsgs < maxIter; numMsgs++ {
			for i := 0; i < numMsgs; i++ {
				events := sdk.Events{}
				for numEvents := 0; numEvents < maxIter; numEvents++ {
					for j := 0; j < numEvents; j++ {
						var attributes []sdk.Attribute
						for numAttributes := 0; numAttributes < maxIter; numAttributes++ {
							for i := 0; i < numAttributes; i++ {
								key := ""
								value := ""
								if !tt.emptyKey {
									key = fmt.Sprintf("key%d", i)
								}
								if !tt.emptyValue {
									value = fmt.Sprintf("value%d", i)
								}
								attributes = append(attributes, sdk.NewAttribute(key, value))
							}
						}
						typeStr := ""
						if !tt.emptyType {
							typeStr = fmt.Sprintf("type%d", i)
						}
						events = append(events, sdk.NewEvent(typeStr, attributes...))
					}
				}

				log := ""
				if !tt.emptyLog {
					log = fmt.Sprintf("log%d", i)
				}
				msgLogs = append(msgLogs, sdk.NewABCIMessageLog(uint32(i), log, events))
			}
		}
		bz, err := cdc.MarshalJSON(msgLogs)

		s.Require().NoError(err)
		s.Require().Equal(string(bz), msgLogs.String())
	}

	var msgLogs sdk.ABCIMessageLogs
	s.Require().Equal("", msgLogs.String())
}

func (s *resultTestSuite) TestNewSearchTxsResult() {
	got := sdk.NewSearchTxsResult(150, 20, 2, 20, []*sdk.TxResponse{})
	s.Require().Equal(&sdk.SearchTxsResult{
		TotalCount: 150,
		Count:      20,
		PageNumber: 2,
		PageTotal:  8,
		Limit:      20,
		Txs:        []*sdk.TxResponse{},
	}, got)
}

func (s *resultTestSuite) TestResponseResultTx() {
	deliverTxResult := abci.ResponseDeliverTx{
		Codespace: "codespace",
		Code:      1,
		Data:      []byte("data"),
		Log:       `[]`,
		Info:      "info",
		GasWanted: 100,
		GasUsed:   90,
	}
	resultTx := &ctypes.ResultTx{
		Hash:     bytes.HexBytes([]byte("test")),
		Height:   10,
		Index:    1,
		TxResult: deliverTxResult,
	}
	logs, err := sdk.ParseABCILogs(`[]`)

	s.Require().NoError(err)

	want := &sdk.TxResponse{
		TxHash:    "74657374",
		Height:    10,
		Codespace: "codespace",
		Code:      1,
		Data:      strings.ToUpper(hex.EncodeToString([]byte("data"))),
		RawLog:    `[]`,
		Logs:      logs,
		Info:      "info",
		GasWanted: 100,
		GasUsed:   90,
		Tx:        nil,
		Timestamp: "timestamp",
		Index:     1,
	}

	s.Require().Equal(want, sdk.NewResponseResultTx(resultTx, nil, "timestamp"))
	s.Require().Equal((*sdk.TxResponse)(nil), sdk.NewResponseResultTx(nil, nil, "timestamp"))
	s.Require().Equal(`code: 1
codespace: codespace
data: "64617461"
events: []
gas_used: "90"
gas_wanted: "100"
height: "10"
index: 1
info: info
logs: []
raw_log: '[]'
timestamp: timestamp
tx: null
txhash: "74657374"
`, sdk.NewResponseResultTx(resultTx, nil, "timestamp").String())
	s.Require().True(sdk.TxResponse{}.Empty())
	s.Require().False(want.Empty())

	resultBroadcastTx := &ctypes.ResultBroadcastTx{
		Code:      1,
		Codespace: "codespace",
		Data:      []byte("data"),
		Log:       `[]`,
		Hash:      bytes.HexBytes([]byte("test")),
	}

	s.Require().Equal(&sdk.TxResponse{
		Code:      1,
		Codespace: "codespace",
		Data:      "64617461",
		RawLog:    `[]`,
		Logs:      logs,
		TxHash:    "74657374",
	}, sdk.NewResponseFormatBroadcastTx(resultBroadcastTx))
	s.Require().Equal((*sdk.TxResponse)(nil), sdk.NewResponseFormatBroadcastTx(nil))
}

func (s *resultTestSuite) TestResponseFormatBroadcastTxCommit() {
	// test nil
	s.Require().Equal((*sdk.TxResponse)(nil), sdk.NewResponseFormatBroadcastTxCommit(nil))

	logs, err := sdk.ParseABCILogs(`[]`)
	s.Require().NoError(err)

	// test checkTx
	checkTxResult := &ctypes.ResultBroadcastTxCommit{
		Height: 10,
		Hash:   bytes.HexBytes([]byte("test")),
		CheckTx: ocabci.ResponseCheckTx{
			Code:      90,
			Data:      nil,
			Log:       `[]`,
			Info:      "info",
			GasWanted: 99,
			GasUsed:   100,
			Codespace: "codespace",
			Events: []abci.Event{
				{
					Type: "message",
					Attributes: []abci.EventAttribute{
						{
							Key:   []byte("action"),
							Value: []byte("foo"),
							Index: true,
						},
					},
				},
			},
		},
	}
	deliverTxResult := &ctypes.ResultBroadcastTxCommit{
		Height: 10,
		Hash:   bytes.HexBytes([]byte("test")),
		DeliverTx: abci.ResponseDeliverTx{
			Code:      90,
			Data:      nil,
			Log:       `[]`,
			Info:      "info",
			GasWanted: 99,
			GasUsed:   100,
			Codespace: "codespace",
			Events: []abci.Event{
				{
					Type: "message",
					Attributes: []abci.EventAttribute{
						{
							Key:   []byte("action"),
							Value: []byte("foo"),
							Index: true,
						},
					},
				},
			},
		},
	}
	want := &sdk.TxResponse{
		Height:    10,
		TxHash:    "74657374",
		Codespace: "codespace",
		Code:      90,
		Data:      "",
		RawLog:    `[]`,
		Logs:      logs,
		Info:      "info",
		GasWanted: 99,
		GasUsed:   100,
		Events: []abci.Event{
			{
				Type: "message",
				Attributes: []abci.EventAttribute{
					{
						Key:   []byte("action"),
						Value: []byte("foo"),
						Index: true,
					},
				},
			},
		},
	}

	s.Require().Equal(want, sdk.NewResponseFormatBroadcastTxCommit(checkTxResult))
	s.Require().Equal(want, sdk.NewResponseFormatBroadcastTxCommit(deliverTxResult))
}

func TestWrapServiceResult(t *testing.T) {
	ctx := sdk.Context{}

	res, err := sdk.WrapServiceResult(ctx, nil, fmt.Errorf("test"))
	require.Nil(t, res)
	require.NotNil(t, err)

	res, err = sdk.WrapServiceResult(ctx, nil, nil)
	require.NotNil(t, res)
	require.Nil(t, err)
	require.Empty(t, res.Events)

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	ctx.EventManager().EmitEvent(sdk.NewEvent("test"))
	res, err = sdk.WrapServiceResult(ctx, nil, nil)
	require.NotNil(t, res)
	require.Nil(t, err)
	require.Len(t, res.Events, 1)

	spot := testdata.Dog{Name: "spot"}
	res, err = sdk.WrapServiceResult(ctx, &spot, nil)
	require.NotNil(t, res)
	require.Nil(t, err)
	require.Len(t, res.Events, 1)
	var spot2 testdata.Dog
	err = proto.Unmarshal(res.Data, &spot2)
	require.NoError(t, err)
	require.Equal(t, spot, spot2)
}

func TestNewResponseFormatBroadcastTx(t *testing.T) {
	hash, err := hex.DecodeString("00000000000000000000000000000000")
	require.NoError(t, err)
	result := ctypes.ResultBroadcastTx{
		Code:      1,
		Data:      []byte("some data"),
		Log:       `[{"log":"","msg_index":1,"success":true}]`,
		Codespace: "codespace",
		Hash:      hash,
	}

	txResponse := sdk.NewResponseFormatBroadcastTx(&result)

	require.NoError(t, err)
	require.Equal(t, result.Code, txResponse.Code)
	require.Equal(t, result.Data.String(), txResponse.Data)
	require.NotEmpty(t, txResponse.Logs)
	require.Equal(t, result.Log, txResponse.RawLog)
	require.Equal(t, result.Codespace, txResponse.Codespace)
	require.Equal(t, result.Hash.String(), txResponse.TxHash)
}
