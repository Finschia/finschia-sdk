package types

func NewCCRef(txhash []byte, msgIdx, batchSize uint32, totalFrames uint64, batchRoot []byte) *CCRef {
	return &CCRef{
		TxHash:      txhash,
		MsgIndex:    msgIdx,
		TotalFrames: totalFrames,
		BatchSize:   batchSize,
		BatchRoot:   batchRoot,
	}
}
