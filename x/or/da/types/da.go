package types

import "time"

func NewCCRef(txhash []byte, msgIdx, batchSize uint32, totalFrames uint64, batchRoot []byte) *CCRef {
	return &CCRef{
		TxHash:      txhash,
		MsgIndex:    msgIdx,
		TotalFrames: totalFrames,
		BatchSize:   batchSize,
		BatchRoot:   batchRoot,
	}
}

func NewSCCRef(totalFrames uint64, batchSize uint32, batchRoot []byte, ISRs [][]byte, blockTime time.Time) *SCCRef {
	return &SCCRef{
		TotalFrames:            totalFrames,
		BatchSize:              batchSize,
		BatchRoot:              batchRoot,
		IntermediateStateRoots: ISRs,
		Deadline:               blockTime,
	}
}
