package sync

import (
	"math"
	"runtime"
	"sync"
	"time"
)

type Mutex struct {
	mu      sync.Mutex
	hold    *int        // Unique value for passing to pprof.
	stack   [32]uintptr // The stack for the current holder.
	start   time.Time   // When the lock was obtained.
	entries int         // Number of entries returned from runtime.Callers.
}

func (m *Mutex) Lock() {
	if contentionOn {
		v := new(int)
		lockBlockers.Add(v, 0)
		m.mu.Lock()
		lockBlockers.Remove(v)
		m.hold = v
		lockHolders.Add(v, 0)
	} else {
		m.mu.Lock()
	}
	if lockTimesOn {
		m.entries = runtime.Callers(2, m.stack[:])
		m.start = time.Now()
	}
}

func (m *Mutex) Unlock() {
	if lockTimesOn {
		d := time.Since(m.start)
		var key [32]uintptr
		copy(key[:], m.stack[:m.entries])
		go func() {
			lockStatsMu.Lock()
			defer lockStatsMu.Unlock()
			v, ok := lockStatsByStack[key]
			if !ok {
				v.min = math.MaxInt64
			}
			if d > v.max {
				v.max = d
			}
			if d < v.min {
				v.min = d
			}
			v.total += d
			v.count++
			lockStatsByStack[key] = v
		}()
	}
	if contentionOn {
		go lockHolders.Remove(m.hold)
	}
	m.mu.Unlock()
}
