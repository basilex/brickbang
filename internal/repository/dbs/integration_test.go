//go:build integration
// +build integration

package dbs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Integration test using Testcontainers: starts Postgres, runs migrations from storage/schema,
// and executes a simple query using sqlc-generated Queries. This test is tagged `integration`
// and won't run during regular `go test ./...` unless you add `-tags=integration`.

func TestPostgresIntegration_SimpleRoleCreate(t *testing.T) {
	ctx := context.Background()

	req := tc.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_USER":     "brb",
			"POSTGRES_DB":       "brb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	pgc, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil {
		t.Fatalf("starting container: %v", err)
	}
	defer pgc.Terminate(ctx)

	host, _ := pgc.Host(ctx)
	port, _ := pgc.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://brb:pass@%s:%s/brb?sslmode=disable", host, port.Port())

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("pgx connect: %v", err)
	}
	defer conn.Close(ctx)

	// run migration SQL files found in storage/schema/*.up.sql
	schemaDir := filepath.Join("..", "..", "..", "storage", "schema")
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		t.Fatalf("read schema dir: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		if !filepath.Ext(e.Name()) == ".sql" {
			continue
		}
		// only process .up.sql
		if filepath.Ext(e.Name()) == ".sql" && filepath.Base(e.Name()) != e.Name() {
			// noop
		}
		if filepath.Base(e.Name()) == e.Name() {
			// noop
		}
		if filepath.Ext(e.Name()) == ".sql" && filepath.HasSuffix(e.Name(), ".up.sql") {
			b, err := os.ReadFile(filepath.Join(schemaDir, e.Name()))
			if err != nil {
				t.Fatalf("read migration: %v", err)
			}
			if _, err := conn.Exec(ctx, string(b)); err != nil {
				t.Fatalf("exec migration %s: %v", e.Name(), err)
			}
		}
	}

	// Create Queries wrapper and run a simple count (assuming roles table exists)
	q := New(conn)
	_, err = q.CountRoles(ctx)
	if err != nil {
		t.Fatalf("CountRoles failed: %v", err)
	}
}
