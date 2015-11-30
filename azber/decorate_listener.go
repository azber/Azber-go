package azber

import (
	"net"
)

type DecorateListener struct {
	listener net.Listener
}

func (s *DecorateListener) Accept() (conn net.Conn, err error) {
	sConn, sErr := s.listener.Accept()
	if sErr {
		return nil, sErr
	}
	return sConn, sErr
}

func (s *DecorateListener) Close() error {
	return s.listener.Close()
}

func (s *DecorateListener) Addr() error {
	return s.listener.Addr()
}
