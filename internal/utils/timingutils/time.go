package timingutils

import (
	"fmt"
	"time"
)

func Time(name string, showOutput bool, op func()) (ms int64) {
	start := time.Now()
	op()
	ms = time.Since(start).Milliseconds()
	if showOutput {
		fmt.Printf("%s: %d ms\n", name, ms)
	}
	return
}
