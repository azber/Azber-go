package azber

import (
	"io"
	"net"
	"strings"
)

type CipherConn struct {
	net.Conn
	rwc io.ReadWriteCloser
}

func (c *CipherConn) Read(data []byte) (int, error) {
	return c.rwc.Read(data)
}

func (c *CipherConn) Write(b []byte) (int, error) {
	return c.rwc.Write(b)
}

func (c *CipherConn) Close() error {
	err := c.Conn.Close()
	c.rwc.Close()
	return err
}

func NewCipherConn(conn net.Conn, cryptMethod string, password []byte) (*CipherConn, error) {
	var rwc io.ReadWriteCloser
	var err error

	switch strings.ToLower(cryptMethod) {
	default:
		rwc = conn
	case "des":
		rwc, err = NewDESCFBCipher(conn, password)
	case "aes-128-cfb":
		rwc, err = NewAESCFBCipher(conn, password, 16)
	case "aes-192-cfb":
		rwc, err = NewAESCFBCipher(conn, password, 24)
	case "aes-256-cfb":
		rwc, err = NewAESCFBCipher(conn, password, 32)
	}
	if err != nil {
		return nil, err
	}

	return &CipherConn{
		Conn: conn,
		rwc:  rwc,
	}, err
}

func NewCipherConnDecorator(cryptoMethod, password string) ConnDecorator {
	return func(conn net.Conn) (net.Conn, error) {
		return NewCipherConn(conn, cryptoMethod, []byte(password))
	}
}
