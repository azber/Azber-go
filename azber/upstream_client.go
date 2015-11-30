package azber

import (
	"golang.org/x/net/proxy"
	"net"
	"sync"
)

type UpstreamDialer struct {
	forwardDialers []proxy.Dialer
	nextRouter     int
	lock           sync.Mutex
}

func NewUpstreamDialer(forwardDialers []proxy.Dialer) (*UpstreamDialer, error) {
	return &UpstreamDialer{
		forwardDialers: forwardDialers,
	}, nil
}

func (u *UpstreamDialer) getNextDialer() proxy.Dialer {
	u.lock.Lock()
	defer u.lock.Unlock()
	index := u.nextRouter
	u.nextRouter++
	if u.nextRouter >= len(u.forwardDialers) {
		u.nextRouter = 0
	}
	if index < len(u.forwardDialers) {
		return u.forwardDialers[index]
	}
	panic("unreached")
}

func (u *UpstreamDialer) Dial(network, address string) (net.Conn, error) {
	router := u.getNextDialer()
	conn, err := router.Dial(network, address)
	if err != nil {
		ErrLog.Println("UpstreamDialer router.Dial failed, err:", err, network, address)
		return nil, err
	}
	return conn
}
