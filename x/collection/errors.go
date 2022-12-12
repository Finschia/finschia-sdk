package collection

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const collectionCodespace = ModuleName

var (
	ErrEmpty              = sdkerrors.Register(collectionCodespace, 2, "empty value")
	ErrDuplicate          = sdkerrors.Register(collectionCodespace, 3, "duplicate value")
	ErrMaxLimit           = sdkerrors.Register(collectionCodespace, 4, "limit exceeded")
	ErrInvalid            = sdkerrors.Register(collectionCodespace, 5, "invalid value")
	ErrNotFound           = sdkerrors.Register(collectionCodespace, 6, "not found")
	ErrAlreadyExists      = sdkerrors.Register(collectionCodespace, 7, "already exists")
	ErrWrongClass         = sdkerrors.Register(collectionCodespace, 8, "class not supports this feature")
	ErrInsufficientTokens = sdkerrors.Register(collectionCodespace, 9, "insufficient tokens")
	ErrCompositionFailed  = sdkerrors.Register(collectionCodespace, 10, "failed composition precondition")

	sdkToGRPC = map[*sdkerrors.Error]codes.Code{
		// this codespace
		ErrEmpty:              codes.InvalidArgument,
		ErrDuplicate:          codes.InvalidArgument,
		ErrMaxLimit:           codes.InvalidArgument,
		ErrInvalid:            codes.InvalidArgument,
		ErrNotFound:           codes.NotFound,
		ErrAlreadyExists:      codes.AlreadyExists,
		ErrWrongClass:         codes.InvalidArgument,
		ErrInsufficientTokens: codes.FailedPrecondition,
		ErrCompositionFailed:  codes.InvalidArgument,

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
