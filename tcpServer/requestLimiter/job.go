package requestlimiter

import "net"

// Job ...
type Job struct {
	c     net.Conn
	host  string
	query string
}

// SetNetConnection ...
func (j *Job) SetNetConnection(c net.Conn) {
	j.c = c
}

// SetHost ...
func (j *Job) SetHost(host string) {
	j.host = host
}

// SetQuery ...
func (j *Job) SetQuery(query string) {
	j.query = query
}

// GetQuery ...
func (j *Job) GetQuery() string {
	return j.query
}

func (j *Job) getHost() string {
	return j.host
}

// WriteResult ...
func (j *Job) WriteResult(result string) {
	if j.c == nil {
		return
	}
	j.c.Write([]byte(string(result) + "\n"))
}
