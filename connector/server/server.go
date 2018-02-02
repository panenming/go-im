package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/panenming/go-im/libs/link"
	"github.com/panenming/go-im/libs/proto"
)

type AddReq struct {
	A, B int
}

type AddRsp struct {
	C int
}

func Serve() {
	go http.ListenAndServe(":5001", nil)
	protocol := proto.Bufio()
	server, err := link.Listen("tcp", "0.0.0.0:5000", protocol, 100 /* sync send */, link.HandlerFunc(serverSessionLoop))
	if err != nil {
		fmt.Println("startserver ", err)
		return
	}

	server.Serve()
}

func serverSessionLoop(session *link.Session) {
	for {
		req, err := session.Receive()

		if err != nil {
			fmt.Println("这时候该session已经关闭...receive err : ", err)
			return
		}

		fmt.Println("req : ", req.(proto.Proto).Operation)

		rsp := &AddRsp{}
		rspData, _ := json.Marshal(rsp)
		p := &proto.Proto{
			HeaderLen: proto.RawHeaderSize,
			Ver:       1,
			Operation: proto.OP_PROTO_FINISH,
			SeqId:     1,
			Body:      rspData,
		}
		err = session.Send(p)
		if err != nil {
			fmt.Println("send err : ", err)
		}

	}
}
