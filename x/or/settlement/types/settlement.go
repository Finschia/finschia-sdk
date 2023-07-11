package types

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

func (c *Challenge) ID() string {
	h := sha256.New()
	h.Write(binary.LittleEndian.AppendUint64([]byte{}, uint64(c.BlockHeight)))
	h.Write([]byte(c.Challenger))
	h.Write([]byte(c.Defender))
	h.Write([]byte(c.RollupName))
	return fmt.Sprintf("%x", h.Sum(nil))
}
