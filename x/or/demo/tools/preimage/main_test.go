package main

import (
	"crypto/sha256"
	"testing"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/stretchr/testify/require"
)

func TestData(t *testing.T) {
	lastBlockHash := sha256Hash([]byte("last_block_hash"))
	lastAppHash := sha256Hash([]byte("last_app_hash"))
	data := sha256Hash([]byte("data"))
	blockBytes := append(append(lastBlockHash, lastAppHash...), data...)

	blockHash := sha256Hash(blockBytes)
	appHash := sha256Hash(lastAppHash, data)

	file, err := preimages("data/height-100.json")
	require.NoError(t, err)

	require.Equal(t, appHash[:], file[preimage.LocalIndexKey(0).PreimageKey()])
	require.Equal(t, blockHash[:], file[preimage.LocalIndexKey(1).PreimageKey()])
	require.Equal(t, blockBytes[:], file[sha256Key(blockHash).PreimageKey()])
}

func TestDefenderData(t *testing.T) {
	lastBlockHash := sha256Hash([]byte("last_block_hash"))
	lastAppHash := sha256Hash([]byte("last_app_hash"))
	data := sha256Hash([]byte("data"))
	mal := [32]byte{}
	blockBytes := append(append(lastBlockHash, lastAppHash...), mal[:]...)

	blockHash := sha256Hash(lastBlockHash, lastAppHash, data)
	appHash := sha256Hash(lastAppHash, mal[:])

	file, err := preimages("data/height-100-dfn.json")
	require.NoError(t, err)

	require.Equal(t, appHash[:], file[preimage.LocalIndexKey(0).PreimageKey()])
	require.Equal(t, blockHash[:], file[preimage.LocalIndexKey(1).PreimageKey()])
	require.Equal(t, blockBytes[:], file[sha256Key(blockHash).PreimageKey()])
}

func sha256Hash(b ...[]byte) []byte {
	h := sha256.New()
	for i := range b {
		h.Write(b[i])
	}
	return h.Sum(nil)
}
