package types_test

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/gogo/protobuf/proto"

	"github.com/Finschia/finschia-rdk/codec/types"
	"github.com/Finschia/finschia-rdk/testutil/testdata"
)

type errOnMarshal struct {
	testdata.Dog
}

var _ proto.Message = (*errOnMarshal)(nil)

var errAlways = fmt.Errorf("always erroring")

func (eom *errOnMarshal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return nil, errAlways
}

var eom = &errOnMarshal{}

// Ensure that returning an error doesn't suddenly allocate and waste bytes.
// See https://github.com/cosmos/cosmos-sdk/issues/8537
func TestNewAnyWithCustomTypeURLWithErrorNoAllocation(t *testing.T) {
	var ms1, ms2 runtime.MemStats

	debug.SetGCPercent(-1) // disable gc. See the comments below for reasons.
	runtime.ReadMemStats(&ms1)
	any, err := types.NewAnyWithValue(eom)
	runtime.ReadMemStats(&ms2)
	debug.SetGCPercent(100) // resume gc
	// Ensure that no fresh allocation was made.
	if diff := ms2.HeapAlloc - ms1.HeapAlloc; diff > 0 {
		// In some cases, `ms1.HeapAlloc` is larger than `ms2.HeapAlloc`.
		// It is probably because the gc worked.
		// That's why we turned off the gc for a while.
		t.Errorf("Unexpected allocation of %d bytes", diff)
	}
	if err == nil {
		t.Fatal("err wasn't returned")
	}
	if any != nil {
		t.Fatalf("Unexpectedly got a non-nil Any value: %v", any)
	}
}

var sink interface{}

func BenchmarkNewAnyWithCustomTypeURLWithErrorReturned(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		any, err := types.NewAnyWithValue(eom)
		if err == nil {
			b.Fatal("err wasn't returned")
		}
		if any != nil {
			b.Fatalf("Unexpectedly got a non-nil Any value: %v", any)
		}
		sink = any
	}
	if sink == nil {
		b.Fatal("benchmark didn't run")
	}
	sink = (interface{})(nil)
}
