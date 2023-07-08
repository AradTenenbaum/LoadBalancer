package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	serverList = []*Server{
		newWeightedServer("server-0", "http://localhost:5000", 2),
		newWeightedServer("server-1", "http://localhost:5001", 3),
		newWeightedServer("server-2", "http://localhost:5002", 2),
		newWeightedServer("server-3", "http://localhost:5003", 1),
		// newWeightedServer("server-4", "http://localhost:5004", 5),
	}
	routeServer     = newRouteServer("server-4", "http://localhost:5004")
	startWeight     = serverList[0].Weight
	lastServedIndex = 0
)

func main() {

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", forwardRequest)
	http.HandleFunc("/route/", forwardRouteRequest)
	log.Print("Starting Load Balancer on 8000...")
	go startHealthCheck()
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func forwardRequest(res http.ResponseWriter, req *http.Request) {
	server, err := getHealthyServer()
	if err != nil {
		http.Error(res, "Couldn't process request: "+err.Error(), http.StatusServiceUnavailable)
		return
	} else {
		log.Printf("Routing the request to the URL: %s", server.URL)
		server.ReverseProxy.ServeHTTP(res, req)
	}
}

func forwardRouteRequest(res http.ResponseWriter, req *http.Request) {
	log.Printf("Routing the request to the URL: %s", routeServer.URL)
	routeServer.ReverseProxy.ServeHTTP(res, req)
}

func getHealthyServer() (*Server, error) {
	for i := 0; i < len(serverList); i++ {
		server := getRandomServer()
		if server.Health {
			return server, nil
		}
	}
	return nil, fmt.Errorf("No healthy server")
}

func getServerRoundRobin() *Server {
	nextIndex := (lastServedIndex + 1) % len(serverList)
	server := serverList[nextIndex]
	lastServedIndex = nextIndex
	return server
}

func getServerWeightedRoundRobin() *Server {
	if startWeight == 0 {
		lastServedIndex = (lastServedIndex + 1) % len(serverList)
		startWeight = serverList[lastServedIndex].Weight
	}
	server := serverList[lastServedIndex]
	startWeight--
	return server
}

func getRandomServer() *Server {
	randomIndex := rand.Intn(len(serverList))
	return serverList[randomIndex]
}
