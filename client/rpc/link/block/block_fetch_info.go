package block

type FetchInfo struct {
	inclusiveFromHeight int64
	exclusiveToHeight   int64
	hasMore             bool
	fetchItemCnt        int8
	fetchItemRange      []int64
}

func NewBlockFetchInfo(inclusiveFromHeight int64, exclusiveToHeight int64, hasMore bool) FetchInfo {
	fetchItemCnt := int8(exclusiveToHeight - inclusiveFromHeight)
	fetchItemRange := make([]int64, fetchItemCnt)
	for i := range fetchItemRange {
		fetchItemRange[i] = inclusiveFromHeight + int64(i)
	}
	return FetchInfo{inclusiveFromHeight, exclusiveToHeight, hasMore,
		fetchItemCnt, fetchItemRange,
	}
}

func NewFetchInfo(latestBlockHeight *int64, fromHeight *int64, fetchSize *int8) (fetchBlockHeight FetchInfo) {
	exclusiveToBlockHeight := *fromHeight + int64(*fetchSize)
	if *latestBlockHeight > exclusiveToBlockHeight-1 {
		fetchBlockHeight = NewBlockFetchInfo(*fromHeight, exclusiveToBlockHeight, true)
	} else if *latestBlockHeight == exclusiveToBlockHeight-1 {
		fetchBlockHeight = NewBlockFetchInfo(*fromHeight, exclusiveToBlockHeight, false)
	} else {
		fetchBlockHeight = NewBlockFetchInfo(*fromHeight, *latestBlockHeight+1, false)
	}
	return
}
