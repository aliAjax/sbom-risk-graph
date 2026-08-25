package repository

import (
	"fmt"
	"sync"
)

type TenantGuard struct {
	mu       sync.RWMutex
	products map[string]string
}

func NewTenantGuard() *TenantGuard { return &TenantGuard{products: make(map[string]string)} }
func (g *TenantGuard) Bind(product, tenant string) error {
	if product == "" || tenant == "" {
		return fmt.Errorf("product and tenant required")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if old, ok := g.products[product]; ok && old != tenant {
		return fmt.Errorf("tenant mismatch")
	}
	if g.products == nil {
		g.products = make(map[string]string)
	}
	g.products[product] = tenant
	return nil
}
func (g *TenantGuard) Allow(product, tenant string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.products[product] == tenant
}

func (g *TenantGuard) Rebind(product, oldTenant, newTenant string) error {
	if product == "" || newTenant == "" || oldTenant == "" {
		return fmt.Errorf("product and tenant required")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.products[product] != oldTenant {
		return fmt.Errorf("tenant recovery failed")
	}
	if g.products == nil {
		g.products = make(map[string]string)
	}
	g.products[product] = newTenant
	return nil
}
