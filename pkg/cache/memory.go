package cache

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"
)

type MemoryCache struct {
	sync.Mutex
	expirations sync.Map
	data        map[string]interface{}
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{data: make(map[string]interface{})}
}

// Get method retrieves a value from the cache, if it exists,
// for a given key. The function locks the cache before
// attempting to retrieve the value to avoid race conditions.
// The returned value is of type interface{} as the cache
// can store any type of value. The second returned boolean
// value indicates if the key was found in the cache or not.
// If the retrieval was successful, the boolean return is 'true',
// if not, it's 'false'.
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	value, ok := c.data[key]
	return value, ok
}

// Delete removes the value associated with the given key from the cache.
// It locks the cache during deletion to prevent concurrent map writes.
func (c *MemoryCache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.data, key)
}

// Clear 清空缓存中的所有键值对。
// 它在清空期间加锁,以防止并发的 map 写入。
func (c *MemoryCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.data = make(map[string]interface{})
}

// Set stores the given value in the cache for the given key and expiration time.
// It locks the cache during the set, and handles expiration by starting a goroutine
// that will delete the key after the expiration time. If expiration is 0, it deletes
// any existing expiration for the key.
func (c *MemoryCache) Set(key string, value interface{}, expiration time.Duration) {
	c.Lock()
	defer c.Unlock()

	c.data[key] = value

	// 如果设置了过期时间，则开启定时器，在到达过期时间时自动删除该键值对
	if expiration > 0 {
		expirationTime := time.Now().Add(expiration)
		c.expirations.Store(key, expirationTime)
		go func() {
			timer := time.NewTimer(expiration)
			<-timer.C

			c.Lock()
			defer c.Unlock()

			// 原子操作检查键值对是否过期
			if expire, ok := c.expirations.Load(key); ok && expire.(time.Time).Equal(expirationTime) {
				delete(c.data, key)
				c.expirations.Delete(key)
			}
		}()
	} else {
		c.expirations.Delete(key)
	}
}

// deepCopy makes a deep copy of the given value using gob encoding/decoding.
// It encodes the value to a bytes buffer, decodes it back to an interface{},
// and returns the decoded value. This creates a deep copy by value instead of
// just copying a reference.
func (c *MemoryCache) deepCopy(value interface{}) interface{} {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	decoder := gob.NewDecoder(buffer)

	// 编码 value
	encoder.Encode(value)

	// 解码 value 的副本
	var copiedValue interface{}
	decoder.Decode(&copiedValue)

	return copiedValue
}
