package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/types"
)

func TestNewStoreKVPairWriteListener(t *testing.T) {
	testWriter := new(bytes.Buffer)
	interfaceRegistry := types.NewInterfaceRegistry()
	testMarshaller := codec.NewProtoCodec(interfaceRegistry)

	wl := NewStoreKVPairWriteListener(testWriter, testMarshaller)

	require.IsType(t, &StoreKVPairWriteListener{}, wl)
	require.Equal(t, testWriter, wl.writer)
	require.Equal(t, testMarshaller, wl.marshaller)
}

func TestOnWrite(t *testing.T) {
	testWriter := new(bytes.Buffer)
	interfaceRegistry := types.NewInterfaceRegistry()
	testMarshaller := codec.NewProtoCodec(interfaceRegistry)

	wl := NewStoreKVPairWriteListener(testWriter, testMarshaller)

	testStoreKey := NewKVStoreKey("test_key")
	testKey := []byte("testing123")
	testValue := []byte("testing321")

	// test set
	err := wl.OnWrite(testStoreKey, testKey, testValue, false)
	require.Nil(t, err)

	outputBytes := testWriter.Bytes()
	outputKVPair := new(StoreKVPair)
	expectedOutputKVPair := &StoreKVPair{
		Key:      testKey,
		Value:    testValue,
		StoreKey: testStoreKey.Name(),
		Delete:   false,
	}
	err = testMarshaller.UnmarshalLengthPrefixed(outputBytes, outputKVPair)
	require.NoError(t, err)
	require.EqualValues(t, expectedOutputKVPair, outputKVPair)
	testWriter.Reset()

	// test delete
	err = wl.OnWrite(testStoreKey, testKey, testValue, true)
	require.Nil(t, err)

	outputBytes = testWriter.Bytes()
	outputKVPair = new(StoreKVPair)
	expectedOutputKVPair = &StoreKVPair{
		Key:      testKey,
		Value:    testValue,
		StoreKey: testStoreKey.Name(),
		Delete:   true,
	}
	err = testMarshaller.UnmarshalLengthPrefixed(outputBytes, outputKVPair)
	require.NoError(t, err)
	require.EqualValues(t, expectedOutputKVPair, outputKVPair)
}
