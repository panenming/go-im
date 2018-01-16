package link

import (
	"sync"
)

/**
Manager 默认有32个SessionMap
SessionMap中以key value形式存储Session
key是Session.id
value是Session
根据Session.id % 32 hash 到不同的SessionMap中
**/
const sessionMapNum = 32

type Manager struct {
	sessionMaps [sessionMapNum]sessionMap
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
}

// session id和对象的对照map,便于查询
type sessionMap struct {
	sync.RWMutex
	sessions map[uint64]*Session
	disposed bool
}

func NewManager() *Manager {
	manager := &Manager{}
	for i := 0; i < len(manager.sessionMaps); i++ {
		manager.sessionMaps[i].sessions = make(map[uint64]*Session)
	}
	return manager
}

func (manager *Manager) Dispose() {
	manager.disposeOnce.Do(func() {
		for i := 0; i < sessionMapNum; i++ {
			smap := &manager.sessionMaps[i]
			smap.Lock()
			smap.disposed = true
			for _, session := range smap.sessions {
				session.Close()
			}
			smap.Unlock()
		}
		manager.disposeWait.Wait()
	})
}

func (manager *Manager) NewSession(codec Codec, sendChanSize int) *Session {
	session := newSession(manager, codec, sendChanSize)
	manager.putSession(session)
	return session
}

func (manager *Manager) putSession(session *Session) {
	smap := &manager.sessionMaps[session.id%sessionMapNum]

	smap.RLock()
	defer smap.RUnlock()
	if smap.disposed {
		session.Close()
		return
	}
	smap.sessions[session.id] = session
	manager.disposeWait.Add(1)
}

func (manager *Manager) GetSession(sessionId uint64) *Session {
	smap := manager.sessionMaps[sessionId%sessionMapNum]
	smap.RLock()
	defer smap.RUnlock()

	session, _ := smap.sessions[sessionId]
	return session
}

func (manager *Manager) delSession(session *Session) {
	smap := &manager.sessionMaps[session.id%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	delete(smap.sessions, session.id)
	manager.disposeWait.Done()
}
