package timingutils

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	ms := Time("test", false, func() { time.Sleep(50 * time.Millisecond) })
	if ms < 50 || ms > 250 {
		t.Errorf("Expected value in range [50,250], got %d", ms)
	}
}
