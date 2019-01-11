package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	port = "8888"
)

func main() {
	// connectionCount := 1

	l, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	log.Println("Start listening port:" + port)

	defer l.Close()

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if sliceQuery, ok := r.URL.Query()["q"]; ok {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(w, "External api response query: "+sliceQuery[0])
		}
	})

	log.Fatal(http.Serve(l, nil))
}
