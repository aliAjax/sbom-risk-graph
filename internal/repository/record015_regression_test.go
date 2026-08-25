package repository

import "testing"

func TestZeroValueTenantGuardDoesNotPanic(t *testing.T) {
	var g TenantGuard
	if err := g.Bind("product", "tenant-a"); err != nil {
		t.Fatal(err)
	}
	if !g.Allow("product", "tenant-a") {
		t.Fatal("binding lost")
	}
}
func TestConflictingBindRollsBackOwnership(t *testing.T) {
	g := NewTenantGuard()
	_ = g.Bind("product", "tenant-a")
	if err := g.Bind("product", "tenant-b"); err == nil {
		t.Fatal("conflict accepted")
	}
	if !g.Allow("product", "tenant-a") || g.Allow("product", "tenant-b") {
		t.Fatal("ownership changed")
	}
}
func TestTenantAuthorizationAfterRecovery(t *testing.T) {
	g := NewTenantService(NewTenantGuard())
	_ = g.BindProduct("product", "tenant-a")
	if err := g.RecoverProduct("product", "tenant-a", "tenant-c"); err != nil {
		t.Fatal(err)
	}
	if !g.Allowed("product", "tenant-c") || g.Allowed("product", "tenant-a") {
		t.Fatal("authorization did not follow recovery")
	}
}
