package link

import (
	"sync"
)

// channel 适用于特殊场景下的推送，比如聊天室等等
// 一个聊天室就是一个channel
type KEY interface{}

type Channel struct {
	mutex    sync.RWMutex
	sessions map[KEY]*Session
	State    interface{}
}

func NewChannel() *Channel {
	return &Channel{
		sessions: make(map[KEY]*Session),
	}
}

func (channel *Channel) Len() int {
	channel.mutex.RLock()
	defer channel.mutex.RUnlock()
	return len(channel.sessions)
}

func (channel *Channel) Fetch(callBack func(*Session)) {
	channel.mutex.RLock()
	defer channel.mutex.RUnlock()
	for _, session := range channel.sessions {
		callBack(session)
	}
}

func (channel *Channel) Get(key KEY) *Session {
	channel.mutex.RLock()
	defer channel.mutex.RUnlock()
	session, _ := channel.sessions[key]
	return session
}

func (channel *Channel) Put(key KEY, session *Session) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	if session, exists := channel.sessions[key]; exists {
		channel.remove(key, session)
	}
	session.AddCloseCallback(channel, key, func() {
		channel.Remove(key)
	})
	channel.sessions[key] = session
}

func (channel *Channel) remove(key KEY, session *Session) {
	session.RemoveCloseCallback(channel, key)
	delete(channel.sessions, key)
}
func (channel *Channel) Remove(key KEY) bool {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	session, exists := channel.sessions[key]
	if exists {
		channel.remove(key, session)
	}
	return exists
}

func (channel *Channel) FetchAndRemove(callback func(*Session)) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	for key, session := range channel.sessions {
		session.RemoveCloseCallback(channel, key)
		delete(channel.sessions, key)
		callback(session)
	}
}

func (channel *Channel) Close() {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	for key, session := range channel.sessions {
		channel.remove(key, session)
	}
}
