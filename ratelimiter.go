package ratelimiter

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type RateLimiter struct {
	maxPerMinute int
	lock         *sync.Mutex
	count        int
	lastPeriod   time.Time
}

func New(maxPerMinute int) *RateLimiter {
	return &RateLimiter{
		maxPerMinute: maxPerMinute,
		lock:         &sync.Mutex{},
	}
}

func (r *RateLimiter) Aquire() {
	r.lock.Lock()

	currentTime := time.Now()

	// this will be the first request
	if r.lastPeriod.IsZero() {
		r.lastPeriod = currentTime
		r.count = 1
		r.lock.Unlock()
		return
	}

	// the minute is up so the timer can be reset
	if currentTime.Truncate(time.Minute).After(r.lastPeriod.Truncate(time.Minute)) {
		r.lastPeriod = currentTime
		r.count = 1
		r.lock.Unlock()
		return
	}

	r.count++

	// this request is the last that can be done this minute, sleep until the minute is up
	if r.count == r.maxPerMinute {
		nextPeriod := currentTime.Add(time.Minute).Truncate(time.Minute)
		waitDuration := nextPeriod.Sub(currentTime)
		logrus.Debugf("ratelimiter: sleeping %s", waitDuration.String())
		time.Sleep(waitDuration)
	}

	r.lock.Unlock()
}
