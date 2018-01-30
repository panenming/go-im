package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/gorilla/websocket"
	"github.com/panenming/go-im/go-drawingdemo/message"
)

var upgrader = websocket.Upgrader{
	// 跨域
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			if client != nil {
				hub.onConnect(client)
			}

		case client := <-hub.unregister:
			if client != nil {
				hub.onDisConnect(client)
			}

		}
	}
}

func (hub *Hub) close() {
	close(hub.register)
	close(hub.unregister)
}

func (hub *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 检测客户端是否支持websocket
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}
	client := newClient(hub, socket)
	hub.clients = append(hub.clients, client)
	hub.register <- client
	client.run()
}

func (hub *Hub) send(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	client.outbound <- data
}

// 广播
func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("client connected : ", client.socket.RemoteAddr())

	users := []message.User{}
	for _, c := range hub.clients {
		users = append(users, message.User{ID: c.id, Color: c.color})
	}

	//用户在线消息发送
	hub.send(message.NewConnected(client.color, users), client)
	// 用户在线广播
	hub.broadcast(message.NewUserJoined(client.id, client.color), client)
}

func (hub *Hub) onDisConnect(client *Client) {
	log.Println("client disconnected: ", client.socket.RemoteAddr())
	client.close()

	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
			break
		}
	}

	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]

	hub.broadcast(message.NewUserLeft(client.id), nil)
}

func (hub *Hub) onMessage(data []byte, client *Client) {
	kind := gjson.GetBytes(data, "kind").Int()
	if kind == message.KindStroke {
		var msg message.Stroke
		if json.Unmarshal(data, &msg) != nil {
			return
		}

		msg.UserID = client.id
		hub.broadcast(msg, client)
	} else if kind == message.KindClear {
		var msg message.Clear
		if json.Unmarshal(data, &msg) != nil {
			return
		}
		msg.UserID = client.id
		hub.broadcast(msg, client)
	}
}
