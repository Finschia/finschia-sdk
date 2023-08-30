package types

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var (
	ErrInvalidCompressedData    = sdkerrors.Register(ModuleName, 1100, "this data cannot be decompressed.")
	ErrInvalidCCBatch           = sdkerrors.Register(ModuleName, 1101, "invalid cc batch")
	ErrInvalidQueueTx           = sdkerrors.Register(ModuleName, 1102, "invalid queue tx")
	ErrCCStateNotFound          = sdkerrors.Register(ModuleName, 1103, "cc state not found")
	ErrCCRefNotFound            = sdkerrors.Register(ModuleName, 1104, "cc reference not found")
	ErrQueueTxStateNotFound     = sdkerrors.Register(ModuleName, 1105, "queue tx state not found")
	ErrQueueTxNotFound          = sdkerrors.Register(ModuleName, 1106, "queue tx not found")
	ErrL2HeightBatchMapNotFound = sdkerrors.Register(ModuleName, 1107, "this l2 height does not match any processed batches")
	ErrSCCStateNotFound         = sdkerrors.Register(ModuleName, 1201, "scc state not found")
	ErrSCCRefNotFound           = sdkerrors.Register(ModuleName, 1202, "scc reference not found")
	ErrInvalidSCCBatch          = sdkerrors.Register(ModuleName, 1203, "invalid scc batch")
)
