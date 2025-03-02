package rsocket

import (
	"net"
	"syscall"
)

type Listener struct {
	addr *net.TCPAddr
	fd   int
}

var _ net.Listener = (*Listener)(nil)

type ListenOption func(ln *Listener) error

func NewListener(localAddr string, backlog int, listenOptions ...ListenOption) (*Listener, error) {
	fd, err := Socket(AF_INET, SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	listener := &Listener{fd: fd}
	for _, fn := range listenOptions {
		if err := fn(listener); err != nil {
			_ = Close(fd)
			return nil, err
		}
	}

	addr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		_ = Close(fd)
		return nil, err
	}
	listener.addr = addr

	sa := &syscall.SockaddrInet4{
		Port: addr.Port,
	}
	copy(sa.Addr[:], addr.IP.To4())

	if err := Bind(fd, sa); err != nil {
		_ = Close(fd)
		return nil, err
	}

	if err := Listen(fd, backlog); err != nil {
		_ = Close(fd)
		return nil, err
	}

	return listener, nil
}

func (l *Listener) FileDescriptor() int {
	return l.fd
}

func (l *Listener) Accept() (net.Conn, error) {
	fd, addr, err := Accept(l.fd)
	if err != nil {
		return nil, err
	}

	socketAddr := addr.(*syscall.SockaddrInet4)
	ip := net.IPv4(socketAddr.Addr[0], socketAddr.Addr[1], socketAddr.Addr[2], socketAddr.Addr[3])
	port := socketAddr.Port

	remoteAddr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}

	return &Conn{
		fd:         fd,
		localAddr:  l.addr,
		remoteAddr: remoteAddr,
	}, nil
}

func (l *Listener) Close() error {
	return Close(l.fd)
}

func (l *Listener) Addr() net.Addr {
	return l.addr
}
