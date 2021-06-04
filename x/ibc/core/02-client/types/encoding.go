package types

import (
	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/x/ibc/core/exported"
)

func GetClientStateUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		var val exported.ClientState
		cdc.UnmarshalInterface(value, &val)
		return val
	}
}

func GetClientStateMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		value, err := cdc.MarshalInterface(obj.(exported.ClientState))
		if err != nil {
			panic(err)
		}
		return value
	}
}

func GetConsensusStateUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		var val exported.ConsensusState
		cdc.UnmarshalInterface(value, &val)
		return val
	}
}

func GetConsensusStateMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		value, err := cdc.MarshalInterface(obj.(exported.ConsensusState))
		if err != nil {
			panic(err)
		}
		return value
	}
}

func GetHeaderUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		var val exported.Header
		cdc.UnmarshalInterface(value, &val)
		return val
	}
}

func GetHeaderMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		value, err := cdc.MarshalInterface(obj.(exported.Header))
		if err != nil {
			panic(err)
		}
		return value
	}
}
