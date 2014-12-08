package sync

import (
	"container/list"
	"runtime"
	"sync"
	"time"
)

var (
	lockHolders  *list.List
	lockBlockers *list.List
	mu           sync.Mutex
)

func init() {
	lockHolders = list.New()
	lockBlockers = list.New()
}

type lockAction struct {
	time.Time
	*Mutex
	Stack string
}

func stack() string {
	var buf [0x1000]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

type Mutex struct {
	mu   sync.Mutex
	hold *list.Element
}

func (m *Mutex) addAction(l *list.List) *list.Element {
	mu.Lock()
	defer mu.Unlock()
	return l.PushBack(lockAction{
		time.Now(),
		m,
		stack(),
	})
}

func (m *Mutex) newAction() *lockAction {
	return &lockAction{
		time.Now(),
		m,
		stack(),
	}
}

func (m *Mutex) Lock() {
	a := m.newAction()
	mu.Lock()
	e := lockBlockers.PushBack(a)
	mu.Unlock()
	m.mu.Lock()
	mu.Lock()
	lockBlockers.Remove(e)
	a.Time = time.Now()
	m.hold = lockHolders.PushBack(a)
	mu.Unlock()
}

func (m *Mutex) Unlock() {
	mu.Lock()
	m.mu.Unlock()
	lockHolders.Remove(m.hold)
	mu.Unlock()
}

type WaitGroup struct {
	sync.WaitGroup
}
