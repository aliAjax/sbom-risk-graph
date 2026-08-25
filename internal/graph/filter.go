package graph

import "example.com/sbom-risk-graph/internal/domain"

func FilterComponents(values []domain.Component, ecosystem, license string) []domain.Component {
	out := make([]domain.Component, 0, len(values))
	for _, value := range values {
		if ecosystem != "" && value.Ecosystem != ecosystem {
			continue
		}
		if license != "" && value.License != license {
			continue
		}
		out = append(out, value)
	}
	return out
}
