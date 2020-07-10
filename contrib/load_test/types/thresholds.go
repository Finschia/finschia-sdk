package types

import "time"

type Thresholds struct {
	Throughput int
	Latency    time.Duration
	TPS        int
}
