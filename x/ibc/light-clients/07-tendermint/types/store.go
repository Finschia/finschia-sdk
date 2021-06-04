package types

import (
	"strings"

	"github.com/line/lbm-sdk/v2/codec"
	types2 "github.com/line/lbm-sdk/v2/store/types"
	sdk "github.com/line/lbm-sdk/v2/types"
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
	clienttypes "github.com/line/lbm-sdk/v2/x/ibc/core/02-client/types"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
	"github.com/line/lbm-sdk/v2/x/ibc/core/exported"
)

// KeyProcessedTime is appended to consensus state key to store the processed time
var KeyProcessedTime = []byte("/processedTime")

func GetConsensusStateUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := ConsensusState{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetConsensusStateMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*ConsensusState))
	}
}

// GetConsensusState retrieves the consensus state from the client prefixed
// store. An error is returned if the consensus state does not exist.
func GetConsensusState(store sdk.KVStore, cdc codec.BinaryMarshaler, height exported.Height) (*ConsensusState, error) {
	val := store.Get(host.ConsensusStateKey(height), GetConsensusStateUnmarshalFunc(cdc))
	consensusState, ok := val.(*ConsensusState)
	if !ok {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidConsensus,
			"invalid consensus type %T, expected %T", consensusState, &ConsensusState{},
		)
	}

	return consensusState, nil
}

// IterateProcessedTime iterates through the prefix store and applies the callback.
// If the cb returns true, then iterator will close and stop.
func IterateProcessedTime(store sdk.KVStore, cb func(key, val []byte) bool) {
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyConsensusStatePrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keySplit := strings.Split(string(iterator.Key()), "/")
		// processed time key in prefix store has format: "consensusState/<height>/processedTime"
		if len(keySplit) != 3 || keySplit[2] != "processedTime" {
			// ignore all consensus state keys
			continue
		}

		if cb(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

// ProcessedTime Store code

// ProcessedTimeKey returns the key under which the processed time will be stored in the client store.
func ProcessedTimeKey(height exported.Height) []byte {
	return append(host.ConsensusStateKey(height), KeyProcessedTime...)
}

// SetProcessedTime stores the time at which a header was processed and the corresponding consensus state was created.
// This is useful when validating whether a packet has reached the specified delay period in the tendermint client's
// verification functions
func SetProcessedTime(clientStore sdk.KVStore, height exported.Height, timeNs uint64) {
	key := ProcessedTimeKey(height)
	val := sdk.Uint64ToBigEndian(timeNs)
	clientStore.Set(key, val, types2.GetBytesMarshalFunc())
}

// GetProcessedTime gets the time (in nanoseconds) at which this chain received and processed a tendermint header.
// This is used to validate that a received packet has passed the delay period.
func GetProcessedTime(clientStore sdk.KVStore, height exported.Height) (uint64, bool) {
	key := ProcessedTimeKey(height)
	bz := clientStore.Get(key, types2.GetBytesUnmarshalFunc())
	if bz == nil {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz.([]byte)), true
}
