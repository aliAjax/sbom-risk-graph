package importer

import "example.com/sbom-risk-graph/internal/domain"

func (b *Batch) Snapshot() []domain.SBOM {
	b.mu.Lock()
	defer b.mu.Unlock()
	return cloneSBOMs(b.current)
}

func cloneSBOMs(items []domain.SBOM) []domain.SBOM {
	cloned := make([]domain.SBOM, len(items))
	for i, item := range items {
		cloned[i] = cloneSBOM(item)
	}
	return cloned
}

func cloneSBOM(item domain.SBOM) domain.SBOM {
	item.Components = append([]domain.Component(nil), item.Components...)
	item.Dependencies = append([]domain.Dependency(nil), item.Dependencies...)
	return item
}
