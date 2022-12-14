package token

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const tokenCodespace = ModuleName

var (
	ErrInvalidContractID          = sdkerrors.Register(tokenCodespace, 2, "invalid contract id")
	ErrContractNotFound           = sdkerrors.Register(tokenCodespace, 3, "contract not found")
	ErrWrongContract              = sdkerrors.Register(tokenCodespace, 4, "contract not supports this feature")
	ErrInvalidPermission          = sdkerrors.Register(tokenCodespace, 5, "invalid permission")
	ErrGrantNotFound              = sdkerrors.Register(tokenCodespace, 6, "grant not found")
	ErrOperatorIsHolder           = sdkerrors.Register(tokenCodespace, 7, "operator and holder should be different")
	ErrAuthorizationNotFound      = sdkerrors.Register(tokenCodespace, 8, "authorization not found")
	ErrAuthorizationAlreadyExists = sdkerrors.Register(tokenCodespace, 9, "authorization already exists")
	ErrInvalidAmount              = sdkerrors.Register(tokenCodespace, 10, "invalid amount")
	ErrInsufficientTokens         = sdkerrors.Register(tokenCodespace, 11, "insufficient tokens")
	ErrInvalidName                = sdkerrors.Register(tokenCodespace, 12, "invalid name")
	ErrInvalidSymbol              = sdkerrors.Register(tokenCodespace, 13, "invalid symbol")
	ErrInvalidImageURI            = sdkerrors.Register(tokenCodespace, 14, "invalid image_uri")
	ErrInvalidMeta                = sdkerrors.Register(tokenCodespace, 15, "invalid meta")
	ErrInvalidDecimals            = sdkerrors.Register(tokenCodespace, 16, "invalid decimals")
	ErrInvalidChanges             = sdkerrors.Register(tokenCodespace, 17, "invalid changes")

	sdkToGRPC = map[*sdkerrors.Error]codes.Code{
		// this codespace
		ErrContractNotFound:           codes.NotFound,
		ErrGrantNotFound:              codes.NotFound,
		ErrAuthorizationNotFound:      codes.NotFound,
		ErrInvalidContractID:          codes.InvalidArgument,
		ErrInvalidPermission:          codes.InvalidArgument,
		ErrOperatorIsHolder:           codes.InvalidArgument,
		ErrInvalidAmount:              codes.InvalidArgument,
		ErrInvalidName:                codes.InvalidArgument,
		ErrInvalidDecimals:            codes.InvalidArgument,
		ErrInvalidImageURI:            codes.InvalidArgument,
		ErrWrongContract:              codes.InvalidArgument,
		ErrInvalidChanges:             codes.InvalidArgument,
		ErrInvalidMeta:                codes.InvalidArgument,
		ErrInsufficientTokens:         codes.FailedPrecondition,
		ErrAuthorizationAlreadyExists: codes.AlreadyExists,

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
