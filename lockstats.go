package sync

import (
	"sort"
	"sync"
	"time"
)

var (
	// Stats on lock usage by call graph.
	lockStatsMu      sync.Mutex
	lockStatsByStack map[lockStackKey]lockStats
)

type (
	lockStats struct {
		min   time.Duration
		max   time.Duration
		total time.Duration
		count lockCount
	}
	lockStackKey = [32]uintptr
	lockCount    = int64
)

func (me *lockStats) MeanTime() time.Duration {
	return me.total / time.Duration(me.count)
}

type stackLockStats struct {
	stack lockStackKey
	lockStats
}

func sortedLockTimes() (ret []stackLockStats) {
	lockStatsMu.Lock()
	for stack, stats := range lockStatsByStack {
		ret = append(ret, stackLockStats{stack, stats})
	}
	lockStatsMu.Unlock()
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].total > ret[j].total
	})
	return
}
