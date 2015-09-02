package synctest

import "sync"

// Locker is an implementation of sync.Locker that notifies callers when
// locks and unlocks happen. otherwise, it behaves identically as a sync.Mutex
type NotifyingLocker struct {
	cond   *sync.Cond
	locked bool
	lck    *sync.Mutex
}

func NewNotifyingLocker() *NotifyingLocker {
	return &NotifyingLocker{cond: sync.NewCond(&sync.Mutex{}), locked: false, lck: &sync.Mutex{}}
}

// NotifyLock returns a channel that will close when n is locked. The channel
// never sends and will be closed immediately if n is already locked
func (n *NotifyingLocker) NotifyLock() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		n.cond.L.Lock()
		for !n.locked {
			n.cond.Wait()
		}
	}()
	return ch
}

// NotifyUnlock returns a channel that will close when n is unlocked. The channel
// never sends and will be closed immediately if n is already unlocked.
func (n *NotifyingLocker) NotifyUnlock() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		n.cond.L.Lock()
		for n.locked {
			n.cond.Wait()
		}
	}()
	return ch
}

func (n *NotifyingLocker) Lock() {
	n.lck.Lock()
	n.cond.L.Lock()
	defer n.cond.L.Unlock()
	n.locked = true
	n.cond.Broadcast()
}

func (n *NotifyingLocker) Unlock() {
	n.lck.Unlock()
	n.cond.L.Lock()
	defer n.cond.L.Unlock()
	n.locked = false
	n.cond.Broadcast()
}
