package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
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
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
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
		resp, err := http.Get(externalAPI)
		if err != nil {
			// handle the case when external api shut down
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		c.Write([]byte(string(body) + "\n"))
	}
	c.Close()
}
