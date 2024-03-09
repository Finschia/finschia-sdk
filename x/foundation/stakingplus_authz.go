package foundation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ Authorization = (*CreateValidatorAuthorization)(nil)

func (a CreateValidatorAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{})
}

func (a CreateValidatorAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (AcceptResponse, error) {
	mCreate, ok := msg.(*stakingtypes.MsgCreateValidator)
	if !ok {
		return AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}

	if mCreate.ValidatorAddress != a.ValidatorAddress {
		return AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("validator address differs from the authorization's")
	}

	return AcceptResponse{Accept: true}, nil
}

func (a CreateValidatorAuthorization) ValidateBasic() error {
	return nil
}
