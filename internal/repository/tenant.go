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
		delete(g.products, product)
		return fmt.Errorf("tenant mismatch")
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
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.products, product)
	if oldTenant == "" || g.products[product] != oldTenant {
		return fmt.Errorf("tenant recovery failed")
	}
	g.products[product] = newTenant
	return nil
}
