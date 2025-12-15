package service

import "testing"

func TestNewGrantServiceReturnsNonNil(t *testing.T) {
    s := NewGrantService(nil)
    if s == nil {
        t.Fatalf("expected non-nil service")
    }
}
