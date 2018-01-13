package main

import (
	"fmt"

	"sync"

	"github.com/micro/go-micro"
	proto "github.com/panenming/go-im/go-microdemo/server/proto"
	"golang.org/x/net/context"
)

func rpcClient(service micro.Service) {
	micro := proto.NewMessageServiceClient("messageMicro", service.Client())
	rsp, err := micro.CreateMsg(context.TODO(), &proto.Message{Id: 1, Fr: "pan"})

	if err != nil {
		fmt.Println(err)
		return
	}
	// Print response
	fmt.Println(rsp.Fr, rsp.Id)
}

func main() {
	service := micro.NewService()
	service.Init()
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			rpcClient(service)
			wg.Done()
		}()
	}
	wg.Wait()
}
