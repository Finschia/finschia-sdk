package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
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
	Items   []*FetchResult `json:"items"`
	HasMore bool           `json:"hasMore"`
}

type FetchResult struct {
	Block     *ctypes.ResultBlock `json:"block"`
	TxResults []*ctypes.ResultTx  `json:"txResults"`
}
