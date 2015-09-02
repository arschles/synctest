package synctest

import (
	"testing"
	"time"
)

const (
	timeout = 1 * time.Second
)

func TestNotifyLock(t *testing.T) {
	nl := NewNotifyingLocker()
	ch := nl.NotifyLock()
	nl.Lock()
	select {
	case <-ch:
	case <-time.After(timeout):
		t.Fatalf("didn't receive on channel after %s", timeout)
	}
}

func TestNotifyLockAlreadyLocked(t *testing.T) {
	nl := NewNotifyingLocker()
	nl.Lock()
	ch := nl.NotifyLock()
	select {
	case <-ch:
	case <-time.After(timeout):
		t.Fatalf("didn't receive on channel after %s", timeout)
	}
}
