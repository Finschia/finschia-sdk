package tmservice

import (
	"context"

	ctypes "github.com/line/ostracon/rpc/core/types"

	"github.com/line/lfb-sdk/client"
)

func getNodeStatus(clientCtx client.Context) (*ctypes.ResultStatus, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}
	return node.Status(context.Background())
}
