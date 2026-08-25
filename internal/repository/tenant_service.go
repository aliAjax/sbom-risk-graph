package repository

type TenantService struct{ guard *TenantGuard }

func NewTenantService(g *TenantGuard) *TenantService { return &TenantService{guard: g} }
func (s *TenantService) BindProduct(product, tenant string) error {
	return s.guard.Bind(product, tenant)
}
func (s *TenantService) RecoverProduct(product, oldTenant, newTenant string) error {
	return s.guard.Bind(product, newTenant)
}
func (s *TenantService) Allowed(product, tenant string) bool { return s.guard.Allow(product, tenant) }
