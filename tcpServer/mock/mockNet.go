package mock

import (
	"net"
	"time"
)

type MockConn struct{}

func (m *MockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}
func (m *MockConn) Write(b []byte) (n int, err error) {
	return 0, nil
}
func (m *MockConn) Close() error {
	return nil
}
func (m *MockConn) LocalAddr() net.Addr {
	return &MockAddr{}
}
func (m *MockConn) RemoteAddr() net.Addr {
	return &MockAddr{}
}
func (m *MockConn) SetDeadline(t time.Time) error {
	return nil
}
func (m *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (m *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type MockAddr struct{}

func (m *MockAddr) Network() string {
	return "mockNetwork"
}

func (m *MockAddr) String() string {
	return "mockString"
}
