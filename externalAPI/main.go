package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	// connectionCount := 1

	l, err := net.Listen("tcp", ":8888")

	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	defer l.Close()

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(3 * time.Second)
		fmt.Fprintf(w, "Welcome to my website!")
	})

	// l = netutil.LimitListener(l, connectionCount)

	log.Fatal(http.Serve(l, nil))
}
