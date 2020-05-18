package gosocket

import (
	"time"
	"sync"
)

var globalTimerPool sync.Pool

// acquireTimer returns a timer with the specified timeout.
func acquireTimer(timeout time.Duration) *time.Timer {
	v := globalTimerPool.Get()
	if v == nil {
		return time.NewTimer(timeout)
	}
	t := v.(*time.Timer)
	t.Reset(timeout)
	return t
}

// releaseTimer puts the given t to the timer pool.
// It will try to stop the timer before pooling, so we
// don't need to call t.Stop() before calling it.
func releaseTimer(t *time.Timer) {
	if !t.Stop() {
		// Drain the channel if it has not been read yet.
		select {
		case <-t.C:
		default:
		}
	}
	globalTimerPool.Put(t)
}

