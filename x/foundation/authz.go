package foundation

import (
	"github.com/cosmos/gogoproto/proto"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// Authorization represents the interface of various Authorization types implemented
// by other modules.
// Caution: It's not for x/authz exec.
type Authorization interface {
	proto.Message

	// MsgTypeURL returns the fully-qualified Msg service method URL (as described in ADR 031),
	// which will process and accept or reject a request.
	MsgTypeURL() string

	// Accept determines whether this grant permits the provided sdk.Msg to be performed,
	// and if so provides an updated authorization instance.
	Accept(ctx sdk.Context, msg sdk.Msg) (AcceptResponse, error)

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error
}

// AcceptResponse instruments the controller of an authz message if the request is accepted
// and if it should be updated or deleted.
type AcceptResponse struct {
	// If Accept=true, the controller can accept and authorization and handle the update.
	Accept bool
	// If Delete=true, the controller must delete the authorization object and release
	// storage resources.
	Delete bool
	// Controller, who is calling Authorization.Accept must check if `Updated != nil`. If yes,
	// it must use the updated version and handle the update on the storage level.
	Updated Authorization
}

var _ Authorization = (*ReceiveFromTreasuryAuthorization)(nil)

func (a ReceiveFromTreasuryAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgWithdrawFromTreasury{})
}

func (a ReceiveFromTreasuryAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (AcceptResponse, error) {
	_, ok := msg.(*MsgWithdrawFromTreasury)
	if !ok {
		return AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}

	return AcceptResponse{Accept: true}, nil
}

func (a ReceiveFromTreasuryAuthorization) ValidateBasic() error {
	return nil
}
