package collections

import (
	"reflect"
	"sync"
)

type Iterator interface {
	HasNext() bool
	Next() interface{}
	Remove()
}

type SyncIdentityMap struct {
	sync.RWMutex
	m map[uintptr]interface{}
}

func NewSyncMap() *SyncIdentityMap {
	return &SyncIdentityMap{
		m: make(map[uintptr]interface{}),
	}
}

func (m *SyncIdentityMap) Get(key interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()
	keyPtr := genKey(key)
	value := m.m[keyPtr]
	return value
}

func genKey(key interface{}) uintptr {
	keyValue := reflect.ValueOf(key)
	return keyValue.Pointer()
}

func (m *SyncIdentityMap) Put(key interface{}, value interface{}) {
	m.Lock()
	defer m.Unlock()
	keyPtr := genKey(key)
	m.m[keyPtr] = value
}

func (m *SyncIdentityMap) Remove(key interface{}) {
	m.Lock()
	defer m.Unlock()
	keyPtr := genKey(key)
	delete(m.m, keyPtr)
}

func (m *SyncIdentityMap) Size() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.m)
}

func (m *SyncIdentityMap) Values() []interface{} {
	m.RLock()
	defer m.RUnlock()
	list := make([]interface{}, 0)
	for _, v := range m.m {
		list = append(list, v)
	}
	return list
}
