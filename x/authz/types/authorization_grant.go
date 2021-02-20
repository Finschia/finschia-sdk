package types

import (
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/codec/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz/exported"
)

// NewAuthorizationGrant returns new AuthrizationGrant
func NewAuthorizationGrant(authorization exported.Authorization, expiration time.Time) (AuthorizationGrant, error) {
	auth := AuthorizationGrant{
		Expiration: expiration,
	}
	msg, ok := authorization.(proto.Message)
	if !ok {
		return AuthorizationGrant{}, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", authorization)
	}

	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return AuthorizationGrant{}, err
	}

	auth.Authorization = any

	return auth, nil
}

var (
	_ types.UnpackInterfacesMessage = &AuthorizationGrant{}
)

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (auth AuthorizationGrant) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var authorization exported.Authorization
	return unpacker.UnpackAny(auth.Authorization, &authorization)
}

// GetAuthorizationGrant returns the cached value from the AuthorizationGrant.Authorization if present.
func (auth AuthorizationGrant) GetAuthorizationGrant() exported.Authorization {
	authorization, ok := auth.Authorization.GetCachedValue().(exported.Authorization)
	if !ok {
		return nil
	}
	return authorization
}
