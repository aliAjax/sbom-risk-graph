package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Status string

const (
	Draft      Status = "draft"
	Validating Status = "validating"
	Valid      Status = "valid"
	Failed     Status = "failed"
	Revoked    Status = "revoked"
	Expired    Status = "expired"
)

type Product struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Revision  uint64    `json:"revision"`
}

func (p Product) Validate() error {
	if p.ID == "" || p.TenantID == "" || strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("product id tenant_id and name are required")
	}
	return nil
}

type Component struct {
	ID        string `json:"id"`
	PURL      string `json:"purl"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Ecosystem string `json:"ecosystem"`
	License   string `json:"license,omitempty"`
}

func (c Component) Key() string { return c.PURL + "@" + c.Version }
func (c Component) Validate() error {
	if c.Name == "" || c.Version == "" {
		return fmt.Errorf("component name and version are required")
	}
	if len(c.Name) > 512 || len(c.Version) > 256 {
		return fmt.Errorf("component field too long")
	}
	return nil
}

type Dependency struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Optional bool   `json:"optional"`
}
type SBOM struct {
	ID           string       `json:"id"`
	ProductID    string       `json:"product_id"`
	Format       string       `json:"format"`
	Serial       string       `json:"serial"`
	Components   []Component  `json:"components"`
	Dependencies []Dependency `json:"dependencies"`
	Digest       string       `json:"digest"`
	Status       Status       `json:"status"`
	ImportedAt   time.Time    `json:"imported_at"`
	Revision     uint64       `json:"revision"`
}

func (s SBOM) Validate() error {
	if s.ID == "" || s.ProductID == "" {
		return fmt.Errorf("sbom id and product_id are required")
	}
	if s.Format != "spdx" && s.Format != "cyclonedx" {
		return fmt.Errorf("format must be spdx or cyclonedx")
	}
	if len(s.Components) > 100000 {
		return fmt.Errorf("too many components")
	}
	return nil
}

type Advisory struct {
	ID            string    `json:"id"`
	ComponentPURL string    `json:"component_purl"`
	AffectedRange string    `json:"affected_range"`
	Severity      float64   `json:"severity"`
	EPSS          float64   `json:"epss"`
	FixedVersion  string    `json:"fixed_version,omitempty"`
	PublishedAt   time.Time `json:"published_at"`
	Withdrawn     bool      `json:"withdrawn"`
}

func (a Advisory) Validate() error {
	if a.ID == "" || a.ComponentPURL == "" || a.AffectedRange == "" {
		return fmt.Errorf("advisory fields are required")
	}
	if a.Severity < 0 || a.Severity > 10 {
		return fmt.Errorf("severity must be 0..10")
	}
	return nil
}

type Policy struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	MaxSeverity     float64  `json:"max_severity"`
	MaxEPSS         float64  `json:"max_epss"`
	DeniedLicenses  []string `json:"denied_licenses"`
	RequireEvidence bool     `json:"require_evidence"`
	Revision        uint64   `json:"revision"`
}

func (p Policy) Validate() error {
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("policy id and name are required")
	}
	return nil
}

func (p Policy) NextRevision(current Policy) (Policy, error) {
	if current.ID != "" && current.ID != p.ID {
		return Policy{}, fmt.Errorf("policy revision belongs to another policy")
	}
	if p.Revision == 0 {
		p.Revision = current.Revision + 1
		return p, nil
	}
	if current.ID != "" && p.Revision <= current.Revision {
		return Policy{}, fmt.Errorf("policy revision %d does not advance current %d", p.Revision, current.Revision)
	}
	return p, nil
}

type Violation struct {
	PolicyID   string  `json:"policy_id"`
	Component  string  `json:"component"`
	AdvisoryID string  `json:"advisory_id,omitempty"`
	Code       string  `json:"code"`
	Message    string  `json:"message"`
	Severity   float64 `json:"severity"`
}
type RiskReport struct {
	ProductID       string      `json:"product_id"`
	GeneratedAt     time.Time   `json:"generated_at"`
	Score           float64     `json:"score"`
	Components      int         `json:"components"`
	Vulnerabilities int         `json:"vulnerabilities"`
	Violations      []Violation `json:"violations"`
	Paths           [][]string  `json:"paths,omitempty"`
}

func (r *RiskReport) Sort() {
	sort.Slice(r.Violations, func(i, j int) bool {
		if r.Violations[i].Severity == r.Violations[j].Severity {
			return r.Violations[i].Component < r.Violations[j].Component
		}
		return r.Violations[i].Severity > r.Violations[j].Severity
	})
}

type Evidence struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"product_id"`
	Kind          string    `json:"kind"`
	Digest        string    `json:"digest"`
	PayloadDigest string    `json:"payload_digest"`
	Signature     string    `json:"signature"`
	Status        Status    `json:"status"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}
type Job struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Type      string    `json:"type"`
	Status    Status    `json:"status"`
	Progress  int       `json:"progress"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
