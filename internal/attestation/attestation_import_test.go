package attestation

import (
	"context"
	"errors"
	"testing"
	"time"
)

type sourceFunc func(context.Context) (Attestation, error)

func (f sourceFunc) Next(ctx context.Context) (Attestation, error) { return f(ctx) }

func TestPipelineFollowsParentCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	pipeline := NewPipeline(parent)
	cancel()
	select {
	case <-pipeline.Done():
	case <-time.After(200 * time.Millisecond):
		pipeline.Close()
		t.Fatal("pipeline did not observe parent cancellation")
	}
}

func TestSourceReceivesCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := readNext(ctx, sourceFunc(func(received context.Context) (Attestation, error) {
		return Attestation{}, received.Err()
	}))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("source did not receive canceled context: %v", err)
	}
}

func TestAbortDiscardsStagedAttestation(t *testing.T) {
	collector := NewCollector()
	collector.Stage(Attestation{ID: "att-12", ProductID: "product-8", Digest: "sha256:abc"})
	collector.Abort()
	if ids := collector.IDs(); len(ids) != 0 {
		t.Fatalf("aborted attestation was committed: %v", ids)
	}
}

func TestImportErrorPreservesCancellation(t *testing.T) {
	err := wrapImportError(context.Canceled)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("import error lost cancellation identity: %v", err)
	}
}
