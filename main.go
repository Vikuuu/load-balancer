package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	lb := NewLoadBalancer("8000")

	servers, err := readYAML(pwd)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	n, err := lb.Read(servers)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	fmt.Fprintf(os.Stdout, "%d number of servers found\n", n)
	fmt.Fprintf(os.Stdout, "%+v\n", lb.servers)
	fmt.Fprintf(os.Stdout, "%v\n", servers)

	handleRedirect := func(w http.ResponseWriter, r *http.Request) { lb.serveProxy(w, r) }
	http.HandleFunc("/", handleRedirect)
	fmt.Println("serving request at localhost:", lb.port)

	// bankai
	go backgroundProcess(lb)

	// start listening for request
	http.ListenAndServe(":"+lb.port, nil)
}

func backgroundProcess(lb *LoadBalancer) {
	time.Sleep(5 * time.Second)
	lb.checkHealth()
}
