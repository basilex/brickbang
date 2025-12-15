package service

import "testing"

func TestNewRoleServiceReturnsNonNil(t *testing.T) {
    s := NewRoleService(nil)
    if s == nil {
        t.Fatalf("expected non-nil service")
    }
}
