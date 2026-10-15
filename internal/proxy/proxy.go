package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/takuya-q/gogate/internal/balancer"
	"github.com/takuya-q/gogate/internal/config"
)

type Gateaway struct {
	config *config.Config
	balancers map[string]*balancer.RoundRobin
}

func NewGateway(cfg *config.Config) *Gateaway {
	balancers := make(map[string]*balancer.RoundRobin)
	for _, route := range cfg.Routes {
		balancers[route.Path] = balancer.NewRoundRobin(route.Backends)
	}
	return &Gateaway{
		config: cfg,
		balancers: balancers,
	}
}

func(g *Gateaway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var targetPath string
	for _, route := range g.config.Routes {
		if strings.HasPrefix(r.URL.Path, route.Path) {
			targetPath = route.Path
			break
		}
	}
	if targetPath == "" {
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	routeBalancer, ok := g.balancers[targetPath]
	if !ok {
		http.Error(w, "Balancer not found", http.StatusInternalServerError)
		return
	}

	backendURLStr := routeBalancer.Next()
	if backendURLStr == "" {
		http.Error(w, "Bad Gateway: No backends", http.StatusBadGateway)
		return
	}

	backendURL, err := url.Parse(backendURLStr)
	if err != nil {
		slog.Error("Failed to parse backend URL", "url", backendURLStr, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Proxying request", "path", r.URL.Path, "backend", backendURLStr)

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	r.Host = backendURL.Host
	proxy.ServeHTTP(w, r)
}