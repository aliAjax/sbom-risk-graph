package parser

import (
	"encoding/json"
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
)

type CycloneDocument struct {
	BomFormat    string `json:"bomFormat"`
	SpecVersion  string `json:"specVersion"`
	SerialNumber string `json:"serialNumber"`
	Components   []struct {
		BOMRef   string `json:"bom-ref"`
		Name     string `json:"name"`
		Version  string `json:"version"`
		PURL     string `json:"purl"`
		Licenses []struct {
			License struct {
				Name string `json:"name"`
			} `json:"license"`
		} `json:"licenses"`
	} `json:"components"`
	Dependencies []struct {
		Ref       string   `json:"ref"`
		DependsOn []string `json:"dependsOn"`
	} `json:"dependencies"`
}

func ParseCyclone(data []byte, id, product string) (domain.SBOM, error) {
	var d CycloneDocument
	if err := json.Unmarshal(data, &d); err != nil {
		return domain.SBOM{}, err
	}
	if d.BomFormat != "CycloneDX" && d.BomFormat != "cyclonedx" {
		return domain.SBOM{}, fmt.Errorf("not CycloneDX")
	}
	s := domain.SBOM{ID: id, ProductID: product, Format: "cyclonedx", Serial: d.SerialNumber, Status: domain.Failed}
	for _, c := range d.Components {
		license := ""
		if len(c.Licenses) > 0 {
			license = c.Licenses[0].License.Name
		}
		s.Components = append(s.Components, domain.Component{ID: c.BOMRef, Name: c.Name, Version: c.Version, PURL: c.PURL, License: license})
	}
	for _, d := range d.Dependencies {
		for _, to := range d.DependsOn {
			s.Dependencies = append(s.Dependencies, domain.Dependency{From: d.Ref, To: to})
		}
	}
	s = NormalizeSBOM(s)
	return s, s.Validate()
}
