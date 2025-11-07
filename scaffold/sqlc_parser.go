package scaffold

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SQLCStruct struct {
	Name   string
	Fields map[string]string // field -> type
}

type SQLCFunction struct {
	Name    string
	Params  []string
	Returns []string
	Doc     string
}

// EntityMeta — структура, которую будем передавать в шаблоны
type EntityMeta struct {
	EntityName string // "Role", "Country", etc.

	// boolean flags for CRUD
	HasCreate     bool
	HasList       bool
	HasGetByID    bool
	HasGetByName  bool
	HasUpdateByID bool
	HasDeleteByID bool
	HasCount      bool

	// names of sqlc-generated Params/Structs (if any)
	CreateParamsName string // e.g. "CreateCountryParams" (may be empty)
	ListParamsName   string // e.g. "ListCountriesParams"
	UpdateParamsName string // e.g. "UpdateRoleByIDParams"
	ModelStructName  string // e.g. "Role"

	// Raw AST-discovered info (useful for subsequent generation)
	AllStructs     []SQLCStruct
	AllFunctions   []SQLCFunction
	OtherFunctions []SQLCFunction
}

//
// Utilities
//

func exprString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprString(t.X)
	case *ast.ArrayType:
		return "[]" + exprString(t.Elt)
	case *ast.MapType:
		return "map[" + exprString(t.Key) + "]" + exprString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return fmt.Sprintf("%T", e)
	}
}

func isExported(name string) bool {
	return name != "" && strings.ToUpper(name[:1]) == name[:1]
}

func isBasicType(t string) bool {
	basics := []string{"string", "int", "int32", "int64", "bool", "float32", "float64", "[]byte", "error", "context.Context"}
	for _, b := range basics {
		if t == b || strings.HasPrefix(t, b) {
			return true
		}
	}
	return false
}

func structExists(list []SQLCStruct, name string) bool {
	for _, s := range list {
		if s.Name == name {
			return true
		}
	}
	return false
}

// Parse a single sqlc-generated file and build an EntityMeta
func ParseSQLCFile(path string) (*EntityMeta, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	meta := &EntityMeta{}
	meta.AllStructs = []SQLCStruct{}
	meta.AllFunctions = []SQLCFunction{}
	meta.OtherFunctions = []SQLCFunction{}

	// collect all type structs
	for _, decl := range node.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						s := SQLCStruct{
							Name:   ts.Name.Name,
							Fields: map[string]string{},
						}
						for _, f := range st.Fields.List {
							ftype := exprString(f.Type)
							if len(f.Names) == 0 {
								s.Fields["<anon>"] = ftype
							} else {
								for _, nm := range f.Names {
									s.Fields[nm.Name] = ftype
								}
							}
						}
						meta.AllStructs = append(meta.AllStructs, s)
					}
				}
			}
		}
	}

	// heuristic: try to detect main model struct (e.g., Role, Country)
	probModel := ""
	for _, st := range meta.AllStructs {
		if isExported(st.Name) && !strings.HasSuffix(st.Name, "Params") && !strings.HasSuffix(st.Name, "Row") && !strings.Contains(st.Name, "Select") {
			if matched, _ := regexp.MatchString(`^[A-Z][a-zA-Z0-9]+$`, st.Name); matched {
				probModel = st.Name
				break
			}
		}
	}
	meta.ModelStructName = probModel

	// parse functions (methods on *Queries)
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil {
			continue
		}
		recv := fn.Recv.List[0].Type
		recvStr := exprString(recv)
		// filter methods on Queries
		if !strings.Contains(recvStr, "Queries") {
			continue
		}

		// extract parameter types
		params := []string{}
		if fn.Type.Params != nil {
			for _, p := range fn.Type.Params.List {
				typ := exprString(p.Type)
				params = append(params, typ)
			}
		}
		// extract return types
		returns := []string{}
		if fn.Type.Results != nil {
			for _, r := range fn.Type.Results.List {
				returns = append(returns, exprString(r.Type))
			}
		}
		doc := ""
		if fn.Doc != nil {
			doc = fn.Doc.Text()
		}

		f := SQLCFunction{
			Name:    fn.Name.Name,
			Params:  params,
			Returns: returns,
			Doc:     doc,
		}
		meta.AllFunctions = append(meta.AllFunctions, f)

		// detect CRUD patterns by function name
		switch {
		case strings.HasPrefix(f.Name, "Create"):
			meta.HasCreate = true
			// detect CreateParamsName if param is some *XxxParams
			for i := 1; i < len(f.Params); i++ {
				t := f.Params[i]
				if !isBasicType(t) && (strings.HasSuffix(t, "Params") || strings.HasSuffix(t, "Param")) {
					meta.CreateParamsName = strings.TrimPrefix(t, "*")
					break
				}
			}
		case strings.HasPrefix(f.Name, "List"):
			meta.HasList = true
			for i := 1; i < len(f.Params); i++ {
				t := f.Params[i]
				if strings.Contains(t, "List") || strings.HasSuffix(t, "Params") {
					meta.ListParamsName = strings.TrimPrefix(t, "*")
					break
				}
			}
		case strings.HasPrefix(f.Name, "Get") && strings.HasSuffix(f.Name, "ByID"):
			meta.HasGetByID = true
		case strings.HasPrefix(f.Name, "Get") && strings.Contains(f.Name, "ByName"):
			meta.HasGetByName = true
		case strings.HasPrefix(f.Name, "Update") && strings.HasSuffix(f.Name, "ByID"):
			meta.HasUpdateByID = true
			for i := 1; i < len(f.Params); i++ {
				t := f.Params[i]
				if strings.Contains(t, "Update") || strings.HasSuffix(t, "Params") {
					meta.UpdateParamsName = strings.TrimPrefix(t, "*")
					break
				}
			}
		case strings.HasPrefix(f.Name, "Delete") && strings.HasSuffix(f.Name, "ByID"):
			meta.HasDeleteByID = true
		case strings.HasPrefix(f.Name, "Count"):
			meta.HasCount = true
		default:
			meta.OtherFunctions = append(meta.OtherFunctions, f)
		}
	}

	// if model unknown try to detect from returns like *Role
	if meta.ModelStructName == "" {
		for _, f := range meta.AllFunctions {
			for _, r := range f.Returns {
				if strings.HasPrefix(r, "*") {
					nm := strings.TrimPrefix(r, "*")
					if structExists(meta.AllStructs, nm) {
						meta.ModelStructName = nm
						break
					}
				}
			}
			if meta.ModelStructName != "" {
				break
			}
		}
	}

	return meta, nil
}

