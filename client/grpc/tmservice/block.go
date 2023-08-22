package tmservice

import (
	"context"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/Finschia/finschia-rdk/client"
)

func getBlock(ctx context.Context, clientCtx client.Context, height *int64) (*ctypes.ResultBlock, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.Block(ctx, height)
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

func GetProtoBlock(ctx context.Context, clientCtx client.Context, height *int64) (tmproto.BlockID, *tmproto.Block, error) {
	block, err := getBlock(ctx, clientCtx, height)
	if err != nil {
		return tmproto.BlockID{}, nil, err
	}
	protoBlock, err := block.Block.ToProto()
	if err != nil {
		return tmproto.BlockID{}, nil, err
	}
	protoBlockID := block.BlockID.ToProto()

	return protoBlockID, protoBlock, nil
}
