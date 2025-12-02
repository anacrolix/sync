package sync

import (
	"sync"
	"testing"
)

func BenchmarkRwMutexContention(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var wg sync.WaitGroup
		var lock RWMutex
		for range 10 {
			wg.Go(func() {
				for range 100 {
					lock.Lock()
					lock.Unlock()
				}
			})
		}
		for range 10 {
			wg.Go(func() {
				for range 100 {
					lock.RLock()
					lock.RUnlock()
				}
			})
		}
		wg.Wait()
	}
}
