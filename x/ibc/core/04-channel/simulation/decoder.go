package simulation

import (
	"bytes"
	"fmt"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/ibc/core/04-channel/types"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding channel type.
func NewDecodeStore(kvA, kvB types2.KOPair) (string, bool) {
	switch {
	case bytes.HasPrefix(kvA.Key, []byte(host.KeyChannelEndPrefix)):
		channelA := *kvA.Value.(*types.Channel)
		channelB := *kvB.Value.(*types.Channel)
		return fmt.Sprintf("Channel A: %v\nChannel B: %v", channelA, channelB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyNextSeqSendPrefix)):
		seqA := sdk.BigEndianToUint64(kvA.Value.([]byte))
		seqB := sdk.BigEndianToUint64(kvB.Value.([]byte))
		return fmt.Sprintf("NextSeqSend A: %d\nNextSeqSend B: %d", seqA, seqB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyNextSeqRecvPrefix)):
		seqA := sdk.BigEndianToUint64(kvA.Value.([]byte))
		seqB := sdk.BigEndianToUint64(kvB.Value.([]byte))
		return fmt.Sprintf("NextSeqRecv A: %d\nNextSeqRecv B: %d", seqA, seqB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyNextSeqAckPrefix)):
		seqA := sdk.BigEndianToUint64(kvA.Value.([]byte))
		seqB := sdk.BigEndianToUint64(kvB.Value.([]byte))
		return fmt.Sprintf("NextSeqAck A: %d\nNextSeqAck B: %d", seqA, seqB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyPacketCommitmentPrefix)):
		return fmt.Sprintf("CommitmentHash A: %X\nCommitmentHash B: %X", kvA.Value.([]byte), kvB.Value.([]byte)), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyPacketAckPrefix)):
		return fmt.Sprintf("AckHash A: %X\nAckHash B: %X", kvA.Value.([]byte), kvB.Value.([]byte)), true

	default:
		return "", false
	}
}
