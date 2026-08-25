package evidence

import (
	"context"
	"errors"
	"testing"

	"example.com/sbom-risk-graph/internal/domain"
)

func registryEvidence(status domain.Status) domain.Evidence {
	return domain.Evidence{ID: "ev-1", ProductID: "product-1", Kind: "sbom", Digest: "digest", Status: status}
}

func TestRevokedEvidenceCannotBecomeValid(t *testing.T) {
	err := ValidateTransition(domain.Revoked, domain.Valid)
	if !errors.Is(err, ErrInvalidEvidenceTransition) {
		t.Fatalf("ValidateTransition error = %v, want ErrInvalidEvidenceTransition", err)
	}
}

func TestRejectedTransitionDoesNotMutateRegistry(t *testing.T) {
	registry := NewRegistry([]domain.Evidence{registryEvidence(domain.Valid)})
	if err := registry.Transition("ev-1", domain.Draft); !errors.Is(err, ErrInvalidEvidenceTransition) {
		t.Fatalf("Transition error = %v, want ErrInvalidEvidenceTransition", err)
	}
	item, ok := registry.Get("ev-1")
	if !ok {
		t.Fatal("evidence disappeared after rejected transition")
	}
	if item.Status != domain.Valid {
		t.Fatalf("status after rejected transition = %q, want %q", item.Status, domain.Valid)
	}
}

type rejectedTransitionRegistry struct{ err error }

func (r rejectedTransitionRegistry) Transition(string, domain.Status) error { return r.err }

func TestPublisherPropagatesTransitionFailure(t *testing.T) {
	publisher := NewPublisher(rejectedTransitionRegistry{err: ErrInvalidEvidenceTransition})
	err := publisher.Publish(context.Background(), "ev-1", domain.Valid)
	if !errors.Is(err, ErrInvalidEvidenceTransition) {
		t.Fatalf("Publish error = %v, want ErrInvalidEvidenceTransition", err)
	}
}

func TestRegistrySnapshotIsDetached(t *testing.T) {
	registry := NewRegistry([]domain.Evidence{registryEvidence(domain.Valid)})
	snapshot := registry.Snapshot()
	changed := snapshot["ev-1"]
	changed.Status = domain.Revoked
	snapshot["ev-1"] = changed
	delete(snapshot, "missing")

	item, ok := registry.Get("ev-1")
	if !ok {
		t.Fatal("evidence disappeared after snapshot mutation")
	}
	if item.Status != domain.Valid {
		t.Fatalf("registry status = %q, want %q", item.Status, domain.Valid)
	}
}
