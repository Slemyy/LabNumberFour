package SpinLock

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type Locker struct {
	_    sync.Mutex
	lock uintptr
}

func (l *Locker) Lock() {
loop:
	if !atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
		runtime.Gosched()
		goto loop
	}
}

func (l *Locker) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}
