package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type Server struct {
	name     string
	weight   int
	addr     string
	proxy    *httputil.ReverseProxy
	requests int
}

func NewServer(name, addr string, weight int) *Server {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}

	return &Server{
		name:   name,
		weight: weight,
		addr:   addr,
		proxy:  httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func (s *Server) IsAlive() bool {
	return true
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

type LoadBalancer struct {
	port    string
	servers []*Server
	down    []*Server
}

func NewLoadBalancer(port string) *LoadBalancer {
	return &LoadBalancer{
		port: port,
	}
}

func (ld *LoadBalancer) Read(servers []string) (int, error) {
	for _, s := range servers {
		sPart := strings.Split(s, " ")
		weight, err := strconv.Atoi(sPart[2])
		if err != nil {
			return len(ld.servers), err
		}
		server := NewServer(sPart[0], sPart[1], weight)

		ld.servers = append(ld.servers, server)
	}

	return len(ld.servers), nil
}

func (lb *LoadBalancer) getNextAvailableServer() *Server {
	return weightedSelection(lb.servers)
}

func (lb *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	server := lb.getNextAvailableServer()
	server.requests++
	log.Println("forwarding request to address ", server.addr)

	server.Serve(w, r)
	server.requests--
}

func (lb *LoadBalancer) checkHealth() {
	for i, server := range lb.servers {
		if !server.IsAlive() {
			lb.down = append(lb.down, server)
			prev := lb.servers[:i]
			next := lb.servers[i+1:]
			copy(lb.servers, prev)
			lb.servers = append(lb.servers, next...)
		}
	}
}
