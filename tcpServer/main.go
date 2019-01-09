package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hotdog132/TCPProxyServer/tcpServer/requestlimiter"
	"github.com/hotdog132/TCPProxyServer/tcpServer/statistics"
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

	lStatics, err := net.Listen("tcp", ":"+"7000")
	if err != nil {
		fmt.Println(err)
		return
	}
	serverStatus := &statistics.ServerStatus{}
	serverStatus.Init(lStatics, jl)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c, jl, serverStatus)
	}
}

func handleConnection(c net.Conn, jl *requestlimiter.JobLimiter, ss *statistics.ServerStatus) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	ss.PeerConnectionCount++
	defer func() {
		c.Close()
		ss.PeerConnectionCount--
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
