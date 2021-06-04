package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types3 "github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/ibc/applications/transfer/types"
	clientsim "github.com/line/lbm-sdk/v2/x/ibc/core/02-client/simulation"
	types4 "github.com/line/lbm-sdk/v2/x/ibc/core/02-client/types"
	keeper2 "github.com/line/lbm-sdk/v2/x/ibc/core/03-connection/keeper"
	connectionsim "github.com/line/lbm-sdk/v2/x/ibc/core/03-connection/simulation"
	channelsim "github.com/line/lbm-sdk/v2/x/ibc/core/04-channel/simulation"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
	types5 "github.com/line/lbm-sdk/v2/x/ibc/light-clients/09-localhost/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding ibc type.
func NewDecodeStore(cdc codec.Marshaler) types2.StoreDecoder {
	return types2.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.HasPrefix(key, host.KeyClientStorePrefix) && bytes.HasSuffix(key, []byte(host.KeyClientState)):
				return types4.GetClientStateMarshalFunc(cdc)
			case bytes.HasPrefix(key, host.KeyClientStorePrefix) && bytes.Contains(key, []byte(host.KeyConsensusStatePrefix)):
				return types4.GetConsensusStateMarshalFunc(cdc)
			case bytes.HasPrefix(key, host.KeyClientStorePrefix) && bytes.HasSuffix(key, []byte(host.KeyConnectionPrefix)):
				return keeper2.GetClientPathsMarshalFunc(cdc)
			case bytes.HasPrefix(key, []byte(host.KeyConnectionPrefix)):
				return keeper2.GetConnectionEndMarshalFunc(cdc)
			case bytes.HasPrefix(key, []byte(host.KeyChannelEndPrefix)):
				return types5.GetChannelMarshalFunc(cdc)
			case bytes.HasPrefix(key, []byte(host.KeyNextSeqSendPrefix)),
				 bytes.HasPrefix(key, []byte(host.KeyNextSeqRecvPrefix)),
			     bytes.HasPrefix(key, []byte(host.KeyNextSeqAckPrefix)),
				 bytes.HasPrefix(key, []byte(host.KeyPacketCommitmentPrefix)),
				 bytes.HasPrefix(key, []byte(host.KeyPacketAckPrefix)):
				return types3.GetBytesMarshalFunc()

			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.HasPrefix(key, host.KeyClientStorePrefix) && bytes.HasSuffix(key, []byte(host.KeyClientState)):
				return types4.GetClientStateUnmarshalFunc(cdc)
			case bytes.HasPrefix(key, host.KeyClientStorePrefix) && bytes.Contains(key, []byte(host.KeyConsensusStatePrefix)):
				return types4.GetConsensusStateUnmarshalFunc(cdc)
			case bytes.HasPrefix(key, host.KeyClientStorePrefix) && bytes.HasSuffix(key, []byte(host.KeyConnectionPrefix)):
				return keeper2.GetClientPathsUnmarshalFunc(cdc)
			case bytes.HasPrefix(key, []byte(host.KeyConnectionPrefix)):
				return keeper2.GetConnectionEndUnmarshalFunc(cdc)
			case bytes.HasPrefix(key, []byte(host.KeyChannelEndPrefix)):
				return types5.GetChannelUnmarshalFunc(cdc)
			case bytes.HasPrefix(key, []byte(host.KeyNextSeqSendPrefix)),
				bytes.HasPrefix(key, []byte(host.KeyNextSeqRecvPrefix)),
				bytes.HasPrefix(key, []byte(host.KeyNextSeqAckPrefix)),
				bytes.HasPrefix(key, []byte(host.KeyPacketCommitmentPrefix)),
				bytes.HasPrefix(key, []byte(host.KeyPacketAckPrefix)):
				return types3.GetBytesUnmarshalFunc()

			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types3.KOPair) string {
			if res, found := clientsim.NewDecodeStore(kvA, kvB); found {
				return res
			}

			if res, found := connectionsim.NewDecodeStore(kvA, kvB); found {
				return res
			}


			if res, found := channelsim.NewDecodeStore(kvA, kvB); found {
				return res
			}

			panic(fmt.Sprintf("invalid %s key prefix: %s", host.ModuleName, string(kvA.Key)))
		},
	}
}
