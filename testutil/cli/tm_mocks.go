package cli

import (
	"context"

	tmbytes "github.com/Finschia/ostracon/libs/bytes"
	rpcclient "github.com/Finschia/ostracon/rpc/client"
	rpcclientmock "github.com/Finschia/ostracon/rpc/client/mock"
	coretypes "github.com/Finschia/ostracon/rpc/core/types"
	tmtypes "github.com/Finschia/ostracon/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-sdk/client"
)

var _ client.TendermintRPC = (*MockTendermintRPC)(nil)

type MockTendermintRPC struct {
	rpcclientmock.Client

	responseQuery abci.ResponseQuery
}

// NewMockTendermintRPC returns a mock TendermintRPC implementation.
// It is used for CLI testing.
func NewMockTendermintRPC(respQuery abci.ResponseQuery) MockTendermintRPC {
	return MockTendermintRPC{responseQuery: respQuery}
}

func (MockTendermintRPC) BroadcastTxSync(context.Context, tmtypes.Tx) (*coretypes.ResultBroadcastTx, error) {
	return &coretypes.ResultBroadcastTx{Code: 0}, nil
}

func (m MockTendermintRPC) ABCIQueryWithOptions(
	_ context.Context,
	_ string,
	_ tmbytes.HexBytes,
	_ rpcclient.ABCIQueryOptions,
) (*coretypes.ResultABCIQuery, error) {
	return &coretypes.ResultABCIQuery{Response: m.responseQuery}, nil
}
