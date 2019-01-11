package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hotdog132/TCPProxyServer/tcpServer/requestlimiter"
	"github.com/hotdog132/TCPProxyServer/tcpServer/statistics"
)

var (
	tcpPort                = "8000"
	apiPort                = "7000"
	externalAPI            = "http://localhost:8888/api"
	connectionTimeoutBySec = 120
	limitQPS               = 30
)

func main() {
	log.Println("TCP server start listening port:" + tcpPort)
	l, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Fatal(err)
		return
	}

	jl := &requestlimiter.JobLimiter{}
	jl.Init(limitQPS)
	jl.SetExternalAPI(externalAPI)

	log.Println("Statistics API start listening port:" + apiPort)
	lStatics, err := net.Listen("tcp", ":"+apiPort)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		l.Close()
		lStatics.Close()
	}()

	serverStatus := &statistics.ServerStatus{}
	serverStatus.Init(lStatics, jl)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(c, jl, serverStatus)
	}
}

func handleConnection(c net.Conn, jl *requestlimiter.JobLimiter, ss *statistics.ServerStatus) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	ss.PeerConnectionCount++
	defer func() {
		log.Println(c.RemoteAddr().String() + " disconnect")
		c.Close()
		ss.PeerConnectionCount--
	}()

	showInstructions(c)

	for {
		c.SetDeadline(time.Now().Add(time.Duration(connectionTimeoutBySec) * time.Second))
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Println(c.RemoteAddr().String() + " " + err.Error())
			return
		}

		query := strings.TrimSpace(string(netData))
		if isValid, errMessage := isValidQuery(query); !isValid {
			c.Write([]byte("Invalid query:" + errMessage.Error() + "\n"))
			continue
		}

		if query == "quit" {
			return
		}

		job := &requestlimiter.Job{}
		job.SetNetConnection(c)
		job.SetHost(c.RemoteAddr().String())
		job.SetQuery(query)

		jl.EnqueueJob(job)
	}
}

func showInstructions(c net.Conn) {
	c.Write([]byte("*** Type any string queries to get informations from external API.\n"))
	c.Write([]byte("*** Queries couldn't contains any special charactors, space or tab, and length of 100 limited.\n"))
	c.Write([]byte("*** Disconnect if the server doesn't received any message within " +
		strconv.Itoa(connectionTimeoutBySec) + " seconds.\n"))
}

func isValidQuery(query string) (bool, error) {
	if len(query) == 0 || len(query) > 100 {
		return false, errors.New("Query length must be bigger than 1 and smaller than 100")
	}

	if strings.ContainsAny(query, " ") || strings.ContainsAny(query, "\t") {
		return false, errors.New("Query contains space or tab")
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	filteredQuery := reg.ReplaceAllString(query, "")

	if len(query) != len(filteredQuery) {
		return false, errors.New("Query contains special charactors")
	}
	return true, nil
}
