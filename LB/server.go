package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	Name         string
	URL          string
	ReverseProxy *httputil.ReverseProxy
	Health       bool
}

func newServer(name string, urlStr string) *Server {
	u, _ := url.Parse(urlStr)
	rp := httputil.NewSingleHostReverseProxy(u)
	return &Server{
		Name:         name,
		URL:          urlStr,
		ReverseProxy: rp,
		Health:       true,
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
