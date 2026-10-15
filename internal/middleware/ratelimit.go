package middleware

import(
	"log"
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(2, 5)
		visitors[ip] = limiter
	}
	return limiter
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		limiter := getVisitor(ip)
		if !limiter.Allow() {
			log.Printf("Blocked request from %s: Too Many Requests", ip)
			http.Error(w, "Too Many Requests: Slow down!", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}