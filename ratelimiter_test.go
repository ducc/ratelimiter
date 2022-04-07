package ratelimiter_test

import (
	"github.com/ducc/ratelimiter"
	"testing"
	"time"
)

func TestAquire(t *testing.T) {
	rl := ratelimiter.NewWithPer(5, time.Second)

	// Iterate 6 times.
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			assertDelay(t, aquireDelay(rl), j == 5)
		}

		// Ensure we're past a second before waiting to do the next itr.
		time.Sleep((time.Second / 2) + time.Second)
	}
}

func assertDelay(t *testing.T, duration time.Duration, hit bool) {
	var gtlt = "<"
	var assertCase = duration < 3*time.Millisecond

	if hit {
		gtlt = ">"
		assertCase = duration > 3*time.Millisecond
	}

	if !assertCase {
		t.Errorf("Assert Delay for Aquire failed: Expected %s%s got %s", gtlt, 3*time.Millisecond, duration)
	}
}

// Time the aquire call
func aquireDelay(rl *ratelimiter.RateLimiter) time.Duration {
	var now = time.Now()
	rl.Aquire()
	return time.Now().Sub(now)
}
