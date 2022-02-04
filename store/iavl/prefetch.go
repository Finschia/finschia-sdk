package iavl

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/line/lbm-sdk/telemetry"
)

var (
	usePrefetch       = 0
	prefetchCommiters int64
	prefetchJobs      chan func()
	prefetchLocks     chan bool
	prefetchToken     chan bool

	prefetchJobsSize = 100000 // should be pending queue * 4
)

func StartPrefetch() {
	usePrefetch = 1
}

func StopPrefetch() {
	usePrefetch = -1
}

func PausePrefetcher() {
	if atomic.AddInt64(&prefetchCommiters, 1) == 1 {
		prefetchToken <- true
	}
}

func ResumePrefetcher() {
	if atomic.AddInt64(&prefetchCommiters, -1) == 0 {
		<-prefetchToken
	}
}

func prefetcher() {
	workers := runtime.NumCPU() / 4
	if workers < 4 {
		workers = 4
	}

	prefetchJobs = make(chan func(), prefetchJobsSize)
	prefetchLocks = make(chan bool, workers)
	prefetchToken = make(chan bool, 1)

	go func() {
		for {
			f := <-prefetchJobs
			if len(prefetchToken) != 0 {
				prefetchToken <- true
				<-prefetchToken
			}
			prefetchLocks <- true
			go func(f func()) {
				f()
				<-prefetchLocks
			}(f)
		}
	}()
}

// Implements type.KVStore.
func (st *Store) Prefetch(key []byte, forSet bool) (hits, misses int, value []byte) {
	if usePrefetch != 1 {
		return
	}
	defer telemetry.MeasureSince(time.Now(), "store", "iavl", "load")
	select {
	case prefetchJobs <- func() {
		defer func() {
			// ignore panic
			recover()
		}()
		st.tree.Prefetch(key, forSet)
	}:
		// good
	default:
		// drop this request
	}
	return
}
