package types

import (
	"fmt"

	sdkerrors "github.com/line/lbm-sdk/types/errors"

	channeltypes "github.com/line/lbm-sdk/x/ibc/core/04-channel/types"
)

const (
	// ackErrorString defines a string constant included in error acknowledgements
	// NOTE: Changing this const is state machine breaking as acknowledgements are written into state
	ackErrorString = "error handling packet on host chain: see events for details"
)

// NewErrorAcknowledgement returns a deterministic error string which may be used in
// the packet acknowledgement.
func NewErrorAcknowledgement(err error) channeltypes.Acknowledgement {
	// the ABCI code is included in the abcitypes.ResponseDeliverTx hash
	// constructed in Ostracon and is therefore determinstic
	_, code, _ := sdkerrors.ABCIInfo(err, false) // discard non-deterministic codespace and log values

	errorString := fmt.Sprintf("ABCI code: %d: %s", code, ackErrorString)

	return channeltypes.NewErrorAcknowledgement(errorString)
}
