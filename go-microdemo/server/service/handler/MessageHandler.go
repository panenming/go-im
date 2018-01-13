package handler

import (
	"fmt"

	proto "github.com/panenming/go-im/go-microdemo/server/proto"
	"golang.org/x/net/context"
)

type MessageService struct{}

func (msgService *MessageService) CreateMsg(ctx context.Context, req *proto.Message, rsp *proto.Message) error {
	fmt.Println(req)
	fmt.Println(req.Id)
	*rsp = proto.Message{
		Id: 32131,
	}
	return nil
}
