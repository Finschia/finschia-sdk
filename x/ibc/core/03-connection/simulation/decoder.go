package simulation

import (
	"bytes"
	"fmt"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/line/lbm-sdk/v2/x/ibc/core/03-connection/types"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding connection type.
func NewDecodeStore(kvA, kvB types2.KOPair) (string, bool) {
	switch {
	case bytes.HasPrefix(kvA.Key, host.KeyClientStorePrefix) && bytes.HasSuffix(kvA.Key, []byte(host.KeyConnectionPrefix)):
		clientConnectionsA := *kvA.Value.(*types.ClientPaths)
		clientConnectionsB := *kvB.Value.(*types.ClientPaths)
		return fmt.Sprintf("ClientPaths A: %v\nClientPaths B: %v", clientConnectionsA, clientConnectionsB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyConnectionPrefix)):
		connectionA := *kvA.Value.(*types.ConnectionEnd)
		connectionB := *kvB.Value.(*types.ConnectionEnd)
		return fmt.Sprintf("ConnectionEnd A: %v\nConnectionEnd B: %v", connectionA, connectionB), true

	default:
		return "", false
	}
}
