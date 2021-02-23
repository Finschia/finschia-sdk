package contract

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract/internal/keeper"
	"github.com/line/lbm-sdk/x/contract/internal/types"
)

func ValidateContractIDBasic(contract Msg) error {
	if !keeper.VerifyContractID(contract.GetContractID()) {
		return sdkerrors.Wrapf(types.ErrInvalidContractID, "ContractID: %s", contract.GetContractID())
	}
	return nil
}
