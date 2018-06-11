// Package sync is an extension of the stdlib "sync" package. It has extra
// functionality that helps debug the use of synchronization primitives. The
// package should be importable in place of "sync". The extra functionality
// can be enabled by calling Enable() or passing a non-empty PPROF_SYNC
// environment variable to the process.
//
// Several profiles are exposed on the default HTTP muxer (and to
// "/debug/pprof" when "net/http/pprof" is imported by the process).
// "lockHolders" lists the stack traces of goroutines that called Mutex.Lock
// that haven't subsequently been Unlocked. "lockBlockers" contains goroutines
// that are waiting to obtain locks. "/debug/lockTimes" or PrintLockTimes()
// shows the longest time a lock is held for each stack trace.
//
// Note that currently RWMutex is treated like a Mutex when the package is
// enabled.
package sync

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"
	"text/tabwriter"

	"github.com/anacrolix/missinggo"
)

var (
	// Protects initialization and enabling of the package.
	enableMu sync.Mutex
	// Whether any of this package is to be active.
	enabled = false
	// Current lock holders.
	lockHolders *pprof.Profile
	// Those blocked on acquiring a lock.
	lockBlockers *pprof.Profile
)

// Writes out the longest time a Mutex remains locked for each stack trace
// that locks a Mutex.
func PrintLockTimes(w io.Writer) {
	lockTimes := sortedLockTimes()
	tw := tabwriter.NewWriter(w, 1, 8, 1, '\t', 0)
	defer tw.Flush()
	w = tw
	for _, elem := range lockTimes {
		fmt.Fprintf(w, "%s (%s * %d [%s, %s])\n", elem.total, elem.MeanTime(), elem.count, elem.min, elem.max)
		missinggo.WriteStack(w, elem.stack[:])
	}
}

func Enable() {
	enableMu.Lock()
	defer enableMu.Unlock()
	if enabled {
		return
	}
	lockStatsByStack = make(map[lockStackKey]lockStats)
	lockHolders = pprof.NewProfile("lockHolders")
	lockBlockers = pprof.NewProfile("lockBlockers")
	http.DefaultServeMux.HandleFunc("/debug/lockTimes", func(w http.ResponseWriter, r *http.Request) {
		PrintLockTimes(w)
	})
	enabled = true
}

func init() {
	if os.Getenv("PPROF_SYNC") != "" {
		Enable()
	}
}

type (
	WaitGroup = sync.WaitGroup
	Cond      = sync.Cond
)
