package contract

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/contract/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/contract/internal/types"
)

func ValidateContractIDBasic(contract Msg) error {
	if !keeper.VerifyContractID(contract.GetContractID()) {
		return sdkerrors.Wrapf(types.ErrInvalidContractID, "ContractID: %s", contract.GetContractID())
	}
	return nil
}
