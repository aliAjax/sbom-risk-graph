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
	return &Batch{current: cloneSBOMs(current)}
}

func (b *Batch) Stage(ctx context.Context, item domain.SBOM) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := item.Validate(); err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	b.staged = append(b.staged, cloneSBOM(item))
	return nil
}

func (b *Batch) Fail(err error) {
	if err == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failed = err
}

func (b *Batch) Commit(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if b.failed != nil {
		return fmt.Errorf("batch import failed: %w", b.failed)
	}
	b.current = cloneSBOMs(b.staged)
	b.staged = nil
	return nil
}
