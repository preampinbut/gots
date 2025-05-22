package util

import "go/ast"

// StructInfo holds a parsed struct's name and fields.
type StructInfo struct {
	Name   string
	Fields []*ast.Field
}

type AliasInfo struct {
	Name       string
	Underlying ast.Expr
}
