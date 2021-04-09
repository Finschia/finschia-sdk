package types

import (
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz/exported"
)

var (
	_ sdk.MsgRequest = &MsgGrantAuthorizationRequest{}
	_ sdk.MsgRequest = &MsgRevokeAuthorizationRequest{}
	_ sdk.MsgRequest = &MsgExecAuthorizedRequest{}

	_ types.UnpackInterfacesMessage = &MsgGrantAuthorizationRequest{}
	_ types.UnpackInterfacesMessage = &MsgExecAuthorizedRequest{}
)

// NewMsgGrantAuthorization creates a new MsgGrantAuthorization
//nolint:interfacer
func NewMsgGrantAuthorization(granter sdk.AccAddress, grantee sdk.AccAddress, authorization exported.Authorization, expiration time.Time) (*MsgGrantAuthorizationRequest, error) {
	m := &MsgGrantAuthorizationRequest{
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		Expiration: expiration,
	}
	err := m.SetAuthorization(authorization)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetSigners implements Msg
func (msg MsgGrantAuthorizationRequest) GetSigners() []sdk.AccAddress {
	err := sdk.ValidateAccAddress(msg.Granter)
	if err != nil {
		panic(err)
	}
	granter := sdk.AccAddress(msg.Granter)
	return []sdk.AccAddress{granter}
}

// ValidateBasic implements Msg
func (msg MsgGrantAuthorizationRequest) ValidateBasic() error {
	err := sdk.ValidateAccAddress(msg.Granter)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid granter address")
	}
	granter := sdk.AccAddress(msg.Granter)

	err = sdk.ValidateAccAddress(msg.Grantee)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid granter address")
	}
	grantee := sdk.AccAddress(msg.Grantee)

	if granter.Equals(grantee) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "granter and grantee cannot be same")
	}

	if msg.Expiration.Unix() < time.Now().Unix() {
		return sdkerrors.Wrap(ErrInvalidExpirationTime, "Time can't be in the past")
	}

	authorization, ok := msg.Authorization.GetCachedValue().(exported.Authorization)
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expected %T, got %T", (exported.Authorization)(nil), msg.Authorization.GetCachedValue())
	}
	return authorization.ValidateBasic()
}

// GetGrantAuthorization returns the cache value from the MsgGrantAuthorization.Authorization if present.
func (msg *MsgGrantAuthorizationRequest) GetGrantAuthorization() exported.Authorization {
	authorization, ok := msg.Authorization.GetCachedValue().(exported.Authorization)
	if !ok {
		return nil
	}
	return authorization
}

// SetAuthorization converts Authorization to any and adds it to MsgGrantAuthorization.Authorization.
func (msg *MsgGrantAuthorizationRequest) SetAuthorization(authorization exported.Authorization) error {
	m, ok := authorization.(proto.Message)
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrPackAny, "can't proto marshal %T", m)
	}
	any, err := types.NewAnyWithValue(m)
	if err != nil {
		return err
	}
	msg.Authorization = any
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgExecAuthorizedRequest) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	for _, x := range msg.Msgs {
		var msgExecAuthorized sdk.MsgRequest
		err := unpacker.UnpackAny(x, &msgExecAuthorized)
		if err != nil {
			return err
		}
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgGrantAuthorizationRequest) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var authorization exported.Authorization
	return unpacker.UnpackAny(msg.Authorization, &authorization)
}

// NewMsgRevokeAuthorization creates a new MsgRevokeAuthorization
//nolint:interfacer
func NewMsgRevokeAuthorization(granter sdk.AccAddress, grantee sdk.AccAddress, methodName string) MsgRevokeAuthorizationRequest {
	return MsgRevokeAuthorizationRequest{
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		MethodName: methodName,
	}
}

// GetSigners implements Msg
func (msg MsgRevokeAuthorizationRequest) GetSigners() []sdk.AccAddress {
	err := sdk.ValidateAccAddress(msg.Granter)
	if err != nil {
		panic(err)
	}
	granter := sdk.AccAddress(msg.Granter)
	return []sdk.AccAddress{granter}
}

// ValidateBasic implements MsgRequest.ValidateBasic
func (msg MsgRevokeAuthorizationRequest) ValidateBasic() error {
	err := sdk.ValidateAccAddress(msg.Granter)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid granter address")
	}
	granter := sdk.AccAddress(msg.Granter)

	err = sdk.ValidateAccAddress(msg.Grantee)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid grantee address")
	}
	grantee := sdk.AccAddress(msg.Grantee)

	if granter.Equals(grantee) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "granter and grantee cannot be same")
	}

	if msg.MethodName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing method name")
	}

	return nil
}

// NewMsgExecAuthorized creates a new MsgExecAuthorized
//nolint:interfacer
func NewMsgExecAuthorized(grantee sdk.AccAddress, msgs []sdk.ServiceMsg) MsgExecAuthorizedRequest {
	msgsAny := make([]*types.Any, len(msgs))
	for i, msg := range msgs {
		bz, err := proto.Marshal(msg.Request)
		if err != nil {
			panic(err)
		}

		anyMsg := &types.Any{
			TypeUrl: msg.MethodName,
			Value:   bz,
		}

		msgsAny[i] = anyMsg
	}

	return MsgExecAuthorizedRequest{
		Grantee: grantee.String(),
		Msgs:    msgsAny,
	}
}

// GetServiceMsgs returns the cache values from the MsgExecAuthorized.Msgs if present.
func (msg MsgExecAuthorizedRequest) GetServiceMsgs() ([]sdk.ServiceMsg, error) {
	msgs := make([]sdk.ServiceMsg, len(msg.Msgs))
	for i, msgAny := range msg.Msgs {
		msgReq, ok := msgAny.GetCachedValue().(sdk.MsgRequest)
		if !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "messages contains %T which is not a sdk.MsgRequest", msgAny)
		}
		srvMsg := sdk.ServiceMsg{
			MethodName: msgAny.TypeUrl,
			Request:    msgReq,
		}

		msgs[i] = srvMsg
	}

	return msgs, nil
}

// GetSigners implements Msg
func (msg MsgExecAuthorizedRequest) GetSigners() []sdk.AccAddress {
	err := sdk.ValidateAccAddress(msg.Grantee)
	if err != nil {
		panic(err)
	}
	grantee := sdk.AccAddress(msg.Grantee)
	return []sdk.AccAddress{grantee}
}

// ValidateBasic implements Msg
func (msg MsgExecAuthorizedRequest) ValidateBasic() error {
	err := sdk.ValidateAccAddress(msg.Grantee)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid grantee address")
	}

	if len(msg.Msgs) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "messages cannot be empty")
	}

	return nil
}
