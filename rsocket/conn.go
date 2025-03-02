package rsocket

import (
	"net"
	"syscall"
	"time"
)

type Conn struct {
	fd         int
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr
}

var _ net.Conn = (*Conn)(nil)

type DialOption func(conn *Conn) error

func WithLocal(localAddr string) DialOption {
	return func(conn *Conn) error {
		addr, err := net.ResolveTCPAddr("tcp", localAddr)
		if err != nil {
			return err
		}

		srcAddr := net.ParseIP(addr.IP.String())
		sa := &syscall.SockaddrInet4{
			Port: addr.Port,
		}
		copy(sa.Addr[:], srcAddr.To4())

		if err := Bind(conn.fd, sa); err != nil {
			return err
		}

		conn.localAddr = addr
		return nil
	}
}

func Dial(remoteAddr string, dialOptions ...DialOption) (*Conn, error) {
	fd, err := Socket(AF_INET, SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	conn := &Conn{fd: fd}
	for _, fn := range dialOptions {
		if err := fn(conn); err != nil {
			_ = Close(fd)
			return nil, err
		}
	}

	addr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		_ = Close(fd)
		return nil, err
	}
	conn.remoteAddr = addr

	sa := &syscall.SockaddrInet4{
		Port: addr.Port,
	}
	copy(sa.Addr[:], addr.IP.To4())

	if err := Connect(fd, sa); err != nil {
		_ = Close(fd)
		return nil, err
	}

	return conn, nil
}

func (c *Conn) FileDescriptor() int {
	return c.fd
}

func (c *Conn) Read(b []byte) (int, error) {
	return Read(c.fd, b)
}

func (c *Conn) Write(b []byte) (int, error) {
	return Write(c.fd, b)
}

func (c *Conn) Close() error {
	return Close(c.fd)
}

func (c *Conn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *Conn) SetDeadline(time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(time.Time) error {
	return nil
}
