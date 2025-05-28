// Copyright 2025 ovs-cni authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"time"
)

// Cache is the ovs bridges cache
type Cache struct {
	lastRefreshTime time.Time
	bridges         map[string]bool
}

// Refresh updates the cached bridges and refresh time
func (c *Cache) Refresh(freshBridges map[string]bool) {
	c.bridges = freshBridges
	c.lastRefreshTime = time.Now()
}

// LastRefreshTime returns the last time the cache was updated
func (c *Cache) LastRefreshTime() time.Time {
	return c.lastRefreshTime
}

// Bridges return the cached bridges
func (c Cache) Bridges() map[string]bool {
	bridgesCopy := make(map[string]bool)
	for bridge, exist := range c.bridges {
		bridgesCopy[bridge] = exist
	}
	return bridgesCopy
}
