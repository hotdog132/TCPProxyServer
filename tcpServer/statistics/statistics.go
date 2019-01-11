package statistics

import (
	"fmt"
	"net"
	"net/http"

	"github.com/hotdog132/TCPProxyServer/tcpServer/requestlimiter"
)

// ServerStatus Statistics of the server
type ServerStatus struct {
	PeerConnectionCount int
	jobLimiter          *requestlimiter.JobLimiter
}

// Init Initial ServerStatus
func (ss *ServerStatus) Init(l net.Listener, jl *requestlimiter.JobLimiter) {
	ss.initStaticAPIRouter(l)
	ss.jobLimiter = jl
}

func (ss *ServerStatus) initStaticAPIRouter(l net.Listener) {
	// fs := http.FileServer(http.Dir("/var/www/tpl"))
	// http.Handle("/*filepath", fs)
	http.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		<head>
		<meta http-equiv="refresh" content="2">
		</head>
		TCP server status (Refresh the page every 2 sec) <br /><br />
		Current connection count: %d <br />
		-------------------------- <br />
		Request per second: %d <br />
		-------------------------- <br />
		Processed request count: %d <br />
		-------------------------- <br />
		Remaining jobs: %d`, ss.PeerConnectionCount, ss.jobLimiter.GetRequestRatePerSec(),
			ss.jobLimiter.GetProcessedJobCount(), ss.jobLimiter.GetRemainingJobCount())
	})

	go http.Serve(l, nil)
}
