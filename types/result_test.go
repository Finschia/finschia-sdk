package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/line/lbm-sdk/codec"
)

func TestParseABCILog(t *testing.T) {
	logs := `[{"log":"","msg_index":1,"success":true}]`

	res, err := ParseABCILogs(logs)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, res[0].Log, "")
	require.Equal(t, res[0].MsgIndex, uint16(1))
}

func TestABCIMessageLog(t *testing.T) {
	events := Events{NewEvent("transfer", NewAttribute("sender", "foo"))}
	msgLog := NewABCIMessageLog(0, "", events)

	msgLogs := ABCIMessageLogs{msgLog}
	bz, err := codec.Cdc.MarshalJSON(msgLogs)
	require.NoError(t, err)
	require.Equal(t, string(bz), msgLogs.String())
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

	txResponse := NewResponseFormatBroadcastTx(&result)

	require.NoError(t, err)
	require.Equal(t, result.Code, txResponse.Code)
	require.Equal(t, result.Data.String(), txResponse.Data)
	require.NotEmpty(t, txResponse.Logs)
	require.Equal(t, result.Log, txResponse.RawLog)
	require.Equal(t, result.Codespace, txResponse.Codespace)
	require.Equal(t, result.Hash.String(), txResponse.TxHash)
}
