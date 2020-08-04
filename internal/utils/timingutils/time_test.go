package timingutils

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	ms := Time("test", false, func() { time.Sleep(50 * time.Millisecond) })
	if ms < 50 || ms > 100 {
		t.Errorf("Expected value in range [50,100], got %d", ms)
	}
}
