package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Codec interface {
	MarshalJSONIndent(res interface{}) ([]byte, error)
	MarshalJSON(res interface{}) ([]byte, error)
	UnmarshalBinaryLengthPrefixed(txBytes []byte) (tx auth.StdTx, err error)
}
type Wrapper struct {
	cdc *codec.Codec
}

func NewCodecWrapper(cdc *codec.Codec) *Wrapper {
	wrapper := &Wrapper{cdc}
	return wrapper
}

func (cw *Wrapper) MarshalJSONIndent(res interface{}) ([]byte, error) {
	return cw.cdc.MarshalJSONIndent(res, "", "  ")
}
func (cw *Wrapper) MarshalJSON(res interface{}) ([]byte, error) {
	return cw.cdc.MarshalJSON(res)
}
func (cw *Wrapper) UnmarshalBinaryLengthPrefixed(txBytes []byte) (tx auth.StdTx, err error) {
	if err := cw.cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx); err != nil {
		return tx, err
	}
	return
}

type HasMoreResponseWrapper struct {
	Items   []*FetchResultWithTxRes `json:"items"`
	HasMore bool                    `json:"has_more"`
}

type FetchResultWithTxRes struct {
	ResultBlock *ctypes.ResultBlock `json:"result_block"`
	TxResponses []*sdk.TxResponse   `json:"tx_responses"`
}
