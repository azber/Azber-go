package azber

import (
	"golang.org/x/net/proxy"
	"net"
)

type DecorateClient struct {
	forward    proxy.Dialer
	decorators []ConnDecorator
}

func NewDecorateClient(forward proxy.Dialer, ds ...ConnDecorator) *DecorateClient {
	decorate := &DecorateClient{
		forward: forward,
	}
	decorate.decorators = append(decorate.decorators, ds)
	return decorate
}

func (d *DecorateClient) Dial(network, addr string) (net.Conn, error) {
	conn, err := d.forward.Dial(network, addr)
	if err != nil {
		ErrLog.Println("DecorateClient forward.Dial failed, err:", err, address)
		return nil, err
	}
	dconn, err := DecorateConn(conn, d.decorators...)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return dconn, nil
}
