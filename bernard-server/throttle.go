package main

import (
	"sync"
	"time"
)

//TODO move to config
const ATTEMPTS_ALLOWED = 5
const TIMEOUT_TIME = time.Minute

type ThrottleList struct {
	connAttempts map[string][]time.Time
	lock         sync.Mutex
}

//Keep track of failed connection attempts
func (tl *ThrottleList) FailedAttempt(addr string) {
	tl.lock.Lock()
	defer tl.lock.Unlock()

	attempts := tl.connAttempts[addr]
	if len(attempts) == ATTEMPTS_ALLOWED {
		attempts = append(attempts[1:], time.Now())
	} else {
		attempts = append(attempts, time.Now())
	}

	tl.connAttempts[addr] = attempts
}

//Check if the remote address has been throttled
func (tl *ThrottleList) IsThrottled(addr string) bool {
	tl.lock.Lock()
	defer tl.lock.Unlock()

	attempts := tl.connAttempts[addr]

	timeout := time.Now().Add(-TIMEOUT_TIME)
	failedCount := 0
	for _, attempt := range attempts {
		if attempt.After(timeout) {
			failedCount += 1
		}
	}

	return failedCount >= ATTEMPTS_ALLOWED
}
