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
	http.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		Current connection count: %d
		--------------------------
		Request per second:%d
		--------------------------
		Processed request count:%d
		--------------------------
		Remaining jobs:%d`, ss.PeerConnectionCount, ss.jobLimiter.GetRequestRatePerSec(),
			ss.jobLimiter.GetProcessedJobCount(), ss.jobLimiter.GetRemainingJobCount())
	})

	go http.Serve(l, nil)
}
