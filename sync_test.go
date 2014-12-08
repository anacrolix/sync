package sync

import (
	"sync"
	"testing"
)

func TestLog(t *testing.T) {
	var mu Mutex
	mu.Lock()
	mu.Unlock()
}

func TestUnlockUnlocked(t *testing.T) {
	var mu sync.Mutex
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("should have panicked")
		}
	}()
	mu.Unlock()
}
