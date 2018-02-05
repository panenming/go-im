package main

import (
	"fmt"

	//	"github.com/panenming/go-im/connector/client"
	"github.com/panenming/go-im/connector/server"
)

func main() {
	// 使用vscode 编辑
	fmt.Println("启动服务器")
	server.Serve()
}
