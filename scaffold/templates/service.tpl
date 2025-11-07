package service

import (
    "context"

    "{{ .ModulePath }}/storage/dbs"
)

// Интерфейс по конвенции I...
type I{{ .Entity }}Service interface {
    {{ if .HasCreate }}Create(ctx context.Context, params dbs.{{ .CreateParams }}) (*dbs.{{ .Entity }}, error){{ end }}
    {{ if .HasList }}List(ctx context.Context, params dbs.{{ .ListParams }}) ([]*dbs.{{ .Entity }}, error){{ end }}
    {{ if .HasGetByID }}GetByID(ctx context.Context, id string) (*dbs.{{ .Entity }}, error){{ end }}
    {{ if .HasGetByName }}GetByName(ctx context.Context, name string) (*dbs.{{ .Entity }}, error){{ end }}
    {{ if .HasUpdateByID }}UpdateByID(ctx context.Context, params dbs.{{ .UpdateParams }}) (*dbs.{{ .Entity }}, error){{ end }}
    {{ if .HasDeleteByID }}DeleteByID(ctx context.Context, id string) error{{ end }}
    {{ if .HasCount }}Count(ctx context.Context) (int64, error){{ end }}
}

// Реализация сервиса
type {{ .Entity }}Service struct {
    q *dbs.Queries
}

func New{{ .Entity }}Service(q *dbs.Queries) I{{ .Entity }}Service {
    return &{{ .Entity }}Service{q: q}
}

{{ if .HasCreate }}
func (s *{{ .Entity }}Service) Create(ctx context.Context, params dbs.{{ .CreateParams }}) (*dbs.{{ .Entity }}, error) {
    return s.q.Create{{ .Entity }}(ctx, params)
}
{{ end }}

{{ if .HasList }}
func (s *{{ .Entity }}Service) List(ctx context.Context, params dbs.{{ .ListParams }}) ([]*dbs.{{ .Entity }}, error) {
    return s.q.List{{ .Entity }}s(ctx, params)
}
{{ end }}

{{ if .HasGetByID }}
func (s *{{ .Entity }}Service) GetByID(ctx context.Context, id string) (*dbs.{{ .Entity }}, error) {
    return s.q.Get{{ .Entity }}ByID(ctx, id)
}
{{ end }}

{{ if .HasGetByName }}
func (s *{{ .Entity }}Service) GetByName(ctx context.Context, name string) (*dbs.{{ .Entity }}, error) {
    return s.q.Get{{ .Entity }}ByName(ctx, name)
}
{{ end }}

{{ if .HasUpdateByID }}
func (s *{{ .Entity }}Service) UpdateByID(ctx context.Context, params dbs.{{ .UpdateParams }}) (*dbs.{{ .Entity }}, error) {
    return s.q.Update{{ .Entity }}ByID(ctx, params)
}
{{ end }}

{{ if .HasDeleteByID }}
func (s *{{ .Entity }}Service) DeleteByID(ctx context.Context, id string) error {
    return s.q.Delete{{ .Entity }}ByID(ctx, id)
}
{{ end }}

{{ if .HasCount }}
func (s *{{ .Entity }}Service) Count(ctx context.Context) (int64, error) {
    return s.q.Count{{ .Entity }}(ctx)
}
{{ end }}
