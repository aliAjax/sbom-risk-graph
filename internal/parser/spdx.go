package parser

import (
	"encoding/json"
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
)

type SPDXDocument struct {
	SPDXVersion string `json:"spdxVersion"`
	SPDXID      string `json:"SPDXID"`
	Packages    []struct {
		SPDXID           string `json:"SPDXID"`
		Name             string `json:"name"`
		Version          string `json:"versionInfo"`
		DownloadLocation string `json:"downloadLocation"`
		License          string `json:"licenseConcluded"`
	} `json:"packages"`
	Relationships []struct {
		From string `json:"spdxElementId"`
		To   string `json:"relatedSpdxElement"`
		Type string `json:"relationshipType"`
	} `json:"relationships"`
}

func ParseSPDX(data []byte, id, product string) (domain.SBOM, error) {
	var d SPDXDocument
	if err := json.Unmarshal(data, &d); err != nil {
		return domain.SBOM{}, err
	}
	if d.SPDXVersion == "" {
		return domain.SBOM{}, fmt.Errorf("missing SPDX version")
	}
	s := domain.SBOM{ID: id, ProductID: product, Format: "spdx", Serial: d.SPDXID, Status: domain.Draft}
	for _, p := range d.Packages {
		s.Components = append(s.Components, domain.Component{ID: p.SPDXID, Name: p.Name, Version: p.Version, PURL: p.DownloadLocation, License: p.License})
	}
	for _, r := range d.Relationships {
		if r.Type == "DEPENDS_ON" {
			s.Dependencies = append(s.Dependencies, domain.Dependency{From: r.From, To: r.To})
		}
	}
	s = NormalizeSBOM(s)
	return s, s.Validate()
}
