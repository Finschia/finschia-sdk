package tendermint

import (
	"github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
)

// Name returns the IBC client name
func Name() string {
	return types.SubModuleName
}
