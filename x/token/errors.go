package token

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const tokenCodespace = ModuleName

var (
	ErrEmpty              = sdkerrors.Register(tokenCodespace, 2, "empty value")
	ErrDuplicate          = sdkerrors.Register(tokenCodespace, 3, "duplicate value")
	ErrMaxLimit           = sdkerrors.Register(tokenCodespace, 4, "limit exceeded")
	ErrInvalid            = sdkerrors.Register(tokenCodespace, 5, "invalid value")
	ErrNotFound           = sdkerrors.Register(tokenCodespace, 6, "not found")
	ErrAlreadyExists      = sdkerrors.Register(tokenCodespace, 7, "already exists")
	ErrInsufficientTokens = sdkerrors.Register(tokenCodespace, 8, "insufficient tokens")

	sdkToGRPC = map[*sdkerrors.Error]codes.Code{
		// this codespace
		ErrEmpty:              codes.InvalidArgument,
		ErrDuplicate:          codes.InvalidArgument,
		ErrMaxLimit:           codes.InvalidArgument,
		ErrInvalid:            codes.InvalidArgument,
		ErrNotFound:           codes.NotFound,
		ErrAlreadyExists:      codes.AlreadyExists,
		ErrInsufficientTokens: codes.FailedPrecondition,

		// sdk codespace
		sdkerrors.ErrInvalidAddress: codes.InvalidArgument,
		sdkerrors.ErrInvalidType:    codes.InvalidArgument,
		sdkerrors.ErrUnauthorized:   codes.PermissionDenied,
	}
)

func SDKErrorToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	for sdkerror, grpcCode := range sdkToGRPC {
		if sdkerror.Is(err) {
			return status.Error(grpcCode, sdkerror.Error())
		}
	}

	panic("unknown error")
}
