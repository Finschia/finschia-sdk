package client

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/v2/x/ibc/core/02-client/client/cli"
	"github.com/line/lbm-sdk/v2/x/ibc/core/02-client/types"
)

// Name returns the IBC client name
func Name() string {
	return types.SubModuleName
}

// GetQueryCmd returns no root query command for the IBC client
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterQueryService registers the gRPC query service for IBC client.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
