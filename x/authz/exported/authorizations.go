package exported

import (
	"github.com/gogo/protobuf/proto"

	ocproto "github.com/line/ostracon/proto/ostracon/types"

	sdk "github.com/line/lbm-sdk/types"
)

// Authorization represents the interface of various Authorization types.
type Authorization interface {
	proto.Message

	// MethodName returns the fully-qualified Msg service method name as described in ADR 031.
	MethodName() string

	// Accept determines whether this grant permits the provided sdk.ServiceMsg to be performed, and if
	// so provides an upgraded authorization instance.
	Accept(msg sdk.ServiceMsg, block ocproto.Header) (updated Authorization, delete bool, err error)
}
