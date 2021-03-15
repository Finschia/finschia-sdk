package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/types/kv"
	"github.com/line/lbm-sdk/v2/x/ibc/core/04-channel/types"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding channel type.
func NewDecodeStore(cdc codec.BinaryMarshaler, kvA, kvB kv.Pair) (string, bool) {
	switch {
	case bytes.HasPrefix(kvA.Key, []byte(host.KeyChannelEndPrefix)):
		var channelA, channelB types.Channel
		cdc.MustUnmarshalBinaryBare(kvA.Value, &channelA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &channelB)
		return fmt.Sprintf("Channel A: %v\nChannel B: %v", channelA, channelB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyNextSeqSendPrefix)):
		seqA := sdk.BigEndianToUint64(kvA.Value)
		seqB := sdk.BigEndianToUint64(kvB.Value)
		return fmt.Sprintf("NextSeqSend A: %d\nNextSeqSend B: %d", seqA, seqB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyNextSeqRecvPrefix)):
		seqA := sdk.BigEndianToUint64(kvA.Value)
		seqB := sdk.BigEndianToUint64(kvB.Value)
		return fmt.Sprintf("NextSeqRecv A: %d\nNextSeqRecv B: %d", seqA, seqB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyNextSeqAckPrefix)):
		seqA := sdk.BigEndianToUint64(kvA.Value)
		seqB := sdk.BigEndianToUint64(kvB.Value)
		return fmt.Sprintf("NextSeqAck A: %d\nNextSeqAck B: %d", seqA, seqB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyPacketCommitmentPrefix)):
		return fmt.Sprintf("CommitmentHash A: %X\nCommitmentHash B: %X", kvA.Value, kvB.Value), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyPacketAckPrefix)):
		return fmt.Sprintf("AckHash A: %X\nAckHash B: %X", kvA.Value, kvB.Value), true

	default:
		return "", false
	}
}
