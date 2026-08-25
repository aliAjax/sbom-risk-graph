package parser

import (
	"encoding/json"
	"example.com/sbom-risk-graph/internal/domain"
	"example.com/sbom-risk-graph/pkg/canon"
	"example.com/sbom-risk-graph/pkg/hashing"
	"fmt"
	"strings"
	"time"
)

type Raw struct {
	SPDXID  string `json:"SPDXID"`
	Name    string `json:"name"`
	Version string `json:"versionInfo"`
	License string `json:"licenseConcluded"`
	PURL    string `json:"purl"`
}
type RawDocument struct {
	SPDXID        string `json:"SPDXID"`
	Name          string `json:"name"`
	Components    []Raw  `json:"packages"`
	Relationships []struct {
		From string `json:"spdxElementId"`
		To   string `json:"relatedSpdxElement"`
		Type string `json:"relationshipType"`
	} `json:"relationships"`
	BOMRef            string `json:"bom-ref"`
	CycloneComponents []struct {
		BOMRef  string `json:"bom-ref"`
		Name    string `json:"name"`
		Version string `json:"version"`
		PURL    string `json:"purl"`
		License string `json:"license"`
	} `json:"components"`
	Dependencies []struct {
		Ref       string   `json:"ref"`
		DependsOn []string `json:"dependsOn"`
	} `json:"dependencies"`
}

func Parse(data []byte, format, productID, id string) (domain.SBOM, error) {
	if len(data) > 64<<20 {
		return domain.SBOM{}, fmt.Errorf("SBOM exceeds 64MiB")
	}
	var raw RawDocument
	if err := json.Unmarshal(data, &raw); err != nil {
		return domain.SBOM{}, fmt.Errorf("decode %s: %w", format, err)
	}
	s := domain.SBOM{ID: id, ProductID: productID, Format: strings.ToLower(format), Serial: raw.SPDXID, Status: domain.Draft, ImportedAt: time.Now().UTC()}
	if s.Format == "spdx" {
		for _, c := range raw.Components {
			s.Components = append(s.Components, domain.Component{ID: c.SPDXID, PURL: c.PURL, Name: c.Name, Version: c.Version, License: c.License})
		}
		for _, r := range raw.Relationships {
			if r.Type == "DEPENDS_ON" {
				s.Dependencies = append(s.Dependencies, domain.Dependency{From: r.From, To: r.To})
			}
		}
	} else if s.Format == "cyclonedx" {
		for _, c := range raw.CycloneComponents {
			s.Components = append(s.Components, domain.Component{ID: c.BOMRef, PURL: c.PURL, Name: c.Name, Version: c.Version, License: c.License})
		}
		for _, d := range raw.Dependencies {
			for _, to := range d.DependsOn {
				s.Dependencies = append(s.Dependencies, domain.Dependency{From: d.Ref, To: to})
			}
		}
	} else {
		return domain.SBOM{}, fmt.Errorf("unsupported format")
	}
	if err := s.Validate(); err != nil {
		return domain.SBOM{}, err
	}
	canonical, err := canon.JSON(s)
	if err != nil {
		return domain.SBOM{}, err
	}
	s.Digest = hashing.Hex(canonical)
	if len(s.Components) > 0 {
		s.Status = domain.Draft
	} else {
		s.Status = domain.Failed
	}
	return s, nil
}
func Canonical(value any) ([]byte, string, error) {
	b, err := canon.JSON(value)
	if err != nil {
		return nil, "", err
	}
	return b, hashing.Hex(b), nil
}
