package main

import (
	"fmt"

	"github.com/junkeWu/GoGame/common/log"
)

func main() {
	fmt.Println("start service")
	log.Config("./log/biz_server.log")
	log.Info("Hello")
}
