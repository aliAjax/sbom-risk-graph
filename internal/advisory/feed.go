package advisory

import (
	"encoding/json"
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
)

func ParseFeed(data []byte) ([]domain.Advisory, error) {
	var items []domain.Advisory
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	if len(items) > 100000 {
		return nil, fmt.Errorf("feed too large")
	}
	for _, a := range items {
		if err := a.Validate(); err != nil {
			return nil, err
		}
	}
	return items, nil
}
