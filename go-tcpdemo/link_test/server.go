package main

import (
	"log"

	"github.com/panenming/go-im/go-tcpdemo/link_test/codec"
	"github.com/panenming/go-im/libs/link"
)

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

	server, err := link.Listen("tcp", "0.0.0.0:9000", json, 100 /* sync send */, link.HandlerFunc(serverSessionLoop))
	checkErr("startserver ", err)
	server.Serve()
}

func serverSessionLoop(session *link.Session) {
	for {
		req, err := session.Receive()

		checkErr("server receive ", err)
		//fmt.Println("请求的sessionid=", session.ID())

		err = session.Send(&AddRsp{
			req.(*AddReq).A + req.(*AddReq).B,
		})
		checkErr("server send ", err)
	}
}

func checkErr(location string, err error) {
	if err != nil {
		log.Fatal(location, err.Error())
	}
}
