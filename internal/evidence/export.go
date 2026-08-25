package evidence

import (
	"encoding/json"
	"example.com/sbom-risk-graph/internal/domain"
)

type Package struct {
	Evidence domain.Evidence     `json:"evidence"`
	Audit    []domain.AuditEvent `json:"audit"`
	Root     string              `json:"root"`
}

func Export(e domain.Evidence, a []domain.AuditEvent) []byte {
	ids := make([]string, 0, len(a))
	for _, v := range a {
		ids = append(ids, v.Hash)
	}
	b, _ := json.Marshal(Package{Evidence: e, Audit: a, Root: Root(ids)})
	return b
}
