package tmservice

import (
	"context"

	ctypes "github.com/line/ostracon/rpc/core/types"

	"github.com/line/lfb-sdk/client"
)

func getBlock(clientCtx client.Context, height *int64) (*ctypes.ResultBlock, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.Block(context.Background(), height)
}

func getBlockByHash(clientCtx client.Context, hash []byte) (*ctypes.ResultBlock, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.BlockByHash(context.Background(), hash)
}

func getBlockResultsByHeight(clientCtx client.Context, height *int64) (*ctypes.ResultBlockResults, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.BlockResults(context.Background(), height)
}
