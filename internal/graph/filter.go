package graph

import "example.com/sbom-risk-graph/internal/domain"

func FilterComponents(values []domain.Component, ecosystem, license string) []domain.Component {
	out := values[:0]
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
