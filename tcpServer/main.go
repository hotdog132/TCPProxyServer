package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hotdog132/TCPProxyServer/tcpServer/requestlimiter"
)

const (
	port        = "8000"
	externalAPI = "http://localhost:8888/api"
)

func main() {
	// arguments := os.Args
	// if len(arguments) == 1 {
	// 	fmt.Println("Please provide a port number!")
	// 	return
	// }

	// PORT := ":" + arguments[1]
	log.Println("Start listening port:" + port)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	jl := &requestlimiter.JobLimiter{}
	jl.Init(1)
	jl.SetExternalAPI(externalAPI)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c, jl)
	}
}

func handleConnection(c net.Conn, jl *requestlimiter.JobLimiter) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())

	defer func() {
		c.Close()
	}()

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		query := strings.TrimSpace(string(netData))
		if query == "quit" {
			break
		}

		job := &requestlimiter.Job{}
		job.SetNetConnection(c)
		job.SetHost(c.RemoteAddr().String())
		job.SetQuery(query)

		jl.EnqueueJob(job)

	}
}
