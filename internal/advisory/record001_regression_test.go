package advisory

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/sbom-risk-graph/internal/domain"
)

func advisoryItem(id string) domain.Advisory {
	return domain.Advisory{ID: id, ComponentPURL: "pkg:golang/example/module", AffectedRange: ">=1.0.0 <1.1.0"}
}

func TestSnapshotCacheOwnsReplacementInput(t *testing.T) {
	cache := &SnapshotCache{}
	input := []domain.Advisory{advisoryItem("ADV-1")}
	cache.Replace(input)
	input[0].ID = "changed"
	_, got := cache.Snapshot()
	if got[0].ID != "ADV-1" {
		t.Fatalf("caller mutated stored snapshot: %q", got[0].ID)
	}
}

func TestSnapshotCacheReturnsDetachedSnapshot(t *testing.T) {
	cache := &SnapshotCache{}
	cache.Replace([]domain.Advisory{advisoryItem("ADV-2")})
	_, first := cache.Snapshot()
	first[0].ID = "changed"
	_, second := cache.Snapshot()
	if second[0].ID != "ADV-2" {
		t.Fatalf("returned slice escaped cache: %q", second[0].ID)
	}
}

type loaderFunc func(context.Context) ([]domain.Advisory, error)

func (f loaderFunc) Load(ctx context.Context) ([]domain.Advisory, error) { return f(ctx) }

func TestRefreshPropagatesCancellationToLoader(t *testing.T) {
	cache := &SnapshotCache{}
	started := make(chan struct{})
	seenCancel := make(chan struct{})
	loader := loaderFunc(func(ctx context.Context) ([]domain.Advisory, error) {
		close(started)
		select {
		case <-ctx.Done():
			close(seenCancel)
			return nil, ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return nil, errors.New("loader did not receive cancellation")
		}
	})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- (Refresher{Cache: cache, Loader: loader}).Refresh(ctx)
	}()
	<-started
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled refresh, got %v", err)
	}
	select {
	case <-seenCancel:
	default:
		t.Fatal("loader did not observe caller cancellation")
	}
}

func TestRefreshFailureKeepsPreviousSnapshot(t *testing.T) {
	cache := &SnapshotCache{}
	cache.Replace([]domain.Advisory{advisoryItem("old")})
	loaderErr := errors.New("feed interrupted")
	loader := loaderFunc(func(context.Context) ([]domain.Advisory, error) {
		return []domain.Advisory{advisoryItem("partial")}, loaderErr
	})
	if err := (Refresher{Cache: cache, Loader: loader}).Refresh(context.Background()); !errors.Is(err, loaderErr) {
		t.Fatalf("expected loader error, got %v", err)
	}
	generation, got := cache.Snapshot()
	if generation != 1 || len(got) != 1 || got[0].ID != "old" {
		t.Fatalf("failed refresh published partial batch: generation=%d items=%v", generation, got)
	}
}

func TestRefreshCanceledAfterLoadDoesNotPublish(t *testing.T) {
	cache := &SnapshotCache{}
	cache.Replace([]domain.Advisory{advisoryItem("old")})
	ctx, cancel := context.WithCancel(context.Background())
	loader := loaderFunc(func(context.Context) ([]domain.Advisory, error) {
		cancel()
		return []domain.Advisory{advisoryItem("new")}, nil
	})
	if err := (Refresher{Cache: cache, Loader: loader}).Refresh(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation after load, got %v", err)
	}
	generation, got := cache.Snapshot()
	if generation != 1 || len(got) != 1 || got[0].ID != "old" {
		t.Fatalf("canceled refresh advanced cache: generation=%d items=%v", generation, got)
	}
}
