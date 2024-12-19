package _map

import (
	"fmt"
	"hash/fnv"
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

// SplitLockMap 基于分片加锁的线程安全map
type SplitLockMap struct {
	slices   []SplitSlice          // 切片数组
	size     int                   // 切片的数量
	hashFunc func(key string) uint // 哈希函数来确定key所在哪个切片中
}

// SplitSlice map里面的单个切片
type SplitSlice struct {
	m     map[string]string
	mutex sync.RWMutex
}

func NewSplitLockMap(size int) *SplitLockMap {
	slices := make([]SplitSlice, size)
	for i := 0; i < size; i++ {
		slices[i] = SplitSlice{m: make(map[string]string)}
	}
	return &SplitLockMap{slices: slices, size: size, hashFunc: fnvHash64}
}

func fnvHash64(key string) uint {
	h := fnv.New64()
	h.Write([]byte(key))
	return uint(h.Sum64())
}

func (s *SplitLockMap) Set(key, value string) {
	index := s.hashFunc(key) % uint(s.size)
	s.slices[index].mutex.Lock()
	defer s.slices[index].mutex.Unlock()
	s.slices[index].m[key] = value
}

func (s *SplitLockMap) Get(key string) (string, bool) {
	index := s.hashFunc(key) % uint(s.size)
	s.slices[index].mutex.RLock()
	defer s.slices[index].mutex.RUnlock()
	val, ok := s.slices[index].m[key]
	return val, ok
}

func (s *SplitLockMap) Delete(key string) {
	index := s.hashFunc(key) % uint(s.size)
	s.slices[index].mutex.Lock()
	defer s.slices[index].mutex.Unlock()
	delete(s.slices[index].m, key)
}

func (s *SplitLockMap) Len() int {
	var length int
	for _, slice := range s.slices {
		slice.mutex.RLock()
		length += len(slice.m)
		slice.mutex.RUnlock()
	}
	return length
}

func (s *SplitLockMap) PrintMap() {
	for _, slice := range s.slices {
		slice.mutex.RLock()
		for k, v := range slice.m {
			fmt.Println(k, v)
		}
		slice.mutex.RUnlock()
	}
}
