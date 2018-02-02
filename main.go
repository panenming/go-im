package main

import (
	"fmt"

	//	"github.com/panenming/go-im/connector/client"
	"github.com/panenming/go-im/connector/server"
)

func main() {
	fmt.Println("启动服务器")
	server.Serve()
}
