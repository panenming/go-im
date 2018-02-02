package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/panenming/go-im/libs/link"
	"github.com/panenming/go-im/libs/proto"
)

type AddReq struct {
	A, B int
}

type AddRsp struct {
	C int
}

func main() {

	addr := "127.0.0.1:5000"
	pro := proto.Bufio()
	clientSessionLoop(addr, pro)
}

func clientSessionLoop(addr string, pro link.Protocol) {
	log.Println("打印返回值！")
	client, err := link.Dial("tcp", addr, pro, 0)
	checkErr(" client start err ：", err)
	for i := 0; i < 3; i++ {

		req := &AddReq{
			A: 1,
			B: 1,
		}
		reqData, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
		p := &proto.Proto{
			HeaderLen: proto.RawHeaderSize,
			Ver:       1,
			Operation: proto.OP_PROTO_READY,
			SeqId:     1,
			Body:      reqData,
		}
		log.Println("发送消息start-----")
		err = client.Send(p)
		log.Println("发送消息end-----")

		checkErr("client send ", err)
		time.Sleep(1 * time.Second)
		rsp, err := client.Receive()

		if err != nil {
			log.Println("receive : ", err)
			return
		}

		log.Println("rsp : ", rsp.(proto.Proto).Operation)
	}

}

func checkErr(location string, err error) {
	if err != nil {
		log.Fatal(location, err.Error())
	}
}
