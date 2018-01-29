package main

import (
	"github.com/gorilla/websocket"
	"github.com/panenming/go-im/go-drawingdemo/utils"
	uuid "github.com/satori/go.uuid"
)

// client 实体
type Client struct {
	id       string
	hub      *Hub
	color    string
	socket   *websocket.Conn
	outbound chan []byte
}

func newClient(hub *Hub, socket *websocket.Conn) *Client {
	uuid, _ := uuid.NewV4()
	return &Client{
		id:       uuid.String(),
		color:    utils.GenerateColor(),
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (client *Client) read() {
	defer func() {
		client.hub.unregister <- client
	}()

	for {
		_, data, err := client.socket.ReadMessage()
		if err != nil {
			continue
		}
		client.hub.onMessage(data, client)
	}
}

func (client *Client) write() {
	for {
		select {
		case data, ok := <-client.outbound:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.socket.WriteMessage(websocket.TextMessage, data)
		}
	}
}

// run method
func (client Client) run() {
	go client.read()
	go client.write()
}

// close method
func (client Client) close() {
	client.socket.Close()
	close(client.outbound)
}
