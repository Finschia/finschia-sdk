package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/store/types"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
	"github.com/line/lbm-sdk/v2/x/ibc/core/exported"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding client type.
func NewDecodeStore(kvA, kvB types.KOPair) (string, bool) {
	switch {
	case bytes.HasPrefix(kvA.Key, host.KeyClientStorePrefix) && bytes.HasSuffix(kvA.Key, []byte(host.KeyClientState)):
		clientStateA := kvA.Value.(exported.ClientState)
		clientStateB := kvB.Value.(exported.ClientState)
		return fmt.Sprintf("ClientState A: %v\nClientState B: %v", clientStateA, clientStateB), true

	case bytes.HasPrefix(kvA.Key, host.KeyClientStorePrefix) && bytes.Contains(kvA.Key, []byte(host.KeyConsensusStatePrefix)):
		consensusStateA := kvA.Value.(exported.ConsensusState)
		consensusStateB := kvB.Value.(exported.ConsensusState)
		return fmt.Sprintf("ConsensusState A: %v\nConsensusState B: %v", consensusStateA, consensusStateB), true

	default:
		return "", false
	}
}
