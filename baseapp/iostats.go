// +build linux

package baseapp

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/iavl"
	"github.com/line/ostracon/libs/log"
)

// log iavl, db & disk stats if IAVL_DEV env. var. is set
var (
	logIoDevNop                                                  string  = " nop"
	logIoDev                                                     *string = nil
	tStats                                                       time.Time
	iavlSetCount, iavlSetBytes, iavlDelCount                     uint64
	diskReadCount, diskReadBytes, diskWriteCount, diskWriteBytes uint64
)

// read disk io stats from /proc/diskstats in linux
func getDeviceIoStats(device string) (readCount, readBytes, writeCount, writeBytes uint64, err error) {
	inf, err := os.Open("/proc/diskstats")
	if err != nil {
		return
	}
	defer inf.Close()

	in := bufio.NewReader(inf)
	r := regexp.MustCompile(`\s+`)
	for {
		line, err2 := in.ReadString('\n')
		if err2 != nil {
			if err2 == io.EOF {
				break
			}
			err = err2
			return
		}
		parts := r.Split(strings.TrimSpace(line), -1)
		if len(parts) < 10 {
			continue
		}
		if len(device) > 0 && device != parts[2] {
			continue
		} else if len(device) == 0 && parts[1] != "0" {
			continue
		}

		rCount, e1 := strconv.ParseUint(parts[3], 10, 64)
		rSectors, e2 := strconv.ParseUint(parts[5], 10, 64)
		wCount, e3 := strconv.ParseUint(parts[7], 10, 64)
		wSectors, e4 := strconv.ParseUint(parts[9], 10, 64)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			continue
		}

		readCount += rCount
		readBytes += rSectors * 512
		writeCount += wCount
		writeBytes += wSectors * 512
	}
	return readCount, readBytes, writeCount, writeBytes, err
}

func logIoStats(logger log.Logger) {
	if logIoDev == &logIoDevNop {
		return
	} else if logIoDev == nil {
		// not set yet
		if dev := os.Getenv("IAVL_DEV"); len(dev) > 0 {
			logIoDev = &dev
		} else {
			logIoDev = &logIoDevNop
			return
		}
	}

	now := time.Now()
	dt := now.Sub(tStats)
	tStats = now

	tCount, newSetBytes, newDelCount := iavl.GetDBStats()
	newReadCount, newReadBytes, newWriteCount, newWriteBytes, _ := getDeviceIoStats(*logIoDev)
	diffSetCount, diffSetBytes, diffDelCount := newSetCount-iavlSetCount, newSetBytes-iavlSetBytes, newDelCount-iavlDelCount
	_, diffReadBytes, _, diffWriteBytes := newReadCount-diskReadCount, newReadBytes-diskReadBytes, newWriteCount-diskWriteCount, newWriteBytes-diskWriteBytes
	iavlSetCount, iavlSetBytes, iavlDelCount = newSetCount, newSetBytes, newDelCount
	diskReadCount, diskReadBytes, diskWriteCount, diskWriteBytes = newReadCount, newReadBytes, newWriteCount, newWriteBytes

	// read & write bytes per second
	rps, wps := int64(0), int64(0)
	if dt.Milliseconds() > 0 && dt.Seconds() < 1000 {
		rps = int64(diffReadBytes) * 1000 / dt.Milliseconds()
		wps = int64(diffWriteBytes) * 1000 / dt.Milliseconds()
	}

	// write amplification factors
	wampf := float64(0)
	if diffSetBytes > 0 {
		wampf = float64(diffWriteBytes*1000/diffSetBytes) / 1000.0
	}

	logger.Info("io_stats", "set_count", diffSetCount,
		"set_bytes", diffSetBytes, "del_count", diffDelCount,
		"read-bytes", diffReadBytes, "rps", rps,
		"write_bytes", diffWriteBytes, "wps", wps,
		"write_amp_f", wampf, "ellapsed", float64(dt.Milliseconds())/1000.0)
}
