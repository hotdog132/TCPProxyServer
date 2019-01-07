package requestlimiter

import "net"

// Job ...
type Job struct {
	c    net.Conn
	host string
}

// SetNetConnection ...
func (j *Job) SetNetConnection(c net.Conn) {
	j.c = c
}

// SetHost ...
func (j *Job) SetHost(host string) {
	j.host = host
}

func (j *Job) getHost() string {
	return j.host
}

// WriteResult ...
func (j *Job) WriteResult(result string) {
	j.c.Write([]byte(string(result) + "\n"))
}
