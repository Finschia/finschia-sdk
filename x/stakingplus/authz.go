package stakingplus

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
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
	if _, err := sdk.AccAddressFromBech32(a.ValidatorAddress); err != nil {
		return err
	}

	return nil
}
