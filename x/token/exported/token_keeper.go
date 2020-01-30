package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

type TokenKeeper interface {
	GetModuleAddress() sdk.AccAddress
	SetToken(sdk.Context, types.Token) sdk.Error
	GetToken(sdk.Context, string) (types.Token, sdk.Error)
	GetAllTokens(sdk.Context) types.Tokens
	IterateTokens(sdk.Context, func(types.Token) bool)
	AddPermission(sdk.Context, sdk.AccAddress, types.PermissionI)
	RemovePermission(sdk.Context, sdk.AccAddress, types.PermissionI) sdk.Error
	HasPermission(sdk.Context, sdk.AccAddress, types.PermissionI) bool
	InheritPermission(sdk.Context, sdk.AccAddress, sdk.AccAddress)
	GrantPermission(sdk.Context, sdk.AccAddress, sdk.AccAddress, types.PermissionI) sdk.Error
	MintToken(sdk.Context, sdk.Coin, sdk.AccAddress) sdk.Error
	BurnToken(sdk.Context, sdk.Coin, sdk.AccAddress) sdk.Error
	GetSupply(sdk.Context, string) (sdk.Int, sdk.Error)
	ModifyTokenURI(sdk.Context, sdk.AccAddress, string, string) sdk.Error
	TransferFT(sdk.Context, sdk.AccAddress, sdk.AccAddress, string, sdk.Int) sdk.Error
	TransferCFT(sdk.Context, sdk.AccAddress, sdk.AccAddress, string, string, sdk.Int) sdk.Error
	TransferNFT(sdk.Context, sdk.AccAddress, sdk.AccAddress, string) sdk.Error
	TransferCNFT(sdk.Context, sdk.AccAddress, sdk.AccAddress, string, string) sdk.Error
	Attach(sdk.Context, sdk.AccAddress, string, string, string) sdk.Error
	Detach(sdk.Context, sdk.AccAddress, sdk.AccAddress, string, string) sdk.Error
}
