// Copyright 2026, TeamDev. All rights reserved.
//
// Redistribution and use in source and/or binary forms, with or without
// modification, must retain the above copyright notice and the following
// disclaimer.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package fragmentation

import (
	"container/list"
	"sync"
)

// cache stores a limited number of least-recently-used values by string key.
type cache[T any] struct {
	sync.Mutex
	limit   int
	values  map[string]T
	entries map[string]*list.Element
	order   *list.List
}

// newCache creates a cache with least-recently-used eviction.
func newCache[T any](limit int) *cache[T] {
	return &cache[T]{
		limit:   limit,
		values:  make(map[string]T),
		entries: make(map[string]*list.Element),
		order:   list.New(),
	}
}

// get returns a cached value and marks it as recently used.
func (c *cache[T]) get(key string) (T, bool) {
	c.Lock()
	defer c.Unlock()

	value, found := c.values[key]
	if found {
		c.markUsed(key)
	}

	return value, found
}

// set stores a value and evicts the least recently used value when the cache exceeds its limit.
func (c *cache[T]) set(key string, value T) {
	c.Lock()
	defer c.Unlock()

	c.values[key] = value
	if entry, found := c.entries[key]; found {
		c.order.MoveToBack(entry)
		return
	}

	c.entries[key] = c.order.PushBack(key)
	if len(c.values) <= c.limit {
		return
	}

	c.evictOldest()
}

// clear removes all cached values.
func (c *cache[T]) clear() {
	c.Lock()
	defer c.Unlock()

	c.values = make(map[string]T)
	c.entries = make(map[string]*list.Element)
	c.order.Init()
}

// markUsed moves a cache key to the most recently used position.
func (c *cache[T]) markUsed(key string) {
	if entry, found := c.entries[key]; found {
		c.order.MoveToBack(entry)
	}
}

// evictOldest removes the least recently used cache entry.
func (c *cache[T]) evictOldest() {
	oldestEntry := c.order.Front()
	if oldestEntry == nil {
		return
	}

	oldestKey := oldestEntry.Value.(string)
	c.order.Remove(oldestEntry)
	delete(c.entries, oldestKey)
	delete(c.values, oldestKey)
}
