package types

import (
	"bytes"
	"testing"

	sdk "github.com/line/lfb-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestGetContractByCodeIDSecondaryIndexPrefix(t *testing.T) {
	specs := map[string]struct {
		src uint64
		exp []byte
	}{
		"small number": {src: 1,
			exp: []byte{6, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		"big number": {src: 1 << (8 * 7),
			exp: []byte{6, 1, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got := GetContractByCodeIDSecondaryIndexPrefix(spec.src)
			assert.Equal(t, spec.exp, got)
		})
	}
}

func TestGetContractByCreatedSecondaryIndexKey(t *testing.T) {
	e := ContractCodeHistoryEntry{
		CodeID:  1,
		Updated: &AbsoluteTxPosition{2 + 1<<(8*7), 3 + 1<<(8*7)},
	}
	addr := sdk.BytesToAccAddress(bytes.Repeat([]byte{4}, sdk.BytesAddrLen))
	got := GetContractByCreatedSecondaryIndexKey(addr, e)
	exp := []byte{6, // prefix
		0, 0, 0, 0, 0, 0, 0, 1, // codeID
		1, 0, 0, 0, 0, 0, 0, 2, // height
		1, 0, 0, 0, 0, 0, 0, 3, // index
		0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x71, 0x73, 0x7a, 0x71, 0x67, 0x70, 0x71, 0x79, 0x71, 0x73, 0x7a, 0x71, 0x67, 0x70, 0x71, 0x79, 0x71, 0x73, 0x7a, 0x71, 0x67, 0x70, 0x71, 0x79, 0x71, 0x73, 0x7a, 0x71, 0x67, 0x70, 0x71, 0x79, 0x68, 0x65, 0x70, 0x61, 0x73, 0x6c, // address
	}
	assert.Equal(t, exp, got)
}
