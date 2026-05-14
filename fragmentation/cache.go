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

// cache stores a limited number of least-recently-used values by key.
type cache[K comparable, V any] struct {
	sync.Mutex
	limit   int
	loader  func(K) (V, error)
	values  map[K]V
	entries map[K]*list.Element
	order   *list.List
}

// newCache creates a cache with a loader and least-recently-used eviction.
func newCache[K comparable, V any](limit int, loader func(K) (V, error)) *cache[K, V] {
	return &cache[K, V]{
		limit:   limit,
		loader:  loader,
		values:  make(map[K]V),
		entries: make(map[K]*list.Element),
		order:   list.New(),
	}
}

// get returns a cached value or loads it when missing.
func (c *cache[K, V]) get(key K) (V, error) {
	c.Lock()
	value, found := c.values[key]
	if found {
		c.markUsed(key)
		c.Unlock()

		return value, nil
	}
	c.Unlock()

	value, err := c.loader(key)
	if err != nil {
		return value, err
	}

	c.Lock()
	defer c.Unlock()
	c.storeLoaded(key, value)

	return value, nil
}

// clear removes all cached values.
func (c *cache[K, V]) clear() {
	c.Lock()
	defer c.Unlock()

	c.values = make(map[K]V)
	c.entries = make(map[K]*list.Element)
	c.order.Init()
}

// storeLoaded stores a loaded value and evicts the least recently used value when needed.
func (c *cache[K, V]) storeLoaded(key K, value V) {
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

// markUsed moves a cache key to the most recently used position.
func (c *cache[K, V]) markUsed(key K) {
	if entry, found := c.entries[key]; found {
		c.order.MoveToBack(entry)
	}
}

// evictOldest removes the least recently used cache entry.
func (c *cache[K, V]) evictOldest() {
	oldestEntry := c.order.Front()
	if oldestEntry == nil {
		return
	}

	oldestKey := oldestEntry.Value.(K)
	c.order.Remove(oldestEntry)
	delete(c.entries, oldestKey)
	delete(c.values, oldestKey)
}
