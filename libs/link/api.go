package link

import (
	"io"
	"net"
	"strings"
	"time"
)

// 协议以及编码定义的接口
// 定义 protocol中NewCodec io.ReadWriter，这样可以直接将网络的包装进行发送
type Codec interface {
	Receive() (interface{}, error)
	Send(interface{}) error
	Close() error
}

type Protocol interface {
	NewCodec(rw io.ReadWriter) (Codec, error)
}

type ProtocolFunc func(rw io.ReadWriter) (Codec, error)

func (pf ProtocolFunc) NewCodec(rw io.ReadWriter) (Codec, error) {
	return pf(rw)
}

type ClearSendChan interface {
	ClearSendChan(<-chan interface{})
}

func Listen(network, addr string, protocol Protocol, sendChanSize int, handler Handler) (*Server, error) {
	listener, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	return NewServer(listener, protocol, sendChanSize, handler), nil
}

func Dial(network, addr string, protocol Protocol, sendChanSize int) (*Session, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	codec, err := protocol.NewCodec(conn)
	if err != nil {
		return nil, err
	}
	return NewSession(codec, sendChanSize), nil
}

func DailTimeout(network, addr string, timeout time.Duration, protocol Protocol, sendChanSize int) (*Session, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	codec, err := protocol.NewCodec(conn)
	if err != nil {
		return nil, err
	}
	return NewSession(codec, sendChanSize), nil
}

func Accept(listener net.Listener) (net.Conn, error) {
	var tempDelay time.Duration
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}
		return conn, nil
	}
}
