package util

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

func RenderInterfaces(structs []StructInfo, allAliases map[string]ast.Expr) string {
	var b strings.Builder
	for _, s := range structs {
		b.WriteString(fmt.Sprintf("interface %s {\n", s.Name))
		for _, field := range s.Fields {
			var name string
			optional := ""
			if len(field.Names) > 0 {
				name = field.Names[0].Name
			} else {
				continue // embedded/anonymous field with no name
			}
			tsType := exprToTS(field.Type, allAliases)

			// Check for JSON tag
			if field.Tag != nil {
				tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
				if jsonTag := tag.Get("json"); jsonTag != "" {
					parts := strings.Split(jsonTag, ",")
					if parts[0] != "" && parts[0] != "-" {
						name = parts[0]
					}
					for _, p := range parts[1:] {
						if p == "omitempty" {
							optional = "?"
						}
					}
				}
			}

			b.WriteString(fmt.Sprintf("  %s%s: %s;\n", name, optional, tsType))
		}
		b.WriteString("}\n\n")
	}
	return strings.TrimSpace(b.String())
}

func ReferencedIdentName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return ReferencedIdentName(t.X)
	case *ast.ArrayType:
		return ReferencedIdentName(t.Elt)
	case *ast.SelectorExpr:
		return t.Sel.Name // e.g., time.Time → "Time"
	default:
		return ""
	}
}

func exprToTS(expr ast.Expr, aliases map[string]ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// resolve alias first
		resolved := ResolveAlias(t.Name, aliases)
		// if resolved is ident, map basic
		if ident, ok := resolved.(*ast.Ident); ok {
			return mapBasicGoToTS(ident.Name)
		}
		// otherwise fallback to any (or handle more cases)
		return "any"
	case *ast.ArrayType:
		return exprToTS(t.Elt, aliases) + "[]"
	case *ast.StarExpr:
		return exprToTS(t.X, aliases) // dereference pointer
	case *ast.SelectorExpr:
		return t.Sel.Name // handle imported type like `time.Time`
	case *ast.StructType:
		// Anonymous inline struct
		return inlineStructToTS(t, aliases)
	default:
		return "any"
	}
}

func inlineStructToTS(st *ast.StructType, aliases map[string]ast.Expr) string {
	var b strings.Builder
	b.WriteString("{ ")
	for _, field := range st.Fields.List {
		name := field.Names[0].Name
		fieldType := exprToTS(field.Type, aliases)
		b.WriteString(fmt.Sprintf("%s: %s; ", name, fieldType))
	}
	b.WriteString("}")
	return b.String()
}

func mapBasicGoToTS(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "interface{}":
		return "any"
	default:
		return goType // assume it's a user-defined type or imported
	}
}
