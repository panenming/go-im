package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/panenming/go-im/go-drawingdemo/message"
	uuid "github.com/satori/go.uuid"
)

type Room struct {
	hub    *Hub
	roomId string
}

// 已经创建的room
type Rooms struct {
	rooms map[string]*Room
}

func RoomsInit(size int) *Rooms {
	rooms := &Rooms{
		rooms: make(map[string]*Room, size),
	}
	return rooms
}

func (rooms *Rooms) newRoom(hub *Hub) *Room {
	uuid, _ := uuid.NewV4()
	roomId := uuid.String()
	room := &Room{
		hub:    hub,
		roomId: roomId,
	}
	rooms.rooms[roomId] = room
	return room
}

func (rooms *Rooms) RoomCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	log.Println("开始创建room")
	// TODO 需要限制room的个数
	hub := newHub()
	go hub.run()

	room := rooms.newRoom(hub)
	if room.roomId != "" { // room 创建成功，返回roomId

		// 每隔60s扫描room是否无效进行删除
		go func() {
			for {
				time.Sleep(60 * time.Second)
				rooms.destoryRoom(room.roomId)

			}
		}()
		roomCreate := &message.RoomCreate{
			Kind:   message.KindRoomCreate,
			RoomID: room.roomId,
		}
		data, _ := json.Marshal(roomCreate)
		w.Write(data)
	} else {
		roomCreateError := &message.RoomCreate{
			Kind: message.KindRoomCreateError,
		}
		data, _ := json.Marshal(roomCreateError)
		w.Write(data)
	}
}

func (rooms *Rooms) RoomHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roomId := r.Form.Get("roomId")
	if roomId == "" {
		// 需要返回到创建room的页面
		return
	}

	// 用户已经有room,查询room是否已经创建
	room := rooms.rooms[roomId]
	if room == nil {
		// 找不到room直接返回
		return
	}

	room.hub.handleWebSocket(w, r)

}

func (rooms *Rooms) RoomListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	log.Println("开始获取room")
	roomLsit := &message.RoomList{
		Kind: 0,
	}
	roomIds := []string{}
	for roomId, _ := range rooms.rooms {
		roomIds = append(roomIds, roomId)
	}

	roomLsit.RoomIds = roomIds
	data, _ := json.Marshal(roomLsit)
	w.Write(data)
}

// 删除无效的room
func (rooms *Rooms) destoryRoom(roomId string) {
	room := rooms.rooms[roomId]
	if room != nil {
		clientSize := len(room.hub.clients)
		if clientSize <= 0 {
			// 无效的room删除
			log.Println("roomId:", roomId, "无效进行删除")
			room.hub.close()
			delete(rooms.rooms, roomId)
		}
	}
}
