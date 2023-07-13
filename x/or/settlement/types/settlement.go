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
