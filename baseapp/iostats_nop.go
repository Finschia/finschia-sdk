//go:build !linux
// +build !linux

package baseapp

import "github.com/line/ostracon/libs/log"

func logIoStats(logger log.Logger) {
}
