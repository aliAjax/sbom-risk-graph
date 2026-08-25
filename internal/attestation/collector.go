package attestation

import (
	"fmt"
	"sort"
	"sync"
)

type Collector struct {
	mu        sync.RWMutex
	staged    map[string]Attestation
	committed map[string]Attestation
}

func NewCollector() *Collector {
	return &Collector{staged: make(map[string]Attestation), committed: make(map[string]Attestation)}
}

func (c *Collector) Stage(item Attestation) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.staged[item.ID] = item
}

func (c *Collector) Commit() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.staged) == 0 {
		return fmt.Errorf("no staged attestations")
	}
	for id, item := range c.staged {
		c.committed[id] = item
		delete(c.staged, id)
	}
	return nil
}

func (c *Collector) Abort() {
	c.mu.Lock()
	defer c.mu.Unlock()
	clear(c.staged)
}

func (c *Collector) IDs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ids := make([]string, 0, len(c.committed))
	for id := range c.committed {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
