package main

import (
	"errors"
	"github.com/azber/Azber-go/azber"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

func main() {

	upstream := azber.Upstream{
		Type:     "shadowsocks",
		Crypto:   "aes-256-cfb",
		Password: "ijdIM@j83!dj.Udi",
		Address:  "133.130.99.18:34781",
	}

	pac := azber.PAC{
		Address:  "127.0.0.1:50000",
		Proxy:    "127.0.0.1:40000",
		SOCKS5:   "127.0.0.1:8000",
		Upstream: upstream,
	}

	c := azber.Config{
		PAC: pac,
	}

	router, _ := BuildUpstreamRouter(c)
	runSOCKS5Server(router)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill, os.Interrupt)
	<-sigChan
}

func BuildUpstreamRouter(conf azber.Proxy) (proxy.Dialer, error) {
	var allForward []proxy.Dialer
	for _, upstream := range conf.Upstreams {
		var forward proxy.Dialer
		var err error
		forward, _ = azber.NewDecorateDirect(conf.DNSCacheTimeout)
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

func BuildUpstream(upstream azber.Upstream, forward proxy.Dialer) (proxy.Dialer, error) {
	cipherDecorator := azber.NewCipherConnDecorator(upstream.Crypto, upstream.Password)
	forward = azber.NewDecorateClient(forward, cipherDecorator)

	switch strings.ToLower(upstream.Type) {
	case "socks5":
		{
			return azber.NewSocks5Client("tcp", upstream.Address, "", "", forward)
		}
	case "shadowsocks":
		{
			return azber.NewShadowsocksClient("tcp", upstream.Address, forward)
		}
	}
	return nil, errors.New("unknown upstream type" + upstream.Type)
}

func runSOCKS5Server(forward proxy.Dialer) {
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
