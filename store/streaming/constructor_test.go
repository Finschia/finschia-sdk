package streaming

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-rdk/codec"
	codecTypes "github.com/Finschia/finschia-rdk/codec/types"
	"github.com/Finschia/finschia-rdk/store/streaming/file"
	"github.com/Finschia/finschia-rdk/store/types"
	sdk "github.com/Finschia/finschia-rdk/types"
)

type fakeOptions struct{}

func (f *fakeOptions) Get(string) interface{} { return nil }

var (
	mockOptions       = new(fakeOptions)
	mockKeys          = []types.StoreKey{sdk.NewKVStoreKey("mockKey1"), sdk.NewKVStoreKey("mockKey2")}
	interfaceRegistry = codecTypes.NewInterfaceRegistry()
	testMarshaller    = codec.NewProtoCodec(interfaceRegistry)
)

func TestStreamingServiceConstructor(t *testing.T) {
	_, err := NewServiceConstructor("unexpectedName")
	require.NotNil(t, err)

	constructor, err := NewServiceConstructor("file")
	require.Nil(t, err)
	var expectedType ServiceConstructor
	require.IsType(t, expectedType, constructor)

	serv, err := constructor(mockOptions, mockKeys, testMarshaller)
	require.Nil(t, err)
	require.IsType(t, &file.StreamingService{}, serv)
	listeners := serv.Listeners()
	for _, key := range mockKeys {
		_, ok := listeners[key]
		require.True(t, ok)
	}
}
