package balancer

import (
	"log"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL   string
	Alive bool
	mux   sync.RWMutex
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

type RoundRobin struct {
	backends []*Backend
	current  uint64
}

func NewRoundRobin(urls []string) *RoundRobin {
	var backends []*Backend
	for _, u := range urls {
		backends = append(backends, &Backend{URL: u, Alive: true})
	}
	rr := &RoundRobin{backends: backends}
	go rr.healthCheckLoop()
	return rr
}

func (r *RoundRobin) Next() string {
	for i := 0; i < len(r.backends); i++ {
		next := atomic.AddUint64(&r.current, 1)
		idx := (next - 1) % uint64(len(r.backends))
		if r.backends[idx].IsAlive() {
			return r.backends[idx].URL
		}
	}
	return ""
}

func (r *RoundRobin) healthCheckLoop() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		for _, b := range r.backends {
			u, err := url.Parse(b.URL)
			if err != nil {
				continue
			}
            
			conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
			if err != nil {
				if b.IsAlive() {
					log.Printf("Backend %s is DOWN", b.URL)
					b.SetAlive(false)
				}
			} else {
				conn.Close()
				if !b.IsAlive() {
					log.Printf("Backend %s is UP", b.URL)
					b.SetAlive(true)
				}
			}
		}
	}
}
