package evidence

import (
	"example.com/sbom-risk-graph/internal/domain"
	"reflect"
	"testing"
)

func TestAuditSnapshotUnaffectedByAppend(t *testing.T) {
	a := &Audit{events: make([]domain.AuditEvent, 1, 4)}
	a.events[0] = domain.AuditEvent{ID: "first", Hash: "h1"}
	snapshot := a.All()
	if cap(snapshot) != len(snapshot) {
		t.Fatalf("published snapshot retains writable audit capacity: len=%d cap=%d", len(snapshot), cap(snapshot))
	}
	a.Append("second", "create", "subject")
	if snapshot[0].ID != "first" {
		t.Fatalf("append changed published snapshot: %#v", snapshot)
	}
}

func TestExportDoesNotReorderAuditSnapshot(t *testing.T) {
	events := []domain.AuditEvent{{ID: "z", Hash: "z"}, {ID: "a", Hash: "a"}}
	want := append([]domain.AuditEvent(nil), events...)
	_ = Export(domain.Evidence{ID: "e1"}, events)
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("export reordered caller events: got=%v want=%v", events, want)
	}
}

func TestMerkleRootDoesNotMutateInput(t *testing.T) {
	values := []string{"z", "a", "m"}
	want := append([]string(nil), values...)
	_ = Root(values)
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("root calculation reordered input: got=%v want=%v", values, want)
	}
}

func TestProofDoesNotShareCallerStorage(t *testing.T) {
	values := []string{"a", "b", "c"}
	want := append([]string(nil), values...)
	proof, ok := Proof(values, "b")
	if !ok {
		t.Fatal("proof missing")
	}
	proof[0] = "changed"
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("proof shares caller storage: got=%v want=%v", values, want)
	}
}

func TestAuditSealDoesNotShareChainState(t *testing.T) {
	first := domain.AuditEvent{ID: "one", Action: "create", Subject: "a"}.Seal()
	second := domain.AuditEvent{ID: "two", Action: "create", Subject: "b"}.Seal()
	if first.PreviousHash != "" || second.PreviousHash != "" {
		t.Fatalf("independent seals inherited chain state: first=%q second=%q", first.PreviousHash, second.PreviousHash)
	}
}
