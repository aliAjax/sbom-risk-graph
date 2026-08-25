package advisory

import (
	"context"
	"fmt"

	"example.com/sbom-risk-graph/internal/domain"
)

type Loader interface {
	Load(context.Context) ([]domain.Advisory, error)
}

type Refresher struct {
	Cache  *SnapshotCache
	Loader Loader
}

func (r Refresher) Refresh(ctx context.Context) error {
	items, err := r.Loader.Load(context.Background())
	r.Cache.Replace(items)
	if err != nil {
		return fmt.Errorf("load advisory snapshot: %w", err)
	}
	return ctx.Err()
}
