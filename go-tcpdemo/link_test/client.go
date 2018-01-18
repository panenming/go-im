package main

import (
	"log"

	"sync"

	"strconv"

	"github.com/panenming/go-im/go-tcpdemo/link_test/codec"
	"github.com/panenming/go-im/libs/link"
)

// 并发1w
const batch int = 10000

type AddReq struct {
	A, B int
}

type AddRsp struct {
	C int
}

func main() {
	json := codec.Json()

	json.Register(AddReq{})
	json.Register(AddRsp{})

	addr := "127.0.0.1:9000"

	clientSessionLoop(addr, json)
}

func clientSessionLoop(addr string, json *codec.JsonProtocol) {
	var wg sync.WaitGroup
	for i := 0; i < batch; i++ {
		wg.Add(1)
		client, err := link.Dial("tcp", addr, json, 0)
		checkErr(strconv.Itoa(i)+"client start ", err)
		err = client.Send(&AddReq{
			i, i,
		})
		checkErr(strconv.Itoa(i)+"client send ", err)
		//log.Printf("Send: %d + %d", i, i)

		_, err = client.Receive()
		checkErr(strconv.Itoa(i)+"client receive ", err)
		//log.Printf("Receive: %d", rsp.(*AddRsp).C)
		log.Println("count=", i)
		wg.Done()
	}
	wg.Wait()
}

func checkErr(location string, err error) {
	if err != nil {
		log.Fatal(location, err.Error())
	}
}
