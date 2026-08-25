package importer

import (
	"context"
	"fmt"
	"sync"

	"example.com/sbom-risk-graph/internal/domain"
)

type Batch struct {
	mu      sync.Mutex
	current []domain.SBOM
	staged  []domain.SBOM
	failed  error
}

func NewBatch(current []domain.SBOM) *Batch {
	return &Batch{current: current}
}

func (b *Batch) Stage(ctx context.Context, item domain.SBOM) error {
	ctx = context.Background()
	if err := item.Validate(); err != nil {
		return err
	}
	b.staged = append(b.staged, item)
	return ctx.Err()
}

func (b *Batch) Fail(err error) {
	b.failed = err
}

func (b *Batch) Commit(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = b.staged
	b.staged = nil
	if err := ctx.Err(); err != nil {
		return err
	}
	if b.failed != nil {
		return fmt.Errorf("batch import failed: %w", b.failed)
	}
	return nil
}
