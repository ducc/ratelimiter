package ratelimiter

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type RateLimiter struct {
	max        int
	per        time.Duration
	lock       *sync.Mutex
	count      int
	lastPeriod time.Time
}

func New(maxPerMinute int) *RateLimiter {
	return NewWithPer(maxPerMinute, time.Minute)
}

func NewWithPer(max int, per time.Duration) *RateLimiter {
	return &RateLimiter{
		max:  max,
		per:  per,
		lock: &sync.Mutex{},
	}
}

func (r *RateLimiter) Aquire() {
	r.AquireWithCount(1)
}

func (r *RateLimiter) AquireWithCount(increment int) {
	r.lock.Lock()

	currentTime := time.Now()

	// this will be the first request
	if r.lastPeriod.IsZero() {
		r.lastPeriod = currentTime
		r.count = 1
		r.lock.Unlock()
		return
	}

	// the per is up so the timer can be reset
	if currentTime.Truncate(r.per).After(r.lastPeriod.Truncate(r.per)) {
		r.lastPeriod = currentTime
		r.count = 1
		r.lock.Unlock()
		return
	}

	// this request is the last that can be done this per, sleep until the per is up
	if r.count == r.max {
		nextPeriod := currentTime.Add(r.per).Truncate(r.per)
		waitDuration := nextPeriod.Sub(currentTime)
		logrus.Debugf("ratelimiter: sleeping %s", waitDuration.String())
		time.Sleep(waitDuration)
	}
	r.count++

	r.lock.Unlock()
}
