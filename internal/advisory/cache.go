package advisory

import (
	"sync"

	"example.com/sbom-risk-graph/internal/domain"
)

type SnapshotCache struct {
	mu         sync.RWMutex
	generation uint64
	items      []domain.Advisory
}

func (c *SnapshotCache) Replace(items []domain.Advisory) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generation++
	c.items = append([]domain.Advisory(nil), items...)
	return c.generation
}

func (c *SnapshotCache) Snapshot() (uint64, []domain.Advisory) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.generation, append([]domain.Advisory(nil), c.items...)
}
