package proto

import (
	"io"

	"github.com/panenming/go-im/libs/proto/bufio"

	"github.com/panenming/go-im/libs/link"
)

type bufioProtocol struct{}

func Bufio() link.Protocol {
	return &bufioProtocol{}
}

func (b *bufioProtocol) NewCodec(rw io.ReadWriter) (cc link.Codec, err error) {
	cc = &bufioCodec{
		rw:     rw,
		reader: bufio.NewReader(rw),
		writer: bufio.NewWriter(rw),
	}
	return
}

type bufioCodec struct {
	rw     io.ReadWriter
	reader *bufio.Reader
	writer *bufio.Writer
}

func (c *bufioCodec) Send(msg interface{}) error {
	if msg == nil {
		return nil
	}
	p := msg.(*Proto)
	err := p.WriteTo(c.writer)
	if err != nil {
		return err
	}
	return nil
}

func (c *bufioCodec) Receive() (interface{}, error) {
	p := emptyProto
	err := p.ReadFr(c.reader)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (c *bufioCodec) Close() error {
	if closer, ok := c.rw.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
