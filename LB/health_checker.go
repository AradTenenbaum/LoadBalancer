package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func startHealthCheck() {
	s := gocron.NewScheduler(time.Local)
	for _, host := range serverList {
		s.Every(2).Second().Do(func(s *Server) {
			healthy := s.checkHealth()
			if healthy {
				log.Printf("%s is healthy", s.Name)
			} else {
				log.Printf("%s is unhealthy", s.Name)
			}
		}, host)
	}
	s.StartAsync()
}
