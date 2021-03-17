package tmservice

import (
	"context"

	ctypes "github.com/line/ostracon/rpc/core/types"

	"github.com/line/lbm-sdk/v2/client"
)

func getNodeStatus(clientCtx client.Context) (*ctypes.ResultStatus, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}
	return node.Status(context.Background())
}
