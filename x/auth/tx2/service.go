package tx2

import (
	"context"

	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Finschia/finschia-rdk/client"
	"github.com/Finschia/finschia-rdk/client/grpc/tmservice"
	codectypes "github.com/Finschia/finschia-rdk/codec/types"
	sdk "github.com/Finschia/finschia-rdk/types"
	sdkerrors "github.com/Finschia/finschia-rdk/types/errors"
	pagination "github.com/Finschia/finschia-rdk/types/query"
	txtypes "github.com/Finschia/finschia-rdk/types/tx"
	tx2types "github.com/Finschia/finschia-rdk/types/tx2"
)

type tx2Server struct {
	clientCtx         client.Context
	interfaceRegistry codectypes.InterfaceRegistry
}

func NewTx2Server(clientCtx client.Context, interfaceRegistry codectypes.InterfaceRegistry) tx2types.ServiceServer {
	return tx2Server{
		clientCtx:         clientCtx,
		interfaceRegistry: interfaceRegistry,
	}
}

var _ tx2types.ServiceServer = tx2Server{}

// protoTxProvider is a type which can provide a proto transaction. It is a
// workaround to get access to the wrapper TxBuilder's method GetProtoTx().
// ref: https://github.com/cosmos/cosmos-sdk/issues/10347
type protoTxProvider interface {
	GetProtoTx() *txtypes.Tx
}

func (s tx2Server) GetBlockWithTxs(ctx context.Context, req *tx2types.GetBlockWithTxsRequest) (*tx2types.GetBlockWithTxsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()

	if req.Height < 1 || req.Height > currentHeight {
		return nil, sdkerrors.ErrInvalidHeight.Wrapf("requested height %d but height must not be less than 1 "+
			"or greater than the current height %d", req.Height, currentHeight)
	}

	blockID, block, err := tmservice.GetProtoBlock(ctx, s.clientCtx, &req.Height)
	if err != nil {
		return nil, err
	}

	var offset, limit uint64
	if req.Pagination != nil {
		offset = req.Pagination.Offset
		limit = req.Pagination.Limit
	} else {
		offset = 0
		limit = pagination.DefaultLimit
	}

	blockTxs := block.Data.Txs
	blockTxsLn := uint64(len(blockTxs))
	txs := make([]*txtypes.Tx, 0, limit)
	if offset >= blockTxsLn && blockTxsLn != 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("out of range: cannot paginate %d txs with offset %d and limit %d", blockTxsLn, offset, limit)
	}
	decodeTxAt := func(i uint64) error {
		tx := blockTxs[i]
		txb, err := s.clientCtx.TxConfig.TxDecoder()(tx)
		if err != nil {
			return err
		}
		p, ok := txb.(protoTxProvider)
		if !ok {
			return sdkerrors.ErrTxDecode.Wrapf("could not cast %T to %T", txb, txtypes.Tx{})
		}
		txs = append(txs, p.GetProtoTx())
		return nil
	}
	if req.Pagination != nil && req.Pagination.Reverse {
		for i, count := offset, uint64(0); i > 0 && count != limit; i, count = i-1, count+1 {
			if err = decodeTxAt(i); err != nil {
				return nil, err
			}
		}
	} else {
		for i, count := offset, uint64(0); i < blockTxsLn && count != limit; i, count = i+1, count+1 {
			if err = decodeTxAt(i); err != nil {
				return nil, err
			}
		}
	}

	return &tx2types.GetBlockWithTxsResponse{
		Txs:     txs,
		BlockId: &blockID,
		Block:   block,
		Pagination: &pagination.PageResponse{
			Total: blockTxsLn,
		},
	}, nil
}

// RegisterTxService registers the tx service on the gRPC router.
func RegisterTxService(
	qrt gogogrpc.Server,
	clientCtx client.Context,
	interfaceRegistry codectypes.InterfaceRegistry,
) {
	tx2types.RegisterServiceServer(
		qrt,
		NewTx2Server(clientCtx, interfaceRegistry),
	)
}

// RegisterGRPCGatewayRoutes mounts the tx service's GRPC-gateway routes on the
// given Mux.
func RegisterGRPCGatewayRoutes(clientConn gogogrpc.ClientConn, mux *runtime.ServeMux) {
	if err := tx2types.RegisterServiceHandlerClient(context.Background(), mux, tx2types.NewServiceClient(clientConn)); err != nil {
		panic(err)
	}
}
