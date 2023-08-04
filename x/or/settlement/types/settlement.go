package types

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
)

const (
	SectionN = int64(16) // 16sect search
)

func (c *Challenge) ID() string {
	h := sha256.New()
	h.Write(binary.LittleEndian.AppendUint64([]byte{}, uint64(c.BlockHeight)))
	h.Write([]byte(c.Challenger))
	h.Write([]byte(c.Defender))
	h.Write([]byte(c.RollupName))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c Challenge) IsSearching() bool {
	return c.L+1 != c.R
}

func (c Challenge) GetSteps() []uint64 {
	sub := sdk.NewIntFromUint64(c.R - c.L)
	steps := []uint64{}

	for i := int64(1); i < SectionN; i++ {
		step := sub.MulRaw(i).QuoRaw(SectionN).Uint64() + c.L
		if c.L == step {
			continue
		}

		if len(steps) > 0 && steps[len(steps)-1] == step {
			continue
		}

		steps = append(steps, step)
	}
	return steps
}

func (c Challenge) CurrentResponder() string {
	if len(c.AssertedStateHashes) == len(c.DefendedStateHashes) {
		return c.Challenger
	}
	return c.Defender
}

func DecodeState(input []byte) (*State, error) {
	if len(input) != 226 {
		return nil, fmt.Errorf("invalid state bytes length, length %d", len(input))
	}
	state := State{}
	memoryRoot := input[:32]
	state.MemRoot = memoryRoot
	state.PreimageKey = input[32:64]
	state.PreimageOffset = binary.BigEndian.Uint32(input[64:68])
	state.Pc = binary.BigEndian.Uint32(input[68:72])
	state.NextPc = binary.BigEndian.Uint32(input[72:76])
	state.Lo = binary.BigEndian.Uint32(input[76:80])
	state.Hi = binary.BigEndian.Uint32(input[80:84])
	state.Heap = binary.BigEndian.Uint32(input[84:88])
	state.ExitCode = uint32(input[88])
	if input[89] == 1 {
		state.Exited = true
	} else if input[89] == 0 {
		state.Exited = false
	} else {
		return nil, fmt.Errorf("invalid state bytes data, index 89")
	}
	state.Step = binary.BigEndian.Uint64(input[90:98])
	for i := 0; i < 32; i++ {
		index := i*4 + 98
		register := binary.BigEndian.Uint32(input[index : index+4])
		state.Registers = append(state.Registers, register)
	}

	return &state, nil
}
