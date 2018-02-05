package connpool

import (
	"errors"
	"net"
)

// tcp 连接池
var (
	ErrClosed = errors.New("pool is closed")
)

type Pool interface {
	Get() (net.Conn, error)
	Close()
	Len() int
}
