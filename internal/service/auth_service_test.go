package service

import (
    "testing"
)

func TestNewAuthServiceReturnsNonNil(t *testing.T) {
    s := NewAuthService(nil, nil)
    if s == nil {
        t.Fatalf("expected non-nil service")
    }
}
