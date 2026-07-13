/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package fragmentation

import "sync"

// cache is a limited collection of recently used values by key.
type cache[K comparable, V any] struct {
	// Mutex guards all cache state.
	sync.Mutex

	// limit is the maximum number of retained values.
	limit int

	// loader resolves values that are not cached.
	loader func(K) (V, error)

	// values contains cached values indexed by key.
	values map[K]V

	// order tracks keys from least to most recently used.
	order []K
}

// newCache creates a cache with a loader and least-recently-used eviction.
func newCache[K comparable, V any](limit int, loader func(K) (V, error)) *cache[K, V] {
	return &cache[K, V]{
		limit:  limit,
		loader: loader,
		values: make(map[K]V),
	}
}

// get returns a cached value or loads it when missing.
//
//nolint:ireturn // The cache is generic, so returning V preserves the stored value type.
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

// storeLoaded stores a loaded value and evicts the least recently used value when needed.
func (c *cache[K, V]) storeLoaded(key K, value V) {
	if _, found := c.values[key]; found {
		c.values[key] = value
		c.markUsed(key)

		return
	}

	c.values[key] = value
	c.order = append(c.order, key)
	if len(c.values) <= c.limit {
		return
	}

	c.evictOldest()
}

// markUsed moves a cache key to the most recently used position.
func (c *cache[K, V]) markUsed(key K) {
	for index, existingKey := range c.order {
		if existingKey == key {
			c.order = append(c.order[:index], c.order[index+1:]...)
			c.order = append(c.order, key)

			return
		}
	}
}

// evictOldest removes the least recently used cache entry.
func (c *cache[K, V]) evictOldest() {
	oldestKey := c.order[0]
	c.order = c.order[1:]
	delete(c.values, oldestKey)
}
