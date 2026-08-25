package httpapi

import (
	"encoding/json"
	"example.com/sbom-risk-graph/internal/domain"
	"example.com/sbom-risk-graph/internal/observability"
	"example.com/sbom-risk-graph/internal/parser"
	"example.com/sbom-risk-graph/internal/repository"
	"example.com/sbom-risk-graph/pkg/metrics"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	service *repository.Service
	metrics *metrics.Counters
	tracer  *observability.Tracer
}

func New(s *repository.Service, m *metrics.Counters) *Handler {
	return &Handler{service: s, metrics: m, tracer: observability.NewTracer()}
}
func (h *Handler) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", h.health)
	m.HandleFunc("/readyz", h.health)
	m.HandleFunc("/metrics", h.metricsHandler)
	m.HandleFunc("/api/v1/products", h.products)
	m.HandleFunc("/api/v1/products/", h.product)
	m.HandleFunc("/api/v1/sboms", h.sboms)
	m.HandleFunc("/api/v1/advisories", h.advisories)
	m.HandleFunc("/api/v1/policies", h.policies)
	m.HandleFunc("/api/v1/graph/paths", h.paths)
	return Middleware(m)
}
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	write(w, 200, map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
}
func (h *Handler) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	_, _ = w.Write([]byte(metrics.Render(h.metrics)))
}

type productReq struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
}

func (h *Handler) products(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, 200, domain.Page[domain.Product]{Items: h.service.Repo.Products(), Count: len(h.service.Repo.Products())})
		return
	}
	var req productReq
	if err := decode(r, &req); err != nil {
		h.fail(w, r, 400, "INVALID_ARGUMENT", err.Error())
		return
	}
	p := domain.Product{ID: req.ID, TenantID: req.TenantID, Name: req.Name, Version: req.Version, CreatedAt: time.Now().UTC()}
	if err := h.service.AddProduct(p); err != nil {
		h.domainError(w, r, err)
		return
	}
	write(w, 201, p)
}
func (h *Handler) product(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	if id == "" {
		h.fail(w, r, 404, "NOT_FOUND", "product not found")
		return
	}
	if strings.HasSuffix(id, "/risk") {
		id = strings.TrimSuffix(id, "/risk")
		p, ok := h.service.Repo.Product(id)
		if !ok {
			h.fail(w, r, 404, "NOT_FOUND", "product not found")
			return
		}
		_ = p
		h.metrics.Evaluations.Add(1)
		write(w, 200, h.service.Risk(id))
		return
	}
	p, ok := h.service.Repo.Product(id)
	if !ok {
		h.fail(w, r, 404, "NOT_FOUND", "product not found")
		return
	}
	write(w, 200, p)
}

type sbomReq struct {
	ID        string          `json:"id"`
	ProductID string          `json:"product_id"`
	Format    string          `json:"format"`
	Document  json.RawMessage `json:"document"`
}

func (h *Handler) sboms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.fail(w, r, 405, "METHOD_NOT_ALLOWED", "POST required")
		return
	}
	var req sbomReq
	if err := decode(r, &req); err != nil {
		h.fail(w, r, 400, "INVALID_ARGUMENT", err.Error())
		return
	}
	b, err := parser.Parse(req.Document, req.Format, req.ProductID, req.ID)
	if err != nil {
		h.metrics.ParserErrors.Add(1)
		h.domainError(w, r, err)
		return
	}
	if err := h.service.AddSBOM(b); err != nil {
		h.domainError(w, r, err)
		return
	}
	h.metrics.Imports.Add(1)
	write(w, 201, b)
}
func (h *Handler) advisories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.fail(w, r, 405, "METHOD_NOT_ALLOWED", "POST required")
		return
	}
	var a domain.Advisory
	if err := decode(r, &a); err != nil {
		h.fail(w, r, 400, "INVALID_ARGUMENT", err.Error())
		return
	}
	if err := h.service.AddAdvisory(a); err != nil {
		h.domainError(w, r, err)
		return
	}
	write(w, 201, a)
}
func (h *Handler) policies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.fail(w, r, 405, "METHOD_NOT_ALLOWED", "POST required")
		return
	}
	var p domain.Policy
	if err := decode(r, &p); err != nil {
		h.fail(w, r, 400, "INVALID_ARGUMENT", err.Error())
		return
	}
	if err := h.service.AddPolicy(p); err != nil {
		h.domainError(w, r, err)
		return
	}
	write(w, 201, p)
}
func (h *Handler) paths(w http.ResponseWriter, r *http.Request) {
	product := r.URL.Query().Get("product_id")
	root := r.URL.Query().Get("root")
	g, ok := h.service.Graph.Graph(product)
	if !ok {
		h.fail(w, r, 404, "NOT_FOUND", "graph not found")
		return
	}
	write(w, 200, map[string]any{"paths": g.PathsFrom(root, 32)})
}
func decode(r *http.Request, target any) error {
	defer r.Body.Close()
	data, err := io.ReadAll(io.LimitReader(r.Body, 64<<20+1))
	if err != nil {
		return err
	}
	if len(data) > 64<<20 {
		return fmt.Errorf("request exceeds 64MiB")
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	return dec.Decode(target)
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errorResponse struct{ Code, Message, RequestID string }

func (h *Handler) fail(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	write(w, status, errorResponse{code, message, requestID(r.Context())})
}
func (h *Handler) domainError(w http.ResponseWriter, r *http.Request, err error) {
	h.fail(w, r, 400, "INVALID_ARGUMENT", err.Error())
}
