package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
)

func main() {
	po := preimage.NewOracleClient(preimage.ClientPreimageChannel())

	expected := *(*[32]byte)(po.Get(preimage.LocalIndexKey(0)))

	blockhash := *(*[32]byte)(po.Get(preimage.LocalIndexKey(1)))
	blockBytes := po.Get(sha256Key(*(*[32]byte)(blockhash[:32])))
	block := DecodeBlock(blockBytes)

	// compute apphash
	h := sha256.New()
	h.Write(block.LastAppHash)
	h.Write(block.Data)
	appHash := h.Sum(nil)

	// result
	if bytes.Equal(expected[:], appHash) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// app_hash = hash(pre_app_hash, data)
type Block struct {
	LastBlockHash []byte // 256bit
	LastAppHash   []byte // 256bit
	Data          []byte
}

func (b Block) Encode() []byte {
	bs := []byte{}
	bs = append(bs, b.LastBlockHash...)
	bs = append(bs, b.LastAppHash...)
	bs = append(bs, b.Data...)
	return bs
}

func (b Block) Hash() []byte {
	h := sha256.New()
	h.Write(b.Encode())
	return h.Sum(nil)
}

func DecodeBlock(bs []byte) Block {
	b := Block{}
	b.LastBlockHash = bs[:32]
	b.LastAppHash = bs[32:64]
	b.Data = bs[64:]
	return b
}

// custom preimage key
type sha256Key [32]byte

const sha256KeyType = 100

func (s sha256Key) PreimageKey() (out [32]byte) {
	out = s                      // copy the keccak hash
	out[0] = byte(sha256KeyType) // apply prefix
	return
}

func (s sha256Key) String() string {
	return "0x" + hex.EncodeToString(s[:])
}

func (s sha256Key) TerminalString() string {
	return "0x" + hex.EncodeToString(s[:])
}
