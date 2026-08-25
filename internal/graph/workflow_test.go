package graph

import (
	"errors"
	"testing"

	"example.com/sbom-risk-graph/internal/domain"
)

func TestGraphWorkflow(t *testing.T) {
	invalid := &domain.Graph{
		Nodes: map[string]domain.Component{},
		Edges: map[string][]string{"missing": {"target"}},
	}

	t.Run("validation remains classifiable", func(t *testing.T) {
		if err := Validate(invalid); !errors.Is(err, ErrInvalidGraph) {
			t.Fatalf("validation error was not preserved: %v", err)
		}
	})

	t.Run("checked expand rejects nil graph", func(t *testing.T) {
		if _, err := ExpandChecked(nil, "root", 2, 10); !errors.Is(err, ErrInvalidGraph) {
			t.Fatalf("expand error was not preserved: %v", err)
		}
	})

	t.Run("checked snapshot rejects dangling edge", func(t *testing.T) {
		if _, err := BuildSnapshotChecked("product-1", invalid); !errors.Is(err, ErrInvalidGraph) {
			t.Fatalf("snapshot error was not preserved: %v", err)
		}
	})

	t.Run("failed replacement is atomic", func(t *testing.T) {
		store := New()
		original := domain.SBOM{
			ID:        "sbom-1",
			ProductID: "product-1",
			Format:    "spdx",
			Components: []domain.Component{{
				ID: "component-1", PURL: "pkg:golang/original", Name: "original", Version: "1.0.0",
			}},
		}
		if err := store.PutSBOM(original); err != nil {
			t.Fatal(err)
		}
		invalidReplacement := original
		invalidReplacement.Components = []domain.Component{{
			ID: "component-2", PURL: "pkg:golang/broken", Version: "2.0.0",
		}}
		if err := store.PutSBOM(invalidReplacement); err == nil {
			t.Fatal("invalid replacement unexpectedly succeeded")
		}
		got, ok := store.GetSBOM(original.ID)
		if !ok || got.Components[0].Name != "original" {
			t.Fatalf("failed replacement changed committed sbom: %#v", got)
		}
	})
}
