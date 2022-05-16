package foundation

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

var _ authz.Authorization = (*ReceiveFromTreasuryAuthorization)(nil)

func (a ReceiveFromTreasuryAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgWithdrawFromTreasury{})
}

func (a ReceiveFromTreasuryAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	_, ok := msg.(*MsgWithdrawFromTreasury)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}

	return authz.AcceptResponse{Accept: true}, nil
}

func (a ReceiveFromTreasuryAuthorization) ValidateBasic() error {
	return nil
}

var _ authz.Authorization = (*CreateValidatorAuthorization)(nil)

func (a CreateValidatorAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{})
}

func (a CreateValidatorAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	mCreate, ok := msg.(*stakingtypes.MsgCreateValidator)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}

	if mCreate.ValidatorAddress != a.ValidatorAddress {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("validator address differs from the authorization's")
	}

	return authz.AcceptResponse{Accept: true}, nil
}

func (a CreateValidatorAuthorization) ValidateBasic() error {
	if err := sdk.ValidateValAddress(a.ValidatorAddress); err != nil {
		return err
	}

	return nil
}
