package cash

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Cash struct {
	Storage    map[string]interface{}
	Mutex      *sync.RWMutex
	TTL        map[string]time.Time
	LastAccess map[string]time.Time
	Capacity   int
}

type KeyNotFoundError struct {
	Key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key '%s' not found in cache", e.Key)
}

func getDefaultCapacity() int {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	freeMemory := memStats.Frees
	capacity := int(freeMemory / 1024 / 10)

	if capacity < 100 {
		return 100
	}
	return capacity
}

func NewCash() *Cash {
	capacity := getDefaultCapacity()
	return &Cash{
		Storage:    make(map[string]interface{}),
		Mutex:      &sync.RWMutex{},
		TTL:        make(map[string]time.Time),
		LastAccess: make(map[string]time.Time),
		Capacity:   capacity,
	}
}

func (c *Cash) Set(key string, value interface{}, ttl time.Duration) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if len(c.Storage) >= c.Capacity {
		c.evict()
	}

	c.Storage[key] = value
	c.TTL[key] = time.Now().Add(ttl)
	c.LastAccess[key] = time.Now()
}

func (c *Cash) Get(key string) (interface{}, error) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	value, exists := c.Storage[key]
	if !exists {
		return nil, &KeyNotFoundError{Key: key}
	}

	if ttl, ok := c.TTL[key]; ok && time.Now().After(ttl) {
		c.deleteKey(key)
		return nil, &KeyNotFoundError{Key: key}
	}

	c.LastAccess[key] = time.Now()
	return value, nil
}

func (c *Cash) DeleteKey(key string) {
	delete(c.Storage, key)
	delete(c.TTL, key)
	delete(c.LastAccess, key)
}

func (c *Cash) evict() {
	var oldestKey string
	var oldestTime time.Time

	for key, accessTime := range c.LastAccess {
		if oldestKey == "" || accessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = accessTime
		}
	}
	if oldestKey != "" {
		c.deleteKey(oldestKey)
	}
}

func main() {
	cash := NewCash()

	cash.Set("key1", "value1", 5*time.Second)
	cash.Set("key2", "value2", 10*time.Second)

	value, err := cash.Get("key1")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Key1: %v\n", value)
	}

	time.Sleep(6 * time.Second)
	value, err = cash.Get("key1")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Key1: %v\n", value)
	}

	for i := 3; i <= cash.Capacity+3; i++ {
		cash.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 10*time.Second)
	}

	value, err = cash.Get("key2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Key2: %v\n", value)
	}
}
