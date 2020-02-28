package contract

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/contract/internal/keeper"
	"github.com/line/link/x/contract/internal/types"
)

func ValidateContractIDBasic(contract Msg) sdk.Error {
	if !keeper.VerifyContractID(contract.GetContractID()) {
		return types.ErrInvalidContractID(types.DefaultCodespace, contract.GetContractID())
	}
	return nil
}
