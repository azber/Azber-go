package main

import (
	"Azber-go/azber"
	"golang.org/x/net/proxy"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":10800")
	if err != nil {
		return
	}
	defer listener.Close()

	if server, err := azber.NewSocks5Server(proxy.Direct); err == nil {
		server.Serve(listener)
	}
}
