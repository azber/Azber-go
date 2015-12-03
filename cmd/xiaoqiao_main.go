package main

import (
	"github.com/azber/Azber-go/xiaoqiao"
	"os"
	"os/signal"
	"fmt"
)

func main() {
	fmt.Println("start")
	for i := 20000; i < 35000; i++ {
		go begin(i)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill, os.Interrupt)
	<-sigChan
}

func begin(port int){
	xiao, err := xiaoqiao.NewXiaoqiao(port)
	if err != nil {
		return
	}
	xiao.Proxy()
}