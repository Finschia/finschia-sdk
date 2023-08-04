package mips

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
	"github.com/ethereum/go-ethereum/common"
)

const MemProofSize = 28 * 32

type StepWitness struct {
	// encoded state witness
	State []byte

	MemProof []byte

	PreimageKey    [32]byte // zeroed when no pre-image is accessed
	PreimageValue  []byte   // including the 8-byte length prefix
	PreimageOffset uint32
}

func WitnessState(p *types.Witness) *State {
	reg := [32]uint32{}
	for i := range reg {
		reg[i] = p.State.Registers[i]
	}
	return &State{
		Memory: &WitnessMemory{
			initialRoot: *(*[32]byte)(p.State.MemRoot),
			root:        *(*[32]byte)(p.State.MemRoot),
			proofs: [2][MemProofSize]byte{
				*(*[MemProofSize]byte)(p.Proofs[:MemProofSize]),
				*(*[MemProofSize]byte)(p.Proofs[MemProofSize : MemProofSize*2]),
			},
			pc: p.State.Pc,
		},
		PreimageKey:    common.BytesToHash(p.State.PreimageKey),
		PreimageOffset: p.State.PreimageOffset,
		PC:             p.State.Pc,
		NextPC:         p.State.NextPc,
		LO:             p.State.Lo,
		HI:             p.State.Hi,
		Heap:           p.State.Heap,
		ExitCode:       uint8(p.State.ExitCode),
		Exited:         p.State.Exited,
		Step:           p.State.Step,
		Registers:      reg,
	}
}

func WitnessStep(witness *types.Witness, challenge *types.Challenge) (*State, error) {
	// verify witness
	state := WitnessState(witness)
	preSteped := sha256.Sum256(state.EncodeWitness())
	if !bytes.Equal(challenge.AssertedStateHashes[challenge.L], preSteped[:]) {
		return nil, types.ErrInvalidWitness
	}

	// run single step
	images := map[[32]byte][]byte{}
	localPreimageKey := func(k uint64) (out [32]byte) {
		localKeyType := 1
		out[0] = byte(localKeyType)
		binary.BigEndian.PutUint64(out[24:], k)
		return
	}
	// TODO:
	// - get defender assertion from da module
	// - get block hash from da module
	images[localPreimageKey(0)], _ = hex.DecodeString("7d682198accf63e24956cd7883c8704154d968b24aaa7b1050f97d8204e11b1a") // defender assertion, this is malicious state, see demo/tools/data/height-100-def.json
	images[localPreimageKey(1)], _ = hex.DecodeString("5e33433d3f48974b10083d3ff6799e9d448b47406168ff2ee9da2e1e2035101b") // block hash

	oracle := NewPreimageOracle(images, witness)
	var stdOutBuf, stdErrBuf bytes.Buffer
	us := NewInstrumentedState(state, oracle, io.MultiWriter(&stdOutBuf, os.Stdout), io.MultiWriter(&stdErrBuf, os.Stderr))
	_, err := us.Step(false)
	if err != nil {
		return nil, err
	}

	return state, nil
}

type IMemory interface {
	AllocPage(pageIndex uint32) *CachedPage
	ForEachPage(fn func(pageIndex uint32, page *Page) error) error
	GetMemory(addr uint32) uint32
	Invalidate(addr uint32)
	MarshalJSON() ([]byte, error)
	MerkleProof(addr uint32) (out [MemProofSize]byte)
	MerkleRoot() [32]byte
	MerkleizeSubtree(gindex uint64) [32]byte
	PageCount() int
	ReadMemoryRange(addr uint32, count uint32) io.Reader
	SetMemory(addr uint32, v uint32)
	SetMemoryRange(addr uint32, r io.Reader) error
	UnmarshalJSON(data []byte) error
	Usage() string
}

type WitnessMemory struct {
	initialRoot [32]byte
	root        [32]byte
	proofs      [2][MemProofSize]byte
	pc          uint32
}

func (m *WitnessMemory) AllocPage(pageIndex uint32) *CachedPage {
	panic("cannot call 'AllocPage' on an witness memory")
}

func (m *WitnessMemory) ForEachPage(fn func(pageIndex uint32, page *Page) error) error {
	panic("cannot call 'ForEachPage' on an witness memory")
}

