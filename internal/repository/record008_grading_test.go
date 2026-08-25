package repository

import (
	"errors"
	"example.com/sbom-risk-graph/internal/domain"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func failPersistence(t *testing.T, store *Store) {
	t.Helper()
	store.path = filepath.Join(t.TempDir(), "missing", "state.json")
}

func TestProductPersistFailureRollsBackState(t *testing.T) {
	store, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	failPersistence(t, store)
	err = store.SaveProduct(domain.Product{ID: "product-a", TenantID: "tenant-a", Name: "widget"})
	if err == nil || !errors.Is(err, fs.ErrNotExist) || !strings.Contains(err.Error(), "save product") {
		t.Fatalf("unexpected persistence error: %v", err)
	}
	if _, ok := store.Snapshot().Products["product-a"]; ok {
		t.Fatal("failed product write remained visible in memory")
	}
}

func TestSBOMPersistFailureRollsBackState(t *testing.T) {
	store, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	failPersistence(t, store)
	err = store.SaveSBOM(domain.SBOM{ID: "sbom-a", ProductID: "product-a"})
	if err == nil || !errors.Is(err, fs.ErrNotExist) || !strings.Contains(err.Error(), "save sbom") {
		t.Fatalf("unexpected persistence error: %v", err)
	}
	if _, ok := store.Snapshot().SBOMs["sbom-a"]; ok {
		t.Fatal("failed sbom write remained visible in memory")
	}
}

func TestAdvisoryPersistFailureDoesNotPublishIndex(t *testing.T) {
	store, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	failPersistence(t, store)
	service := NewService(store)
	advisory := domain.Advisory{ID: "ADV-1", ComponentPURL: "pkg:generic/lib", AffectedRange: ">=1.0.0 <2.0.0"}
	if err := service.AddAdvisory(advisory); err == nil {
		t.Fatal("advisory write unexpectedly succeeded")
	}
	component := domain.Component{PURL: advisory.ComponentPURL, Name: "lib", Version: "1.5.0"}
	if matches := service.Advisory.Match(component); len(matches) != 0 {
		t.Fatalf("failed advisory was published to matcher: %#v", matches)
	}
}

func TestPolicyPersistFailureDoesNotPublishEvaluatorState(t *testing.T) {
	store, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	failPersistence(t, store)
	service := NewService(store)
	policy := domain.Policy{ID: "policy-a", Name: "strict", MaxSeverity: 1}
	if err := service.AddPolicy(policy); err == nil {
		t.Fatal("policy write unexpectedly succeeded")
	}
	if _, ok := service.Policies.Get("policy-a"); ok {
		t.Fatal("failed policy was published to evaluator store")
	}
}
