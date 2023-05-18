package or

import (
	"context"
	"fmt"
	"time"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/types"
	da "github.com/Finschia/finschia-sdk/x/or/da/types"
	ostypes "github.com/Finschia/ostracon/types"
)

// This source is based on Rollkit's implementation for Celestia.

// StatusCode is a type for DA layer return status.
// TODO: define an enum of different non-happy-path cases
// that might need to be handled by Rollkit independent of
// the underlying DA chain.
type StatusCode uint64

// Data Availability return codes.
const (
	StatusUnknown StatusCode = iota
	StatusSuccess
	StatusNotFound
	StatusError
)

// BaseResult contains basic information returned by DA layer.
type BaseResult struct {
	// Code is to determine if the action succeeded.
	Code StatusCode
	// Message may contain DA layer specific information (like DA block height/hash, detailed error message, etc)
	Message string
	// DAHeight informs about a height on Data Availability Layer for given result.
	DAHeight uint64
}

// ResultSubmitBlock contains information returned from DA layer after block submission.
type ResultSubmitBlock struct {
	BaseResult
	// Not sure if this needs to be bubbled up to other
	// parts of Rollkit.
	// Hash hash.Hash
}

// ResultCheckBlock contains information about block availability, returned from DA layer client.
type ResultCheckBlock struct {
	BaseResult
	// DataAvailable is the actual answer whether the block is available or not.
	// It can be true if and only if Code is equal to StatusSuccess.
	DataAvailable bool
}

// ResultRetrieveBlocks contains batch of blocks returned from DA layer client.
type ResultRetrieveBlocks struct {
	BaseResult
	// Block is the full block retrieved from Data Availability Layer.
	// If Code is not equal to StatusSuccess, it has to be nil.
	Blocks []*ostypes.Block
}

// DataAvailabilityLayerClient use celestia-node public API.
type DataAvailabilityLayerClient struct {
	client da.MsgClient
}

func NewClient(ctx client.Context) DataAvailabilityLayerClient {
	return DataAvailabilityLayerClient{
		client: da.NewMsgClient(ctx),
	}
}

// Config stores Celestia DALC configuration parameters.
type Config struct {
	BaseURL  string        `json:"base_url"`
	Timeout  time.Duration `json:"timeout"`
	Fee      int64         `json:"fee"`
	GasLimit uint64        `json:"gas_limit"`
}

// Init initializes DataAvailabilityLayerClient instance.
func (c *DataAvailabilityLayerClient) Init() error {
	return nil
}

// Start prepares DataAvailabilityLayerClient to work.
func (c *DataAvailabilityLayerClient) Start() error {
	return nil
}

// Stop stops DataAvailabilityLayerClient.
func (c *DataAvailabilityLayerClient) Stop() error {
	return nil
}

// SubmitBlock submits a block to DA layer.
func (c *DataAvailabilityLayerClient) SubmitBlock(ctx context.Context, block *ostypes.Block) ResultSubmitBlock {

	// TODO Fields must be fixed.
	msg := da.MsgAppendCTCBatch{
		RollupName: "",
		Batch: da.CTCBatch{
			ShouldStartAtElement: types.NewInt(0),
			BatchContexts: []*da.CTCBatchContext{
				{
					NumSequencedTxs:       0,
					NumSubsequentQueueTxs: 0,
					Timestamp:             time.Now(),
					L1Height:              0,
				},
			},
			Elements: []*da.CTCBatchElement{
				{
					Timestamp:  time.Now(),
					L1Height:   0,
					Txraw:      []byte{},
					QueueIndex: 0,
					L2Height:   0,
				},
			},
			Compression: da.OptionEmpty,
		},
	}

	_, err := c.client.AppendCTCBatch(ctx, &msg)
	if err != nil {
		return ResultSubmitBlock{
			BaseResult: BaseResult{
				Code:    StatusError,
				Message: err.Error(),
			},
		}
	}

	// TODO What does the response include?
	return ResultSubmitBlock{
		BaseResult: BaseResult{
			Code:     StatusSuccess,
			Message:  fmt.Sprintf("tx hash: " /* TODO: %s", res.TxHash */),
			DAHeight: uint64(0 /* TODO: res.Height */),
		},
	}
}

// CheckBlockAvailability queries DA layer to check data availability of block at given height.
func (c *DataAvailabilityLayerClient) CheckBlockAvailability(ctx context.Context, dataLayerHeight uint64) ResultCheckBlock {
	// TODO The gRPC API for DA layer isn't yet defined.
	panic("not implemented")
}

// RetrieveBatch gets a batch of blocks from DA layer.
func (c *DataAvailabilityLayerClient) RetrieveBatch(ctx context.Context, dataLayerHeight uint64) ResultRetrieveBlocks {
	// TODO The gRPC API for DA layer isn't yet defined.
	panic("not implemented")
}
