package types

import (
	tmtypes "github.com/tendermint/tendermint/types"
)

type FetchInfo struct {
	InclusiveFromHeight int64
	ExclusiveToHeight   int64
	HasMore             bool
	FetchItemCnt        int64
	FetchItemRange      []int64
}

func NewBlockFetchInfo(inclusiveFromHeight int64, exclusiveToHeight int64, hasMore bool) FetchInfo {
	fetchItemCnt := exclusiveToHeight - inclusiveFromHeight
	fetchItemRange := make([]int64, fetchItemCnt)
	for i := range fetchItemRange {
		fetchItemRange[i] = inclusiveFromHeight + int64(i)
	}
	return FetchInfo{inclusiveFromHeight, exclusiveToHeight, hasMore,
		fetchItemCnt, fetchItemRange,
	}
}

func NewFetchInfo(latestBlockHeight, fromHeight, fetchSize int64) FetchInfo {
	var fetchBlockHeight FetchInfo
	switch exclusiveToBlockHeight := fromHeight + fetchSize; {
	case latestBlockHeight > exclusiveToBlockHeight-1:
		fetchBlockHeight = NewBlockFetchInfo(fromHeight, exclusiveToBlockHeight, true)
	case latestBlockHeight == exclusiveToBlockHeight-1:
		fetchBlockHeight = NewBlockFetchInfo(fromHeight, exclusiveToBlockHeight, false)
	default:
		fetchBlockHeight = NewBlockFetchInfo(fromHeight, latestBlockHeight+1, false)
	}
	return fetchBlockHeight
}

type HasMoreResponseWrapper struct {
	Items   []*ResultBlockWithTxResponses `json:"items"`
	HasMore bool                          `json:"has_more"`
}

type ResultBlockWithTxResponses struct {
	ResultBlock *ResultBlock `json:"result_block"`
	TxResponses []TxResponse `json:"tx_responses"`
}

type ResultBlock struct {
	BlockSize int             `json:"block_size"`
	BlockID   tmtypes.BlockID `json:"block_id"`
	Block     *tmtypes.Block  `json:"block"`
}
