package collection

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const collectionCodespace = ModuleName

var (
	ErrTokenNotFound              = sdkerrors.Register(collectionCodespace, 3, "token not found")
	ErrNotMintable                = sdkerrors.Register(collectionCodespace, 4, "not mintable")
	ErrInvalidName                = sdkerrors.Register(collectionCodespace, 5, "invalid name")
	ErrInvalidTokenID             = sdkerrors.Register(collectionCodespace, 6, "invalid token id")
	ErrInvalidDecimals            = sdkerrors.Register(collectionCodespace, 7, "invalid decimals")
	ErrBadUseCase                 = sdkerrors.Register(collectionCodespace, 8, "bad use case")
	ErrInvalidBaseImgURI          = sdkerrors.Register(collectionCodespace, 10, "invalid base_img_uri")
	ErrInvalidClassID             = sdkerrors.Register(collectionCodespace, 12, "invalid class id")
	ErrContractNotFound           = sdkerrors.Register(collectionCodespace, 15, "contract not found")
	ErrClassNotFound              = sdkerrors.Register(collectionCodespace, 17, "class not found")
	ErrGrantNotFound              = sdkerrors.Register(collectionCodespace, 21, "grant not found")
	ErrParentNotFound             = sdkerrors.Register(collectionCodespace, 23, "parent not found")
	ErrWrongClass                 = sdkerrors.Register(collectionCodespace, 26, "class not supports this feature")
	ErrInvalidComposition         = sdkerrors.Register(collectionCodespace, 27, "invalid nft composition")
	ErrOperatorIsHolder           = sdkerrors.Register(collectionCodespace, 29, "operator and holder should be different")
	ErrAuthorizationNotFound      = sdkerrors.Register(collectionCodespace, 30, "authorization not found")
	ErrAuthorizationAlreadyExists = sdkerrors.Register(collectionCodespace, 31, "authorization already exists")
	ErrInvalidCoins               = sdkerrors.Register(collectionCodespace, 35, "invalid coins")
	ErrInvalidChanges             = sdkerrors.Register(collectionCodespace, 36, "invalid changes")
	ErrInvalidModificationTarget  = sdkerrors.Register(collectionCodespace, 39, "invalid modification target")
	ErrInsufficientTokens         = sdkerrors.Register(collectionCodespace, 41, "insufficient tokens")
	ErrInvalidMeta                = sdkerrors.Register(collectionCodespace, 43, "invalid meta")
	ErrInvalidPermission          = sdkerrors.Register(collectionCodespace, 49, "invalid permission")
	ErrGrantAlreadyExists         = sdkerrors.Register(collectionCodespace, 50, "grant already exists")
	ErrInvalidMintNFTParams       = sdkerrors.Register(collectionCodespace, 51, "invalid mint nft params")
	ErrInvalidContractID          = sdkerrors.Register(collectionCodespace, 52, "invalid contract id")
	ErrEmptyTokenIDs              = sdkerrors.Register(collectionCodespace, 53, "empty token ids")

	sdkToGRPC = map[*sdkerrors.Error]codes.Code{
		// this codespace
		ErrTokenNotFound:              codes.NotFound,
		ErrNotMintable:                codes.Unimplemented,
		ErrInvalidName:                codes.InvalidArgument,
		ErrInvalidTokenID:             codes.InvalidArgument,
		ErrInvalidDecimals:            codes.InvalidArgument,
		ErrBadUseCase:                 codes.InvalidArgument,
		ErrInvalidBaseImgURI:          codes.InvalidArgument,
		ErrInvalidClassID:             codes.InvalidArgument,
		ErrContractNotFound:           codes.NotFound,
		ErrClassNotFound:              codes.NotFound,
		ErrGrantNotFound:              codes.NotFound, // PermissionDenied
		ErrParentNotFound:             codes.NotFound,
		ErrWrongClass:                 codes.InvalidArgument,
		ErrInvalidComposition:         codes.FailedPrecondition,
		ErrOperatorIsHolder:           codes.InvalidArgument,
		ErrAuthorizationNotFound:      codes.NotFound, // PermissionDenied
		ErrAuthorizationAlreadyExists: codes.AlreadyExists,
		ErrInvalidCoins:               codes.InvalidArgument,
		ErrInvalidChanges:             codes.InvalidArgument,
		ErrInvalidModificationTarget:  codes.InvalidArgument,
		ErrInsufficientTokens:         codes.FailedPrecondition,
		ErrInvalidMeta:                codes.InvalidArgument,
		ErrInvalidPermission:          codes.InvalidArgument,
		ErrGrantAlreadyExists:         codes.AlreadyExists,
		ErrInvalidMintNFTParams:       codes.InvalidArgument,
		ErrInvalidContractID:          codes.InvalidArgument,
		ErrEmptyTokenIDs:              codes.InvalidArgument,

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
