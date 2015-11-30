package main

import (
	"errors"
	"github.com/azber/Azber-go/azber"
	"github.com/eahydra/socks"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

func main() {
	c := &azber.Proxy{}

	router := BuildUpstreamRouter(c)
	runSOCKS5Server(router)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill, os.Interrupt)
	<-sigChan
}

func BuildUpstreamRouter(conf azber.Proxy) socks.Dialer {
	var allForward []socks.Dialer
	for _, upstream := range conf.Upstreams {
		var forward socks.Dialer
		var err error
		forward = azber.NewDecorateDirect(conf.DNSCacheTimeout)
		forward, err = BuildUpstream(upstream, forward)
		if err != nil {
			azber.ErrLog.Println("failed to BuildUpstream, err:", err)
			continue
		}
		allForward = append(allForward, forward)
	}
	if len(allForward) == 0 {
		router, _ := azber.NewDecorateDirect(conf.DNSCacheTimeout)
		allForward = append(allForward, router)
	}
	return azber.NewUpstreamDialer(allForward)
}

func BuildUpstream(upstream azber.Upstream, forward socks.Dialer) (socks.Dialer, error) {
	cipherDecorator := azber.NewCipherConnDecorator(upstream.Crypto, upstream.Password)
	forward = azber.NewDecorateClient(forward, cipherDecorator)

	switch strings.ToLower(upstream.Type) {
	case "socks5":
		{
			return socks.NewSocks5Client("tcp", upstream.Address, "", "", forward)
		}
	case "shadowsocks":
		{
			return socks.NewShadowSocksClient("tcp", upstream.Address, forward)
		}
	}
	return nil, errors.New("unknown upstream type" + upstream.Type)
}

func runSOCKS5Server(forward socks.Dialer) {
	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Println("net.Listen failed, err:", err, ":7777")
		return
	}
	cipherDecorator := azber.NewCipherConnDecorator("aes-256-cfb", "1234567890")
	listener = azber.NewDecorateListener(listener, cipherDecorator)
	socks5Svr, err := azber.NewSocks5Server(forward)
	if err != nil {
		listener.Close()
		log.Println("socks.NewSocks5Server failed, err:", err)
		return
	}
	go func() {
		defer listener.Close()
		socks5Svr.Serve(listener)
	}()
}
