package main

import (
	"log"
	"net/http"
)

const (
	roomsize = 10 * 1000 //最多创建的room数量
)

func main() {

	rooms := RoomsInit(roomsize)

	// ws连接room
	http.HandleFunc("/ws", rooms.RoomHandler)
	// 创建room
	http.HandleFunc("/create", rooms.RoomCreateHandler)
	// 获取roomId列表
	http.HandleFunc("/roomlist", rooms.RoomListHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
