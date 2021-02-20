package types

import (
	ocproto "github.com/line/ostracon/proto/ostracon/types"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/authz/exported"
)

var (
	_ exported.Authorization = &GenericAuthorization{}
)

// NewGenericAuthorization creates a new GenericAuthorization object.
func NewGenericAuthorization(methodName string) *GenericAuthorization {
	return &GenericAuthorization{
		MessageName: methodName,
	}
}

// MethodName implements Authorization.MethodName.
func (cap GenericAuthorization) MethodName() string {
	return cap.MessageName
}

// Accept implements Authorization.Accept.
func (cap GenericAuthorization) Accept(msg sdk.ServiceMsg, block ocproto.Header) (updated exported.Authorization, delete bool, err error) {
	return &cap, false, nil
}
