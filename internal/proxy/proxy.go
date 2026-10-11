package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/takuya-q/gogate/internal/config"
)

type Gateaway struct {
	config *config.Config
}

func NewGateway(cfg *config.Config) *Gateaway {
	return &Gateaway{
		config: cfg,
	}
}

func(g *Gateaway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var targetRoute *config.Route
	for _, route := range g.config.Routes {
		if strings.HasPrefix(r.URL.Path, route.Path){
			tmp := route
			targetRoute = &tmp
			break
		}
	}
	if targetRoute == nil {
		log.Printf("Route not found for path: %s", r.URL.Path)
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}
	if len(targetRoute.Backends) == 0 {
		log.Printf("No backends available for route: %s", targetRoute.Path)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	backendURLStr := targetRoute.Backends[0]

	backendURL, err := url.Parse(backendURLStr)
	if err != nil {
		log.Printf("Failed to parse backend URL %s: %v", backendURLStr, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Proxying request: %s -> %s", r.URL.Path, backendURLStr)
	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	r.Host = backendURL.Host
	proxy.ServeHTTP(w, r)
}