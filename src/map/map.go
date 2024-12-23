package _map

import (
	"fmt"
	"sync"
)

// RWMap 基于读写锁实现的线程安全map
type RWMap struct {
	m    map[string]string
	lock sync.RWMutex
}

func NewRWMap() *RWMap {
	return &RWMap{m: make(map[string]string)}
}

func (m *RWMap) Set(key, value string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.m[key] = value
}

func (m *RWMap) Get(key string) (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	val, ok := m.m[key]
	return val, ok
}

func (m *RWMap) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.m, key)
}

func (m *RWMap) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.m)
}

func (m *RWMap) PrintMap() {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for k, v := range m.m {
		fmt.Println(k, v)
	}
}

// ChannelMap 基于channel实现的线程安全map
type ChannelMap struct {
	m     map[string]string
	mutex chan struct{} // 基于通道实现的锁,保证了在同一时间只有一个goroutine可以访问map
}

func NewChannelMap() *ChannelMap {
	return &ChannelMap{m: make(map[string]string), mutex: make(chan struct{}, 1)}
}

func (cm *ChannelMap) Lock() {
	cm.mutex <- struct{}{}
}

func (cm *ChannelMap) Unlock() {
	<-cm.mutex
}

func (cm *ChannelMap) Set(key, value string) {
	cm.Lock()
	defer cm.Unlock()
	cm.m[key] = value
}

func (cm *ChannelMap) Get(key string) (string, bool) {
	cm.Lock()
	defer cm.Unlock()
	if cm.m == nil {
		return "", false
	}
	value, ok := cm.m[key]
	return value, ok
}

func (cm *ChannelMap) Delete(key string) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.m, key)
}

func (cm *ChannelMap) Len() int {
	cm.Lock()
	defer cm.Unlock()
	return len(cm.m)
}

func (cm *ChannelMap) PrintMap() {
	cm.Lock()
	defer cm.Unlock()
	for k, v := range cm.m {
		fmt.Println(k, v)
	}
}
