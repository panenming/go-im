package connpool

import (
	"errors"
	"net"
	"sync"
)

type PoolConn struct {
	net.Conn
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
		return nil
	}
	return p.c.put(p.Conn)
}

func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.unusable = true
}

func (c *channelPool) wrapConn(conn net.Conn) net.Conn {
	p := &PoolConn{
		c: c,
	}
	p.Conn = conn
	return p
}

func (c *channelPool) put(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conns == nil {
		return conn.Close()
	}

	select {
	case c.conns <- conn:
		return nil
	default:
		return conn.Close()
	}
}

func (c *channelPool) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	conns := c.conns
	c.conns = nil

	if conns == nil {
		return
	}
	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func (c *channelPool) Len() int {
	conns, _ := c.getConnsAndFactory()
	return len(conns)
}
