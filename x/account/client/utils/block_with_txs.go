package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/line/lbm-sdk/v2/x/account/client/types"

	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func LatestBlockHeight(cliCtx context.CLIContext) (int64, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return -1, err
	}

	// Get the latest block
	latestBlock, err := node.Block(nil)
	if err != nil {
		return -1, err
	}

	return latestBlock.Block.Height, nil
}

func BlockWithTxResponses(cliCtx context.CLIContext, latestBlockHeight, fromBlockHeight, fetchSize int64) (blockWithRxResultsWrapper *types.HasMoreResponseWrapper, err error) {
	fbh := types.NewFetchInfo(latestBlockHeight, fromBlockHeight, fetchSize)
	results := make([]*types.ResultBlockWithTxResponses, len(fbh.FetchItemRange))
	for idx, height := range fbh.FetchItemRange {
		block, err := getBlock(cliCtx, height)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while fetching a block by blockHeight(%d), err(%s)", height, err)
		}
		txs, err := getTxs(cliCtx, height)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while fetching a block by blockHeight(%d), err(%s)", height, err)
		}
		results[idx] = &types.ResultBlockWithTxResponses{
			ResultBlock: &types.ResultBlock{
				BlockSize: block.Block.Size(),
				BlockID:   block.BlockID,
				Block:     block.Block,
			},
			TxResponses: txs,
		}
	}

	return &types.HasMoreResponseWrapper{
		Items:   results,
		HasMore: fbh.HasMore,
	}, nil
}

func getBlock(cliCtx context.CLIContext, height int64) (*ctypes.ResultBlock, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}
	resultBlock, err := node.Block(&height)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		check, err := cliCtx.Verify(resultBlock.Block.Height)
		if err != nil {
			return nil, err
		}

		if err := tmliteProxy.ValidateHeader(&resultBlock.Block.Header, check); err != nil {
			return nil, err
		}

		if err = tmliteProxy.ValidateBlock(resultBlock.Block, check); err != nil {
			return nil, err
		}
	}
	return resultBlock, nil
}

func getTxs(cliCtx context.CLIContext, height int64) ([]types.TxResponse, error) {
	const defaultLimit = 100
	// nolint:prealloc
	txResponses := []types.TxResponse{}

	nextTxPage := 1
	for {
		searchResult, err := QueryTxsByEvents(cliCtx, []string{fmt.Sprintf("tx.height=%d", height)}, nextTxPage, defaultLimit)
		if err != nil {
			return nil, err
		}

		txResponses = append(txResponses, searchResult.Txs...)

		nextTxPage++
		if nextTxPage > searchResult.PageTotal {
			break
		}
	}
	return txResponses, nil
}
