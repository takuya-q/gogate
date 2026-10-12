package balancer

import (
	"sync/atomic"
)

type RoundRobin struct {
	backends []string
	current uint64
}

func NewRoundRobin(backends []string) *RoundRobin {
	return &RoundRobin{
		backends: backends,
	}
}

func (r *RoundRobin) Next() string{
	if len(r.backends) == 0 {
		return ""
	}
	next := atomic.AddUint64(&r.current, 1)
	idx := (next - 1) % uint64(len(r.backends))
	return r.backends[idx]
}