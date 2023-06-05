package keeper

import (
	"bytes"
	"compress/zlib"
	sdkerror "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
	"io"
)

//nolint:unused
func (k Keeper) appendSequencerBatch() {
	panic("implement me")
}

func (k Keeper) DecompressCCBatch(origin types.CompressedCCBatch) (*types.CCBatch, error) {
	switch origin.Compression {
	case types.OptionZLIB:
		b := bytes.NewReader(origin.Data)
		r, err := zlib.NewReader(b)
		defer r.Close()
		if err != nil {
			return nil, types.ErrInvalidCompressedData.Wrap(err.Error())
		}
		buf := make([]byte, types.DefaultCCBatchMaxBytes)
		n, err := r.Read(buf)
		out := buf[:n]
		if err != nil && err != io.EOF {
			return nil, err
		}
		batch := new(types.CCBatch)
		k.cdc.MustUnmarshal(out, batch)

		return batch, nil

	case types.OptionZSTD:
		return nil, types.ErrInvalidCompressedData.Wrapf("compression %s not supported", origin.Compression)
	case types.OptionEmpty:
		return nil, types.ErrInvalidCompressedData
	default:
		return nil, sdkerror.ErrInvalidRequest.Wrapf("no compression option provided")
	}
}
