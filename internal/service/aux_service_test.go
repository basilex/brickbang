package service

import (
	"strings"
	"testing"
	"time"
)

func TestHealth(t *testing.T) {
	s := NewAuxService(map[string]string{"k": "v"})
	got := s.Health()
	if got["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", got)
	}
}

func TestMetadata(t *testing.T) {
	meta := map[string]string{"version": "1.2.3"}
	s := NewAuxService(meta)
	got := s.Metadata()
	if got["version"] != "1.2.3" {
		t.Fatalf("expected metadata preserved, got %v", got)
	}
}

func TestUptime(t *testing.T) {
	// construct auxService with deterministic start time in the past
	as := &auxService{metadata: map[string]string{"a": "b"}, startTime: time.Now().Add(-2 * time.Hour)}
	got := as.Uptime()
	if uptime, ok := got["uptime"]; !ok || !strings.Contains(uptime, "h") {
		t.Fatalf("expected uptime to contain hours, got %v", got)
	}
}
