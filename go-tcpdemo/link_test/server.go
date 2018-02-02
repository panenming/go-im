package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

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
	go http.ListenAndServe(":5001", nil)
	json := codec.Json()

	json.Register(AddReq{})
	json.Register(AddRsp{})

	server, err := link.Listen("tcp", "0.0.0.0:5000", json, 100 /* sync send */, link.HandlerFunc(serverSessionLoop))
	if err != nil {
		fmt.Println("startserver ", err)
		return
	}

	server.Serve()
}

func serverSessionLoop(session *link.Session) {
	for {
		req, err := session.Receive()

		//fmt.Println("请求的sessionid=", session.ID())

		if err != nil {
			fmt.Println("server receive---", err.Error())
			continue
		}
		err = session.Send(&AddRsp{
			req.(*AddReq).A,
		})

		if err != nil {
			fmt.Println("server send---", err.Error())
			continue
		}
	}
}
