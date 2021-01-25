package types

import (
	"reflect"

	ocproto "github.com/line/ostracon/proto/ostracon/types"

	sdk "github.com/line/lbm-sdk/types"
	bank "github.com/line/lbm-sdk/x/bank/types"
)

var (
	_ Authorization = &SendAuthorization{}
)

// NewSendAuthorization creates a new SendAuthorization object.
func NewSendAuthorization(spendLimit sdk.Coins) *SendAuthorization {
	return &SendAuthorization{
		SpendLimit: spendLimit,
	}
}

// MethodName implements Authorization.MethodName.
func (authorization SendAuthorization) MethodName() string {
	return "/lbm.bank.v1.Msg/Send"
}

// Accept implements Authorization.Accept.
func (authorization SendAuthorization) Accept(msg sdk.ServiceMsg, block ocproto.Header) (allow bool, updated Authorization, delete bool) {
	if reflect.TypeOf(msg.Request) == reflect.TypeOf(&bank.MsgSend{}) {
		msg, ok := msg.Request.(*bank.MsgSend)
		if ok {
			limitLeft, isNegative := authorization.SpendLimit.SafeSub(msg.Amount)
			if isNegative {
				return false, nil, false
			}
			if limitLeft.IsZero() {
				return true, nil, true
			}

			return true, &SendAuthorization{SpendLimit: limitLeft}, false
		}
	}
	return false, nil, false
}
