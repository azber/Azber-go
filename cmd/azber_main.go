package main

import (
	"Azber-go/azber"
	"github.com/eahydra/socks"
	"net"
)

func main() {

}

func runSOCKS5Server(conf Proxy, forward socks.Dialer) {
	listener, err := net.Listen("tcp", "10802")
	if err != nil {
		ErrLog.Println("net.Listen failed, err:", err, "10802")
		return
	}
	cipherDecorator := NewCipherConnDecorator(conf.Crypto, conf.Password)
	listener = NewDecorateListener(listener, cipherDecorator)
	socks5Svr, err := azber.NewSocks5Server(forward)
	if err != nil {
		listener.Close()
		ErrLog.Println("socks.NewSocks5Server failed, err:", err)
		return
	}
	go func() {
		defer listener.Close()
		socks5Svr.Serve(listener)
	}()
}
