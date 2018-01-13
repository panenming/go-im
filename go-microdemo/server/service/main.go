package main

import (
	"fmt"

	"github.com/panenming/go-im/go-microdemo/server/service/handler"

	"github.com/micro/go-micro"
)

func main() {
	fmt.Println("Hello")

	service := micro.NewService(
		micro.Name("messageMicro"),
		micro.Version("last"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),
	)
	service.Init()

	server := service.Server()
	server.Handle(server.NewHandler(&handler.MessageService{}))

	err := service.Run()
	if err != nil {
		fmt.Println(err)
	}

}
