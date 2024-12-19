/*
@author: fengxu
@since: 2024-12-19
@desc: //TODO
*/

package _map

import (
	"fmt"
	"sync"
	"testing"
)

// TestRWMap 测试 RWMap 的基本功能
func TestRWMap(t *testing.T) {
	m := NewRWMap()

	// 测试 Set 和 Get 方法
	m.Set("key1", "value1")
	if val, ok := m.Get("key1"); !ok || val != "value1" {
		t.Errorf("Expected 'value1', got '%v'", val)
	}

	// 测试 Delete 方法
	m.Delete("key1")
	if _, ok := m.Get("key1"); ok {
		t.Error("Expected key1 to be deleted")
	}

	// 测试 Len 方法
	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}

// TestChannelMap 测试 ChannelMap 的基本功能
func TestChannelMap(t *testing.T) {
	m := NewChannelMap()

	// 测试 Set 和 Get 方法
	m.Set("key1", "value1")
	if val, ok := m.Get("key1"); !ok || val != "value1" {
		t.Errorf("Expected 'value1', got '%v'", val)
	}

	// 测试 Delete 方法
	m.Delete("key1")
	if _, ok := m.Get("key1"); ok {
		t.Error("Expected key1 to be deleted")
	}

	// 测试 Len 方法
	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}

// TestSplitLockMap 测试 SplitLockMap 的基本功能
func TestSplitLockMap(t *testing.T) {
	m := NewSplitLockMap(8)

	// 测试 Set 和 Get 方法
	m.Set("key1", "value1")
	if val, ok := m.Get("key1"); !ok || val != "value1" {
		t.Errorf("Expected 'value1', got '%v'", val)
	}

	// 测试 Delete 方法
	m.Delete("key1")
	if _, ok := m.Get("key1"); ok {
		t.Error("Expected key1 to be deleted")
	}

	// 测试 Len 方法
	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}

// TestConcurrentRWMap 测试 RWMap 在并发环境下的表现
func TestConcurrentRWMap(t *testing.T) {
	m := NewRWMap()
	var wg sync.WaitGroup
	const goroutines = 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			m.Set(key, "value")
			if val, ok := m.Get(key); !ok || val != "value" {
				t.Errorf("Failed to get value for key: %s", key)
			}
			m.Delete(key)
		}(i)
	}
	wg.Wait()

	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}

// TestConcurrentChannelMap 测试 ChannelMap 在并发环境下的表现
func TestConcurrentChannelMap(t *testing.T) {
	m := NewChannelMap()
	var wg sync.WaitGroup
	const goroutines = 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			m.Set(key, "value")
			if val, ok := m.Get(key); !ok || val != "value" {
				t.Errorf("Failed to get value for key: %s", key)
			}
			m.Delete(key)
		}(i)
	}
	wg.Wait()

	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}

// TestConcurrentSplitLockMap 测试 SplitLockMap 在并发环境下的表现
func TestConcurrentSplitLockMap(t *testing.T) {
	m := NewSplitLockMap(8)
	var wg sync.WaitGroup
	const goroutines = 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			m.Set(key, "value")
			if val, ok := m.Get(key); !ok || val != "value" {
				t.Errorf("Failed to get value for key: %s", key)
			}
			m.Delete(key)
		}(i)
	}
	wg.Wait()

	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}

// TestSplitLockMapLen 并发操作后检查长度是否正确
func TestSplitLockMapLen(t *testing.T) {
	m := NewSplitLockMap(8)
	var wg sync.WaitGroup
	const goroutines = 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			m.Set(key, "value")
		}(i)
	}
	wg.Wait()

	if m.Len() != goroutines {
		t.Errorf("Expected length %d, got %d", goroutines, m.Len())
	}

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			m.Delete(key)
		}(i)
	}
	wg.Wait()

	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}
