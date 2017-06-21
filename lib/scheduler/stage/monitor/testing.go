package monitor

import (
	"sync"
	"time"
)

type TestMonitor struct {
	PushC, PullC               map[string]uint32
	PushErrC, PullErrC         uint32
	PushLatencyC, PullLatencyC time.Duration
	OffersC, OffersDeclinedC   uint32
	// Key is role'+'name
	ResourcesC map[string]float64
	sync.Mutex
}

func NewTestMonitor() *TestMonitor {
	return &TestMonitor{
		PushC:      make(map[string]uint32),
		PullC:      make(map[string]uint32),
		ResourcesC: make(map[string]float64),
	}
}

func (t *TestMonitor) Push(name string) {
	t.Lock()
	t.PushC[name] += 1
	t.Unlock()
}

func (t *TestMonitor) Pull(name string) {
	t.Lock()
	t.PullC[name] += 1
	t.Unlock()
}

func (t *TestMonitor) PushErr(err error) {
	t.Lock()
	t.PushErrC += 1
	t.Unlock()
}

func (t *TestMonitor) PullErr(err error) {
	t.Lock()
	t.PullErrC += 1
	t.Unlock()
}

func (t *TestMonitor) PushLatency(time time.Duration) {
	if time < 0 {
		panic("Latency must be > 0")
	}
	t.Lock()
	t.PushLatencyC += time
	t.Unlock()
}

func (t *TestMonitor) PullLatency(time time.Duration) {
	if time < 0 {
		panic("Latency must be > 0")
	}
	t.Lock()
	t.PullLatencyC += time
	t.Unlock()
}

func (t *TestMonitor) OffersReceived(count uint32) {
	t.Lock()
	t.OffersC += count
	t.Unlock()
}

func (t *TestMonitor) OffersDeclined(count uint32) {
	t.Lock()
	t.OffersDeclinedC += count
	t.Unlock()
}

func (t *TestMonitor) ResourceOffered(name, role string, value float64) {
	t.Lock()
	t.ResourcesC[name+":"+role] += value
	t.Unlock()
}
