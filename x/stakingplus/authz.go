package stakingplus

import (
	sdk "github.com/Finschia/finschia-rdk/types"
	sdkerrors "github.com/Finschia/finschia-rdk/types/errors"
	"github.com/Finschia/finschia-rdk/x/foundation"
	stakingtypes "github.com/Finschia/finschia-rdk/x/staking/types"
)

var _ foundation.Authorization = (*CreateValidatorAuthorization)(nil)

func (a CreateValidatorAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{})
}

func (a CreateValidatorAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (foundation.AcceptResponse, error) {
	mCreate, ok := msg.(*stakingtypes.MsgCreateValidator)
	if !ok {
		return foundation.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}

	if mCreate.ValidatorAddress != a.ValidatorAddress {
		return foundation.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("validator address differs from the authorization's")
	}

	return foundation.AcceptResponse{Accept: true}, nil
}

func (a CreateValidatorAuthorization) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(a.ValidatorAddress); err != nil {
		return err
	}

	return nil
}
