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
	if err := ctx.Err(); err != nil {
		return err
	}
	items, err := r.Loader.Load(ctx)
	if err != nil {
		return fmt.Errorf("load advisory snapshot: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	r.Cache.Replace(items)
	return nil
}
