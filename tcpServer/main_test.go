package main

import (
	"testing"
	"time"
	"github.com/hotdog132/TCPProxyServer/tcpServer/mock"
	"github.com/hotdog132/TCPProxyServer/tcpServer/requestlimiter"
	"github.com/hotdog132/TCPProxyServer/tcpServer/statistics"
	"github.com/magiconair/properties/assert"
)

func TestIsValidQuery(t *testing.T) {
	testStrings := []string{
		"abcd,cdef",
		"abcd cdef",
		"abcd	cdef",
		"abcd!@ef",
		getLongString(101),
		"",
	}
	for _, testString := range testStrings {
		testInvalidString(t, testString)
	}
}
func TestPeerConnectionTimeout(t *testing.T) {
	c := &mock.MockConn{}
	jl := &requestlimiter.JobLimiter{}
	ss := &statistics.ServerStatus{}
	connectionTimeoutBySec = 2
	handleConnection(c, jl, ss)
	time.Sleep(3 * time.Second)
	assert.Equal(t, 0, ss.PeerConnectionCount, "Connection count must be 0")
}

func testInvalidString(t *testing.T, testString string) {
	isValid, _ := isValidQuery(testString)
	assert.Equal(t, false, isValid, "\""+testString+"\" must be invalid string")
}

func getLongString(length int) string {
	longString := ""
	for length > 0 {
		longString += "a"
		length--
	}
	return longString
}
