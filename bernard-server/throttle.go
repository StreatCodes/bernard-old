package main

import (
	"sync"
	"time"
)

type ThrottleList struct {
	connAttempts    map[string][]time.Time
	lock            sync.Mutex
	attemptsAllowed int
	timeoutTime     time.Duration
}

func NewThrottleList(attemptsAllowed int, timeoutTime time.Duration) ThrottleList {
	return ThrottleList{
		connAttempts:    make(map[string][]time.Time),
		lock:            sync.Mutex{},
		attemptsAllowed: attemptsAllowed,
		timeoutTime:     timeoutTime,
	}
}

//Keep track of failed connection attempts
func (tl *ThrottleList) FailedAttempt(addr string) {
	tl.lock.Lock()
	defer tl.lock.Unlock()

	attempts := tl.connAttempts[addr]
	if len(attempts) == tl.attemptsAllowed {
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

	timeout := time.Now().Add(-tl.timeoutTime)
	failedCount := 0
	for _, attempt := range attempts {
		if attempt.After(timeout) {
			failedCount += 1
		}
	}

	return failedCount >= tl.attemptsAllowed
}
