package importer

import "example.com/sbom-risk-graph/internal/domain"

func (b *Batch) Snapshot() []domain.SBOM {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.current
}

func cloneSBOMs(items []domain.SBOM) []domain.SBOM {
	return append([]domain.SBOM(nil), items...)
}

func cloneSBOM(item domain.SBOM) domain.SBOM {
	return item
}
