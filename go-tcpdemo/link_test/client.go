package main

import (
	"log"

	"sync"

	"strconv"

	"github.com/panenming/go-im/go-tcpdemo/link_test/codec"
	"github.com/panenming/go-im/libs/link"
)

// 并发100w
const batch int = 4000

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
	log.Println("打印返回值！")
	for i := 0; i < batch; i++ {
		wg.Add(1)
		go func() {
			client, err := link.Dial("tcp", addr, json, 20)
			checkErr(strconv.Itoa(i)+"client start ", err)
			//time.Sleep(2 * 60 * time.Second)
			err = client.Send(&AddReq{
				i, i,
			})
			checkErr(strconv.Itoa(i)+"client send ", err)
			//log.Printf("Send: %d + %d", i, i)
			rsp, err := client.Receive()
			wg.Done()
			if err != nil {
				log.Println("receive : ", err)
				return
			}
			log.Printf("Receive: %d", rsp.(*AddRsp).C)
		}()

		//log.Println("count=", i)
	}
	wg.Wait()
}

func checkErr(location string, err error) {
	if err != nil {
		log.Fatal(location, err.Error())
	}
}
