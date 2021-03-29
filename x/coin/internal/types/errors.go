package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const Codespace = ModuleName

var (
	ErrNoInputs                    = sdkerrors.Register(Codespace, 1, "no inputs to send transaction")
	ErrNoOutputs                   = sdkerrors.Register(Codespace, 2, "no outputs to send transaction")
	ErrInputOutputMismatch         = sdkerrors.Register(Codespace, 3, "sum inputs != sum outputs")
	ErrSendDisabled                = sdkerrors.Register(Codespace, 4, "send transactions are disabled")
	ErrCanNotTransferToBlacklisted = sdkerrors.Register(Codespace, 5, "Cannot transfer to safety box addresses")
	ErrRequestGetsLimit            = sdkerrors.Register(Codespace, 6, "the gets should be limited")
	ErrInvalidDenom                = sdkerrors.Register(Codespace, 7, "invalid denom")
)
