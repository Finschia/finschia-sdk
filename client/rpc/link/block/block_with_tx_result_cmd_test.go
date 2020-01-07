package block

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/golang/mock/gomock"
	. "github.com/line/link/client/rpc/mock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"math"
	"strconv"
	"strings"
	"testing"
)

const ballotX = "\u2717"
const failMsg = "The code did not panic" + ballotX

func TestValidateCmd(t *testing.T) {
	rootCmd := &cobra.Command{
		Use: "linkcli",
	}

	queryCmd := &cobra.Command{
		Use: "query",
	}
	rootCmd.AddCommand(queryCmd)

	blockWithTxResultCmd := WithTxResultCommand(codec.New())
	queryCmd.AddCommand(blockWithTxResultCmd)

	tests := []struct {
		reason  string
		args    []string
		wantErr bool
	}{
		{"misspelled command", []string{"block_with_tx_result"}, true},
		{"no command provided", []string{}, false},
		{"help flag", []string{"block-with-tx-result", "--help"}, false},
		{"shorthand help flag", []string{"block-with-tx-result", "-h"}, false},
	}

	for _, tt := range tests {
		err := client.ValidateCmd(blockWithTxResultCmd, tt.args)
		require.Equal(t, tt.wantErr, err != nil, tt.reason)
	}
}

func TestProcess(t *testing.T) {

	t.Log("normal case", checkMark)
	{
		args, mockCliContext, mockClient, mockTendermint, mockCodec, rb := prepareCMD(t)

		mockCliContext.EXPECT().GetNode().Return(mockClient, nil).Times(2)
		mockCliContext.EXPECT().TrustNode().Return(false).Times(1)
		check := tmtypes.SignedHeader{}
		mockClient.EXPECT().Block(gomock.Any()).Return(rb, nil).Times(2)
		mockCliContext.EXPECT().Verify(rb.Block.Height).Return(check, nil).Times(1)
		mockTendermint.EXPECT().ValidateBlockMeta(rb.BlockMeta, check).Return(nil).Times(1)
		mockTendermint.EXPECT().ValidateBlock(rb.Block, check).Return(nil).Times(1)
		mockCliContext.EXPECT().Indent().Return(true).Times(1)
		data := []byte(`"hello"`)
		mockCodec.EXPECT().MarshalJSONIndent(gomock.Any()).Return(data, nil).Times(1)
		mockCodec.EXPECT().MarshalJSON(gomock.Any()).Times(0)

		actual, err := process(args, &Util{lcliCtx: mockCliContext, ltmtl: mockTendermint, lcdc: mockCodec})
		require.Equal(t, data, actual)
		require.Equal(t, nil, err)
	}
	t.Log("error case", checkMark)
	{
		args, mockCliContext, _, mockTendermint, mockCodec, _ := prepareCMD(t)

		expectedErr := fmt.Errorf("could not create client")
		mockCliContext.EXPECT().GetNode().Return(nil, expectedErr).Times(1)
		actual, err := process(args, &Util{lcliCtx: mockCliContext, ltmtl: mockTendermint, lcdc: mockCodec})
		require.Equal(t, expectedErr, err)
		require.Equal(t, "", string(actual))
	}
}

func prepareCMD(t *testing.T) ([]string, *MockCLIContext, *MockClient, *MockTendermint, *MockCodec, *ctypes.ResultBlock) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fromBlockHeight := int64(1)
	args := []string{strconv.Itoa(int(fromBlockHeight)), "1"}
	mockCliContext := NewMockCLIContext(ctrl)
	mockClient := NewMockClient(ctrl)
	mockTendermint := NewMockTendermint(ctrl)
	mockCodec := NewMockCodec(ctrl)
	rb := &ctypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: fromBlockHeight},
		},
		BlockMeta: &tmtypes.BlockMeta{
			BlockID: tmtypes.BlockID{},
			Header:  tmtypes.Header{},
		},
	}
	return args, mockCliContext, mockClient, mockTendermint, mockCodec, rb
}

var cmdTestTable = []struct {
	msg             string
	args            []string
	fromBlockHeight int64
	fetchSize       int8
	failMsg         string
}{
	{msg: "parse fromBlockHeight fetchSizeInt8",
		args:            []string{"1", "20"},
		fromBlockHeight: 1,
		fetchSize:       DefaultBlockFetchSize,
	},
	{msg: "maximum fetchsize(2nd param) will be adjusted to DefaultBlockFetchSize automatically",
		args:            []string{"1", "21"},
		fromBlockHeight: 1,
		fetchSize:       DefaultBlockFetchSize,
	},
	{msg: "fetchsize(2nd param) is optional",
		args:            []string{"1"},
		fromBlockHeight: 1,
		fetchSize:       DefaultBlockFetchSize,
	},
	{msg: "overflow blockheight",
		args:            []string{"922337203685471422412475807", "21"},
		fromBlockHeight: 1,
		failMsg:         failMsg,
	},
	{msg: "fetchsize(2nd param) must be in the data range of int8",
		args:            []string{"1", strconv.Itoa(math.MaxInt16)},
		fromBlockHeight: 1,
		failMsg:         failMsg,
	},
	{msg: "more than expected args",
		args:    []string{"1", strconv.Itoa(math.MaxInt16), "3"},
		failMsg: failMsg,
	},
	{msg: "blank args",
		args:    nil,
		failMsg: failMsg,
	},
}

func TestParseCmdParams(t *testing.T) {
	for _, testingOption := range cmdTestTable {
		t.Logf("%s, args %s, fromBlockHeight %d", testingOption.msg, strings.Join(testingOption.args, ","),
			testingOption.fromBlockHeight)
		{
			if testingOption.failMsg != "" {
				assert.Panics(t, func() {
					parseCmdParams(testingOption.args)
				}, failMsg)
			} else {
				fromBlockHeight, fetchSizeInt8 := parseCmdParams(testingOption.args)
				require.Equal(t, testingOption.fromBlockHeight, fromBlockHeight)
				require.Equal(t, testingOption.fetchSize, fetchSizeInt8)
			}
		}
	}
}
