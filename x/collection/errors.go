package collection

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const collectionCodespace = ModuleName

var (
	ErrInvalidContractID          = sdkerrors.Register(collectionCodespace, 2, "invalid contract id")
	ErrContractNotFound           = sdkerrors.Register(collectionCodespace, 3, "contract not found")
	ErrInvalidClassID             = sdkerrors.Register(collectionCodespace, 4, "invalid class id")
	ErrWrongClass                 = sdkerrors.Register(collectionCodespace, 5, "class not supports this feature")
	ErrClassNotFound              = sdkerrors.Register(collectionCodespace, 6, "class not found")
	ErrInvalidTokenID             = sdkerrors.Register(collectionCodespace, 7, "invalid token id")
	ErrTokenNotFound              = sdkerrors.Register(collectionCodespace, 8, "token not found")
	ErrInvalidPermission          = sdkerrors.Register(collectionCodespace, 9, "invalid permission")
	ErrGrantNotFound              = sdkerrors.Register(collectionCodespace, 10, "grant not found")
	ErrGrantAlreadyExists         = sdkerrors.Register(collectionCodespace, 11, "grant already exists")
	ErrOperatorIsHolder           = sdkerrors.Register(collectionCodespace, 12, "operator and holder should be different")
	ErrAuthorizationNotFound      = sdkerrors.Register(collectionCodespace, 13, "authorization not found")
	ErrAuthorizationAlreadyExists = sdkerrors.Register(collectionCodespace, 14, "authorization already exists")
	ErrInvalidCoins               = sdkerrors.Register(collectionCodespace, 15, "invalid coins")
	ErrInsufficientTokens         = sdkerrors.Register(collectionCodespace, 16, "insufficient tokens")
	ErrInvalidName                = sdkerrors.Register(collectionCodespace, 17, "invalid name")
	ErrInvalidBaseImgURI          = sdkerrors.Register(collectionCodespace, 18, "invalid base_img_uri")
	ErrInvalidMeta                = sdkerrors.Register(collectionCodespace, 19, "invalid meta")
	ErrInvalidDecimals            = sdkerrors.Register(collectionCodespace, 20, "invalid decimals")
	ErrNotMintable                = sdkerrors.Register(collectionCodespace, 21, "not mintable")
	ErrBadUseCase                 = sdkerrors.Register(collectionCodespace, 22, "bad use case")
	ErrInvalidChanges             = sdkerrors.Register(collectionCodespace, 23, "invalid changes")
	ErrInvalidModificationTarget  = sdkerrors.Register(collectionCodespace, 24, "invalid modification target")
	ErrParentNotFound             = sdkerrors.Register(collectionCodespace, 25, "parent not found")
	ErrInvalidComposition         = sdkerrors.Register(collectionCodespace, 26, "invalid nft composition")
	ErrInvalidMintNFTParams       = sdkerrors.Register(collectionCodespace, 27, "invalid mint nft params")
	ErrEmptyTokenIDs              = sdkerrors.Register(collectionCodespace, 28, "empty token ids")

	sdkToGRPC = map[*sdkerrors.Error]codes.Code{
		// this codespace
		ErrContractNotFound:           codes.NotFound,
		ErrClassNotFound:              codes.NotFound,
		ErrTokenNotFound:              codes.NotFound,
		ErrGrantNotFound:              codes.NotFound, // PermissionDenied
		ErrAuthorizationNotFound:      codes.NotFound, // PermissionDenied
		ErrParentNotFound:             codes.NotFound,
		ErrInvalidContractID:          codes.InvalidArgument,
		ErrInvalidClassID:             codes.InvalidArgument,
		ErrInvalidTokenID:             codes.InvalidArgument,
		ErrInvalidPermission:          codes.InvalidArgument,
		ErrOperatorIsHolder:           codes.InvalidArgument,
		ErrInvalidCoins:               codes.InvalidArgument,
		ErrInvalidName:                codes.InvalidArgument,
		ErrInvalidDecimals:            codes.InvalidArgument,
		ErrInvalidBaseImgURI:          codes.InvalidArgument,
		ErrWrongClass:                 codes.InvalidArgument,
		ErrBadUseCase:                 codes.InvalidArgument,
		ErrInvalidChanges:             codes.InvalidArgument,
		ErrInvalidModificationTarget:  codes.InvalidArgument,
		ErrInvalidMeta:                codes.InvalidArgument,
		ErrInvalidMintNFTParams:       codes.InvalidArgument,
		ErrEmptyTokenIDs:              codes.InvalidArgument,
		ErrInsufficientTokens:         codes.FailedPrecondition,
		ErrInvalidComposition:         codes.FailedPrecondition,
		ErrGrantAlreadyExists:         codes.AlreadyExists,
		ErrAuthorizationAlreadyExists: codes.AlreadyExists,
		ErrNotMintable:                codes.Unimplemented,

		// sdk codespace
		sdkerrors.ErrInvalidAddress: codes.InvalidArgument,
		sdkerrors.ErrInvalidType:    codes.InvalidArgument,
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
