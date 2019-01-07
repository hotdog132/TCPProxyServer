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
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "quit" {
			break
		}

		// result := strconv.Itoa(rand.Intn(100)) + "\n"

		// external api call
		// resp, err := http.Get(externalAPI)
		// if err != nil {
		// 	// handle the case when external api shut down
		// }
		// defer resp.Body.Close()
		// body, err := ioutil.ReadAll(resp.Body)

		// c.Write([]byte(string(body) + "\n"))

		job := &requestlimiter.Job{}
		job.SetNetConnection(c)
		job.SetHost(c.RemoteAddr().String())
		jl.EnqueueJob(job)

	}
	c.Close()
}
