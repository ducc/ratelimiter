package ratelimiter

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var DefaultWaitCallback = func(waitDuration time.Duration) {
	logrus.Debugf("ratelimiter: sleeping %s", waitDuration.String())
	time.Sleep(waitDuration)
}

type RateLimiter struct {
	max        int
	per        time.Duration
	lock       *sync.Mutex
	count      int
	lastPeriod time.Time
	callback   func(waitDuration time.Duration)
}

func New(maxPerMinute int) *RateLimiter {
	return NewWithPer(maxPerMinute, time.Minute)
}

func NewWithPer(max int, per time.Duration) *RateLimiter {
	return NewWithPerAndCallback(max, per, DefaultWaitCallback)
}

func NewWithPerAndCallback(max int, per time.Duration, callback func(waitDuration time.Duration)) *RateLimiter {
	return &RateLimiter{
		max:      max,
		per:      per,
		lock:     &sync.Mutex{},
		callback: callback,
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
		if r.callback != nil {
			r.callback(waitDuration)
		}
	}
	r.count++

	r.lock.Unlock()
}
