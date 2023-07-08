package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Server struct {
	Name         string
	URL          string
	ReverseProxy *httputil.ReverseProxy
	Health       bool
	Weight       int
}

func newServer(name string, urlStr string) *Server {
	return newWeightedServer(name, urlStr, 1)
}

func newWeightedServer(name string, urlStr string, weight int) *Server {
	u, _ := url.Parse(urlStr)
	rp := httputil.NewSingleHostReverseProxy(u)
	return &Server{
		Name:         name,
		URL:          urlStr,
		ReverseProxy: rp,
		Health:       true,
		Weight:       weight,
	}
}

func newRouteServer(name string, urlStr string) *Server {
	u, _ := url.Parse(urlStr)
	return &Server{
		Name: name,
		URL:  urlStr,
		ReverseProxy: &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				// Remove the route from the original request path
				if strings.HasPrefix(req.URL.Path, "/route") {
					req.URL.Path = strings.TrimPrefix(req.URL.Path, "/route")
				}

				// Set the proxy target to the original server
				req.URL.Scheme = u.Scheme
				req.URL.Host = u.Host
			},
		},
		Health: true,
		Weight: 1,
	}
}

func (s *Server) checkHealth() bool {
	resp, err := http.Head(s.URL)
	if err != nil {
		s.Health = false
		return s.Health
	}
	if resp.StatusCode != http.StatusOK {
		s.Health = false
		return s.Health
	}
	s.Health = true
	return s.Health
}