func (m *WitnessMemory) GetMemory(addr uint32) uint32 {
	val, ok := m.verifyProof(addr)
	if !ok {
		panic("invalid proof")
	}
	return val
}

func (m *WitnessMemory) Invalidate(addr uint32) {
	panic("cannot call 'Invalidate' on an witness memory")
}

func (m *WitnessMemory) MarshalJSON() ([]byte, error) {
	panic("cannot call 'MarshalJSON' on an witness memory")
}

func (m *WitnessMemory) MerkleProof(addr uint32) (out [MemProofSize]byte) {
	panic("cannot call 'MerkleProof' on an witness memory")
}

func (m *WitnessMemory) MerkleRoot() [32]byte {
	return m.root
}

func (m *WitnessMemory) MerkleizeSubtree(gindex uint64) [32]byte {
	panic("cannot call 'MerkleizeSubtree' on an witness memory")
}

func (m *WitnessMemory) PageCount() int {
	panic("cannot call 'PageCount' on an witness memory")
}

func (m *WitnessMemory) ReadMemoryRange(addr uint32, count uint32) io.Reader {
	return &emptyReader{}
}

func (m *WitnessMemory) SetMemory(addr uint32, v uint32) {
	_, ok := m.verifyProof(addr)
	if !ok {
		panic("invalid proof")
	}
	proof := m.proof(addr)
	root := *(*[32]byte)(proof[:32])
	binary.BigEndian.PutUint32(root[addr%32:addr%32+4], v)
	path := addr >> 5
	for i := 1; i < 28; i++ {
		sib := *(*[32]byte)(proof[i*32 : (i+1)*32])
		if (path >> (i - 1) & 1) == 0 {
			root = HashPair(root, sib)
		} else {
			root = HashPair(sib, root)
		}
	}
	m.root = root
}

func (m *WitnessMemory) SetMemoryRange(addr uint32, r io.Reader) error {
	panic("cannot call 'SetMemoryRange' on an witness memory")
}

func (m *WitnessMemory) UnmarshalJSON(data []byte) error {
	panic("cannot call 'UnmarshalJSON' on an witness memory")
}

func (m *WitnessMemory) Usage() string {
	panic("cannot call 'Usage' on an witness memory")
}

func (m *WitnessMemory) proof(addr uint32) [MemProofSize]byte {
	index := 0
	if m.pc != addr {
		index = 1
	}
	return m.proofs[index]
}

func (m *WitnessMemory) verifyProof(addr uint32) (uint32, bool) {
	proof := m.proof(addr)
	value := proof[:32][addr%32 : addr%32+4]
	root := *(*[32]byte)(proof[:32])
	path := addr >> 5
	for i := 1; i < 28; i++ {
		sib := *(*[32]byte)(proof[i*32 : (i+1)*32])
		if (path >> (i - 1) & 1) == 0 {
			root = HashPair(root, sib)
		} else {
			root = HashPair(sib, root)
		}
	}

	if !bytes.Equal(m.initialRoot[:], root[:]) {
		return 0, false
	}

	return binary.BigEndian.Uint32(value), true
}

type emptyReader struct {
}

func (emptyReader) Read(p []byte) (n int, err error) { return 0, nil }

func NewPreimageOracle(images map[[32]byte][]byte, witness *types.Witness) PreimageOracle {
	// TODO:
	// - verify preimage_value = hash(preimage_key)
	// - prevent value of local key rewriten
	images[*(*[32]byte)(witness.PreimageKey)] = witness.PreimageValue[8:] // remove prefix appended at https://github.com/ethereum-optimism/cannon/blob/b39dc7fb3a5bc1c5235351dc0aec175f299727ad/mipsevm/mips.go#L14-L17

	oracle := &simpleOracle{
		getPreimage: func(k [32]byte) []byte {
			p, ok := images[k]
			if !ok {
				panic(fmt.Sprintf("missing pre-image %s", k))
			}
			return p
		},
	}
	return oracle
}

type simpleOracle struct {
	hint        func(v []byte)
	getPreimage func(k [32]byte) []byte
}

func (t *simpleOracle) Hint(v []byte) {
	t.hint(v)
}

func (t *simpleOracle) GetPreimage(k [32]byte) []byte {
	return t.getPreimage(k)
}
