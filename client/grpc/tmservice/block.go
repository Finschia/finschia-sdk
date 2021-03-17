package tmservice

import (
	"context"

	ctypes "github.com/line/ostracon/rpc/core/types"

	"github.com/line/lbm-sdk/v2/client"
)

func getBlock(clientCtx client.Context, height *int64) (*ctypes.ResultBlock, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.Block(context.Background(), height)
}
