package grpc

const (
	// GRPCBlockHeightHeader is the gRPC header for block height.
	GRPCBlockHeightHeader = "x-cosmos-block-height"
	// GRPCCheckStateHeader is the gRPC header for mempool state. Assign "on" to this header when you want to query checkState values.
	// If you use both GRPCBlockHeightHeader and GRPCCheckStateHeader, GRPCCheckStateHeader would be ignored.
	GRPCCheckStateHeader = "x-lbm-checkstate"
)