// Walk directory and parse all .sql.go
func ParseAllSQLC(folder string) ([]*EntityMeta, error) {
	var metas []*EntityMeta
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".sql.go") {
			m, err := ParseSQLCFile(path)
			if err != nil {
				return fmt.Errorf("parse %s: %w", path, err)
			}
			metas = append(metas, m)
		}
		return nil
	})
	return metas, err
}

// Filter: return entities that look like a full CRUD reference data (heuristic)
func FilterFullCrudEntities(metas []*EntityMeta) []*EntityMeta {
	var res []*EntityMeta
	for _, m := range metas {
		// Basic rule: must have Create, List, GetByID, UpdateByID, DeleteByID
		if m.HasCreate && m.HasList && m.HasGetByID && m.HasUpdateByID && m.HasDeleteByID {
			// Ensure model name exists or we can still bind (we accept missing model if at least params exist)
			res = append(res, m)
		}
	}
	return res
}

// CLI entrypoint for parser
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: go run scaffold/sqlc_parser.go <storage/dbs path>")
	}
	path := os.Args[1]
	metas, err := ParseAllSQLC(path)
	if err != nil {
		log.Fatalf("parse error: %v", err)
	}
	fmt.Printf("Parsed %d .sql.go files\n\n", len(metas))

	full := FilterFullCrudEntities(metas)
	fmt.Printf("Detected %d full-CRUD entities (heuristic)\n\n", len(full))

	for _, m := range full {
		fmt.Println("==== ENTITY ====")
		fmt.Printf("EntityName: %s\n", m.ModelStructName)
		if m.ModelStructName == "" {
			// fallback: try to infer name from functions (CreateXxx -> Xxx)
			for _, f := range m.AllFunctions {
				if strings.HasPrefix(f.Name, "Create") {
					mName := strings.TrimPrefix(f.Name, "Create")
					m.ModelStructName = mName
					break
				}
			}
		}
		fmt.Printf("ModelStructName: %s\n", m.ModelStructName)
		fmt.Printf("HasCreate: %v CreateParams: %s\n", m.HasCreate, m.CreateParamsName)
		fmt.Printf("HasList: %v ListParams: %s\n", m.HasList, m.ListParamsName)
		fmt.Printf("HasGetByID: %v\n", m.HasGetByID)
		fmt.Printf("HasGetByName: %v\n", m.HasGetByName)
		fmt.Printf("HasUpdateByID: %v UpdateParams: %s\n", m.HasUpdateByID, m.UpdateParamsName)
		fmt.Printf("HasDeleteByID: %v\n", m.HasDeleteByID)
		fmt.Printf("HasCount: %v\n", m.HasCount)
		fmt.Println("Other functions (sample):")
		for _, f := range m.OtherFunctions {
			fmt.Printf(" - %s params=%v returns=%v\n", f.Name, f.Params, f.Returns)
		}
		fmt.Println()
	}
}
