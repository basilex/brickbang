package main

import (
	"bytes"
	"fmt"

	// "go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"brickbang/scaffold"
)

// Template file names (in scaffold/templates)
var tplFiles = []string{
	"module.tpl",
	"mapper.tpl",
	"service.tpl",
	"repository.tpl",
	"controller.tpl",
}

// Target paths relative to project root
var outPaths = map[string]string{
	"module.tpl":     "internal/module/%s_module.go",
	"mapper.tpl":     "internal/mapper/%s_mapper.go",
	"service.tpl":    "internal/service/%s_service.go",
	"repository.tpl": "internal/repository/%s_repository.go",
	"controller.tpl": "internal/controller/%s_controller.go",
}

type TemplateData struct {
	ModulePath string

	// Entity naming
	Entity         string // singular, e.g. Role, Country
	EntityLower    string // lower-case e.g. role, country
	EntityPlural   string // naive plural e.g. Roles, Countries (Entity + "s")
	EntityPluralLo string // lowercase plural

	// Flags and param struct names discovered by parser
	HasCreate     bool
	HasList       bool
	HasGetByID    bool
	HasGetByName  bool
	HasUpdateByID bool
	HasDeleteByID bool
	HasCount      bool

	CreateParams string // e.g. CreateRoleParams (may be "")
	ListParams   string // e.g. ListRolesParams
	UpdateParams string // e.g. UpdateRoleByIDParams

	ModelStruct string // e.g. Role

	// Raw meta in case template needs more info
	Meta *scaffold.EntityMeta
}

func ensureDirForFile(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0o755)
}

func renderTemplateFile(tplPath string, outPath string, data TemplateData) error {
	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", tplPath, err)
	}

	tpl, err := template.New(filepath.Base(tplPath)).Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).Parse(string(tplBytes))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tplPath, err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template %s: %w", tplPath, err)
	}

	if err := ensureDirForFile(outPath); err != nil {
		return err
	}

	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write output %s: %w", outPath, err)
	}
	return nil
}

func naivePlural(s string) string {
	// very naive pluralization: if ends with 'y' -> replace y->ies else add 's'
	if s == "" {
		return ""
	}
	last := s[len(s)-1]
	if last == 'y' || last == 'Y' {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

func generateFromMeta(meta *scaffold.EntityMeta, modulePath string) error {
	// choose entity name: prefer ModelStructName, fallback to infer from functions
	entity := meta.ModelStructName
	if entity == "" {
		// try to pick name from Create function name
		for _, f := range meta.AllFunctions {
			if strings.HasPrefix(f.Name, "Create") {
				entity = strings.TrimPrefix(f.Name, "Create")
				break
			}
		}
	}
	if entity == "" {
		return fmt.Errorf("cannot determine entity name for meta: %#v", meta)
	}

	// prepare template data
	td := TemplateData{
		ModulePath:     modulePath,
		Entity:         entity,
		EntityLower:    strings.ToLower(entity),
		EntityPlural:   naivePlural(entity),
		EntityPluralLo: strings.ToLower(naivePlural(entity)),
		HasCreate:      meta.HasCreate,
		HasList:        meta.HasList,
		HasGetByID:     meta.HasGetByID,
		HasGetByName:   meta.HasGetByName,
		HasUpdateByID:  meta.HasUpdateByID,
		HasDeleteByID:  meta.HasDeleteByID,
		HasCount:       meta.HasCount,
		CreateParams:   meta.CreateParamsName,
		ListParams:     meta.ListParamsName,
		UpdateParams:   meta.UpdateParamsName,
		ModelStruct:    meta.ModelStructName,
		Meta:           meta,
	}

	// loop over template files
	for _, tpl := range tplFiles {
		tplPath := filepath.Join("scaffold", "templates", tpl)
		outPattern, ok := outPaths[tpl]
		if !ok {
			return fmt.Errorf("no out path for template %s", tpl)
		}
		outPath := fmt.Sprintf(outPattern, td.EntityLower)
		if err := renderTemplateFile(tplPath, outPath, td); err != nil {
			return fmt.Errorf("render %s -> %s: %w", tplPath, outPath, err)
		}
		log.Printf("Generated: %s\n", outPath)
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run scaffold/gen_from_sqlc.go <storage/dbs path> <module/path>")
		fmt.Println("Example: go run scaffold/gen_from_sqlc.go ./storage/dbs brickbang")
		os.Exit(1)
	}
	dbsPath := os.Args[1]
	modulePath := os.Args[2]

	metas, err := scaffold.ParseAllSQLC(dbsPath)
	if err != nil {
		log.Fatalf("parse sqlc: %v", err)
	}

	full := scaffold.FilterFullCrudEntities(metas)
	if len(full) == 0 {
		log.Fatalf("no full-CRUD entities found")
	}

	for _, m := range full {
		if err := generateFromMeta(m, modulePath); err != nil {
			log.Fatalf("generate for %s: %v", m.ModelStructName, err)
		}
	}

	log.Printf("Done: generated %d entities\n", len(full))
}
