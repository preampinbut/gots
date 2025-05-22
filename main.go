package main

import (
	"flag"
	"fmt"
	"github.com/preampinbut/gots/util"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	// Define flags
	outFile := flag.String("output", "", "write output to a file instead of stdout")
	flag.StringVar(outFile, "o", "", "shorthand for --output")

	verbose := flag.Bool("verbose", false, "enable verbose output")
	flag.BoolVar(verbose, "v", false, "shorthand for --verbose")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage: gots [-o --output output.ts] [-v --verbose] <input.go>")
	}
	inputFile := flag.Arg(0)

	ts, err := TranspileFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	if *outFile != "" {
		if *verbose {
			log.Printf("📄 Transpiling file: %s\n", inputFile)
		}

		err := os.WriteFile(*outFile, []byte(ts), 0644)
		if err != nil {
			log.Fatalf("failed to write output to %s: %v", *outFile, err)
		}
		if *verbose {
			log.Printf("✅ Output written to: %s\n", *outFile)
		}
	} else {
		fmt.Println(ts)
	}
}

func TranspileFile(filename string) (string, error) {
	fset := token.NewFileSet()
	dir := filepath.Dir(filename)

	// Step 1: Parse all structs in the directory
	allStructs := make(map[string]*ast.StructType)
	allAliases := make(map[string]ast.Expr)

	pkgs, err := parser.ParseDir(fset, dir, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec := spec.(*ast.TypeSpec)

					switch t := typeSpec.Type.(type) {
					case *ast.StructType:
						allStructs[typeSpec.Name.Name] = t
					default:
						allAliases[typeSpec.Name.Name] = t
					}
				}
			}
		}
		break
	}

	// Step 2: Parse the input file to find used top-level structs
	inputFileAst, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}

	used := map[string]bool{}
	var queue []string

	// collect top-level structs from input file
	for _, decl := range inputFileAst.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec := spec.(*ast.TypeSpec)
			name := typeSpec.Name.Name
			if _, found := allStructs[name]; found {
				used[name] = true
				queue = append(queue, name)
			}
		}
	}

	// Step 3: Recursively find referenced struct types
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		st := allStructs[name]
		for _, field := range st.Fields.List {
			t := field.Type
			ref := util.ReferencedIdentName(t)
			if ref != "" && allStructs[ref] != nil && !used[ref] {
				used[ref] = true
				queue = append(queue, ref)
			}
		}
	}

	var names []string
	for name := range used {
		names = append(names, name)
	}
	sort.Strings(names)

	// Step 4: Generate TS interfaces
	var structs []util.StructInfo
	for _, name := range names {
		structs = append(structs, util.StructInfo{
			Name:   name,
			Fields: allStructs[name].Fields.List,
		})
	}

	return util.RenderInterfaces(structs, allAliases), nil
}
