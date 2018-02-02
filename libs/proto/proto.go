package proto

import (
	"encoding/json"
	"errors"
	"fmt"

	"time"

	"github.com/panenming/go-im/libs/proto/bufio"
	"github.com/panenming/go-im/libs/proto/encoding/binary"

	"github.com/gorilla/websocket"
)

// tcp 最大传输数据量
const (
	MaxBodySize = int32(1 << 10)
)

// 整个接收到的数据
const (
	PackSize      = 4
	HeaderSize    = 2
	VerSize       = 2
	OperationSize = 4
	SeqIdSize     = 4
	RawHeaderSize = PackSize + HeaderSize + VerSize + OperationSize + SeqIdSize
	MaxPackSize   = MaxBodySize + int32(RawHeaderSize) // 最大承载的数据量
	// 获取相关值的偏移量
	PackOffset      = 0
	HeaderOffset    = PackOffset + PackSize
	VerOffset       = HeaderOffset + HeaderSize
	OperationOffset = VerOffset + VerSize
	SeqIdOffset     = OperationOffset + OperationSize
)

type Proto struct {
	HeaderLen int16           `json:"-"`    // header length
	Ver       int16           `json:"ver"`  // protocol version
	Operation int32           `json:"op"`   // operation for request
	SeqId     int32           `json:"seq"`  // sequence number chosen by client
	Body      json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte)
	Time      time.Time       `json:"-"`    // proto send time
}

var (
	emptyProto    = Proto{}
	emptyJSONBody = []byte("{}")

	ErrProtoPackLen   = errors.New("消息包长度解析出错")
	ErrProtoHeaderLen = errors.New("信息头部解析出错")
)

var (
	ProtoReady  = &Proto{Operation: OP_PROTO_READY}
	ProtoFinish = &Proto{Operation: OP_PROTO_FINISH}
)

func (p *Proto) Reset() {
	*p = emptyProto
}

func (p *Proto) String() string {
	return fmt.Sprintf("\n-------- proto --------\nver: %d\nop: %d\nseq: %d\nbody: %v\n-----------------------", p.Ver, p.Operation, p.SeqId, p.Body)
}

// 封装writer
func (p *Proto) WriteTo(wr *bufio.Writer) (err error) {
	defer wr.Flush()
	var (
		buf     []byte
		packLen int32
	)
	packLen = RawHeaderSize + int32(len(p.Body))
	p.HeaderLen = RawHeaderSize
	if buf, err = wr.Peek(RawHeaderSize); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[PackOffset:], packLen)
	binary.BigEndian.PutInt16(buf[HeaderOffset:], p.HeaderLen)
	binary.BigEndian.PutInt16(buf[VerOffset:], p.Ver)
	binary.BigEndian.PutInt32(buf[OperationOffset:], p.Operation)
	binary.BigEndian.PutInt32(buf[SeqIdOffset:], p.SeqId)
	if p.Body != nil {
		_, err = wr.Write(p.Body)
	}
	return
}

// 从io中解析数据
func (p *Proto) ReadFr(rr *bufio.Reader) (err error) {
	var (
		bodyLen int
		packLen int32
		buf     []byte
	)
	if buf, err = rr.Pop(RawHeaderSize); err != nil {
		return
	}
	packLen = binary.BigEndian.Int32(buf[PackOffset:HeaderOffset])
	p.HeaderLen = binary.BigEndian.Int16(buf[HeaderOffset:VerOffset])
	p.Ver = binary.BigEndian.Int16(buf[VerOffset:OperationOffset])
	p.Operation = binary.BigEndian.Int32(buf[OperationOffset:SeqIdOffset])
	p.SeqId = binary.BigEndian.Int32(buf[SeqIdOffset:])
	if packLen > MaxPackSize {
		return ErrProtoPackLen
	}
	if p.HeaderLen != RawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(p.HeaderLen)); bodyLen > 0 {
		p.Body, err = rr.Pop(bodyLen)
	} else {
		p.Body = nil
	}
	return
}

func (p *Proto) ReadWebsocket(wr *websocket.Conn) (err error) {
	err = wr.ReadJSON(p)
	return
}

func (p *Proto) WriteWebsocket(wr *websocket.Conn) (err error) {
	if p.Body == nil {
		p.Body = emptyJSONBody
	}
	// [{"ver":1,"op":8,"seq":1,"body":{}}, {"ver":1,"op":3,"seq":2,"body":{}}]
	err = wr.WriteJSON([]*Proto{p})
	return
}
