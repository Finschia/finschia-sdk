package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/wasm/internal/types"
)

var (
	gzipIdent = []byte("\x1F\x8B\x08")
	wasmIdent = []byte("\x00\x61\x73\x6D")
)

// List of CLI flags
const (
	FlagPage       = "page"
	FlagLimit      = "limit"
	FlagPageKey    = "page-key"
	FlagOffset     = "offset"
	FlagCountTotal = "count-total"
)

// IsGzip returns checks if the file contents are gzip compressed
func IsGzip(input []byte) bool {
	return bytes.Equal(input[:3], gzipIdent)
}

// IsWasm checks if the file contents are of wasm binary
func IsWasm(input []byte) bool {
	return bytes.Equal(input[:4], wasmIdent)
}

// GzipIt compresses the input ([]byte)
func GzipIt(input []byte) ([]byte, error) {
	// Create gzip writer.
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(input)
	if err != nil {
		return nil, err
	}
	err = w.Close() // You must close this first to flush the bytes to the buffer.
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// AddPaginationFlagsToCmd adds common pagination flags to cmd
func AddPaginationFlagsToCmd(cmd *cobra.Command, query string) {
	cmd.Flags().Uint64(FlagPage, 1, fmt.Sprintf("pagination page of %s to query for. This sets offset to a multiple of limit", query))
	cmd.Flags().String(FlagPageKey, "", fmt.Sprintf("pagination page-key of %s to query for", query))
	cmd.Flags().Uint64(FlagOffset, 0, fmt.Sprintf("pagination offset of %s to query for", query))
	cmd.Flags().Uint64(FlagLimit, 100, fmt.Sprintf("pagination limit of %s to query for", query))
	cmd.Flags().Bool(
		FlagCountTotal, false, fmt.Sprintf("count total number of records in %s to query for", query))
}

// BigEndianToUint64 returns an uint64 from big endian encoded bytes. If encoding
// is empty, zero is returned.
// This function is included in cosmos-sdk v0.40.0
// Once cosmos-sdk is updated, use the sdk functions.
func BigEndianToUint64(bz []byte) uint64 {
	if len(bz) == 0 {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

// ReadPageRequest reads and builds the necessary page request flags for pagination.
func ReadPageRequest(flagSet *pflag.FlagSet) (*types.PageRequest, error) {
	pageKey, err := flagSet.GetString(FlagPageKey)
	if err != nil {
		return nil, err
	}
	offset, err := flagSet.GetUint64(FlagOffset)
	if err != nil {
		return nil, err
	}
	limit, err := flagSet.GetUint64(FlagLimit)
	if err != nil {
		return nil, err
	}
	countTotal, err := flagSet.GetBool(FlagCountTotal)
	if err != nil {
		return nil, err
	}
	page, err := flagSet.GetUint64(FlagPage)
	if err != nil {
		return nil, err
	}

	return NewPageRequest(pageKey, offset, limit, page, countTotal)
}

func NewPageRequest(pageKey string, offset, limit, page uint64, countTotal bool) (*types.PageRequest, error) {
	if page > 1 && offset > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "page and offset cannot be used together")
	}

	if page > 1 {
		offset = (page - 1) * limit
	}

	return &types.PageRequest{
		Key:        []byte(pageKey),
		Offset:     offset,
		Limit:      limit,
		CountTotal: countTotal,
	}, nil
}
