package util

import "go/ast"

func ResolveAlias(name string, aliases map[string]ast.Expr) ast.Expr {
	visited := map[string]bool{}
	curr := name
	for {
		if visited[curr] {
			break // avoid infinite loop
		}
		visited[curr] = true

		alias, ok := aliases[curr]
		if !ok {
			break
		}

		ident, ok := alias.(*ast.Ident)
		if !ok {
			// alias to complex type - stop here
			return alias
		}

		curr = ident.Name
	}

	// Return simple ident for the last resolved type name
	return &ast.Ident{Name: curr}
}
