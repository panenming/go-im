package link

import (
	"net"
)

// 定义和server先关的属性和方法
type Handler interface {
	HandleSession(*Session)
}
type Server struct {
	manager      *Manager
	listener     net.Listener
	protocol     Protocol
	handler      Handler
	sendChanSize int
}

var _ Handler = HandlerFunc(nil)

type HandlerFunc func(*Session)

func (f HandlerFunc) HandleSession(session *Session) {
	f(session)
}

func NewServer(listener net.Listener, protocol Protocol, sendChanSize int, handler Handler) *Server {
	return &Server{
		manager:      NewManager(),
		listener:     listener,
		protocol:     protocol,
		handler:      handler,
		sendChanSize: sendChanSize,
	}
}

func (server *Server) Listener() net.Listener {
	return server.listener
}

func (server *Server) Serve() error {
	for {
		conn, err := Accept(server.listener)
		if err != nil {
			return err
		}

		go func() {
			codec, err := server.protocol.NewCodec(conn)
			if err != nil {
				conn.Close()
				return
			}
			session := server.manager.NewSession(codec, server.sendChanSize)
			server.handler.HandleSession(session)
		}()
	}
}

func (server *Server) GetSession(sessionId uint64) *Session {
	return server.manager.GetSession(sessionId)
}

func (server *Server) Stop() {
	server.listener.Close()
	server.manager.Dispose()
}
